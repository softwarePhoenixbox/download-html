package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	// Verificar que se pasó la URL
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run . <URL>")
		os.Exit(1)
	}

	url := os.Args[1]

	// Hacer la petición GET
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error al conectar:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Crear archivo de salida
	out, err := os.Create("descarga.html")
	if err != nil {
		fmt.Println("Error al crear archivo:", err)
		os.Exit(1)
	}
	defer out.Close()

	// Copiar el contenido al archivo
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Error al guardar:", err)
		os.Exit(1)
	}

	fmt.Println("✅ Descargado en: descarga.html")
}