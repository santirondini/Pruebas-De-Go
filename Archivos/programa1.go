package main

import (
	"fmt"
	"log/slog"
	"os"
	// "time"
)

func main() {
	 
	fmt.Println("Hello, world!")
	var numero int

	f, err := os.Create("archivoLog.log")
	if err != nil {
		slog.Error("Error al crear el archivo", "error", err)
		os.Exit(1)
	}
	defer f.Close()

	// Configura slog para escribir en el archivo de log
	logger := slog.New(slog.NewTextHandler(f, nil))
	slog.SetDefault(logger)

	fmt.Print("Ingrese un número: ")
	_, err = fmt.Scan(&numero)
	if numero == 0 {
		slog.Error("El numero es 0 idiota", "error", err)
		os.Exit(1)
	} else  {
		slog.Info("El numero es distinto de 0", "numero", numero)
	}
}
