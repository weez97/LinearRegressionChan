package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func main() {
	// Puerto en el que escuchará el servidor
	port := ":12345"

	// Iniciar servidor
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
	defer ln.Close()
	fmt.Printf("Servidor escuchando en el puerto %s...\n", port)

	// Aceptar conexiones entrantes
	conn, err := ln.Accept()
	if err != nil {
		log.Fatalf("Error al aceptar la conexión: %v", err)
	}
	defer conn.Close()
	fmt.Println("Cliente conectado:", conn.RemoteAddr())

	// Leer el número total de líneas a procesar del cliente
	reader := bufio.NewReader(conn)
	numLinesStr, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error al leer el número total de líneas: %v", err)
	}
	numLinesStr = strings.TrimSpace(numLinesStr)
	numLines, err := strconv.Atoi(numLinesStr)
	if err != nil {
		log.Fatalf("Error al convertir el número total de líneas: %v", err)
	}

	// Barra de progreso para simular el procesamiento
	bar := progressbar.Default(int64(numLines))

	// Variables para la regresión lineal
	var sumX, sumY, sumXY, sumXX float64

	// Leer datos del cliente y procesarlos
	for i := 0; i < numLines; i++ {
		// Leer línea de datos del cliente
		dataStr, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error al leer los datos del cliente: %v", err)
		}
		dataStr = strings.TrimSpace(dataStr)
		data := strings.Split(dataStr, ",")

		// Convertir datos a números
		x, err := strconv.ParseFloat(data[0], 64)
		if err != nil {
			log.Fatalf("Error al convertir x a float64: %v", err)
		}
		y, err := strconv.ParseFloat(data[1], 64)
		if err != nil {
			log.Fatalf("Error al convertir y a float64: %v", err)
		}

		// Actualizar sumas para la regresión lineal
		sumX += x
		sumY += y
		sumXY += x * y
		sumXX += x * x

		// Actualizar barra de progreso
		bar.Add(1)
	}

	// Calcular los coeficientes de la regresión lineal
	n := float64(numLines)
	b := (n*sumXY - sumX*sumY) / (n*sumXX - sumX*sumX)
	a := (sumY - b*sumX) / n

	// Construir y enviar respuesta al cliente
	regression := fmt.Sprintf("%.6f %.6f", a, b)
	_, err = conn.Write([]byte(regression + "\n"))
	if err != nil {
		log.Fatalf("Error al enviar la regresión al cliente: %v", err)
	}

	// Imprimir mensaje de finalización en el servidor
	fmt.Println("Procesamiento completado, regresión calculada:", regression)
}
