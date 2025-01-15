package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// TxLog represents a transaction log entry.
type TxLog struct {
	TxIdx     int
	Timestamp int64
	Height    int
}

// BlockLog represents a block log entry.
type BlockLog struct {
	Timestamp int64
	Height    int
}

// parseTxLogs parses the tx_log.txt file to extract transaction logs.
func parseTxLogs(filePath string) ([]TxLog, error) {
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

	return txLogs, nil
}

// parseBlockLogs parses the block log files to extract block logs.
func parseBlockLogs(logDir string) (map[int]int64, error) {
	files, err := filepath.Glob(filepath.Join(logDir, "output*.log"))
	if err != nil {
		return nil, fmt.Errorf("failed to find block log files: %v", err)
	}

	blockLogs := make(map[int]int64)
	blockLogRegex := regexp.MustCompile(`(\d+)\s+.*height=(\d+)`)

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return nil, fmt.Errorf("failed to open block log file %s: %v", file, err)
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			// Remove color codes
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

	return blockLogs, nil
}

// calculateLatency calculates latency for each transaction and writes it to latency.txt.
func calculateLatency(txLogs []TxLog, blockLogs map[int]int64, outputFile string) error {
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
			fmt.Printf("No block timestamp found for height %d (TxIdx: %d)\n", tx.Height, tx.TxIdx)
			continue
		}

		latency := blockTimestamp - tx.Timestamp
		fmt.Fprintf(writer, "TxIdx: %d, Timestamp: %d, Height: %d, Latency: %d ms\n", tx.TxIdx, tx.Timestamp, tx.Height, latency)
	}

	return nil
}

func main() {
	txLogFile := "tx_log.txt"
	logDir := "./"
	latencyFile := "latency.txt"

	// Parse transaction logs
	txLogs, err := parseTxLogs(txLogFile)
	if err != nil {
		fmt.Printf("Failed to parse transaction logs: %v\n", err)
		return
	}

	// Parse block logs
	blockLogs, err := parseBlockLogs(logDir)
	if err != nil {
		fmt.Printf("Failed to parse block logs: %v\n", err)
		return
	}

	// Calculate and write latency
	if err := calculateLatency(txLogs, blockLogs, latencyFile); err != nil {
		fmt.Printf("Failed to calculate latency: %v\n", err)
		return
	}

	fmt.Printf("Latency calculated and written to %s\n", latencyFile)
}
