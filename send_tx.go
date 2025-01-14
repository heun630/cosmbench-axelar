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
	encodedTxDir = "/data/axelar/cosmbench-axelar/axelar-cosmbench_encoded_txs" // Directory containing encoded transactions
	HOSTS        = []string{"127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1"} // Node IPs
	REST_PORTS   = []string{"22200", "22201", "22202", "22203"}                 // Node REST API ports
	InputTPS     int                                                            // Transactions per second
	runTime      int                                                            // Runtime in seconds
	numTxs       int                                                            // Total number of transactions
)

type TxData struct {
	TxBytes string `json:"tx_bytes"` // Encoded transaction data
	Mode    string `json:"mode"`     // Broadcast mode
}

// Reads encoded transactions from the specified directory
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

// Extracts the latest height from the log file
func extractHeightFromLog(logFileName string) (string, error) {
	file, err := os.Open(logFileName)
	if err != nil {
		return "", fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	heightRegex := regexp.MustCompile(`committed state.*?height=\^\[0m([0-9]+)`) // Correct regex to extract height
	latestHeight := "0"

	for scanner.Scan() {
		line := scanner.Text()
		if matches := heightRegex.FindStringSubmatch(line); matches != nil {
			latestHeight = matches[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read log file: %v", err)
	}

	return latestHeight, nil
}

// Sends a single transaction and logs its height
func sendTransaction(txIdx int, tx string, wg *sync.WaitGroup, fileMutex *sync.Mutex, logFile *os.File) {
	defer wg.Done()

	host := HOSTS[txIdx%len(HOSTS)]
	port := REST_PORTS[txIdx%len(REST_PORTS)]
	url := fmt.Sprintf("http://%s:%s/cosmos/tx/v1beta1/txs", host, port)

	requestData := TxData{
		TxBytes: tx,
		Mode:    "BROADCAST_MODE_BLOCK", // Block mode to ensure transaction is committed
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("[TxIdx %d] Response read error: %v\n", txIdx, err)
		return
	}

	var responseMap map[string]interface{}
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		fmt.Printf("[TxIdx %d] JSON unmarshal error: %v\n", txIdx, err)
		return
	}

	timestamp := time.Now().UnixMilli()
	if txResponse, ok := responseMap["tx_response"].(map[string]interface{}); ok {
		if code, ok := txResponse["code"].(float64); ok && code != 0 {
			fmt.Printf("[TxIdx %d] Transaction failed with code: %.0f, log: %s\n", txIdx, code, txResponse["raw_log"])
			return
		}

		if h, ok := txResponse["height"].(string); ok && h != "0" {
			fileMutex.Lock()
			defer fileMutex.Unlock()
			fmt.Fprintf(logFile, "txIdx: %d time: %d height: %s\n", txIdx, timestamp, h)
			fmt.Printf("[TxIdx %d] Response: %s\n", txIdx, string(body))
			return
		}
	}

	// Fetch latest height from logs if height is missing
	latestHeight, err := extractHeightFromLog("output0.log")
	if err != nil {
		fmt.Printf("[TxIdx %d] Failed to extract height from log: %v\n", txIdx, err)
		return
	}

	fileMutex.Lock()
	defer fileMutex.Unlock()
	fmt.Fprintf(logFile, "txIdx: %d time: %d height: %s\n", txIdx, timestamp, latestHeight)
	fmt.Printf("[TxIdx %d] Response: %s\n", txIdx, string(body))
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

	fmt.Printf("Starting with TPS: %d, RunTime: %d seconds\n", InputTPS, runTime)

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

		if sentTxs >= numTxs {
			break
		}
	}

	fmt.Printf("All transactions sent (%d total). Logs saved to tx_log.txt\n", sentTxs)
}
