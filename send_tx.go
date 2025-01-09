package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 환경 변수 로드 함수
func loadEnv(files ...string) (map[string]string, error) {
	env := make(map[string]string)
	for _, file := range files {
		cmd := exec.Command("bash", "-c", "source "+file+" && env")
		cmd.Stderr = os.Stderr

		out, err := cmd.Output()
		if err != nil {
			return nil, fmt.Errorf("%s 로드 실패: %v", file, err)
		}

		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				env[parts[0]] = parts[1]
			}
		}
	}
	return env, nil
}

var HOSTS []string
var REST_PORTS []string
var numTxs int   // 총 트랜잭션 수
var InputTPS int // TPS
var runTime int  // 실행 시간
var encodedTxDir string
var coinUnit string
var sendAmount int

// TxData는 트랜잭션 데이터를 저장하는 구조체
type TxData struct {
	TxBytes string `json:"tx_bytes"`
	Mode    string `json:"mode"`
}

// readEncodedTxs는 디렉토리에서 인코딩된 트랜잭션 데이터를 읽어옵니다
func readEncodedTxs(dir string) ([]string, error) {
	pattern := filepath.Join(dir, "*")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("파일 검색 실패: %v", err)
	}

	var txs []string
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("파일 읽기 실패 (%s): %v", file, err)
		}
		txs = append(txs, string(bytes.TrimSpace(content)))
	}
	numTxs = len(txs)
	return txs, nil
}

// sendTransaction은 단일 트랜잭션을 지정된 노드로 전송합니다
func sendTransaction(txIdx int, tx string, wg *sync.WaitGroup) {
	defer wg.Done()
	host := HOSTS[txIdx%len(HOSTS)]
	port := REST_PORTS[txIdx%len(REST_PORTS)]
	url := fmt.Sprintf("http://%s:%s/cosmos/tx/v1beta1/txs", host, port)

	requestData := TxData{
		TxBytes: tx,
		Mode:    "BROADCAST_MODE_ASYNC",
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Printf("[TxSequence %d, Host %s] JSON 변환 실패: %v\n", txIdx, host, err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("[TxSequence %d, Host %s] 요청 생성 실패: %v\n", txIdx, host, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[TxSequence %d, Host %s] 요청 전송 실패: %v\n", txIdx, host, err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Printf("[TxSequence %d, Host %s] 응답 읽기 실패: %v\n", txIdx, host, err)
		return
	}

	fmt.Printf("[TxSequence %d, Host %s, Port %s] 응답: %s\n", txIdx, host, port, string(body))
}

func main() {
	// `run_env.sh`와 `env.sh` 파일을 로드합니다
	env, err := loadEnv("run_env.sh", "env.sh")
	if err != nil {
		fmt.Printf("환경 변수 로드 실패: %v\n", err)
		return
	}

	// 노드 수와 디렉토리 정보 동적으로 설정
	nodeCount, err := strconv.Atoi(env["NODE_COUNT"])
	if err != nil {
		fmt.Printf("NODE_COUNT 값이 올바르지 않습니다: %v\n", err)
		return
	}

	encodedTxDir = env["ENCODED_TX_ROOT_DIR"]
	coinUnit = env["UNIT"]

	// 송금 금액 설정
	sendAmount, err = strconv.Atoi(env["SEND_AMOUNT"])
	if err != nil {
		fmt.Printf("SEND_AMOUNT 값이 올바르지 않습니다: %v\n", err)
		return
	}

	// HOSTS 초기화
	for i := 0; i < nodeCount; i++ {
		HOSTS = append(HOSTS, env[fmt.Sprintf("PRIVATE_HOSTS[%d]", i)])
	}

	// REST_PORTS 초기화 (API_PORTS 사용)
	for i := 0; i < nodeCount; i++ {
		REST_PORTS = append(REST_PORTS, env[fmt.Sprintf("API_PORTS[%d]", i)])
	}

	if len(os.Args) != 3 {
		fmt.Println("사용법: go run send_tx.go [TPS] [RunTime]")
		return
	}

	InputTPS, err = strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("TPS 값이 올바르지 않습니다: %v\n", err)
		return
	}

	runTime, err = strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("실행시간 값이 올바르지 않습니다: %v\n", err)
		return
	}

	txs, err := readEncodedTxs(encodedTxDir)
	if err != nil {
		fmt.Printf("트랜잭션 데이터 읽기 실패: %v\n", err)
		return
	}

	fmt.Printf("총 트랜잭션 수: %d\n", numTxs)

	var wg sync.WaitGroup
	sentTxs := 0

	for i := 0; i < runTime && sentTxs < numTxs; i++ {
		startTime := time.Now()
		remainingTxs := numTxs - sentTxs
		txsToSend := InputTPS
		if remainingTxs < InputTPS {
			txsToSend = remainingTxs
		}

		for j := 0; j < txsToSend; j++ {
			wg.Add(1)
			go sendTransaction(sentTxs+j, txs[sentTxs+j], &wg)
		}

		wg.Wait()
		sentTxs += txsToSend

		elapsedTime := time.Since(startTime).Milliseconds()
		if elapsedTime < 1000 {
			time.Sleep(time.Duration(1000-elapsedTime) * time.Millisecond)
		}

		if sentTxs >= numTxs {
			break
		}
	}
	fmt.Printf("모든 트랜잭션 전송 완료 (총 %d개)\n", sentTxs)
}
