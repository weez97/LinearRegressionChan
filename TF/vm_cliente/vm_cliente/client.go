package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type RegressionResult struct {
	Intercept float64 `json:"intercept"`
	Slope     float64 `json:"slope"`
}

func main() {
	serverIP := "http://192.168.18.97:8080"

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the URL of the CSV file: ")
	csvURL, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	// Remove newline character from the input
	csvURL = csvURL[:len(csvURL)-1]

	resp, err := http.Get(fmt.Sprintf("%s/regression?url=%s", serverIP, csvURL))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: Server returned status code %d\n", resp.StatusCode)
		return
	}

	var result RegressionResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		panic(err)
	}

	fmt.Printf("Intercept: %.6f\n", result.Intercept)
	fmt.Printf("Slope: %.6f\n", result.Slope)
}
