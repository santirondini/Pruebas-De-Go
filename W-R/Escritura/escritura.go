package main

import (
	"fmt"
	"math/bits"
	"unsafe"
)

type Alumno struct {
	Nombre  string
	Edad    int
	Email   string
	Carrera string
}

var santino = Alumno{
	Nombre:  "Santino", 
	Edad:    20,
	Email:   "santino@gmail.com",
	Carrera: "Ingenieria en Sistemas",
}


func main() {
	memoriaPrincipal := make([]byte, 4096)


	fmt.Println("Tamaño de memoriaPrincipal: ", len(memoriaPrincipal))
	
	// Usar la variable santino para evitar el error de variable no utilizada
	fmt.Printf("Alumno: %+v\n", santino)

	fmt.Println("Tamaño de santino: ", unsafe.Sizeof(santino)) // Tamaño en bytes de la estructura Alumno
	
	fmt.Println("Tamaño de variables dentro de santino:")

	fmt.Println("Tamaño de Nombre: ", unsafe.Sizeof(santino.Nombre)) // Tamaño en bytes del string Nombre
	fmt.Println("Tamaño de Edad: ", unsafe.Sizeof(santino.Edad))     // Tamaño en bytes del int Edad
	fmt.Println("Tamaño de Email: ", unsafe.Sizeof(santino.Email))   // Tamaño en bytes del string Email
	fmt.Println("Tamaño de Carrera: ", unsafe.Sizeof(santino.Carrera)) // Tamaño en bytes del string Carrera
	fmt.Println("Tamaño de struct: ", unsafe.Sizeof(Alumno{})) // Tamaño en bytes de la estructura Alumno

	fmt.Println("Tamaño de variables dentro de santino en bytes:")
	fmt.Println("Tamaño de Nombre: " , len("Santino")) 
	fmt.Println("Tamaño de Edad: ", unsafe.Sizeof(20)) 
	fmt.Println("Tamaño de Email: ", len("santino@gmai.com"))
	fmt.Println("Tamaño de Carrera: ", len("Ingenieria en Sistemas"))



	fmt.Println("Tamaño de los tipos de variables dentro de santino:") 
	
	
	fmt.Println("Dirección de memoria de santino: ", &santino) // Imprime la dirección de memoria de la variable santino

	fmt.Println("Dirección de memoria del primer elemento de memoriaPrincipal: ", &memoriaPrincipal[0])
	fmt.Println("Meto a santino dentro de la memoriaPrincipal")
	// copy(memoriaPrincipal[0:unsafe.Sizeof(santino)], santino) // Copia el nombre del alumno en memoriaPrincipal
}