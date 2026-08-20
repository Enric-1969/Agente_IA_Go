package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ObtenerConfigRutas lee las variables del .env o asigna los valores por defecto.
func ObtenerConfigRutas() ConfigRutas {
	rutaDocs := os.Getenv("RUTA_DOCS_ORIGEN")
	if rutaDocs == "" {
		rutaDocs = "./documentos" // Según explicación en .env.example
	}

	rutaVector := os.Getenv("RUTA_VECTORSTORES")
	if rutaVector == "" {
		rutaVector = "./mis_vectorstores"
	}

	return ConfigRutas{
		RutaDocsOrigen:   filepath.Clean(rutaDocs),
		RutaVectorstores: filepath.Clean(rutaVector),
	}
}

// ValidarRutaVerificacion comprueba si una carpeta existe y es accesible.
func ValidarRutaVerificacion(ruta string) bool {
	info, err := os.Stat(ruta)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// ListarArchivosDisponibles escanea de forma recursiva y filtra los formatos procesables por la IA.
func ListarArchivosDisponibles(ruta string) ([]string, error) {
	var archivos []string

	// Extensiones permitidas según las carpetas de conocimiento de Files_IA
	extPermitidas := map[string]bool{
		// Documentación / Texto general
		".md": true, ".pdf": true, ".txt": true,
		// Go & Python
		".go": true, ".py": true, ".ipynb": true,
		// Excel, Word, Macros (VBA) y módulos exportados
		".xlsx": true, ".xls": true, ".xlsm": true, ".doc": true, ".docx": true, ".bas": true, ".cls": true,
		// WordPress, Web y bases de datos
		".php": true, ".js": true, ".css": true, ".html": true, ".sql": true, ".json": true,
	}

	err := filepath.WalkDir(ruta, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if extPermitidas[ext] {
				relPath, errRel := filepath.Rel(ruta, path)
				if errRel == nil {
					archivos = append(archivos, relPath)
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error al recorrer la ruta %s: %w", ruta, err)
	}

	return archivos, nil
}
