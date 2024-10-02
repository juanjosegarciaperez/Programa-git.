package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf/v2"
)

type Producto struct {
	nombre   string
	precio   float64
	cantidad int
}

func main() {

	var productos []Producto
	reader := bufio.NewReader(os.Stdin)
	var total float64 = 0.0

	for {
		fmt.Print("Ingrese el nombre del producto (o 'salir' para terminar): ")
		nombre, _ := reader.ReadString('\n')
		nombre = strings.TrimSpace(nombre)

		if strings.ToLower(nombre) == "salir" {
			break
		}

		fmt.Print("Ingrese el precio del produnto: ")
		precioInput, _ := reader.ReadString('\n')
		precioInput = strings.TrimSpace(precioInput)
		precio, err := strconv.ParseFloat(precioInput, 64)
		if err != nil {
			fmt.Println("Error: el precio no es valido")
			continue
		}

		fmt.Print("Ingrese la cantidad del producto: ")
		cantidadInput, _ := reader.ReadString('\n')
		cantidadInput = strings.TrimSpace(cantidadInput)
		cantidad, err := strconv.Atoi(cantidadInput)
		if err != nil {
			fmt.Println("Error: la cantidad no es valida")
			continue
		}

		producto := Producto{
			nombre:   nombre,
			precio:   precio,
			cantidad: cantidad,
		}
		productos = append(productos, producto)

		subtotal := producto.precio * float64(producto.cantidad)
		total += subtotal

		fmt.Println("Producto agregado con exito")
	}

	fmt.Println("\nLista de productos ingresados: ")
	for _, producto := range productos {
		subtotal := producto.precio * float64(producto.cantidad)
		fmt.Printf("Producto: %s | Precio %.2f | Cantidad %d | Subtotal %.2f\n", producto.nombre, producto.precio, producto.cantidad, subtotal)
	}
	fmt.Printf("\nTotal general: %.2f\n", total)

	fmt.Print("\n¿Desea continuar con la operacion del PDF? (s/n): ")
	confirmacion, _ := reader.ReadString('\n')
	confirmacion = strings.TrimSpace(confirmacion)

	if strings.ToLower(confirmacion) != "s" {
		fmt.Println("Generacion de PDF cancelada")
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Factura")

	fechaActual := time.Now().Format("02/01/2006 15:04:05")

	fmt.Print("Ingrese el nombre del cliente: ")
	scanner := bufio.NewReader(os.Stdin)
	nombreCliente, _ := scanner.ReadString('\n')

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Cliente: %s", nombreCliente))
	pdf.Ln(8)

	pdf.Cell(40, 10, fmt.Sprintf("Fecha: %s", fechaActual))
	pdf.Ln(8)

	drawTable(pdf, productos, total)
	// Guardar el PDF
	err := pdf.OutputFileAndClose("factura.pdf")
	if err != nil {
		fmt.Println("Error al crear el PDF:", err)
		return
	}

	fmt.Println("Factura creada con éxito: factura.pdf")
}

func drawTable(pdf *gofpdf.Fpdf, productos []Producto, total float64) {
	// Establecer el ancho de las columnas
	colWidth := []float64{80, 40, 40, 40}

	// Encabezados de la tabla
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(colWidth[0], 10, "Producto")
	pdf.Cell(colWidth[1], 10, "Precio")
	pdf.Cell(colWidth[2], 10, "Cantidad")
	pdf.Cell(colWidth[3], 10, "Subtotal")
	pdf.Ln(10)

	// Dibuja la línea horizontal debajo de los encabezados
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())

	// Establecer la fuente para los productos
	pdf.SetFont("Arial", "", 12)

	// Agregar productos a la tabla
	for _, producto := range productos {
		subtotal := producto.precio * float64(producto.cantidad)
		pdf.Cell(colWidth[0], 10, producto.nombre)
		pdf.Cell(colWidth[1], 10, fmt.Sprintf("%.2f", producto.precio))
		pdf.Cell(colWidth[2], 10, fmt.Sprintf("%d", producto.cantidad))
		pdf.Cell(colWidth[3], 10, fmt.Sprintf("%.2f", subtotal))
		pdf.Ln(10)

		// Dibuja la línea horizontal al final de la fila
		pdf.Line(10, pdf.GetY(), 200, pdf.GetY())
	}

	// Dibuja la línea horizontal al final de la tabla
	pdf.Line(10, pdf.GetY(), 200, pdf.GetY())

	// Agregar el total al PDF
	pdf.Ln(5) // Salto de línea
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(colWidth[0], 10, "Total:")
	pdf.Cell(colWidth[1], 10, fmt.Sprintf("%.2f", total))

}
