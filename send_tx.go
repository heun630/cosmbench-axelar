package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
	"time"
)

// Configuration
var (
	encodedTxDir = "/data/axelar/cosmbench-axelar/axelar-cosmbench_encoded_txs"
	HOSTS        = []string{"127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1"}
	REST_PORTS   = []string{"22200", "22201", "22202", "22203"}
	InputTPS     int
	runTime      int
	numTxs       int
)

type TxData struct {
	TxBytes string `json:"tx_bytes"`
	Mode    string `json:"mode"`
}

func readEncodedTxs(dir string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return nil, fmt.Errorf("failed to find files: %v", err)
	}

	txs := make([]string, 0, len(files))
	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read file (%s): %v", file, err)
		}
		txs = append(txs, string(bytes.TrimSpace(content)))
	}
	numTxs = len(txs)
	return txs, nil
}

// Remove ANSI escape codes from a string
func removeANSI(text string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansiRegex.ReplaceAllString(text, "")
}

// Extracts the latest height from the last 10 lines of the log file
func extractHeightFromLog(logFileName string) (string, error) {
	file, err := os.Open(logFileName)
	if err != nil {
		return "", fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > 10 {
			lines = lines[1:] // Keep only the last 10 lines
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read log file: %v", err)
	}

	heightRegex := regexp.MustCompile(`height=([0-9]+)`)
	for i := len(lines) - 1; i >= 0; i-- {
		cleanedLine := removeANSI(lines[i])
		if matches := heightRegex.FindStringSubmatch(cleanedLine); matches != nil {
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("no height found in log file")
}

func sendTransaction(txIdx int, tx string, wg *sync.WaitGroup, fileMutex *sync.Mutex, logFile *os.File) {
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
		fmt.Printf("[TxIdx %d] JSON marshal error: %v\n", txIdx, err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("[TxIdx %d] Request creation error: %v\n", txIdx, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[TxIdx %d] HTTP request error: %v\n", txIdx, err)
		return
	}
	defer resp.Body.Close()

	timestamp := time.Now().UnixMilli()
	latestHeight, err := extractHeightFromLog("output0.log")
	if err != nil {
		fmt.Printf("[TxIdx %d] Failed to extract height from log: %v\n", txIdx, err)
		return
	}

	fileMutex.Lock()
	defer fileMutex.Unlock()
	fmt.Fprintf(logFile, "txIdx: %d time: %d height: %s\n", txIdx, timestamp, latestHeight)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run send_tx.go [TPS] [RunTime]")
		return
	}

	var err error
	InputTPS, err = strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Invalid TPS value: %v\n", err)
		return
	}

	runTime, err = strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Invalid RunTime value: %v\n", err)
		return
	}

	txs, err := readEncodedTxs(encodedTxDir)
	if err != nil {
		fmt.Printf("Error reading transactions: %v\n", err)
		return
	}

	fmt.Printf("Loaded %d transactions\n", numTxs)

	logFile, err := os.Create("tx_log.txt")
	if err != nil {
		fmt.Printf("Error creating log file: %v\n", err)
		return
	}
	defer logFile.Close()

	var fileMutex sync.Mutex
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
			go sendTransaction(sentTxs+j, txs[sentTxs+j], &wg, &fileMutex, logFile)
		}

		wg.Wait()
		sentTxs += txsToSend

		elapsed := time.Since(startTime).Milliseconds()
		if elapsed < 1000 {
			time.Sleep(time.Duration(1000-elapsed) * time.Millisecond)
		}
	}

	fmt.Printf("All transactions sent (%d total). Logs saved to tx_log.txt\n", sentTxs)
}
