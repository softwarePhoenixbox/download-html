package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	fmt.Print("URL: ")
	reader := bufio.NewReader(os.Stdin)
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url)

	if url == "" {
		fmt.Println("No escribiste ninguna URL")
		return
	}

	// Crear carpeta para la descarga
	nombreCarpeta := "descarga"
	os.MkdirAll(nombreCarpeta, 0755)

	// Descargar HTML
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	htmlContent, _ := io.ReadAll(resp.Body)

	// Guardar HTML
	os.WriteFile(filepath.Join(nombreCarpeta, "pagina.html"), htmlContent, 0644)
	fmt.Println("✅ HTML descargado")

	// Buscar imágenes en el HTML
	imgRegex := regexp.MustCompile(`<img[^>]+src="([^"]+)"`)
	matches := imgRegex.FindAllStringSubmatch(string(htmlContent), -1)

	fmt.Printf("🔍 Encontradas %d imágenes\n", len(matches))

	// Descargar cada imagen
	for i, match := range matches {
		if len(match) > 1 {
			imgURL := match[1]

			// Si es URL relativa, hacerla absoluta
			if !strings.HasPrefix(imgURL, "http") {
				baseURL := url[:strings.LastIndex(url, "/")+1]
				imgURL = baseURL + imgURL
			}

			// Extraer nombre del archivo
			partes := strings.Split(imgURL, "/")
			nombreImg := partes[len(partes)-1]
			if nombreImg == "" {
				nombreImg = fmt.Sprintf("imagen_%d.jpg", i)
			}

			// Descargar imagen
			imgResp, err := http.Get(imgURL)
			if err != nil {
				fmt.Printf("❌ Error en %s: %v\n", nombreImg, err)
				continue
			}

			imgData, _ := io.ReadAll(imgResp.Body)
			imgResp.Body.Close()

			ruta := filepath.Join(nombreCarpeta, nombreImg)
			os.WriteFile(ruta, imgData, 0644)
			fmt.Printf("✅ Imagen %d/%d: %s\n", i+1, len(matches), nombreImg)
		}
	}

	fmt.Println("\n🎉 Todo guardado en carpeta 'descarga'")
}
