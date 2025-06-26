package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)


var entradasPorPagina = 3 // Cantidad de entradas por pagina
var numeroDeNiveles = 2 // Cantidad de niveles de la tabla de paginas
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
var TDPMultinivel map[uint]*Tabla 

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


func CrearTablaDePaginas(pid uint, tamanio int) {
	log.Printf("Creando tabla de paginas para el PID %d", pid)
	paginasRestantes := CantidadDePaginasDeProceso(tamanio)
	tabla := CreaTablaJerarquica(pid, numeroDeNiveles, &paginasRestantes) // Crea una tabla jerárquica para el PID
	TDPMultinivel[pid] = tabla // Asigna la tabla al mapa TDP
}

func GetTDP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	tablaJSON, err := json.MarshalIndent(TDPMultinivel, "", "  ")
	if err != nil {
		http.Error(w, "Error al serializar la tabla de páginas", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(tablaJSON)
}

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



func MMU(pid uint, direccionLogica int) int {

	// configMemoria, err := PedirConfigMemoria()
	// if err != nil {
	// 	log.Println("No se pudo obtener la configuración de memoria: %v", err)
	// 	return -1
	// }
	
	desplazamiento := Desplazamiento(direccionLogica, tamanioPagina)
	tabla := TDPMultinivel[pid] // Obtengo la tabla de páginas del PID

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

	log.Println("Iniciando el sistema de paginación multinivel...")

	// Inicializar las estructuras necesarias
	Ocupadas = make(map[uint]FrameInfo)
	TDPMultinivel = make(map[uint]*Tabla)

	for i := 0; i < tamanioMemoria/tamanioPagina; i++ {
		Ocupadas[uint(i)] = FrameInfo{EstaOcupado: false, PID: -1}
	}

	MarcarPrimerosNOcupados(12, 111)

	paginas := CantidadDePaginasDeProceso(128)
	TDPMultinivel[15] = CreaTablaJerarquica(15, numeroDeNiveles, &paginas) // Crea una tabla jerárquica para el PID 15)
	MostrarTablaArbol(TDPMultinivel[15], "", true) // Muestra la tabla jerárquica creada

	fmt.Println("MMU PARA DIRECION 0 = ", MMU(15, 0)) // Prueba de MMU para la dirección lógica 0 del PID 15
	
	MostrarFramesOcupados()

	http.HandleFunc("/tdp", GetTDP) // Ruta para obtener la tabla de páginas en formato JSON
	http.ListenAndServe(":8080", nil) // Inicia el servidor HTTP para servir la tabla de páginas

}
