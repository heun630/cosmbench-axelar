package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
)

func convertJSONToCSV(jsonFile string, csvFile string) error {
	// Open the JSON file
	file, err := os.Open(jsonFile)
	if err != nil {
		return fmt.Errorf("failed to open JSON file: %v", err)
	}
	defer file.Close()

	// Decode JSON into a generic interface
	var data interface{}
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

	// Determine the type of JSON data
	value := reflect.ValueOf(data)

	if value.Kind() == reflect.Slice {
		// Handle JSON arrays
		dataArray, ok := data.([]interface{})
		if !ok {
			return fmt.Errorf("failed to cast JSON array")
		}

		if len(dataArray) == 0 {
			fmt.Printf("[INFO] JSON file %s contains no data.\n", jsonFile)
			return nil
		}

		// Extract headers from the first element
		firstElement, ok := dataArray[0].(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed to parse first element as object")
		}

		var headers []string
		for key := range firstElement {
			headers = append(headers, key)
		}
		if err := writer.Write(headers); err != nil {
			return fmt.Errorf("failed to write headers to CSV file: %v", err)
		}

		// Write data rows
		for _, item := range dataArray {
			record, ok := item.(map[string]interface{})
			if !ok {
				return fmt.Errorf("failed to parse item as object")
			}

			var row []string
			for _, header := range headers {
				row = append(row, fmt.Sprintf("%v", record[header]))
			}
			if err := writer.Write(row); err != nil {
				return fmt.Errorf("failed to write row to CSV file: %v", err)
			}
		}
	} else if value.Kind() == reflect.Map {
		// Handle JSON objects
		dataMap, ok := data.(map[string]interface{})
		if !ok {
			return fmt.Errorf("failed to cast JSON object")
		}

		// Write headers (keys of the object)
		var headers []string
		var row []string
		for key, val := range dataMap {
			headers = append(headers, key)
			row = append(row, fmt.Sprintf("%v", val))
		}
		if err := writer.Write(headers); err != nil {
			return fmt.Errorf("failed to write headers to CSV file: %v", err)
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row to CSV file: %v", err)
		}
	} else {
		return fmt.Errorf("unsupported JSON format")
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
