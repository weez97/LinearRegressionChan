package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type RegressionResult struct {
	Intercept float64 `json:"intercept"`
	Slope     float64 `json:"slope"`
}

func fetchCSV(url string) ([][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch CSV: %s", resp.Status)
	}

	// Read the CSV file
	reader := csv.NewReader(resp.Body)
	reader.LazyQuotes = true // Allow lazy quotes to handle improperly formatted CSV files

	// Read all records
	var data [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV: %v", err)
		}

		// Remove BOM (Byte Order Mark) if present in the first field of the first record
		if len(data) == 0 && len(record) > 0 && strings.HasPrefix(record[0], "\uFEFF") {
			record[0] = record[0][3:] // Remove the BOM (3 bytes for UTF-8)
		}

		// Append the record to data
		data = append(data, record)
	}

	return data, nil
}

func linearRegression(data [][]string) (float64, float64, error) {
	var sumX, sumY, sumXY, sumXX float64
	n := float64(len(data))

	for _, record := range data {
		x, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Printf("Skipping invalid x value: %v", record[0])
			continue // Skip invalid x values
		}

		y, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Printf("Skipping invalid y value: %v", record[1])
			continue // Skip invalid y values
		}

		sumX += x
		sumY += y
		sumXY += x * y
		sumXX += x * x
	}

	slope := (n*sumXY - sumX*sumY) / (n*sumXX - sumX*sumX)
	intercept := (sumY - slope*sumX) / n
	return intercept, slope, nil
}

func regressionHandler(w http.ResponseWriter, r *http.Request) {
	csvURL := r.URL.Query().Get("url")
	if csvURL == "" {
		http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
		return
	}

	data, err := fetchCSV(csvURL)
	if err != nil {
		log.Printf("Error fetching CSV: %v", err)
		http.Error(w, fmt.Sprintf("Failed to fetch CSV: %v", err), http.StatusInternalServerError)
		return
	}

	intercept, slope, err := linearRegression(data)
	if err != nil {
		log.Printf("Error performing linear regression: %v", err)
		http.Error(w, fmt.Sprintf("Failed to perform linear regression: %v", err), http.StatusInternalServerError)
		return
	}

	result := RegressionResult{Intercept: intercept, Slope: slope}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, fmt.Sprintf("Failed to encode JSON response: %v", err), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/regression", regressionHandler)
	log.Println("Server is listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
