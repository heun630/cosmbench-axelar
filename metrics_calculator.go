package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// TxLog represents a transaction log entry
type TxLog struct {
	TxIdx     int
	Timestamp int64
	Height    int
}

// parseTxLogs reads tx_log.txt and extracts transaction information
func parseTxLogs(filePath string) ([]TxLog, error) {
	fmt.Println("[INFO] Parsing transaction logs...")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open tx log file: %v", err)
	}
	defer file.Close()

	var txLogs []TxLog
	txLogRegex := regexp.MustCompile(`txIdx:\s+(\d+)\s+timestamp:\s+(\d+)\s+height:\s+(\d+)`)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		matches := txLogRegex.FindStringSubmatch(line)
		if len(matches) > 0 {
			txIdx, _ := strconv.Atoi(matches[1])
			timestamp, _ := strconv.ParseInt(matches[2], 10, 64)
			height, _ := strconv.Atoi(matches[3])
			txLogs = append(txLogs, TxLog{TxIdx: txIdx, Timestamp: timestamp, Height: height})
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read tx log file: %v", err)
	}

	fmt.Printf("[INFO] Parsed %d transactions from %s\n", len(txLogs), filePath)
	return txLogs, nil
}

// parseBlockLogs reads output*.log files and extracts block information
func parseBlockLogs(logDir string) (map[int]int64, error) {
	fmt.Println("[INFO] Parsing block logs...")
	files, err := filepath.Glob(filepath.Join(logDir, "output*.log"))
	if err != nil || len(files) == 0 {
		return nil, fmt.Errorf("failed to find block log files: %v", err)
	}

	blockLogs := make(map[int]int64)
	blockLogRegex := regexp.MustCompile(`(\d+)\s+.*height=(\d+).*`)

	for _, file := range files {
		fmt.Printf("[INFO] Processing block log file: %s\n", file)
		f, err := os.Open(file)
		if err != nil {
			return nil, fmt.Errorf("failed to open block log file %s: %v", file, err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			// Remove color codes if any
			colorCodeRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
			cleanedLine := colorCodeRegex.ReplaceAllString(line, "")

			matches := blockLogRegex.FindStringSubmatch(cleanedLine)
			if len(matches) > 0 {
				timestamp, _ := strconv.ParseInt(matches[1], 10, 64)
				height, _ := strconv.Atoi(matches[2])
				blockLogs[height] = timestamp
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("failed to read block log file %s: %v", file, err)
		}
	}

	fmt.Printf("[INFO] Parsed %d block entries\n", len(blockLogs))
	return blockLogs, nil
}

// calculateLatency computes the latency for each transaction and writes to output file
func calculateLatency(txLogs []TxLog, blockLogs map[int]int64, outputFile string) error {
	fmt.Println("[INFO] Calculating latency...")
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create latency file: %v", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, tx := range txLogs {
		blockTimestamp, exists := blockLogs[tx.Height]
		if !exists {
			fmt.Printf("[WARN] No block timestamp found for height %d (TxIdx: %d)\n", tx.Height, tx.TxIdx)
			continue
		}

		latency := blockTimestamp - tx.Timestamp
		fmt.Printf("[INFO] TxIdx: %d, Latency: %d ms\n", tx.TxIdx, latency)
		fmt.Fprintf(writer, "TxIdx: %d, Timestamp: %d, Height: %d, Latency: %d ms\n", tx.TxIdx, tx.Timestamp, tx.Height, latency)
	}

	fmt.Printf("[INFO] Latency calculation completed. Results saved to %s\n", outputFile)
	return nil
}

func main() {
	txLogFile := "tx_log.txt"
	logDir := "./"
	outputFile := "latency.txt"

	fmt.Println("[INFO] Starting latency calculation tool...")

	// Parse transaction logs
	txLogs, err := parseTxLogs(txLogFile)
	if err != nil {
		fmt.Printf("[ERROR] Failed to parse tx logs: %v\n", err)
		return
	}

	// Parse block logs
	blockLogs, err := parseBlockLogs(logDir)
	if err != nil {
		fmt.Printf("[ERROR] Failed to parse block logs: %v\n", err)
		return
	}

	// Calculate latency and write to file
	if err := calculateLatency(txLogs, blockLogs, outputFile); err != nil {
		fmt.Printf("[ERROR] Failed to calculate latency: %v\n", err)
		return
	}

	fmt.Println("[INFO] Latency calculation completed successfully.")
}
