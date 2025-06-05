package main

import (
	"net/http"
	"os"
	"io"
	"log"
	"encoding/json"
)

type InfoAlumno struct {
	Nombre string `json:"nombre"`
	Edad   int    `json:"edad"`
	Email  string `json:"email"`
	Carrera string `json:"carrera"`
}

func ConfigurarLogger() {
	logFile, err := os.OpenFile("moduloY.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil { 
		panic(err) 
	}
	mw := io.MultiWriter(os.Stdout, logFile)  
	log.SetOutput(mw) 
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var alumno InfoAlumno
	err := json.NewDecoder(r.Body).Decode(&alumno)

	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}
	log.Printf("Recibido: %+v\n", alumno)
	w.WriteHeader(http.StatusOK)
}

func ObtenerAlumnos() []InfoAlumno {
	alumnos := []InfoAlumno{
		{
			Nombre:  "Ana Gomez",
			Edad:    22,
			Email:   "anagomez@gmail.com",
			Carrera: "Ingenieria Civil",
		},
		{
			Nombre:  "Carlos Lopez",
			Edad:    21,
			Email:   "carloslopez@gmail.com",
			Carrera: "Ingenieria Industrial",
		},
		{
			Nombre:  "Maria Rodriguez",
			Edad:    23,
			Email:   "mariarodriguez@gmail.com",
			Carrera: "Ingenieria Mecanica",
		},
	}
	return alumnos
}

func ObtenerAlumnosHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	alumnos := ObtenerAlumnos()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(alumnos)
	if err != nil {
		http.Error(w, "Error al codificar los datos", http.StatusInternalServerError)
		return
	}
	log.Printf("Enviados %d alumnos\n", len(alumnos))
}

func main(){
		
	ConfigurarLogger()
	// var alumnos []InfoAlumno
	 
	
	http.HandleFunc("/obtenerAlumnos", ObtenerAlumnosHandler)
	http.HandleFunc("/mandarPaquete", Handler)
	http.ListenAndServe(":8089", nil)
}

