package main

import (
	"log"
	"os"
	"io"
	"time"
	"fmt"
	"strconv"
)

func check(mensaje string, e error) {
	if e != nil {
		log.Println(mensaje, e)
	}
}

func NombreDelArchivoDMP(pid string) string {
	return "" + pid + "-" + time.Now().Format("2006-01-02-15-04-05") + ".dmp"
}

func CreacionArchivoDump(path string, pid uint) (*os.File, error) {
	nombreArchivo := NombreDelArchivoDMP(strconv.Itoa(int(pid)))
	rutaCompleta := fmt.Sprintf("%s/%s", path, nombreArchivo)
	file, err := os.Create(rutaCompleta)
	check("No se pudo crear el archivo Dump con el nombre" , err)
	return file, nil
}

func ConfigurarLogger() {
	logFile, err := os.OpenFile("dump.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil { // si err no es nil, significa que ocurrió un error al intentar abrir/crear el archivo
		panic(err) //detiene la ejecución del programa y muestra el error. El programa no continua si no se puede configuarar el logger
	}
	mw := io.MultiWriter(os.Stdout, logFile) // se crea un escritor Writer que envia los datos a multiples destinos. 
	log.SetOutput(mw) // se configura para que le logger envie mensajes a mw, que a suvez los redirige tanto a consola como a tp0.log
}

func main() {
	ConfigurarLogger() 
	log.Println("Nueva Prueba de log ---- ")
	nombre := NombreDelArchivoDMP("1234")
	log.Println("Nombre del archivo dump: ", nombre)
	file, err := CreacionArchivoDump("/.", 1234)
	check("No se pudo crear el archivo Dump con el nombre" , err)
	defer file.Close()
	log.Println("Archivo dump creado con exito: ", file.Name())
}