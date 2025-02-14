package main

import (
	"flag"
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

	// Definir los comandos disponibles
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	completeCmd := flag.NewFlagSet("complete", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	// Selección de comando
	switch os.Args[1] {
	case "add":
		addTask := addCmd.String("task", "", "Nombre de la tarea a agregar")
		addCmd.Parse(os.Args[2:])
		if *addTask == "" {
			fmt.Println("Debes especificar una tarea con -task")
			os.Exit(1)
		}
		fmt.Println("Tarea agregada:", *addTask)
	case "list":
		listCmd.Parse(os.Args[2:])
		fmt.Println("Mostrando todas las tareas...")
	case "complete":
		completeTaskID := completeCmd.Int("id", -1, "ID de la tarea a marcar como completada")
		completeCmd.Parse(os.Args[2:])
		if *completeTaskID == -1 {
			fmt.Println("Debes especificar un ID de tarea con -id")
			os.Exit(1)
		}
		fmt.Println("Tarea completada con ID:", *completeTaskID)
	case "delete":
		deleteTaskID := deleteCmd.Int("id", -1, "ID de la tarea a eliminar")
		deleteCmd.Parse(os.Args[2:])
		if *deleteTaskID == -1 {
			fmt.Println("Debes especificar un ID de tarea con -id")
			os.Exit(1)
		}
		fmt.Println("Tarea eliminada con ID:", *deleteTaskID)
	default:
		fmt.Println("Comando no reconocido:", os.Args[1])
		fmt.Println("Comandos disponibles: add, list, complete, delete")
		os.Exit(1)
	}
}
