package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type BlockInfo struct {
	Height    int64
	Timestamp int64
}

func parseBlockLogs(logFiles []string) ([]BlockInfo, error) {
	var blocks []BlockInfo
	blockRegex := regexp.MustCompile(`(\d+)\s+.*height=(\d+)`)

	for _, logFile := range logFiles {
		file, err := os.Open(logFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file %s: %v", logFile, err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			matches := blockRegex.FindStringSubmatch(line)
			if matches != nil {
				timestamp, _ := strconv.ParseInt(matches[1], 10, 64)
				height, _ := strconv.ParseInt(matches[2], 10, 64)
				blocks = append(blocks, BlockInfo{Height: height, Timestamp: timestamp})
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("failed to read log file %s: %v", logFile, err)
		}
	}
	return blocks, nil
}

func calculateLatency(txLog string, blockInfo []BlockInfo) (map[int64][]int64, error) {
	file, err := os.Open(txLog)
	if err != nil {
		return nil, fmt.Errorf("failed to open tx log: %v", err)
	}
	defer file.Close()

	latencyMap := make(map[int64][]int64) // Map of block height to latencies
	txRegex := regexp.MustCompile(`txIdx: \d+ timestamp: (\d+) hash: .* height: (\d+)`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := txRegex.FindStringSubmatch(line)
		if matches != nil {
			txTimestamp, _ := strconv.ParseInt(matches[1], 10, 64)
			txHeight, _ := strconv.ParseInt(matches[2], 10, 64)

			for _, block := range blockInfo {
				if block.Height == txHeight {
					latency := block.Timestamp - txTimestamp
					latencyMap[txHeight] = append(latencyMap[txHeight], latency)
					break
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read tx log: %v", err)
	}

	return latencyMap, nil
}

func main() {
	blockLogFiles := []string{"output0.log", "output1.log", "output2.log", "output3.log"}
	txLog := "tx_log.txt"

	// Parse block logs
	blockInfo, err := parseBlockLogs(blockLogFiles)
	if err != nil {
		fmt.Printf("Error parsing block logs: %v\n", err)
		return
	}

	// Calculate latency for transactions
	latencyMap, err := calculateLatency(txLog, blockInfo)
	if err != nil {
		fmt.Printf("Error calculating latency: %v\n", err)
		return
	}

	// Print latencies for each block height
	for height, latencies := range latencyMap {
		fmt.Printf("Block Height: %d\n", height)
		for _, latency := range latencies {
			fmt.Printf("  Latency: %d ms\n", latency)
		}
	}
}
