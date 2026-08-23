package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

// DespacharHerramienta evalúa la intención asignada por el router y ejecuta la función nativa en Go.
func DespacharHerramienta(comando string, rutaDocs string) ResultadoHerramienta {
	switch comando {
	case "listar_documentos":
		return HerramientaListarDocumentos(rutaDocs)
	case "info_sistema":
		return HerramientaInfoSistema()
	default:
		return ResultadoHerramienta{
			Nombre:    comando,
			Exito:     false,
			Contenido: "",
			Error:     fmt.Sprintf("Herramienta '%s' no reconocida por el despachador.", comando),
		}
	}
}

// HerramientaListarDocumentos inspecciona la carpeta de entrada local.
func HerramientaListarDocumentos(rutaDir string) ResultadoHerramienta {
	archivos, err := os.ReadDir(rutaDir)
	if err != nil {
		return ResultadoHerramienta{
			Nombre: "listar_documentos",
			Exito:  false,
			Error:  fmt.Sprintf("No se pudo leer la carpeta '%s': %v", rutaDir, err),
		}
	}

	if len(archivos) == 0 {
		return ResultadoHerramienta{
			Nombre:    "listar_documentos",
			Exito:     true,
			Contenido: fmt.Sprintf("La carpeta '%s' está vacía.", rutaDir),
		}
	}

	res := fmt.Sprintf("Archivos encontrados en '%s':\n", rutaDir)
	for _, file := range archivos {
		if !file.IsDir() {
			info, err := file.Info()
			if err == nil {
				res += fmt.Sprintf("- %s (%d bytes, modificado: %s)\n",
					file.Name(), info.Size(), info.ModTime().Format("2006-01-02 15:04:02"))
			}
		}
	}

	return ResultadoHerramienta{
		Nombre:    "listar_documentos",
		Exito:     true,
		Contenido: res,
	}
}

// HerramientaInfoSistema devuelve métricas básicas de ejecución local del ejecutable Go.
func HerramientaInfoSistema() ResultadoHerramienta {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	info := fmt.Sprintf(
		"Estado del Agente Go:\n"+
			"- Sistema Operativo: %s/%s\n"+
			"- CPUS Cores: %d\n"+
			"- Memoria Asignada (Heap): %.2f MB\n"+
			"- Hora del sistema: %s\n",
		runtime.GOOS, runtime.GOARCH,
		runtime.NumCPU(),
		float64(m.Alloc)/(1024*1024),
		time.Now().Format("2006-01-02 15:04:05"),
	)

	return ResultadoHerramienta{
		Nombre:    "info_sistema",
		Exito:     true,
		Contenido: info,
	}
}
