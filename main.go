package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/eiannone/keyboard"
	"github.com/inancgumus/screen"
)

// Ruta del archivo CSV donde se almacenan las tareas
const filePath = "tasks.csv"

// Estructura para representar una tarea
type Task struct {
	ID        int
	Name      string
	Completed bool
}

var colorInit = lipgloss.Color("#00e4ff")

var initStyle = lipgloss.NewStyle().Align().
	Bold(true).
	Italic(true).
	Foreground(colorInit).
	Padding(1, 3).
	Margin(1).
	Align(lipgloss.Center)
var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
var optionSelectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#38d1e3"))
var optionNotSelectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))
var placeHolderStyle = lipgloss.NewStyle().Italic(true)

func main() {

	options := []string{"Agregar tarea", "Listar tareas", "Completar tarea", "Eliminar tarea", "Salir"}
	selectedIndex := 0

	for {
		// Limpiar la pantalla
		screen.Clear()
		screen.MoveTopLeft()

		fmt.Print(initStyle.Render("T0D0List"), placeHolderStyle.Render("by aruizgg"))
		fmt.Println()

		// Mostrar el menú
		fmt.Println(helpStyle.Render("\nSeleccione una opción:"))
		for i, option := range options {
			if i == selectedIndex {
				fmt.Printf(optionSelectedStyle.Render("> %s"), option)
				fmt.Println()

			} else {
				fmt.Printf(optionNotSelectedStyle.Render("%s"), option)
				fmt.Println()

			}
		}

		// Capturar la tecla presionada
		_, key, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}

		// Manejar la navegación
		switch key {
		case keyboard.KeyArrowUp:
			selectedIndex = (selectedIndex - 1 + len(options)) % len(options)
		case keyboard.KeyArrowDown:
			selectedIndex = (selectedIndex + 1) % len(options)
		case keyboard.KeyEnter:
			// Ejecutar la acción seleccionada
			screen.Clear()
			screen.MoveTopLeft()
			switch selectedIndex {
			case 0:
				handleAddCommand()
			case 1:
				showTaskList()
			case 2:
				handleCompleteCommand()
			case 3:
				handleDeleteCommand()
			case 4:
				return
			}
			fmt.Println("\nPresione cualquier tecla para volver al menú...")
			keyboard.GetSingleKey()
		case keyboard.KeyEsc:
			return
		}
	}
}

func handleAddCommand() {
	fmt.Println("Ingrese el nombre de la tarea:")
	reader := bufio.NewReader(os.Stdin)
	task, _ := reader.ReadString('\n')
	task = strings.TrimSpace(task)
	if task != "" {
		addNewTask(task)
	} else {
		fmt.Println("La tarea no puede estar vacía.")
	}
}

func handleCompleteCommand() {
	fmt.Print("Ingrese el ID de la tarea a completar: ")
	id := selectTask()
	completeTask(id)
}

func handleDeleteCommand() {
	fmt.Print("Ingrese el ID de la tarea a eliminar: ")
	id := selectTask()
	deleteTask(id)
}

func selectTask() int {

	selectedIndex := 0
	tasks := loadTasks()
	seleccionado, fin := false, false
	for !fin {

		// Limpiar la pantalla
		screen.Clear()
		screen.MoveTopLeft()

		fmt.Println("\nSeleccione una opción:")
		for i, task := range tasks {
			var status string
			if task.Completed {
				status = "Completado"
			} else {
				status = "Pendiente"
			}
			if i == selectedIndex {
				fmt.Printf(optionSelectedStyle.Render("> %d. %s - %s\n"), task.ID, task.Name, status)
				fmt.Println()
			} else {
				fmt.Printf(optionNotSelectedStyle.Render("%d. %s - %s\n"), task.ID, task.Name, status)
				fmt.Println()
			}
		}

		// Capturar la tecla presionada
		_, key, err := keyboard.GetSingleKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyEnter {

		}
		// Manejar la navegación
		switch key {
		case keyboard.KeyArrowUp:
			selectedIndex = (selectedIndex - 1 + len(tasks)) % len(tasks)
		case keyboard.KeyArrowDown:
			selectedIndex = (selectedIndex + 1) % len(tasks)
		case keyboard.KeyEnter:
			seleccionado, fin = true, true
		case keyboard.KeyEsc:
			fin = true
		}
	}
	if seleccionado {
		return tasks[selectedIndex].ID
	} else {
		return 0
	}
}

// Agregar una nueva tarea al archivo CSV
func addNewTask(taskName string) {
	tasks := loadTasks()

	// Asignación del primer ID disponible
	maxID := 0
	usedIDs := make(map[int]bool, len(tasks))

	for _, task := range tasks {
		usedIDs[task.ID] = true
		if task.ID > maxID {
			maxID = task.ID
		}
	}

	newID := 1
	for newID <= maxID {
		if !usedIDs[newID] {
			break
		}
		newID++
	}

	if newID > maxID {
		newID = maxID + 1
	}

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

	if taskID == 0 {
		return
	}
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
		}

		// fmt.Printf(" Se va a escribir %d. %s - %s\n", record[0], record[1], record[2])

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

	if taskID == 0 {
		return
	}

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
