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

func parseTxLogs(filePath string) ([]TxLog, int64, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, 0, fmt.Errorf("트랜잭션 로그 파일 열기 실패: %v", err)
	}
	defer file.Close()

	var txLogs []TxLog
	var minTimestamp int64 = -1
	scanner := bufio.NewScanner(file)
	txLogRegex := regexp.MustCompile(`txIdx:\s+(\d+)\s+time:\s+(\d+)`)

	for scanner.Scan() {
		line := scanner.Text()
		match := txLogRegex.FindStringSubmatch(line)
		if len(match) > 0 {
			txIdx, _ := strconv.Atoi(match[1])
			timestamp, _ := strconv.ParseInt(match[2], 10, 64)
			txLogs = append(txLogs, TxLog{TxIdx: txIdx, Timestamp: timestamp})

			// 최소 타임스탬프 계산
			if minTimestamp == -1 || timestamp < minTimestamp {
				minTimestamp = timestamp
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, fmt.Errorf("트랜잭션 로그 파일 읽기 실패: %v", err)
	}

	if minTimestamp == -1 {
		return nil, 0, fmt.Errorf("유효한 타임스탬프를 찾을 수 없습니다")
	}

	return txLogs, minTimestamp, nil
}

func parseAndMergeBlockLogs(logDir string) ([]BlockLog, int64, error) {
	files, err := filepath.Glob(filepath.Join(logDir, "output*.log"))
	if err != nil || len(files) == 0 {
		return nil, 0, fmt.Errorf("블록 로그 파일 검색 실패: %v", err)
	}

	var blockLogs []BlockLog
	var maxTimestamp int64 = 0
	blockLogRegex := regexp.MustCompile(`(\d+)\s+.*committed state.*height=(\d+).*num_txs=(\d+)`)

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return nil, 0, fmt.Errorf("파일 열기 실패 (%s): %v", file, err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			// 컬러 코드 제거
			colorCodeRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
			cleanedLine := colorCodeRegex.ReplaceAllString(line, "")

			match := blockLogRegex.FindStringSubmatch(cleanedLine)
			if len(match) > 0 {
				timestamp, _ := strconv.ParseInt(match[1], 10, 64)
				height, _ := strconv.Atoi(match[2])
				numTxs, _ := strconv.Atoi(match[3])

				// 중복 데이터 방지
				exists := false
				for _, b := range blockLogs {
					if b.Height == height {
						exists = true
						break
					}
				}
				if !exists {
					blockLogs = append(blockLogs, BlockLog{Timestamp: timestamp, Height: height, NumTxs: numTxs})
				}

				// num_txs > 0일 경우, 최대 타임스탬프 갱신
				if numTxs > 0 && timestamp > maxTimestamp {
					maxTimestamp = timestamp
				}
			} else {
				fmt.Printf("매칭 실패 라인: %s\n", cleanedLine)
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, 0, fmt.Errorf("파일 읽기 실패 (%s): %v", file, err)
		}
	}

	if maxTimestamp == 0 {
		return blockLogs, 0, fmt.Errorf("유효한 타임스탬프를 찾을 수 없습니다")
	}

	return blockLogs, maxTimestamp, nil
}

func calculateMetrics(txLogs []TxLog, blockLogs []BlockLog, minTimestamp int64, maxTimestamp int64) (float64, float64, error) {
	if len(txLogs) == 0 || len(blockLogs) == 0 {
		return 0, 0, fmt.Errorf("로그 데이터가 부족합니다")
	}

	// Start and end time
	totalTimeSeconds := float64(maxTimestamp-minTimestamp) / 1000.0

	// TPS 계산
	totalTransactions := 0
	for _, block := range blockLogs {
		totalTransactions += block.NumTxs
	}
	tps := float64(totalTransactions) / totalTimeSeconds

	// Latency 계산
	var maxLatency float64
	for _, tx := range txLogs {
		for _, block := range blockLogs {
			if tx.Timestamp <= block.Timestamp {
				latency := float64(block.Timestamp - tx.Timestamp)
				if latency > maxLatency {
					maxLatency = latency
				}
				break
			}
		}
	}

	return maxLatency, tps, nil
}

// 블록별 트랜잭션 요약 (트랜잭션이 있는 블록만 출력)
func summarizeBlocks(blockLogs []BlockLog) string {
	var summary strings.Builder
	for _, block := range blockLogs {
		if block.NumTxs > 0 {
			summary.WriteString(fmt.Sprintf("Height %d: %d transactions\n", block.Height, block.NumTxs))
		}
	}
	return summary.String()
}

func main() {
	txLogFile := "tx_log.txt" // 트랜잭션 로그 파일
	logDir := "./"            // 블록 로그 파일이 위치한 디렉토리

	// 트랜잭션 로그 파싱
	txLogs, minTimestamp, err := parseTxLogs(txLogFile)
	if err != nil {
		fmt.Printf("트랜잭션 로그 파싱 실패: %v\n", err)
		return
	}

	// 블록 로그 병합 및 파싱
	blockLogs, maxTimestamp, err := parseAndMergeBlockLogs(logDir)
	if err != nil {
		fmt.Printf("블록 로그 병합 및 파싱 실패: %v\n", err)
		return
	}

	// Latency 및 TPS 계산
	maxLatency, tps, err := calculateMetrics(txLogs, blockLogs, minTimestamp, maxTimestamp)
	if err != nil {
		fmt.Printf("지표 계산 실패: %v\n", err)
		return
	}

	// 블록 요약
	blockSummary := summarizeBlocks(blockLogs)

	// 결과 출력
	fmt.Printf("Maximum Latency (ms): %.2f\n", maxLatency)
	fmt.Printf("Throughput (TPS): %.2f\n", tps)
	fmt.Println("\nBlock Summary:")
	fmt.Println(blockSummary)
}
