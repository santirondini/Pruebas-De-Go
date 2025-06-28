package main 

import (
	"log"
	"time"
)

// TLB

type WriteInstruction struct {
	Address int
	Data    string
	PID     uint
}

func InicializarTLB() TLB {
	entradas := make([]EntradaTLB, 0, Config.CacheEntries)
	for i := 0; i < Config.CacheEntries; i++ {
		entradas = append(entradas, EntradaTLB{
			NumeroPagina:        -1, // Inicialmente no hay páginas cargadas
			NumeroFrame:         -1, // Inicialmente no hay frames asignados
			BitPresencia:        false, // Inicialmente no hay páginas presentes
			PID:                 -1, // Inicialmente no hay PID asignado
			InstanteDeReferencia: 0, // Inicialmente no hay instante de referencia
		})
	}
	return TLB{
		Entradas:    entradas,
		MaxEntradas: Config.CacheEntries,
		Algoritmo:   Config.CacheReplacement,
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
	InstanteDeReferencia int  `json:"instante_referencia"` // Marca el instante de referencia para LRU
}

func TLBHabilitada() bool {
	return tlb.MaxEntradas != 0
}

func BuscarDireccion(pagina int) (bool,int) { // devolvemos el frame ya que la pagina esta cargada en el TLB
	i := 0
	for i < len(tlb.Entradas) {
		if tlb.Entradas[i].NumeroPagina == pagina && tlb.Entradas[i].BitPresencia {
			return true,i // La página está en la TLB y es válida
		}
	}
	return false,-1 // La página no está en la TLB o no es válida
}


func AccesoATLB(pid int, nropagina int) int {
	
	if !TLBHabilitada() {
		log.Fatal("TLB no habilitada, no se puede acceder a la TLB")
		return -1 // TLB no habilitada, no se puede acceder a la
	}
	
	bool, indice := BuscarDireccion(nropagina) // Verificamos si la página está en la TLB
	if bool { 
		log.Fatalf("PID: < %d > - TLB HIT - Página: %d", pid, nropagina)
		return tlb.Entradas[indice].NumeroFrame // Si la página está en la TLB, devolvemos el frame y true
	} else {
		log.Println("PID: < %d > - TLB MISS - Página: %d", pid, nropagina)
		return -1
	}
}

func IndiceDeEntradaVictima() int {
		
	if tlb.Algoritmo == "FIFO" {
		return 0 
	} else { // LRU
		tiempoActual := int(time.Now().UnixNano()) // tiempo actual en nanosegundos (convertido a int)
		indice := 0
		victima := 0
		for indice < len(tlb.Entradas) {
			if tlb.Entradas[indice].InstanteDeReferencia < tiempoActual { // Si el instante de referencia es menor al tiempo actual, es una candidata a ser la victima
				victima = indice 
			}
			indice++
		}
		return victima
	}
}

func AgregarEntradaATLB(pid int, nropagina int, nroframe int) {

	nuevaEntrada := EntradaTLB{
		NumeroPagina: nropagina,
		NumeroFrame: nroframe,
		BitPresencia: true, // La pagina esta presente en memoria
		PID: pid, // Asignamos el PID del proceso
		InstanteDeReferencia: int(time.Now().UnixNano()), // Asignamos el instante de referencia actual
	}

	if len(tlb.Entradas) == tlb.MaxEntradas { // si la cantidad de entradas es la maxima => hay que reemplazar
		indiceVictima := IndiceDeEntradaVictima() 
		tlb.Entradas[indiceVictima] = nuevaEntrada // reemplazo  la entrada victima por la nueva entrada
		return 
	} else { // si no esta lleno, agrego la nueva entrada al final
		tlb.Entradas = append(tlb.Entradas, nuevaEntrada)
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
		log.Printf("Entrada %d: PID=%d | Página=%d | Frame=%d | Presente=%v | InstanteReferencia=%d",
			i, entrada.PID, entrada.NumeroPagina, entrada.NumeroFrame, entrada.BitPresencia, entrada.InstanteDeReferencia)
	}
}