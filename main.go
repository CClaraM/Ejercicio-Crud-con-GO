package main // Selecciona la funcion main como punto de arranque del programa

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	_ "github.com/go-sql-driver/mysql"
)

// Variables globales
var (
	datos       [][]string // Datos tabla
	original    [][]string // Datos consulta
	tabla       *widget.Table
	ventana     fyne.Window
	labelEstado *widget.Label
)

func conexion() (*sql.DB, error) {
	// su estructura es usuario:pasword@tcp(host:puerto)/base de datos
	dsn := "root:@tcp(127.0.0.1:3306)/tiendactpi"
	return sql.Open("mysql", dsn)
}

func cargarProductos() {
	db, err := conexion()
	if err != nil {
		labelEstado.SetText("Error de conexión")
		dialog.ShowError(fmt.Errorf("Error de conexión"), ventana)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, nombre, precio, cantidad FROM productos")
	if err != nil {
		labelEstado.SetText("Error al listar")
		dialog.ShowError(fmt.Errorf("Error al listar"), ventana)
		return
	}
	defer rows.Close()

	// Dates iniciales de la tabla Agregamos columna "Eliminar" al final
	datos = [][]string{{"ID", "Nombre", "Precio", "Cantidad", "Eliminar"}}
	original = [][]string{{"ID", "Nombre", "Precio", "Cantidad", "Eliminar"}}

	for rows.Next() {
		var id, cantidad int
		var nombre string
		var precio float64
		rows.Scan(&id, &nombre, &precio, &cantidad)

		fila := []string{
			fmt.Sprintf("%d", id),
			nombre,
			fmt.Sprintf("%.2f", precio),
			fmt.Sprintf("%d", cantidad),
			"", // columna eliminar, no guarda texto
		}
		datos = append(datos, fila)
		original = append(original, append([]string{}, fila...)) // copia
	}

	if tabla != nil {
		tabla.Refresh()
	}
}

// funcion para comparar la fila almacenada y la editada, necesaria para guardar los registros editados
func igualFila(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func main() {
	miApp := app.New()
	ventana = miApp.NewWindow("CRUD Productos - Go + Fyne")

	// Entradas
	inputNombre := widget.NewEntry()
	inputPrecio := widget.NewEntry()
	inputCantidad := widget.NewEntry()
	labelEstado = widget.NewLabel("")

	// Botón Insertar
	btnInsertar := widget.NewButton("Agregar Producto", func() {
		nombre := inputNombre.Text
		precio, err := strconv.ParseFloat(inputPrecio.Text, 64)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Precio inválido"), ventana)
			return
		}
		cantidad, err := strconv.Atoi(inputCantidad.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Cantidad inválida"), ventana)
			return
		}

		db, err := conexion()
		if err != nil {
			dialog.ShowError(fmt.Errorf("Error de conexión"), ventana)
			return
		}
		defer db.Close()

		_, err = db.Exec("INSERT INTO productos (nombre, precio, cantidad) VALUES (?, ?, ?)", nombre, precio, cantidad)
		if err != nil {
			dialog.ShowError(fmt.Errorf("No se pudo insertar: %v", err), ventana)
			return
		}

		inputNombre.SetText("")
		inputPrecio.SetText("")
		inputCantidad.SetText("")
		ventana.Canvas().Focus(inputNombre)

		dialog.ShowInformation("Éxito", "Producto agregado correctamente", ventana)
		cargarProductos()
	})

	// Inicializar datos
	datos = [][]string{{"ID", "Nombre", "Precio", "Cantidad", "Eliminar"}}
	original = [][]string{{"ID", "Nombre", "Precio", "Cantidad", "Eliminar"}}

	// Tabla con Entries editables + botón eliminar
	tabla = widget.NewTable(
		func() (int, int) { return len(datos), len(datos[0]) },
		func() fyne.CanvasObject {
			e := widget.NewEntry()
			b := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
			return container.NewMax(e, b)
		},
		func(id widget.TableCellID, o fyne.CanvasObject) {
			cont := o.(*fyne.Container)
			var e *widget.Entry
			var b *widget.Button
			for _, ch := range cont.Objects {
				switch v := ch.(type) {
				case *widget.Entry:
					e = v
				case *widget.Button:
					b = v
				}
			}
			if e == nil || b == nil {
				return
			}

			row, col := id.Row, id.Col

			// limpiar handler viejo antes de SetText
			// Necesario para fyne para que no recicle datos
			e.OnChanged = nil
			if row < len(datos) && col < len(datos[row]) {
				e.SetText(datos[row][col])
			} else {
				e.SetText("")
			}

			// Cabecera
			if row == 0 {
				e.Disable()
				b.Hide()
				return
			}

			// Última columna = Eliminar
			if col == len(datos[0])-1 {
				e.Hide()
				b.Show()
				r := row
				idProd := datos[r][0]
				nombre := datos[r][1]
				b.OnTapped = func() {
					dialog.ShowConfirm("Eliminar",
						fmt.Sprintf("¿Seguro que deseas eliminar el producto ID %s (%s)?", idProd, nombre),
						func(ok bool) {
							if !ok {
								return
							}
							db, err := conexion()
							if err != nil {
								dialog.ShowError(fmt.Errorf("Error de conexión"), ventana)
								return
							}
							defer db.Close()

							_, err = db.Exec("DELETE FROM productos WHERE id=?", idProd)
							if err != nil {
								dialog.ShowError(fmt.Errorf("Error eliminando: %v", err), ventana)
								return
							}
							dialog.ShowInformation("Éxito", "Producto eliminado", ventana)
							cargarProductos()
						}, ventana)
				}
				return
			}

			// Columna ID = no editable
			if col == 0 {
				e.Show()
				b.Hide()
				e.Disable()
				return
			}

			// Columnas editables
			e.Show()
			b.Hide()
			e.Enable()
			localRow, localCol := row, col
			e.OnChanged = func(s string) {
				if localRow < len(datos) && localCol < len(datos[localRow]) {
					datos[localRow][localCol] = s
				}
			}
		},
	)

	// Anchos de columnas
	tabla.SetColumnWidth(0, 60)  // ID
	tabla.SetColumnWidth(1, 200) // Nombre
	tabla.SetColumnWidth(2, 100) // Precio
	tabla.SetColumnWidth(3, 100) // Cantidad
	tabla.SetColumnWidth(4, 85)  // Eliminar

	tablaScroll := container.NewVScroll(tabla)
	tablaScroll.SetMinSize(fyne.NewSize(510, 320))

	// Botón Listar
	btnListar := widget.NewButton("Listar Productos", func() {
		cargarProductos()
	})

	// Botón Actualizar
	btnActualizar := widget.NewButton("Actualizar Productos", func() {
		cambios := []string{}
		for i := 1; i < len(datos); i++ { // saltar cabecera
			for j := 1; j < len(datos[i])-1; j++ { // saltar ID y Eliminar
				if datos[i][j] != original[i][j] { // agrega una lista de los cambios
					cambios = append(cambios,
						fmt.Sprintf("ID %s: %s → %s",
							datos[i][0], original[i][j], datos[i][j]))
				}
			}
		}

		if len(cambios) == 0 {
			dialog.ShowInformation("Sin cambios", "No hay cambios para guardar", ventana)
			return
		}

		resumen := "Se detectaron los siguientes cambios:\n\n" +
			strings.Join(cambios, "\n")

		dialog.ShowConfirm("Confirmar actualización", resumen, func(ok bool) {
			if !ok {
				cargarProductos()
				return
			}

			db, err := conexion()
			if err != nil {
				dialog.ShowError(fmt.Errorf("Error de conexión"), ventana)
				return
			}
			defer db.Close()

			for i := 1; i < len(datos); i++ { // Recorre el arreglo de valores editados e inserta en la base de datos
				if !igualFila(datos[i], original[i]) {
					id := datos[i][0]
					nombre := datos[i][1]
					precio, _ := strconv.ParseFloat(datos[i][2], 64)
					cantidad, _ := strconv.Atoi(datos[i][3])

					_, err = db.Exec("UPDATE productos SET nombre=?, precio=?, cantidad=? WHERE id=?",
						nombre, precio, cantidad, id)
					if err != nil {
						dialog.ShowError(fmt.Errorf("Error actualizando ID %s: %v", id, err), ventana)
						return
					}
				}
			}

			dialog.ShowInformation("Éxito", "Cambios guardados", ventana)
			cargarProductos()
		}, ventana)
	})

	// Layout
	form := container.NewVBox(
		widget.NewLabel("Nombre:"), inputNombre,
		widget.NewLabel("Precio:"), inputPrecio,
		widget.NewLabel("Cantidad:"), inputCantidad,
		container.NewHBox(btnInsertar, btnListar, btnActualizar),
		tablaScroll,
		labelEstado,
	)

	ventana.SetContent(form)
	ventana.Resize(fyne.NewSize(580, 550))
	ventana.SetFixedSize(true)
	cargarProductos()
	ventana.ShowAndRun()
}
