package main 

import (
	"log"
	"time"
)

// TLB

type WriteInstruction struct {
	LogicAddress int // Lógica
	Data    string
	PID     uint
}

func InicializarTLB() TLB {
	entradas := make([]EntradaTLB, Config.CacheEntries)
	for i := range entradas {
		entradas[i] = EntradaTLB{
			NumeroPagina: -1, // -1 indica entrada vacía
			NumeroFrame:  -1,
			BitPresencia: false,
			PID:          -1,
			Referencia: -1,
			Llegada: -1, // -1 indica que la entrada no ha sido utilizada
		}
	}
	return TLB{
		Entradas: entradas, // Inicializa las entradas de la TLB con el tamaño máximo configurado
		MaxEntradas: Config.CacheEntries,
		Algoritmo:   "LRU",
	}
}
var tlb TLB = InicializarTLB()

type TLB struct {
	Entradas    []EntradaTLB `json:"entradas"`
	MaxEntradas int          `json:"max_entradas"`
	Algoritmo   string       `json:"algoritmo"`
}

type EntradaTLB struct {
	NumeroPagina        int  `json:"numero_pagina"`
	NumeroFrame         int  `json:"numero_frame"`
	BitPresencia        bool `json:"bit_presencia"`    // Indica si el frame esta presente en memoria
	PID                 int  `json:"pid"`              // Identificador del proceso al que pertenece el frame
	Referencia int  `json:"instante_referencia"` // Marca el instante de referencia para LRU
	Llegada int  `json:"llegada"` // Marca el instante de llegada para FIFO
}

func TLBHabilitada() bool {
	return tlb.MaxEntradas != 0
}

func BuscarDireccion(pagina int) (bool,int) { // devolvemos el frame ya que la pagina esta cargada en el TLB
	
	for i := 0; i < len(tlb.Entradas);i++ {
		if tlb.Entradas[i].NumeroPagina == pagina && tlb.Entradas[i].BitPresencia {
			return true,i // La página está en la TLB y es válida
		}
	}
	return false,-1 // La página no está en la TLB o no es válida
}

func ActualizarReferencia(nropagina int) {
	bool, indice := BuscarDireccion(nropagina)
	if !bool { // Si la página no está en la TLB, no se puede actualizar la referencia
		log.Println("No se puede actualizar la referencia, la página no está en la TLB")
		return
	}
	entradaReferenciada := tlb.Entradas[indice]
	entradaReferenciada = EntradaTLB{
		NumeroPagina: tlb.Entradas[indice].NumeroPagina, // Mantenemos el número de página
		NumeroFrame: tlb.Entradas[indice].NumeroFrame, // Mantenemos el número de frame
		BitPresencia: tlb.Entradas[indice].BitPresencia, //
		PID: tlb.Entradas[indice].PID, // Mantenemos el PID
		Llegada: tlb.Entradas[indice].Llegada, // Manten
		Referencia: int(time.Now().UnixNano()), // Actualizamos el instante de referencia
	}
	tlb.Entradas[indice] = entradaReferenciada // Actualizamos la entrada en la TLB
	return 
}

func AccesoATLB(pid int, nropagina int) int {
	
	if !TLBHabilitada() {
		log.Println("TLB no habilitada, no se puede acceder a la TLB")
		return -1 // TLB no habilitada, no se puede acceder a la
	}
	
	bool, indice := BuscarDireccion(nropagina) // Verificamos si la página está en la TLB
	if bool { 
		log.Printf("PID: < %d > - TLB HIT - Página: %d", pid, nropagina)
		ActualizarReferencia(indice) // Actualizamos la referencia de la entrada en la TLB
		return tlb.Entradas[indice].NumeroFrame // Si la página está en la TLB, devolvemos el frame y true
	} else {
		log.Println("PID: < %d > - TLB MISS - Página: %d", pid, nropagina)
		return -1
	}
}

func IndiceDeEntradaVictima(segun func(EntradaTLB) int) int {
		
	victima := tlb.Entradas[0] // Inicializamos la víctima con la primera entrada
	indice := 0
	for i:=0; i < len(tlb.Entradas); i++{
		if segun(victima) > segun(tlb.Entradas[i]) { // Comparamos la entrada actual con la víctima
			victima = tlb.Entradas[i] // Si la entrada actual es más
			indice = i // Actualizamos el índice de la víctima
		}

	}
	return indice 
}


func TLBLleno() bool {
	for i:= 0; i < len(tlb.Entradas); i++ {
		if tlb.Entradas[i].NumeroPagina == -1 { // Si hay una entrada con número de página -1, significa que la TLB no está llena
			return false // La TLB no está llena
		}
	}
	return true 
}

func EntradaTLBValida() int {
	for i := 0; i < len(tlb.Entradas); i++ {
		if tlb.Entradas[i].NumeroPagina == -1 { // Si hay una entrada con número de página -1, significa que la TLB tiene espacio
			return i // Retorna el índice de la entrada válida
		}
	}
	return -1 // Si no hay entradas válidas, retorna -1
}


func AgregarEntradaATLB(pid int, nropagina int, nroframe int) {

	nuevaEntrada := EntradaTLB{
		NumeroPagina: nropagina,
		NumeroFrame: nroframe,
		BitPresencia: true, // La pagina esta presente en memoria
		PID: pid, // Asignamos el PID del proceso
		Referencia: int(time.Now().UnixNano()), // Asignamos el instante de referencia actual
		Llegada: int(time.Now().UnixNano()), // Asignamos el instante de llegada actual
	}

	if TLBLleno() { // si la cantidad de entradas es la maxima => hay que reemplazar
		indiceVictima := IndiceDeEntradaVictima(func(e EntradaTLB) int {
			if Config.CacheReplacement == "FIFO" {
				return e.Llegada
			} else {
				return e.Referencia 
			}
		}) 
		tlb.Entradas[indiceVictima] = nuevaEntrada // reemplazo  la entrada victima por la nueva entrada
		return 
	} else { // si no esta lleno, agrego la nueva entrada al final
		indiceRemplazo := EntradaTLBValida() // Busco una entrada valida para agregar la nueva entrada
		tlb.Entradas[indiceRemplazo] = nuevaEntrada
		return 
	}
}

func MostrarTLB() {
	if len(tlb.Entradas) == 0 {
		log.Println("La TLB está vacía.")
		return
	}
	log.Println("Contenido de la TLB:")
	for i, entrada := range tlb.Entradas {
		log.Printf("Entrada %d: PID=%d | Página=%d | Frame=%d | Presente=%v | InstanteReferencia=%d | Llegada=%d",
			i, entrada.PID, entrada.NumeroPagina, entrada.NumeroFrame, entrada.BitPresencia, entrada.Referencia, entrada.Llegada)
	}
}

func DesalojoTlB(pid uint) {
	for i := 0; i < len(tlb.Entradas); i++ {
		if tlb.Entradas[i].PID == int(pid) { // Verificamos si la entrada pertenece al PID
			tlb.Entradas[i] = EntradaTLB{ // Desalojamos la entrada del PID
				NumeroPagina: -1, // -1 indica que la entrada está vacía
				NumeroFrame:  -1,
				BitPresencia: false,
				PID:          -1,
				Referencia: -1,
				Llegada: -1, // -1 indica que la entrada no ha sido utilizada
			}
		}
	}
}