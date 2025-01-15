package main

//
//import (
//	"bufio"
//	"fmt"
//	"os"
//	"os/exec"
//	"path/filepath"
//	"regexp"
//	"strconv"
//	"strings"
//)
//
//type TxLog struct {
//	TxIdx     int
//	Timestamp int64
//}
//
//type BlockLog struct {
//	Timestamp int64
//	Height    int
//	NumTxs    int
//}
//
//func parseTxLogs(filePath string) ([]TxLog, []int64, error) {
//	file, err := os.Open(filePath)
//	if err != nil {
//		return nil, nil, fmt.Errorf("트랜잭션 로그 파일 열기 실패: %v", err)
//	}
//	defer file.Close()
//
//	var txLogs []TxLog
//	var timestamps []int64
//	scanner := bufio.NewScanner(file)
//	txLogRegex := regexp.MustCompile(`txIdx:\s+(\d+)\s+time:\s+(\d+)`)
//
//	for scanner.Scan() {
//		line := scanner.Text()
//		match := txLogRegex.FindStringSubmatch(line)
//		if len(match) > 0 {
//			txIdx, _ := strconv.Atoi(match[1])
//			timestamp, _ := strconv.ParseInt(match[2], 10, 64)
//			txLogs = append(txLogs, TxLog{TxIdx: txIdx, Timestamp: timestamp})
//			timestamps = append(timestamps, timestamp)
//		}
//	}
//
//	if err := scanner.Err(); err != nil {
//		return nil, nil, fmt.Errorf("트랜잭션 로그 파일 읽기 실패: %v", err)
//	}
//
//	if len(timestamps) == 0 {
//		return nil, nil, fmt.Errorf("유효한 타임스탬프를 찾을 수 없습니다")
//	}
//
//	return txLogs, timestamps, nil
//}
//
//func parseAndMergeBlockLogs(logDir string) ([]BlockLog, int64, error) {
//	files, err := filepath.Glob(filepath.Join(logDir, "output*.log"))
//	if err != nil || len(files) == 0 {
//		return nil, 0, fmt.Errorf("블록 로그 파일 검색 실패: %v", err)
//	}
//
//	var blockLogs []BlockLog
//	var maxTimestamp int64
//	blockLogRegex := regexp.MustCompile(`(\d+)\s+.*committed state.*height=(\d+).*num_txs=(\d+)`)
//
//	for _, file := range files {
//		f, err := os.Open(file)
//		if err != nil {
//			return nil, 0, fmt.Errorf("파일 열기 실패 (%s): %v", file, err)
//		}
//		defer f.Close()
//
//		scanner := bufio.NewScanner(f)
//		for scanner.Scan() {
//			line := scanner.Text()
//			// 컬러 코드 제거
//			colorCodeRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
//			cleanedLine := colorCodeRegex.ReplaceAllString(line, "")
//
//			match := blockLogRegex.FindStringSubmatch(cleanedLine)
//			if len(match) > 0 {
//				timestamp, _ := strconv.ParseInt(match[1], 10, 64)
//				height, _ := strconv.Atoi(match[2])
//				numTxs, _ := strconv.Atoi(match[3])
//
//				blockLogs = append(blockLogs, BlockLog{Timestamp: timestamp, Height: height, NumTxs: numTxs})
//
//				// 최대 타임스탬프 계산 (NumTxs > 0 인 경우만)
//				if numTxs > 0 && timestamp > maxTimestamp {
//					maxTimestamp = timestamp
//				}
//			}
//		}
//
//		if err := scanner.Err(); err != nil {
//			return nil, 0, fmt.Errorf("파일 읽기 실패 (%s): %v", file, err)
//		}
//	}
//
//	return blockLogs, maxTimestamp, nil
//}
//
//func calculateLatency(txLogs []TxLog, blockLogs []BlockLog) error {
//	for _, tx := range txLogs {
//		var blockTimestamp int64 = -1
//
//		// 블록 로그에서 해당 tx가 포함된 height 확인
//		for _, block := range blockLogs {
//			// CLI로 트랜잭션이 해당 블록에 있는지 확인
//			cmd := exec.Command("axelar", "query", "block", strconv.Itoa(block.Height))
//			output, err := cmd.Output()
//			if err != nil {
//				return fmt.Errorf("블록 확인 실패 (height: %d): %v", block.Height, err)
//			}
//
//			// 블록에 tx 포함 여부 확인
//			if strings.Contains(string(output), strconv.Itoa(tx.TxIdx)) {
//				blockTimestamp = block.Timestamp
//				break
//			}
//		}
//
//		if blockTimestamp == -1 {
//			fmt.Printf("트랜잭션 %d를 포함한 블록을 찾을 수 없습니다.\n", tx.TxIdx)
//			continue
//		}
//
//		// Latency 계산
//		latency := blockTimestamp - tx.Timestamp
//		fmt.Printf("Tx %d: Latency = %d ms\n", tx.TxIdx, latency)
//	}
//
//	return nil
//}
//
//func main() {
//	txLogFile := "tx_log.txt"
//	logDir := "./"
//
//	// 트랜잭션 로그 파싱
//	txLogs, _, err := parseTxLogs(txLogFile)
//	if err != nil {
//		fmt.Printf("트랜잭션 로그 파싱 실패: %v\n", err)
//		return
//	}
//
//	// 블록 로그 병합 및 파싱
//	blockLogs, _, err := parseAndMergeBlockLogs(logDir)
//	if err != nil {
//		fmt.Printf("블록 로그 병합 및 파싱 실패: %v\n", err)
//		return
//	}
//
//	// Latency 계산
//	if err := calculateLatency(txLogs, blockLogs); err != nil {
//		fmt.Printf("Latency 계산 실패: %v\n", err)
//	}
//}
