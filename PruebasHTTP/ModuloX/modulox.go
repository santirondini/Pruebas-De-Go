package main

// La idea es que el modulo X le mande un paquete al modulo Y

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"io"
)

type InfoAlumno struct {
	Nombre string `json:"nombre"`
	Edad   int    `json:"edad"`
	Email  string `json:"email"`
	Carrera string `json:"carrera"`
}

func ConfigurarLogger() {
	logFile, err := os.OpenFile("moduloX.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil { 
		panic(err) 
	}
	mw := io.MultiWriter(os.Stdout, logFile)  
	log.SetOutput(mw) 
}


func CrearAlumno() InfoAlumno {
	// Crear un objeto de tipo InfoAlumno
	alumno := InfoAlumno{
		Nombre:  "Juan Perez",
		Edad:    20,
		Email:   "juanperez@gmail.com",
		Carrera: "Ingenieria en Sistemas",
	}
	return alumno
}

func MandarPaquete(){
	alumno := CrearAlumno()
	body, err := json.Marshal(alumno)
	if err != nil {
		fmt.Println("Error al serializar el objeto:", err)
		return
	}
	url := fmt.Sprintf("http://localhost:%d/mandarPaquete", 8089)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error al enviar la solicitud:", err)
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Error al enviar el paquete, código de estado: %d\n", resp.StatusCode)
		return
	}

	log.Println("Paquete enviado correctamente al modulo Y.")
}

func main() {

	ConfigurarLogger() 
	log.Println("Iniciando el modulo X...")
	MandarPaquete()
	log.Println("Paquete enviado al modulo Y.")

	http.ListenAndServe(":8088", nil) 
}