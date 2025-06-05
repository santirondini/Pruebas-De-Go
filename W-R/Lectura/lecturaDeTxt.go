package main

import (
	"log"
	"os"
	"io"
	"strings"
	"bufio"
	"strconv"
)

type Instruccion struct {
	Operacion string 
	Operandos []string
}

func check(mensaje string, e error) {
	if e != nil {
		log.Println(mensaje, e)
	}
}

func LeerArchivoYGuardarInstrucciones(path string) []string {
	
	var instrucciones []string
	
	file , err := os.Open(path)
	check("No se pudo abrir el archivo",err)

	reader := bufio.NewReader(file)
	for { 
		instruccion, err := reader.ReadString('\n') 
		if err != nil {
			if err == io.EOF {   			
				break
				}
		check("Error al leer la instruccion",err) 
		}
		instruccion = strings.TrimSpace(instruccion)
		instrucciones = append(instrucciones, instruccion) 
	}
	defer file.Close() 

	return instrucciones
}

func LeerArchivoYGuardarInstruccionesConScanner(path string) []string {
	
	var instrucciones []string
	
	file , err := os.Open(path)
	check("No se pudo abrir el archivo",err)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() { 
		instruccion := strings.TrimSpace(scanner.Text()) 
		if instruccion != ""{
			instrucciones = append(instrucciones, instruccion)
		}
	}
	if err := scanner.Err(); err != nil {
		check("Error al leer la instruccion",err)
	}
	defer file.Close() 

	return instrucciones
}


func ConfigurarLogger() {
	logFile, err := os.OpenFile("prueba.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil { // si err no es nil, significa que ocurrió un error al intentar abrir/crear el archivo
		panic(err) //detiene la ejecución del programa y muestra el error. El programa no continua si no se puede configuarar el logger
	}
	mw := io.MultiWriter(os.Stdout, logFile) // se crea un escritor Writer que envia los datos a multiples destinos. 
	log.SetOutput(mw) // se configura para que le logger envie mensajes a mw, que a suvez los redirige tanto a consola como a tp0.log
}

func DecodificarInstruccion(instruccion string) Instruccion {
	var instruccionCompleta Instruccion
	partesInstruccion := strings.Fields(instruccion)
	instruccionCompleta.Operacion = partesInstruccion[0]
	instruccionCompleta.Operandos = partesInstruccion[1:]
	// log.Println("Instruccion decodificada:", instruccionCompleta.Operacion)
	return instruccionCompleta
}

func main() {

	// var instruccionEnStruct Instruccion
	// var instruccion string = "JUMP 2 1"

	ConfigurarLogger() // se configura el logger para que envie mensajes a consola y a un archivo
	// log.Println("Instruccion original:", instruccion) // se logea la instruccion original
	// instruccionEnStruct = DecodificarInstruccion(instruccion) // se decodifica la instruccion
	// log.Println("Instruccion decodificada:", instruccionEnStruct.Operacion) // se logea la instruccion decodificada
	// log.Println("Operandos:", instruccionEnStruct.Operandos) // se logean los operandos
	instrucciones := LeerArchivoYGuardarInstruccionesConScanner("santino.txt") // se leen las instrucciones del archivo y se guardan en un slice
	i := 0
	log.Println("Nueva prueba: ")
	log.Println("Cantidad de intrucciones leidas:", len(instrucciones)) // se logea la cantidad de instrucciones leidas
	for _, instruccion := range instrucciones { // por cada instruccion en el slice de instrucciones 
		log.Println("Instruccion "+strconv.Itoa(i), instruccion) // se logea la instruccion
		i++
	}

	var instruccionesDecodificadas []Instruccion
	for _, instr := range instrucciones {
		decodificada := DecodificarInstruccion(instr)
		instruccionesDecodificadas = append(instruccionesDecodificadas, decodificada)
	}
	log.Println("Instrucciones decodificadas:", instruccionesDecodificadas)
}
