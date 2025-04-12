package main 

import (
	"fmt"
)


func main() {

	fmt.Println("Primera Prueba de Go:")
	fmt.Println("")

	var frase string
	fmt.Println("Escribe una frase:")
	fmt.Scanln(&frase)
	fmt.Println("La frase es: ", frase)
}