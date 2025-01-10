package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// TxLog는 트랜잭션 로그 데이터를 저장하는 구조체
type TxLog struct {
	TxIdx     int
	Timestamp int64
}

// BlockLog는 블록 로그 데이터를 저장하는 구조체
type BlockLog struct {
	Timestamp int64
	Height    int
	NumTxs    int
}

// 트랜잭션 로그 파싱
func parseTxLogs(filePath string) ([]TxLog, error) {
	fmt.Printf("트랜잭션 로그 파일 경로: %s\n", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("트랜잭션 로그 파일 열기 실패: %v", err)
	}
	defer file.Close()

	var txLogs []TxLog
	scanner := bufio.NewScanner(file)
	txLogRegex := regexp.MustCompile(`txIdx:\s+(\d+)\s+time:\s+(\d+)`)

	for scanner.Scan() {
		line := scanner.Text()
		match := txLogRegex.FindStringSubmatch(line)
		if len(match) > 0 {
			txIdx, _ := strconv.Atoi(match[1])
			timestamp, _ := strconv.ParseInt(match[2], 10, 64)
			txLogs = append(txLogs, TxLog{TxIdx: txIdx, Timestamp: timestamp})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("트랜잭션 로그 파일 읽기 실패: %v", err)
	}

	return txLogs, nil
}

// 블록 로그 병합 및 파싱
func parseAndMergeBlockLogs(logDir string) ([]BlockLog, error) {
	fmt.Printf("블록 로그 디렉토리 경로: %s\n", logDir) // 디렉토리 확인
	files, err := filepath.Glob(filepath.Join(logDir, "output*.log"))
	if err != nil || len(files) == 0 {
		return nil, fmt.Errorf("블록 로그 파일 검색 실패: %v", err)
	}

	fmt.Printf("발견된 블록 로그 파일: %v\n", files) // 발견된 파일 목록 출력

	var blockLogs []BlockLog
	blockLogRegex := regexp.MustCompile(`(\d+)\s+.*committed state.*height=(\d+).*num_txs=(\d+)[^0]`)

	for _, file := range files {
		fmt.Printf("파싱 중인 파일: %s\n", file) // 현재 처리 중인 파일
		f, err := os.Open(file)
		if err != nil {
			return nil, fmt.Errorf("파일 열기 실패 (%s): %v", file, err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()

			// ANSI 컬러 코드 제거
			colorCodeRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
			cleanedLine := colorCodeRegex.ReplaceAllString(line, "")

			// 매칭 수행
			match := blockLogRegex.FindStringSubmatch(cleanedLine)
			if len(match) > 0 {
				timestamp, _ := strconv.ParseInt(match[1], 10, 64)
				height, _ := strconv.Atoi(match[2])
				numTxs, _ := strconv.Atoi(match[3])
				blockLogs = append(blockLogs, BlockLog{Timestamp: timestamp, Height: height, NumTxs: numTxs})
			} else {
				fmt.Printf("매칭 실패 라인: %s\n", cleanedLine)
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("파일 읽기 실패 (%s): %v", file, err)
		}
	}

	return blockLogs, nil
}

// Latency 및 TPS 계산
func calculateMetrics(txLogs []TxLog, blockLogs []BlockLog) (float64, float64, error) {
	if len(txLogs) == 0 || len(blockLogs) == 0 {
		return 0, 0, fmt.Errorf("로그 데이터가 부족합니다")
	}

	// Start and end time
	startTime := txLogs[0].Timestamp
	endTime := blockLogs[len(blockLogs)-1].Timestamp
	totalTimeSeconds := float64(endTime-startTime) / 1000.0

	// TPS 계산
	totalTransactions := 0
	for _, block := range blockLogs {
		totalTransactions += block.NumTxs
	}
	tps := float64(totalTransactions) / totalTimeSeconds

	// Latency 계산
	var totalLatency int64
	var latencyCount int64

	for _, tx := range txLogs {
		for _, block := range blockLogs {
			if tx.Timestamp <= block.Timestamp {
				totalLatency += block.Timestamp - tx.Timestamp
				latencyCount++
				break
			}
		}
	}

	avgLatency := float64(totalLatency) / float64(latencyCount)
	return avgLatency, tps, nil
}

// 블록별 트랜잭션 요약
func summarizeBlocks(blockLogs []BlockLog) string {
	var summary strings.Builder
	for _, block := range blockLogs {
		summary.WriteString(fmt.Sprintf("Height %d: %d transactions\n", block.Height, block.NumTxs))
	}
	return summary.String()
}

// main 함수
func main() {
	txLogFile := "tx_log.txt" // 트랜잭션 로그 파일
	logDir := "./"            // 블록 로그 파일이 위치한 디렉토리

	// 트랜잭션 로그 파싱
	txLogs, err := parseTxLogs(txLogFile)
	if err != nil {
		fmt.Printf("트랜잭션 로그 파싱 실패: %v\n", err)
		return
	}
	fmt.Printf("파싱된 트랜잭션 로그: %v\n", txLogs)

	// 블록 로그 병합 및 파싱
	blockLogs, err := parseAndMergeBlockLogs(logDir)
	if err != nil {
		fmt.Printf("블록 로그 병합 및 파싱 실패: %v\n", err)
		return
	}
	fmt.Printf("파싱된 블록 로그: %v\n", blockLogs)

	// Latency 및 TPS 계산
	avgLatency, tps, err := calculateMetrics(txLogs, blockLogs)
	if err != nil {
		fmt.Printf("지표 계산 실패: %v\n", err)
		return
	}

	// 블록 요약
	blockSummary := summarizeBlocks(blockLogs)

	// 결과 출력
	fmt.Printf("Average Latency (ms): %.2f\n", avgLatency)
	fmt.Printf("Throughput (TPS): %.2f\n", tps)
	fmt.Println("\nBlock Summary:")
	fmt.Println(blockSummary)
}
