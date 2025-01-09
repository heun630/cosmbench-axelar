package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

/*
** 수정해야 할 configuration 항목
1. encodedTxDir - encoded txs 디렉토리 이름 및 위치
2. HOSTS - 노드 IP (노드는 4대 기준)
3. REST_PORTS - 노드의 REST API 포트 (노드는 4대 기준)
*/
var encodedTxDir = "/data/axelar/cosmbench-axelar/axelar-cosmbench_encoded_txs"

// 노드 설정
var HOSTS = []string{"127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1"}
var REST_PORTS = []string{"22200", "22201", "22202", "22203"}

var numTxs int   // 총 트랜잭션 수
var InputTPS int // 입력 TPS (초당 입력 트랜잭션 수)
var runTime int  // 실행 시간 (초)

// TxData는 트랜잭션 데이터를 저장하는 구조체
type TxData struct {
	TxBytes string `json:"tx_bytes"` // 인코딩된 트랜잭션 데이터
	Mode    string `json:"mode"`     // 브로드캐스트 모드 (ex) BROADCAST_MODE_ASYNC)
}

// readEncodedTxs는 지정된 디렉토리에서 인코딩된 트랜잭션 데이터를 읽어옵니다
func readEncodedTxs(dir string) ([]string, error) {
	// 디렉토리 내의 모든 파일 경로를 가져옵니다
	pattern := filepath.Join(dir, "*")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("파일 검색 실패: %v", err)
	}

	var txs []string

	// 각 파일의 내용을 읽어서 트랜잭션 배열에 추가합니다
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("파일 읽기 실패 (%s): %v", file, err)
		}
		// 파일 내용을 문자열로 변환하여 배열에 추가
		txs = append(txs, string(bytes.TrimSpace(content)))
	}
	numTxs = len(txs)
	return txs, nil
}

// sendTransaction은 단일 트랜잭션을 지정된 노드로 전송합니다
func sendTransaction(txIdx int, tx string, wg *sync.WaitGroup) {
	defer wg.Done()
	// tx 순서(tx sequence index)에 따라 호스트와 포트번호를 라운드로빈 방식으로 선택
	host := HOSTS[txIdx%len(HOSTS)]
	port := REST_PORTS[txIdx%len(REST_PORTS)]
	url := fmt.Sprintf("http://%s:%s/cosmos/tx/v1beta1/txs", host, port)

	// 요청 데이터 준비: 트랜잭션 바이트와 브로드캐스트 모드 설정
	requestData := TxData{
		TxBytes: tx,
		Mode:    "BROADCAST_MODE_ASYNC", // 비동기 브로드캐스트 모드 사용
	}

	// 요청 데이터를 JSON으로 변환
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Printf("[TxSequence %d, Host %s] JSON 변환 실패: %v\n", txIdx, host, err)
		return
	}

	// HTTP POST 요청 생성
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("[TxSequence %d, Host %s] 요청 생성 실패: %v\n", txIdx, host, err)
		return
	}

	// Content-Type 헤더 설정
	req.Header.Set("Content-Type", "application/json")

	// HTTP 클라이언트를 사용하여 요청 전송
	client := &http.Client{}
	println("txIdx:", txIdx, "time:", time.Now().UnixMilli()) // txIDx 코멘트: 고쳐놓기
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[TxSequence %d, Host %s] 요청 전송 실패: %v\n", txIdx, host, err)
		return
	}

	// 응답 처리
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Printf("[TxSequence %d, Host %s] 응답 읽기 실패: %v\n", txIdx, host, err)
		return
	}

	fmt.Printf("[TxSequence %d, Host %s, Port %s] 응답: %s\n", txIdx, host, port, string(body))
}

// go run send_tx.go 100 60  # TPS 100으로 60초 동안 실행
func main() {
	// 커맨드 라인에서 입력 TPS와 실행시간 값을 입력받음
	if len(os.Args) != 3 {
		fmt.Println("사용법: go run send_tx.go [TPS] [RunTime]")
		return
	}

	// TPS 값을 정수로 변환
	var err error
	InputTPS, err = strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("TPS 값이 올바르지 않습니다: %v\n", err)
		return
	}

	// 실행시간 값을 정수로 변환
	runTime, err = strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("실행시간 값이 올바르지 않습니다: %v\n", err)
		return
	}
	fmt.Printf("입력된 TPS: %d, 실행시간: %d초\n", InputTPS, runTime)

	// 트랜잭션 데이터 파일들을 읽어옴
	txs, err := readEncodedTxs(encodedTxDir)
	if err != nil {
		fmt.Printf("트랜잭션 데이터 읽기 실패: %v\n", err)
		return
	}

	fmt.Printf("총 트랜잭션 수: %d\n", numTxs)

	// WaitGroup 생성
	var wg sync.WaitGroup

	// 현재까지 전송한 트랜잭션 수를 추적
	sentTxs := 0

	for i := 0; i < runTime && sentTxs < numTxs; i++ {
		startTime := time.Now()

		// 남은 트랜잭션 수 계산
		remainingTxs := numTxs - sentTxs
		// 이번 반복에서 보낼 트랜잭션 수 결정 (InputTPS와 남은 트랜잭션 수 중 작은 값)
		txsToSend := InputTPS
		if remainingTxs < InputTPS {
			txsToSend = remainingTxs
		}

		// txsToSend 만큼의 트랜잭션 전송
		for j := 0; j < txsToSend; j++ {
			wg.Add(1)
			go sendTransaction(sentTxs+j, txs[sentTxs+j], &wg)
		}

		// 모든 goroutine이 완료될 때까지 대기
		wg.Wait()

		// 전송한 트랜잭션 수 업데이트
		sentTxs += txsToSend

		// 1초에 한 번씩 실행되도록 조절
		elapsedTime := time.Since(startTime).Milliseconds()
		if elapsedTime < 1000 {
			duration := time.Duration(1000-elapsedTime) * time.Millisecond
			println("duration time:", duration)
			time.Sleep(duration)
		}

		// 모든 트랜잭션을 전송했면 종료
		if sentTxs >= numTxs {
			break
		}
	}
	fmt.Printf("모든 트랜잭션 전송 완료 (총 %d개)\n", sentTxs)
}
