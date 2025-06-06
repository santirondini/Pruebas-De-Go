package main

import (
	"fmt"
	"bytes"
)

type Alumno struct {
	Nombre  string
	Edad    int
	Email   string
	Carrera string
}

// var santino = Alumno{
// 	Nombre:  "Santino", 
// 	Edad:    20,
// 	Email:   "santino@gmail.com",
// 	Carrera: "Ingenieria en Sistemas",
// }


func main() {
	memoriaPrincipal := make([]byte, 2048)


	fmt.Println("Tamaño de memoriaPrincipal: ", len(memoriaPrincipal))
	santino := "santino"
	rondini := "rondini" // por letra = 1 byte

	copy(memoriaPrincipal[0:len(santino)], santino)
	copy(memoriaPrincipal[len(santino):len(santino)+len(rondini)], rondini)
	
	fmt.Println("Memoria: ", memoriaPrincipal)
	tamanioTotal := len(santino) + len(rondini)

	// Imprimir memoria en formato ASCII85 - Lectura de byte
	
}

// 115 97 110 116 105 110 111 114 111 110 100 105 105 110 105
	











