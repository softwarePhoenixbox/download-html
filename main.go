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

// Extensiones válidas de imagen
var extensionesValidas = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".bmp":  true,
	".svg":  true,
	".ico":  true,
}

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

	fmt.Printf("🔍 Encontradas %d referencias\n", len(matches))

	descargadas := 0
	saltadas := 0

	// Descargar cada imagen
	for _, match := range matches {
		if len(match) > 1 {
			imgURL := match[1]

			// Si es URL relativa, hacerla absoluta
			if !strings.HasPrefix(imgURL, "http") {
				if strings.HasPrefix(imgURL, "//") {
					imgURL = "https:" + imgURL
				} else {
					baseURL := url[:strings.LastIndex(url, "/")+1]
					imgURL = baseURL + imgURL
				}
			}

			// Extraer nombre del archivo
			partes := strings.Split(imgURL, "/")
			nombreImg := partes[len(partes)-1]

			// Quitar parámetros de URL (después de ?)
			if strings.Contains(nombreImg, "?") {
				nombreImg = strings.Split(nombreImg, "?")[0]
			}

			// VERIFICAR EXTENSIÓN VÁLIDA
			ext := strings.ToLower(filepath.Ext(nombreImg))
			if !extensionesValidas[ext] {
				saltadas++
				continue // Saltar esta imagen
			}

			// Si no tiene nombre válido, generar uno
			if nombreImg == "" || nombreImg == "." {
				nombreImg = fmt.Sprintf("imagen_%d%s", descargadas, ext)
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
			descargadas++
			fmt.Printf("✅ Imagen %d: %s\n", descargadas, nombreImg)
		}
	}

	fmt.Printf("\n🎉 Total: %d descargadas, %d saltadas (sin extensión válida)\n", descargadas, saltadas)
	fmt.Println("📁 Todo guardado en carpeta 'descarga'")
}
