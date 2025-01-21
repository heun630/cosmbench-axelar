package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func convertJSONToCSV(jsonFile string, csvFile string) error {
	// Open the JSON file
	file, err := os.Open(jsonFile)
	if err != nil {
		return fmt.Errorf("failed to open JSON file: %v", err)
	}
	defer file.Close()

	// Decode JSON into a generic slice of maps
	var data []map[string]interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("failed to decode JSON file: %v", err)
	}

	// Create the CSV file
	csvFileHandle, err := os.Create(csvFile)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer csvFileHandle.Close()

	writer := csv.NewWriter(csvFileHandle)
	defer writer.Flush()

	// Write the header row
	if len(data) > 0 {
		var headers []string
		for key := range data[0] {
			headers = append(headers, key)
		}
		if err := writer.Write(headers); err != nil {
			return fmt.Errorf("failed to write headers to CSV file: %v", err)
		}

		// Write the data rows
		for _, record := range data {
			var row []string
			for _, header := range headers {
				value := record[header]
				row = append(row, fmt.Sprintf("%v", value))
			}
			if err := writer.Write(row); err != nil {
				return fmt.Errorf("failed to write data row to CSV file: %v", err)
			}
		}
	}

	fmt.Printf("[INFO] JSON data successfully converted to %s\n", csvFile)
	return nil
}

func main() {
	// Define the directory containing the JSON files
	resultDir := "results"

	// Find all JSON files in the results directory
	jsonFiles, err := filepath.Glob(filepath.Join(resultDir, "*.json"))
	if err != nil {
		fmt.Printf("[ERROR] Failed to find JSON files: %v\n", err)
		return
	}

	if len(jsonFiles) == 0 {
		fmt.Println("[INFO] No JSON files found in the results directory.")
		return
	}

	// Convert each JSON file to a corresponding CSV file
	for _, jsonFile := range jsonFiles {
		csvFile := jsonFile[:len(jsonFile)-len(filepath.Ext(jsonFile))] + ".csv"
		if err := convertJSONToCSV(jsonFile, csvFile); err != nil {
			fmt.Printf("[ERROR] Failed to convert %s to CSV: %v\n", jsonFile, err)
		}
	}

	fmt.Println("[INFO] All JSON files have been successfully converted to CSV.")
}
