package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"log"
	"io"
)

/*

Cosas para testear:

- Agregado de pagina a cache => Testeado
- Remplazo de pagina en cache 
- Verificacion de si una pagina esta en cache => Testeado
- Verificacion de si una pagina fue modificada => Testeado
- Envio de pagina a memoria 
*/

type ConfigStruct struct {
	CacheEntries int // Cantidad de entradas en la caché
	CacheReplacement string // Algoritmo de reemplazo de caché (FIFO, LRU, CLOCK)
	IPMemory string // IP del servidor de memoria
	PortMemory int // Puerto del servidor de memoria
}

var Config ConfigStruct = ConfigStruct{
	CacheEntries: 10, // Por ejemplo, 10 entradas en la caché
	CacheReplacement: "CLOCK", // Algoritmo de reemplazo de caché
	IPMemory: "127.0.0.1", // IP del servidor de memoria
	PortMemory: 8002, // Puerto del servidor de memoria
}

type PaginaCache struct {
	NumeroPagina int // Numero de pagina en la tabla de paginas
	BitPresencia bool // Indica si el frame esta presente en memoria
	BitModificado bool // Indica si el frame ha sido modificado
	BitDeUso bool // Indica si el frame ha sido usado recientemente
	PID int // Identificador del proceso al que pertenece el frame
	Contenido []byte // Contenido de la pagina
}

type CacheStruct struct {
	Paginas []PaginaCache 
	Algoritmo string 
	Clock int // dato para saber donde quedó la "aguja" del clock
}

var Cache CacheStruct = InicializarCache()

func MostrarPaginaCache(pagina PaginaCache) {
	log.Printf("PaginaCache: NumeroPagina: %d, BitPresencia: %t, BitModificado: %t, BitDeUso: %t, PID: %d, Contenido: %s",
		pagina.NumeroPagina, pagina.BitPresencia, pagina.BitModificado, pagina.BitDeUso, pagina.PID, string(pagina.Contenido))
}

func MostrarCache() {
	
	log.Println("Contenido de la caché => ")
	for i:=0; i < Config.CacheEntries; i++ {
		MostrarPaginaCache(Cache.Paginas[i])
	}
}

func InicializarCache() CacheStruct {
	paginas := make([]PaginaCache, Config.CacheEntries) // Slice vacío, capacidad predefinida
	
		for i := 0; i < Config.CacheEntries; i++ {
		paginas[i] = PaginaCache{
			NumeroPagina: -1,
			PID: -1,
			}
		}
	return CacheStruct{
		Paginas: paginas,
		Algoritmo: Config.CacheReplacement,
	}
}

func CacheHabilitado() bool {
	return len(Cache.Paginas) > 0 
}

func FueModificada(pagina PaginaCache) bool {
	return pagina.BitModificado
}

func EstaEnCache(pid uint, nropagina int) bool {
	if !CacheHabilitado() {
		log.Fatalf("Caché no habilitada, no se puede verificar si la página está en caché")
		return false 
	}

	for _, pagina := range Cache.Paginas {
		if pagina.PID == int(pid) && pagina.NumeroPagina == nropagina && pagina.BitPresencia {
			return true // La página está en la caché
		}
	}
	return false 
}

func ObtenerPaginaDeCache(pid uint, nropagina int) int {
	
	if !CacheHabilitado() {
		log.Fatalf("Caché no habilitada, no se puede obtener la página de caché")
		return -1
	}

	for i, pagina := range Cache.Paginas {
		if pagina.PID == int(pid) && pagina.NumeroPagina == nropagina && pagina.BitPresencia {
			log.Println("Página encontrada en caché: PID %d, Página %d", pid, nropagina)
			return i // Retorna la página y su índice en caché
		}
	}
	return -1
}

func MandarDatosAMP(paginas PaginaCache) {
	url := fmt.Sprintf("http://%s:%d/actualizarMP", Config.IPMemory, Config.PortMemory)
	body, err := json.Marshal(paginas)
	if err != nil {
		log.Println("Error al serializar la pagina a JSON:", err)
		return
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Error al enviar la pagina a la memoria:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error al enviar la pagina a la memoria, status code: %d", resp.StatusCode)
		return
	}
	log.Println("Pagina enviada a la memoria correctamente")
}

func PaginasModificadas() []PaginaCache {
	var paginasModificadas []PaginaCache
	for _, pagina := range Cache.Paginas {
		if FueModificada(pagina) {
			paginasModificadas = append(paginasModificadas, pagina)
		}
	}
	return paginasModificadas
}

// Debe venir una request de memoria o kernel
func DesaolojoDeProceso(w http.ResponseWriter, r *http.Request){
	
	modificadas := PaginasModificadas()
	
	if len(modificadas) == 0 {
		w.Write([]byte("No hay paginas modificadas. No se actualiza memoria"))
		w.WriteHeader(http.StatusOK) // No hay paginas modificadas, todo bien
		return 
	}

	for i:=0; i < len(modificadas); i++ {
		// Consulto direccion fisica => TLB
		// contenido := modificadas[i].Contenido
		// Write de su contenido => pegarle al endpoint de memoria wirite
		// eliminar todas las entradas del caché 
		return
	}
}

func EliminarEntradasDeCache(pid uint) {
	log.Printf("Eliminando entradas de caché para el PID %d", pid)
	for i := 0; i < len(Cache.Paginas); i++ {
		if Cache.Paginas[i].PID == int(pid) { // Si el PID coincide, eliminamos la entrada
			Cache.Paginas[i] = PaginaCache{} // Reinicializamos la entrada
			log.Println("Entrada de caché eliminada para el PID %d", pid)
		}
	}
}

func CreacionDePaginaCache(pid uint, nropagina int, contenido []byte) PaginaCache {
	return PaginaCache{
		NumeroPagina: nropagina,
		BitPresencia: true, // La pagina esta presente en memoria
		BitModificado: false, // Inicialmente no ha sido modificada
		BitDeUso: true, // Inicialmente se considera que la pagina ha sido usada
		PID: int(pid), // Asignamos el PID del proceso
		Contenido: contenido, // Asignamos el contenido de la pagina
	}
}

func PedirFrameAMemoria(pid uint, nropagina int) (PaginaCache, error) {
	
	direccionFisica := MMU(pid, nropagina)
	log.Printf("Pidiendo frame a memoria para PID %d, Numero de pagina %d, Direccion fisica %d", pid, nropagina, direccionFisica) 
	
	url := fmt.Sprintf("http://%s:%d/pedirFrame?pid=%d&direccion=%d", Config.IPMemory, Config.PortMemory, pid, direccionFisica)
	log.Printf("GET HECHO PARA PEDIR FRAME")
	
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error al pedir el frame a memoria: %v", err)
		return PaginaCache{}, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error al pedir el frame a memoria, status code: %d", resp.StatusCode)
		return PaginaCache{}, fmt.Errorf("error al pedir el frame a memoria, status code: %d", resp.StatusCode)
	}

	log.Println("Frame recibido de memoria, decodificando...")
	var frame []byte 
	err = json.NewDecoder(resp.Body).Decode(&frame)
	if err != nil {
		log.Fatalf("Error al decodificar el frame: %v", err)
		return PaginaCache{}, err
	}

	paginaCache := CreacionDePaginaCache(pid, nropagina, frame) 

	return paginaCache, nil
}

func AgregarPaginaACache(pagina PaginaCache) {
	
	log.Println("Agregando pagina a cache")
	
	if len(Cache.Paginas) == Config.CacheEntries {
		RemplazarPaginaEnCache(pagina) // Reemplazamos una pagina segun el algoritmo de reemplazo
		if FueModificada(pagina) {
			log.Println("Pagina modificada, escribiendo en memoria")
			MandarDatosAMP(pagina) 
		}
		return 
	} else {
		Cache.Paginas = append(Cache.Paginas, pagina)
		log.Println("Pagina agregada a la Cache") 
		return 
	}
}

func RemplazarPaginaEnCache(pagina PaginaCache) {
	indiceVictima := IndiceDeCacheVictima() // Obtenemos el indice de la pagina victima

	if FueModificada(Cache.Paginas[indiceVictima]) { // Si la pagina victima fue modificada, debemos escribir su contenido en memoria
		log.Println("Pagina modificada, escribiendo en memoria")
		MandarDatosAMP(Cache.Paginas[indiceVictima]) // Enviamos la pagina a memoria
	}
	Cache.Paginas[indiceVictima] = pagina // Reemplazamos la pagina victima por la nueva pagina
	log.Println("Pagina reemplazada en Cache") 
}


func EscribirEnCache(pid uint, adress int, data string) {

	indice := ObtenerPaginaDeCache(pid, adress)
	if indice == -1 {
		log.Fatalf("Error al obtener la pagina de Cache")
		return
	}

	pagina := Cache.Paginas[indice].Contenido
	copy(pagina[adress:], []byte(data)) // Escribimos el contenido en la pagina de Cache
	Cache.Paginas[indice].Contenido = pagina // Actualizamos el contenido de la pagina en Cache
	Cache.Paginas[indice].BitModificado = true // Marcamos la pagina como modificada
	log.Println("Pagina escrita en Cache: PID %d, Direccion %d, Contenido %s", pid, adress, data)
}

func LeerDeCache(pid uint, adress int, tam int) []byte {
	
	indice := ObtenerPaginaDeCache(pid, adress)
	
	if indice == -1 {
		log.Fatalf("Fatalf al obtener la pagina de Cache")
		return nil
	}

	if indice < 0 || indice >= len(Cache.Paginas) {
		log.Fatalf("Indice de pagina fuera de rango: %d", indice)
		return nil 
	}

	pagina := Cache.Paginas[indice]
	if pagina.BitPresencia && pagina.PID == int(pid) {
		contenido := pagina.Contenido[adress:adress+tam] // Leemos el contenido de la pagina en Cache
		return contenido 
	} else {
		log.Fatalf("Pagina no encontrada en Cache o no pertenece al PID %d", pid)
		return nil 
	}
}

// Para CLOCK-M

func IndiceDeCacheVictima() int {

	if Cache.Algoritmo == "CLOCK" {
		for {
			i := Cache.Clock 
			if !Cache.Paginas[i].BitDeUso {
				Cache.Clock = (i + 1) % len(Cache.Paginas) // Avanzamos al siguiente indice circularmente => por si llegamos al final del vector, poder volver al inicio
				return i
			}
			Cache.Paginas[i].BitDeUso = false  // false = 1
			Cache.Clock = (i + 1) % len(Cache.Paginas) // Avanzamos al siguiente indice circularmente => por si llegamos al final del vector, poder volver al inicio
		}  
	} else {
		i := 0
		for i < len(Cache.Paginas) {
			if !Cache.Paginas[i].BitDeUso && !Cache.Paginas[i].BitModificado {
				Cache.Paginas[i].BitDeUso = true 
				return i // Retorna el indice de la primera pagina con bits 00
			} else {
				if !Cache.Paginas[i].BitDeUso && Cache.Paginas[i].BitModificado { 
					Cache.Paginas[i].BitDeUso = true
					return i;
				}
			}
		}
	}
	return -1 // Si no se encuentra una pagina con bits 00, retorna -1
}

func TraducirDireccion(pid uint, direccion int) int {

	paginaLogica := direccion / tamanioPagina 
	offset := Desplazamiento(direccion, tamanioPagina) // Desplazamiento dentro de la página

	// 1. Preguntamos a TLB
	if AccesoATLB(int(pid), paginaLogica) != -1 {
		frame := AccesoATLB(int(pid), paginaLogica) // Obtenemos el frame desde la TLB
		return frame * tamanioPagina + offset // Retornamos la dirección física
	} 

	// 2. Si no está en TLB, buscamos en la tabla de páginas
	direccionFisica := MMU(pid, direccion) // Obtenemos el frame físico correspondiente a la página lógica
	if direccionFisica == -1 {
		log.Printf("Error al traducir la dirección lógica %d para el PID %d", direccion, pid)
		return -1 // Retornamos -1 para indicar que no se pudo traducir la dirección
	}

	frameFisico := direccionFisica / tamanioPagina 

	// HUBO MISS => AGREGAR A TLB
	AgregarEntradaATLB(int(pid), paginaLogica, frameFisico) // Agregamos la entrada a la TLB
	return direccionFisica// Retornamos la dirección física
}

func Write(pid uint, inst WriteInstruction) {
	
	// Verificar si la página esta en Cache
	if EstaEnCache(pid, inst.Address) {
		log.Println("Página encontrada en caché, escribiendo directamente en caché")
		nropagina := inst.Address / tamanioPagina // Obtenemos el número de página
		EscribirEnCache(pid, nropagina, inst.Data) // Escribimos en la caché
		return
	}

	direccionFisica := TraducirDireccion(pid, inst.Address) // Traducimos la dirección lógica a física
	if direccionFisica == -1 {
		log.Fatalf("Error al traducir la dirección lógica %d para el PID %d", inst.Address, pid)
		return
	}

	inst = WriteInstruction{
		Address: direccionFisica, // Asignamos la dirección física
		Data:    inst.Data, // Asignamos los datos a escribir
		PID:     pid, // Asignamos el PID del proceso
	}

	resp := EnviarMensaje(Config.IPMemory, Config.PortMemory, "write" , inst)
	if resp != "OK" {
		log.Fatalf("Error al escribir en memoria para el PID %d, dirección %d", pid, inst.Address)
		return
	}

	log.Println("Escritura exitosa. - Agregando página a caché")
	
	// Si la página no estaba en cache, pedirla a memoria
	pagina, err := PedirFrameAMemoria(pid, inst.Address)
	if err != nil {
		log.Fatalf("Error al pedir el frame a memoria: %v", err)
		return
	}

	log.Println("PAGINA CACHE CREADA: ", pagina)
	AgregarPaginaACache(pagina)
}

type Respuesta struct {
	Mensaje string
}


func EnviarMensaje(ip string, puerto int, endpoint string, mensaje any) string {
	body, err := json.Marshal(mensaje)
	if err != nil {
		log.Fatalf("No se pudo codificar el mensaje (%v)", err)
		return ""
	}

	url := fmt.Sprintf("http://%s:%d/%s", ip, puerto, endpoint)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("No se pudo enviar mensaje a %s:%d/%s (%v)", ip, puerto, endpoint, err)
		return ""
	}
	defer resp.Body.Close()

	var resData Respuesta
	respuesta, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	err = json.Unmarshal(respuesta, &resData)
	if err != nil {
		return ""
	}

	// log.Printf("respuesta del servidor: %s", resp.Status)
	log.Printf("Respuesta =", resData)
	return resData.Mensaje
}