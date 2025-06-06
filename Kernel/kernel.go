package main

import (
	"log"
	"os"
	"io"
)

func ConfigurarLogger() {
	logFile, err := os.OpenFile("kernel.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil { // si err no es nil, significa que ocurrió un error al intentar abrir/crear el archivo
		panic(err) //detiene la ejecución del programa y muestra el error. El programa no continua si no se puede configuarar el logger
	}
	mw := io.MultiWriter(os.Stdout, logFile) // se crea un escritor Writer que envia los datos a multiples destinos. 
	log.SetOutput(mw) // se configura para que le logger envie mensajes a mw, que a suvez los redirige tanto a consola como a tp0.log
}

func main() {
	ConfigurarLogger()
	log.Println("Kernel iniciado")
	
	tablaConInt := []int{1, 2, 3, 4, 5}
	log.Println("Tabla de enteros:", tablaConInt)


	tablaConMap := map[int]int{
		1: 10,
		2: 20,
		3: 30,
		4: 40,
		5: 50,
	}
	log.Println("Tabla con map:", tablaConMap)

	var matrizConInt [3]string[]int
	valor := 1
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			matrizConInt[i][j] = valor
			valor++
		}
	}
	log.Println("Matriz cargada:", matrizConInt)


	matriz := [3][3]int{}
	valor = 1
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			matriz[i][j] = valor
			valor++
		}
	}
	log.Println("Tabla 3x3:", matriz)
}