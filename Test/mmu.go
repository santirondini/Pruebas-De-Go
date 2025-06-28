package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)


// ---------------------------------------- MMU --------------------------------------------- //

// TESTEADO 

func NroPagina(direccionLogica int, pagesize int) int {
	return direccionLogica / pagesize
}

func Desplazamiento(direccionLogica int, pagesize int) int {
	return direccionLogica % pagesize
}

func EntradaNiveln(direccionlogica int, niveles int, idTabla int, pagesize int, cantEntradas int) int {
	return (NroPagina(direccionlogica, (pagesize^(niveles - idTabla))) % cantEntradas )
}

func PedirTablaDePaginas(pid uint) *Tabla {
	url := fmt.Sprintf("http://%s:%d/tabla-paginas?pid=", Config.IPMemory, Config.PortMemory) + strconv.Itoa(int(pid))
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error al solicitar la tabla de páginas: %v", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error al obtener la tabla de páginas, código de estado: %d", resp.StatusCode)
		return nil
	}

	var tabla Tabla
	if err := json.NewDecoder(resp.Body).Decode(&tabla); err != nil {
		log.Fatalf("Error al decodificar la tabla de páginas: %v", err)
		return nil
	}

	return &tabla
}

func MMU(pid uint, direccionLogica int) int {

	desplazamiento := Desplazamiento(direccionLogica, tamanioPagina)
	tabla := PedirTablaDePaginas(pid) // Obtengo la tabla de páginas del PID

	if tabla == nil {	
		log.Printf("No se pudo obtener la tabla de páginas para el PID %d", pid)
		return -1
	}
	
	raiz := tabla
	for nivel := 1; nivel <= numeroDeNiveles; nivel++ {
		entrada := EntradaNiveln(direccionLogica, numeroDeNiveles, nivel, tamanioPagina, entradasPorPagina)

		// Si llegamos al nivel final => queda buscar el frame unicamente 
		if nivel == numeroDeNiveles {

			if entrada >= len(raiz.Valores) || raiz.Valores[entrada] == -1 { // verifico si la entrada es válida
				log.Printf("Dirección lógica %d no está mapeada en la tabla de páginas del PID %d", direccionLogica, pid)
				return -1 // Dirección no mapeada
			}
		frame := raiz.Valores[entrada] // Obtengo el frame correspondiente a la entrada
		return frame*tamanioPagina + desplazamiento // Esto es el frame correspondiente a la dirección lógica 
		}

		// Si estamos en niveles intermedios => seguimos recorriendo la tabla de páginas
		if entrada >= len(raiz.Punteros) || raiz.Punteros[entrada] == nil { 
			log.Printf("Dirección lógica %d no está mapeada en la tabla de páginas del PID %d", direccionLogica, pid)
			return -1 // Dirección no mapeada
		}
		raiz = raiz.Punteros[entrada] // Avanzamos al siguiente nivel de la tabla de páginas
	}

	log.Printf("Error al procesar la dirección lógica %d para el PID %d", direccionLogica, pid)
	return -1 // Si llegamos hasta aca => error en el procesamiento de la dirección lógica
}
