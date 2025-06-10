package main 

import (
	"math/rand"
	"time"
	"fmt"
	"errors"
)

type Tabla struct {
	Punteros []*Tabla
	Valores []int 
}

// Validaciones -----------------------------------------

func ValidarIndicesEnRango(direccion []int, config Config) error {
	for i := 0; i < config.Niveles; i++ {
		if direccion[i] < 0 || direccion[i] >= config.Entradas {
			return fmt.Errorf("índice fuera de rango en el nivel %d: %d", i+1, direccion[i])
		}
	}
	return nil
}

func ValidarPunterosExistentes(tabla *Tabla, direccion []int, config Config) error {
	actual := tabla
	for i := 0; i < config.Niveles-1; i++ {
		idx := direccion[i]
		if idx >= len(actual.Punteros) || actual.Punteros[idx] == nil {
			return fmt.Errorf("puntero nulo en nivel %d, índice %d", i+1, idx)
		}
		actual = actual.Punteros[idx]
	}
	return nil
}

func ValidarValorFinal(tabla *Tabla, direccion []int, config Config) error {
	actual := tabla
	for i := 0; i < config.Niveles-1; i++ {
		actual = actual.Punteros[direccion[i]]
	}
	ultimoIndice := direccion[config.Niveles-1]
	if ultimoIndice >= len(actual.Valores) {
		return fmt.Errorf("índice fuera de rango en nivel final: %d", ultimoIndice)
	}
	return nil
}

func ValidarDesplazamiento(direccion []int, tamañoFrame int, config Config) error {
	desplazamiento := direccion[len(direccion)-1]
	if desplazamiento < 0 || desplazamiento >= tamañoFrame {
		return fmt.Errorf("desplazamiento fuera del rango del frame: %d", desplazamiento)
	}
	return nil
}

func VerificacionTotal(tabla *Tabla, direccion []int, config Config, tamañoFrame int) error {
	if err := ValidarIndicesEnRango(direccion, config); err != nil {
		return err
	}
	if err := ValidarPunterosExistentes(tabla, direccion, config); err != nil {
		return err
	}
	if err := ValidarValorFinal(tabla, direccion, config); err != nil {
		return err
	}
	if err := ValidarDesplazamiento(direccion, tamañoFrame, config); err != nil {
		return err
	}
	return nil
}

// Mostra Tabla -----------------------------------------

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

// Cración de tabla -----------------------------------------

func CreaTablaJerarquica(niveles int, entradasPorNivel int) *Tabla {
	if niveles == 0 || entradasPorNivel == 0 {
		return nil
	}

	tabla := &Tabla{}

	nivelesRestantes := niveles - 1

	if nivelesRestantes == 0 {
		// Nivel hoja: asignar valores aleatorios
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < entradasPorNivel; i++ {
			valorAleatorio := rand.Intn(1000) // rango arbitrario
			tabla.Valores = append(tabla.Valores, valorAleatorio)
		}
	} else {
		// Nivel intermedio: asignar punteros a nuevas tablas
		for i := 0; i < entradasPorNivel; i++ {
			tabla.Punteros = append(tabla.Punteros, CreaTablaJerarquica(nivelesRestantes, entradasPorNivel))
		}
	}

	return tabla
}

func MostrarPorNiveles(tabla *Tabla) {
	niveles := make(map[int][]*Tabla)
	recolectarPorNivel(tabla, 1, niveles)

	flag := false
	for i := 1; ; i++ {
		tabs, existe := niveles[i]
		if !existe {
			break
		}
		fmt.Printf("Nivel %d:\n", i)
		for j, t := range tabs {
			if len(t.Valores) > 0 {
				fmt.Printf("  Tabla %d: Valores: %v\n", j, t.Valores)
			} else {
				if !flag {
					fmt.Printf("  Tabla %d: [Tabla Maestra. N = 1]\n", j)
					flag = true
				} else {
					fmt.Printf("  Tabla %d: [tabla intermedia]\n", j)
				}
			}
		}
	}
}

// recolectarPorNivel es auxiliar para llenar el mapa niveles[nivel] = []*Tabla
func recolectarPorNivel(tabla *Tabla, nivel int, niveles map[int][]*Tabla) {
	if tabla == nil {
		return
	}
	niveles[nivel] = append(niveles[nivel], tabla)
	for _, subtabla := range tabla.Punteros {
		recolectarPorNivel(subtabla, nivel+1, niveles)
	}
}


// Traducción de direcciones -----------------------------------------

type Config struct {
	Niveles   int // Número de niveles en la tabla jerárquica
	Entradas  int // Número de entradas por nivel
}

func TraducirDireccion(config Config, tabla *Tabla, direccion []int) (int, error) {
	
	if len(direccion) != config.Niveles+1 {
		return 0, fmt.Errorf("la dirección debe tener %d niveles + 1 desplazamiento", config.Niveles)
	}

	actual := tabla

	for nivel := 0; nivel < config.Niveles; nivel++ {
		indice := direccion[nivel]
		if indice < 0 || indice >= config.Entradas {
			return 0, fmt.Errorf("índice fuera de rango en nivel %d: %d", nivel+1, indice)
		}

		if nivel == config.Niveles-1 {
			// Último nivel: valores
			if indice >= len(actual.Valores) {
				return 0, fmt.Errorf("índice %d fuera de rango en valores del último nivel", indice)
			}
			return actual.Valores[indice], nil
		} else {
			if indice >= len(actual.Punteros) || actual.Punteros[indice] == nil {
				return 0, fmt.Errorf("puntero nulo o fuera de rango en nivel %d", nivel+1)
			}
			actual = actual.Punteros[indice]
		}
	}

	return 0, errors.New("no se encontró frame")
}

// Main -----------------------------------------

func main() {

	niveles := 3
	entradasPorNivel := 2

	config := Config{
		Niveles:  niveles,
		Entradas: entradasPorNivel,
	}
	
	fmt.Println("")
	fmt.Println("N = niveles = ", niveles)
	fmt.Println("n = entradas por nivel = ", entradasPorNivel)
	fmt.Println("")

	tabla := CreaTablaJerarquica(niveles, entradasPorNivel)
	fmt.Println("Tabla jerarquica en arbol:")
	MostrarTablaArbol(tabla, "",true)

	fmt.Println("")
	fmt.Println("")


	// Dirección en base 0 => como vectores 
	direccion := []int{0,1,0,0}
	fmt.Println("Dirección: ", direccion)
	
	frame, err := TraducirDireccion(config, tabla, direccion)
	
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Frame correspondiente %d\n", frame)
	}
}