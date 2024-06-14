package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Dirección IP y puerto del servidor
	serverAddr := "localhost:12345"

	// Conectar con el servidor
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalf("Error al conectar con el servidor: %v", err)
	}
	defer conn.Close()

	// Descargar el archivo CSV desde la URL
	url := "https://raw.githubusercontent.com/weez97/LinearRegressionChan/main/test.csv"
	csvData, err := downloadCSV(url)
	if err != nil {
		log.Fatalf("Error al descargar el archivo CSV: %v", err)
	}

	// Ignorar la primera línea del CSV si existe
	if len(csvData) > 0 {
		csvData = csvData[1:] // Omitir la primera línea
	}

	// Limitar a las primeras 50,000 líneas del CSV
	if len(csvData) > 50000 {
		csvData = csvData[:50000]
	}

	// Número total de líneas del CSV a calcular
	numLineas := len(csvData)

	// Enviar el número total de líneas al servidor
	_, err = conn.Write([]byte(strconv.Itoa(numLineas) + "\n"))
	if err != nil {
		log.Fatalf("Error al enviar el número total de líneas al servidor: %v", err)
	}

	// Enviar los datos al servidor en lotes de batchSize
	batchSize := 1000 // Tamaño del lote para enviar al servidor
	for i := 0; i < numLineas; i += batchSize {
		end := i + batchSize
		if end > numLineas {
			end = numLineas
		}
		batch := csvData[i:end]

		// Enviar el lote al servidor
		for _, row := range batch {
			data := strings.Join(row, ",") + "\n"
			_, err := conn.Write([]byte(data))
			if err != nil {
				log.Fatalf("Error al enviar datos al servidor: %v", err)
			}
		}

		// Simular procesamiento en el servidor
		time.Sleep(500 * time.Millisecond)

		// Imprimir mensaje cada vez que se complete un cálculo en el servidor
		fmt.Printf("Se terminó un cálculo [%d/%d]\n", end, numLineas)
	}

	// Recibir los resultados del servidor
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatalf("Error al recibir los resultados del servidor: %v", err)
	}

	// Imprimir mensaje final de todos los cálculos terminados
	fmt.Println("\nTodos los cálculos terminados")

	// Imprimir la regresión final calculada
	fmt.Println("La regresión final es:", response)
}

// Función para descargar un archivo CSV desde una URL y leerlo
func downloadCSV(url string) ([][]string, error) {
	// Realizar la petición HTTP GET
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error al descargar el archivo CSV: %v", err)
	}
	defer resp.Body.Close()

	// Leer el contenido del cuerpo de la respuesta
	reader := csv.NewReader(resp.Body)
	reader.FieldsPerRecord = -1 // Permitir un número variable de campos por registro

	// Leer todas las líneas del CSV
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error al leer el archivo CSV: %v", err)
	}

	return records, nil
}
