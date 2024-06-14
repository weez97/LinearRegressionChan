package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type partialCalc struct {
	sumX  float64
	sumY  float64
	sumXY float64
	sumX2 float64
}

func calculatePartialSums(x []float64, y []float64, startIndex int, endIndex int, wg *sync.WaitGroup, results chan partialCalc) {
	defer wg.Done()

	partial := partialCalc{}
	for i := startIndex; i < endIndex; i++ {
		partial.sumX += x[i]
		partial.sumY += y[i]
		partial.sumXY += x[i] * y[i]
		partial.sumX2 += x[i] * x[i]
	}

	results <- partial
}

func concurrentLinearRegression(x []float64, y []float64) (float64, float64) {
	numDataPoints := len(x)
	numRoutines := 4
	results := make(chan partialCalc, numRoutines)
	var wg sync.WaitGroup

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		start := i * (numDataPoints / numRoutines)
		end := (i + 1) * (numDataPoints / numRoutines)
		if i == numRoutines-1 {
			end = numDataPoints
		}
		go calculatePartialSums(x, y, start, end, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	total := partialCalc{}
	for partial := range results {
		total.sumX += partial.sumX
		total.sumY += partial.sumY
		total.sumXY += partial.sumXY
		total.sumX2 += partial.sumX2
	}

	n := float64(numDataPoints)
	coefB := (n*total.sumXY - total.sumX*total.sumY) / (n*total.sumX2 - total.sumX*total.sumX)
	coefA := (total.sumY / n) - coefB*(total.sumX/n)

	return coefA, coefB
}

func fetchDataset(url string) ([]float64, []float64, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	reader := csv.NewReader(response.Body)

	// Read the first record to handle BOM and header
	record, err := reader.Read()
	if err != nil {
		return nil, nil, err
	}

	// Check for BOM and remove if present
	if strings.HasPrefix(record[0], "\ufeff") {
		record[0] = strings.TrimPrefix(record[0], "\ufeff")
	}

	// Verify header
	if record[0] != "x" || record[1] != "y" {
		return nil, nil, fmt.Errorf("unexpected header: %v", record)
	}

	var x []float64
	var y []float64

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}

		xi, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			return nil, nil, err
		}

		yi, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, nil, err
		}

		x = append(x, xi)
		y = append(y, yi)
	}

	return x, y, nil
}

func main() {
	file, err := os.Create("elapsed_times.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	datasetURL := "https://raw.githubusercontent.com/weez97/LinearRegressionChan/main/test.csv" // Updated URL

	for i := 0; i < 5; i++ {
		x, y, err := fetchDataset(datasetURL)
		if err != nil {
			fmt.Println("Error fetching dataset:", err)
			return
		}

		//concurrentLinearRegression(x, y)
		coefA, coefB := concurrentLinearRegression(x, y)

		fmt.Printf("Regression result for iteration %d: y = %.6f + %.6fx\n", i+1, coefA, coefB) // Print regression result
	}
}
