package main

import (
	"fmt"
	"log"
	"time"
)


var entradasPorPagina = 4// Cantidad de entradas por pagina
var numeroDeNiveles = 5  // Cantidad de niveles de la tabla de paginas
var tamanioPagina = 64
var tamanioMemoria = 4096 // Tamaño total de la memoria en bytes
var espacioUsuario = make([]byte,tamanioMemoria)


// ---------------------------------------- Tabla de páginas --------------------------------------------- //

// TESTEADO


func MostrarTablaArbol(tabla *Tabla, prefijo string, esUltimo bool) {
	if tabla == nil {
		return
	}

	conector := "├──"
	if esUltimo {
		conector = "└──"
	}

	if len(tabla.Valores) > 0 {
		fmt.Printf("%s%s Tabla nivel N: %v\n", prefijo, conector, tabla.Valores)
	} else {
		fmt.Printf("%s%s Tabla\n", prefijo, conector)
		for i, subtabla := range tabla.Punteros {
			esUltimaEntrada := i == len(tabla.Punteros)-1
			nuevoPrefijo := prefijo
			if esUltimo {
				nuevoPrefijo += "    "
			} else {
				nuevoPrefijo += "│   "
			}
			MostrarTablaArbol(subtabla, nuevoPrefijo, esUltimaEntrada)
		}
	}
}


type Tabla struct {
	Punteros []*Tabla `json:"tabla"` 
	Valores []int `json:"valores"`
}

type FrameInfo struct{
	EstaOcupado bool `json:"esta_ocupado"` // Indica si el frame está ocupado
	PID int `json:"pid"` // Identificador del proceso al que pertenece el frame
}

var Ocupadas map[uint]FrameInfo

func MostrarOcupadas() {
	fmt.Println("Frames ocupados:")
	for frame, info := range Ocupadas {
		estado := "Libre"
		if info.EstaOcupado {
			estado = fmt.Sprintf("Ocupado (PID %d)", info.PID)
		}
		fmt.Printf("Frame %d: %s\n", frame, estado)
	}
}

func FrameLibre(numero uint) bool {
	return !Ocupadas[numero].EstaOcupado
}

func PrimerFrameLibre(arranque uint) int { // arranque => desde cual frame arranco a buscar
	log.Println("Buscando frame libre...")
	CantidadDeFrames := uint(len(Ocupadas)) // Cantidad de frames que hay en memoria
	for i := arranque; i < uint(CantidadDeFrames); i++ {
		if FrameLibre(i) {
			log.Printf("Frame libre encontrado - NRO Frame %d", i)
			return int(i)
		}
	}
	log.Println("No se encontraron frames libres")
	return -1 // Si no encuentra un frame libre => memoria llena => devuelvo -1
}

func MarcarFrameOcupado(frame uint, pid uint) {
	info := Ocupadas[frame]
	info.EstaOcupado = true // Marca el frame como ocupado
	info.PID = int(pid)          // Asigna el PID del proceso que ocupa el frame
	Ocupadas[frame] = info  // Actualiza el mapa con la información del frame
}

func CreaTablaJerarquica(pid uint, nivelesRestantes int, paginasRestantes *int) *Tabla {

	tabla := &Tabla{}

	if nivelesRestantes == 1 {
		// Nivel hoja: asignar solo las páginas necesarias
		for i := 0; i < entradasPorPagina; i++ {
			if *paginasRestantes > 0 {
				frameLibre := PrimerFrameLibre(uint(i)) // Busca el primer frame libre 
 				MarcarFrameOcupado(uint(frameLibre), pid) // Lo marca como ocupado para el PID
				tabla.Valores = append(tabla.Valores, frameLibre)
				*paginasRestantes--
			} else {
				tabla.Valores = append(tabla.Valores, -1) // El valor -1 respresenta que no esta asignado 
			}
		}
	} else {
		// Nivel intermedio: crear subtablas recursivamente
		for i := 0; i < entradasPorPagina; i++ {
			subtabla := CreaTablaJerarquica(pid, nivelesRestantes - 1,paginasRestantes)
			tabla.Punteros = append(tabla.Punteros, subtabla)
		}
	}
	return tabla
}

func CantidadDePaginasDeProceso(tamanio int) int {
	tamanioPagina := int(tamanioPagina)
	cantPaginas := (tamanio / tamanioPagina) // Cantidad de paginas que se necesitan
	return cantPaginas
}



func MarcarPrimerosNOcupados(n int, pid int) {
	for i := 0; i < n; i++ {
		Ocupadas[uint(i)] = FrameInfo{
			PID:  pid,
			EstaOcupado: true,
		}
	}
}

func MostrarFramesOcupados() {
	fmt.Println("Frames ocupados:")
	if len(Ocupadas) == 0 {
		fmt.Println("  (no hay frames ocupados)")
		return
	}

	for i := 0; i < len(Ocupadas); i++ {
		info := Ocupadas[uint(i)]
		if info.EstaOcupado {
			fmt.Printf("  Frame %d: PID %d\n", i, info.PID)
		} else {
			fmt.Printf("  Frame %d: Libre\n", i)
		}
	}
}


// ---------------------------------------- MAIN --------------------------------------------- //


func main() {

	// MostrarCache()
	
	Write(1001, WriteInstruction{Address: 0, Data: "Santino Rondini", PID: 1001})

	time.Sleep(5 * time.Second) // Espera para que se procese la escritura
	MostrarCache()
	MostrarTLB()

	time.Sleep(5 * time.Second) // Espera para que se procese la escritura
	Write(1001, WriteInstruction{Address: 20, Data: "Facultad de Ingenieria", PID: 1001})

	time.Sleep(10 * time.Second) // Espera para que se procese la escritura
	MostrarCache()
	MostrarTLB()
}
