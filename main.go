package main

import (
	"fmt"
	"os"
)

func main() {

	// Comprueba el número de argumentos
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <comando> [opciones]")
		fmt.Println("Comandos disponibles: add, list, complete, delete")
		os.Exit(1)
	}

	// Selección de comando
	switch os.Args[1] {
	case "add":
		fmt.Println("Comando: Agregar tarea")
	case "list":
		fmt.Println("Comando: Listar tareas")
	case "complete":
		fmt.Println("Comando: Completar tarea")
	case "delete":
		fmt.Println("Comando: Eliminar tarea")
	default:
		fmt.Println("Comando no reconocido:", os.Args[1])
		fmt.Println("Comandos disponibles: add, list, complete, delete")
		os.Exit(1)
	}
}
