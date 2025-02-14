package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Ruta del archivo CSV donde se almacenan las tareas
const filePath = "tasks.csv"

// Estructura para representar una tarea
type Task struct {
	ID        int
	Name      string
	Completed bool
}

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
		addNewTask(*addTask)
	case "list":
		listCmd.Parse(os.Args[2:])
		showTaskList()
	case "complete":
		completeTaskID := completeCmd.Int("id", -1, "ID de la tarea a marcar como completada")
		completeCmd.Parse(os.Args[2:])
		if *completeTaskID == -1 {
			fmt.Println("Debes especificar un ID de tarea con -id")
			os.Exit(1)
		}
		completeTask(*completeTaskID)
	case "delete":
		deleteTaskID := deleteCmd.Int("id", -1, "ID de la tarea a eliminar")
		deleteCmd.Parse(os.Args[2:])
		if *deleteTaskID == -1 {
			fmt.Println("Debes especificar un ID de tarea con -id")
			os.Exit(1)
		}
		deleteTask(*deleteTaskID)
	default:
		fmt.Println("Comando no reconocido:", os.Args[1])
		fmt.Println("Comandos disponibles: add, list, complete, delete")
		os.Exit(1)
	}
}

// Agregar una nueva tarea al archivo CSV
func addNewTask(taskName string) {
	tasks := loadTasks()

	// Crear una nueva tarea con ID autoincremental
	newID := len(tasks) + 1
	newTask := Task{ID: newID, Name: taskName, Completed: false}

	// Guardar la nueva tarea en el archivo
	tasks = append(tasks, newTask)
	saveTasks(tasks)

	fmt.Println("Tarea agregada: " + taskName + " (" + strconv.Itoa(newID) + ")")
}

// Cargar las tareas desde el archivo CSV
func loadTasks() []Task {
	file, err := os.Open(filePath)
	if err != nil {
		return []Task{} // Si el archivo no existe, devolver una lista vacía
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error leyendo el archivo CSV:", err)
		os.Exit(1)
	}

	var tasks []Task
	for _, line := range lines {
		id, _ := strconv.Atoi(line[0])
		completed, _ := strconv.ParseBool(line[2])
		tasks = append(tasks, Task{ID: id, Name: line[1], Completed: completed})
	}

	return tasks
}

// Guardar las tareas en el archivo CSV
func saveTasks(tasks []Task) {
	file, err := os.Create(filePath)

	if err != nil {
		fmt.Println("Error al escribir en el archivo CSV:", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, task := range tasks {
		writer.Write([]string{
			strconv.Itoa(task.ID),
			task.Name,
			strconv.FormatBool(task.Completed),
		})
	}
}

// Muestra la lista de tareas guardada
func showTaskList() {
	tasks := loadTasks()

	if len(tasks) == 0 {
		fmt.Println("No existe ninguna tarea")
	} else {
		for _, task := range tasks {
			var status string
			if task.Completed {
				status = "Completado"
			} else {
				status = "Pendiente"
			}
			fmt.Printf("%d. %s - %s\n", task.ID, task.Name, status)
		}
	}
}

// Marca una tarea como completada en base a su id
func completeTask(taskID int) {
	// Abrir el archivo CSV original
	inputFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error al abrir el archivo csv:", err)
		return
	}
	defer inputFile.Close()

	// Crear un archivo temporal para escribir los cambios
	outputFile, err := os.Create("tasks_temp.csv")
	if err != nil {
		fmt.Println("Error al crear el archivo temporal:", err)
		return
	}
	defer outputFile.Close()

	// Crear un lector y un escritor CSV
	reader := csv.NewReader(inputFile)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// Definir los parámetros de búsqueda y modificación
	nuevoValor := true
	columnaABuscar := 0    // Índice de la columna donde buscar el elemento x
	columnaAModificar := 2 // Índice de la columna a modificar
	found := false

	// Leer y procesar el archivo línea por línea
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Println("Error al leer una línea:", err)
			return
		}

		// Modificar el registro si se encuentra el elemento x
		if record[columnaABuscar] == strconv.Itoa(taskID) {
			found = true
			if record[columnaAModificar] == "false" {
				record[columnaAModificar] = strconv.FormatBool(nuevoValor)
				fmt.Printf("La tarea %s con id %d ha sido completada\n", record[1], taskID)
			} else {
				fmt.Printf("La tarea %s con id %d ya se encuentra completada\n", record[1], taskID)
			}
			break
		}

		// Escribir el registro (modificado o no) en el archivo temporal
		if err := writer.Write(record); err != nil {
			fmt.Println("Error al escribir una línea:", err)
			return
		}
	}

	// Asegurar que todos los datos se escriban en el archivo
	writer.Flush()

	// Cerrar ambos archivos
	inputFile.Close()
	outputFile.Close()

	// Reemplazar el archivo original con el archivo temporal
	if err := os.Rename("tasks_temp.csv", filePath); err != nil {
		fmt.Println("Error al reemplazar el archivo original:", err)
		return
	}

	if !found {
		fmt.Printf("La tarea con id %d no existe\n", taskID)
	}

}

// Marca una tarea como completada en base a su id
func deleteTask(taskID int) {
	// Abrir el archivo CSV original
	inputFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error al abrir el archivo csv:", err)
		return
	}
	defer inputFile.Close()

	// Crear un archivo temporal para escribir los cambios
	outputFile, err := os.Create("tasks_temp.csv")
	if err != nil {
		fmt.Println("Error al crear el archivo temporal:", err)
		return
	}
	defer outputFile.Close()

	// Crear un lector y un escritor CSV
	reader := csv.NewReader(inputFile)
	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	// Definir los parámetros de búsqueda y modificación
	columnaABuscar := 0 // Índice de la columna donde buscar el elemento x
	found := false

	// Leer y procesar el archivo línea por línea
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Println("Error al leer una línea:", err)
			return
		}

		// Modificar el registro si se encuentra el elemento x
		if record[columnaABuscar] != strconv.Itoa(taskID) {
			// Escribir el registro (modificado o no) en el archivo temporal
			if err := writer.Write(record); err != nil {
				fmt.Println("Error al escribir una línea:", err)
				return
			}
		} else {
			found = true
			fmt.Printf("La tarea %s con id %d ha sido eliminada\n", record[1], taskID)
		}
	}

	// Asegurar que todos los datos se escriban en el archivo
	writer.Flush()

	// Cerrar ambos archivos
	inputFile.Close()
	outputFile.Close()

	// Reemplazar el archivo original con el archivo temporal
	if err := os.Rename("tasks_temp.csv", filePath); err != nil {
		fmt.Println("Error al reemplazar el archivo original:", err)
		return
	}

	if !found {
		fmt.Printf("La tarea con id %d no existe\n", taskID)
	}

}
