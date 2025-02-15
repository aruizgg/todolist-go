# T0D0List

Una aplicación CLI de gestión de tareas desarrollada en Go.

## Descripción
T0D0List es una aplicación de línea de comandos para gestionar tareas, desarrollada en Go. Ofrece una interfaz de usuario intuitiva y atractiva gracias al uso de la biblioteca `lipgloss` para el estilizado de la consola.

## Características
- Agregar nuevas tareas
- Listar tareas existentes
- Marcar tareas como completadas
- Eliminar tareas
- Interfaz de usuario interactiva con navegación por teclado
- Almacenamiento persistente de tareas en un archivo CSV

## Requisitos
- Go 1.x o superior
- Bibliotecas externas:
  - [`github.com/charmbracelet/lipgloss`](https://github.com/charmbracelet/lipgloss)
  - [`github.com/eiannone/keyboard`](https://github.com/eiannone/keyboard)
  - [`github.com/inancgumus/screen`](https://github.com/inancgumus/screen)

## Instalación
Clona el repositorio:
```bash
git clone https://github.com/tu-usuario/T0D0List.git
cd T0D0List
```

Instala las dependencias:
```bash
go mod tidy
```

Compila y ejecuta la aplicación:
```bash
go run main.go
```

## Uso
Al ejecutar la aplicación, se mostrará un menú interactivo con las siguientes opciones:

1. **Agregar tarea**
2. **Listar tareas**
3. **Completar tarea**
4. **Eliminar tarea**
5. **Salir**

Utiliza las teclas de flecha para navegar por el menú y presiona `Enter` para seleccionar una opción.

## Estructura del Proyecto
```
T0D0List/
├── main.go
├── go.mod
├── go.sum
└── tasks.csv
```
- **`main.go`**: Contiene el código fuente de la aplicación.
- **`tasks.csv`**: Archivo donde se almacenan las tareas.

## Autor
Desarrollado por **aruizgg**.

