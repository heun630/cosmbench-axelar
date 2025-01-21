package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

type TxLog struct {
	TxIdx     int    `json:"txIdx"`
	Timestamp int64  `json:"timestamp"`
	Height    int    `json:"height"`
	TxHash    string `json:"txHash"`
}

type TPSData struct {
	TotalTxs         int     `json:"totalTxs"`
	FirstTxTimestamp int64   `json:"firstTxTimestamp"`
	LastTxTimestamp  int64   `json:"lastTxTimestamp"`
	TotalElapsedTime int64   `json:"totalElapsedTime"`
	TPS              float64 `json:"tps"`
}

type BlockTransactionCount struct {
	Height           int `json:"height"`
	TransactionCount int `json:"transactionCount"`
}

type LatencyEntry struct {
	TxIdx     int   `json:"txIdx"`
	Timestamp int64 `json:"timestamp"`
	Height    int   `json:"height"`
	Latency   int64 `json:"latency"`
}

// Function to convert JSON data to CSV and JSON
func convertJSON(inputFile, csvOutput, jsonOutput string, headers []string, parseFunc func([]byte) ([][]string, []byte, error)) error {
	// Read the input JSON file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file %s: %v", inputFile, err)
	}

	// Parse the data into records and JSON bytes
	records, newJSON, err := parseFunc(data)
	if err != nil {
		return fmt.Errorf("failed to parse data: %v", err)
	}

	// Write the CSV file
	csvFile, err := os.Create(csvOutput)
	if err != nil {
		return fmt.Errorf("failed to create CSV file %s: %v", csvOutput, err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Write CSV headers
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers to CSV: %v", err)
	}

	// Write CSV records
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record to CSV: %v", err)
		}
	}

	// Write the new JSON file
	if err := os.WriteFile(jsonOutput, newJSON, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file %s: %v", jsonOutput, err)
	}

	fmt.Printf("Successfully converted %s to %s and %s\n", inputFile, csvOutput, jsonOutput)
	return nil
}

// Parse functions for each file type
func parseTxLog(data []byte) ([][]string, []byte, error) {
	var txLogs []TxLog
	if err := json.Unmarshal(data, &txLogs); err != nil {
		return nil, nil, err
	}

	// Prepare CSV records and new JSON
	records := [][]string{}
	for _, log := range txLogs {
		records = append(records, []string{
			fmt.Sprintf("%d", log.TxIdx),
			fmt.Sprintf("%d", log.Timestamp),
			fmt.Sprintf("%d", log.Height),
			log.TxHash,
		})
	}

	newJSON, err := json.MarshalIndent(txLogs, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	return records, newJSON, nil
}

func parseTPSData(data []byte) ([][]string, []byte, error) {
	var tpsData TPSData
	if err := json.Unmarshal(data, &tpsData); err != nil {
		return nil, nil, err
	}

	records := [][]string{
		{
			fmt.Sprintf("%d", tpsData.TotalTxs),
			fmt.Sprintf("%d", tpsData.FirstTxTimestamp),
			fmt.Sprintf("%d", tpsData.LastTxTimestamp),
			fmt.Sprintf("%d", tpsData.TotalElapsedTime),
			fmt.Sprintf("%.2f", tpsData.TPS),
		},
	}

	newJSON, err := json.MarshalIndent(tpsData, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	return records, newJSON, nil
}

func parseBlockTransactionCount(data []byte) ([][]string, []byte, error) {
	var blockCounts []BlockTransactionCount
	if err := json.Unmarshal(data, &blockCounts); err != nil {
		return nil, nil, err
	}

	records := [][]string{}
	for _, count := range blockCounts {
		records = append(records, []string{
			fmt.Sprintf("%d", count.Height),
			fmt.Sprintf("%d", count.TransactionCount),
		})
	}

	newJSON, err := json.MarshalIndent(blockCounts, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	return records, newJSON, nil
}

func parseLatencyData(data []byte) ([][]string, []byte, error) {
	var latencies []LatencyEntry
	if err := json.Unmarshal(data, &latencies); err != nil {
		return nil, nil, err
	}

	records := [][]string{}
	for _, latency := range latencies {
		records = append(records, []string{
			fmt.Sprintf("%d", latency.TxIdx),
			fmt.Sprintf("%d", latency.Timestamp),
			fmt.Sprintf("%d", latency.Height),
			fmt.Sprintf("%d", latency.Latency),
		})
	}

	newJSON, err := json.MarshalIndent(latencies, "", "  ")
	if err != nil {
		return nil, nil, err
	}

	return records, newJSON, nil
}

func main() {
	// Define files to convert
	files := []struct {
		input      string
		csvOutput  string
		jsonOutput string
		headers    []string
		parseFunc  func([]byte) ([][]string, []byte, error)
	}{
		{
			input:      "results/tx_log.json",
			csvOutput:  "results/tx_log.csv",
			jsonOutput: "results/tx_log_converted.json",
			headers:    []string{"TxIdx", "Timestamp", "Height", "TxHash"},
			parseFunc:  parseTxLog,
		},
		{
			input:      "results/tps.json",
			csvOutput:  "results/tps.csv",
			jsonOutput: "results/tps_converted.json",
			headers:    []string{"TotalTxs", "FirstTxTimestamp", "LastTxTimestamp", "TotalElapsedTime", "TPS"},
			parseFunc:  parseTPSData,
		},
		{
			input:      "results/block_transactions.json",
			csvOutput:  "results/block_transactions.csv",
			jsonOutput: "results/block_transactions_converted.json",
			headers:    []string{"Height", "TransactionCount"},
			parseFunc:  parseBlockTransactionCount,
		},
		{
			input:      "results/latency.json",
			csvOutput:  "results/latency.csv",
			jsonOutput: "results/latency_converted.json",
			headers:    []string{"TxIdx", "Timestamp", "Height", "Latency"},
			parseFunc:  parseLatencyData,
		},
	}

	for _, file := range files {
		if err := convertJSON(file.input, file.csvOutput, file.jsonOutput, file.headers, file.parseFunc); err != nil {
			fmt.Printf("[ERROR] Failed to convert %s: %v\n", file.input, err)
		}
	}
}
