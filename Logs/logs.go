package main

import (
	"log"
	"os"
	"io"
	// "fmt"
	"strconv"
	"strings"
	
)

type MetricaP struct {
	PID int 
	AccesosATabladePaginas int 
	InstruccionesSolicitadas int
	BajadasAlSwap int
	SubidasAMP int 
	LecturasDeMemoria int
	EscriturasEnMemoria int
} 

type Instruccion struct {
	Operacion string 
	Operandos []string
}

func ConfigurarLogger() {
	logFile, err := os.OpenFile("pruebasDeLogs.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil { // si err no es nil, significa que ocurrió un error al intentar abrir/crear el archivo
		panic(err) //detiene la ejecución del programa y muestra el error. El programa no continua si no se puede configuarar el logger
	}
	mw := io.MultiWriter(os.Stdout, logFile) // se crea un escritor Writer que envia los datos a multiples destinos. 
	log.SetOutput(mw) // se configura para que le logger envie mensajes a mw, que a suvez los redirige tanto a consola como a tp0.log
}

func MensajeEscritura(pid int, dirFisica int, tamanio int) string {
	return "## PID: <" + strconv.Itoa(pid) + "> - Escritura - Dir. Física: <" + strconv.Itoa(dirFisica) + "> - Tamaño: <" + strconv.Itoa(tamanio) + ">"
}

func MensajeMemoryDump(pid int) string {
	return "## PID: <" + strconv.Itoa(pid) + "> - Memory Dump solicitado"
}

func MensajeObtenerInstruccion(pid int, pc int, instruccion Instruccion) string {
	return "## PID: <" + strconv.Itoa(pid) + "> - Obtener instrucción: <" + strconv.Itoa(pc) + "> - Instrucción: <" + instruccion.Operacion + "> < " + strings.Join(instruccion.Operandos, ", ") + ">"
}

func MensajeCreacionDeProceso(pid int, tamanio uint) string {
	return "## PID: <" + strconv.Itoa(pid) + "> - Proceso Creado - Tamaño: <" + strconv.Itoa(int(tamanio)) + ">"
}

func MensajeDestruccionDeProceso(pid int, metrica MetricaP) string {
	return "## PID: <" + strconv.Itoa(pid) + "> - Proceso Destruido - Metricas - Acc.T.Pag: <" + strconv.Itoa(metrica.AccesosATabladePaginas) + ">; Inst.Sol.: <" + strconv.Itoa(metrica.InstruccionesSolicitadas) + ">; SWAP: <" + strconv.Itoa(metrica.BajadasAlSwap) + ">; Mem.Prin.: <" +
	 strconv.Itoa(metrica.SubidasAMP) + ">; Lec.Mem.: <" + strconv.Itoa(metrica.LecturasDeMemoria) + ">; Esc.Mem.: <" + strconv.Itoa(metrica.EscriturasEnMemoria) + ">"
}



func main() {
	ConfigurarLogger()
	log.Println("Nueva Prueba de log ---- ")
	pid := 1234
	dirFisica := 5678
	tamanio := 1024
	instruccion := Instruccion{
		Operacion: "LOAD",
		Operandos: []string{"R1", "0x1000"},
	}
	metrica := MetricaP{
		PID:                    pid,
		AccesosATabladePaginas: 10,
		InstruccionesSolicitadas: 20,
		BajadasAlSwap:           2,
		SubidasAMP:              3,
		LecturasDeMemoria:       15,
		EscriturasEnMemoria:     5,
	}

	log.Println(MensajeEscritura(pid, dirFisica, tamanio))
	log.Println(MensajeMemoryDump(pid))
	log.Println(MensajeObtenerInstruccion(pid, 0, instruccion))
	log.Println(MensajeCreacionDeProceso(pid, uint(tamanio)))
	log.Println(MensajeDestruccionDeProceso(pid, metrica))
	log.Println("Fin de la prueba de log ---- ")


	
}