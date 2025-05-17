package main

import (
	"log"
	"os"
	"io"
	"fmt"
)

func ConfigurarLogger() {
	logFile, err := os.OpenFile("CPU.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil { // si err no es nil, significa que ocurrió un error al intentar abrir/crear el archivo
		panic(err) //detiene la ejecución del programa y muestra el error. El programa no continua si no se puede configuarar el logger
	}
	mw := io.MultiWriter(os.Stdout, logFile) // se crea un escritor Writer que envia los datos a multiples destinos. 
	log.SetOutput(mw) // se configura para que le logger envie mensajes a mw, que a suvez los redirige tanto a consola como a tp0.log
}

func main() {
	ConfigurarLogger()
	log.Println("CPU iniciado") 
	fmt.Println("NASHE") // se envia un mensaje al logger. Este mensaje se mostrara tanto en consola como en el archivo tp0.log
}