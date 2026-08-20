package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// --- UTILITIES: MANIFEST Y HASH ---

func GenerarHashMD5(texto string) string {
	hash := md5.Sum([]byte(texto))
	return hex.EncodeToString(hash[:])
}

func ObtenerHashArchivo(ruta string) (string, error) {
	info, err := os.Stat(ruta)
	if err != nil {
		return "", err
	}
	clave := fmt.Sprintf("%s_%d_%d", info.Name(), info.ModTime().UnixNano(), info.Size())
	return GenerarHashMD5(clave), nil
}

func CargarManifest(ruta string) map[string]string {
	manifest := make(map[string]string)
	data, err := os.ReadFile(ruta)
	if err == nil {
		_ = json.Unmarshal(data, &manifest)
	}
	return manifest
}

func GuardarManifest(ruta string, manifest map[string]string) error {
	_ = os.MkdirAll(filepath.Dir(ruta), 0755)
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ruta, data, 0644)
}

// --- FILTROS Y VALIDACIÓN DE ARCHIVOS ---

func esExtensionSoportada(ext string) bool {
	soportadas := map[string]bool{
		".pdf": true, ".txt": true, ".md": true, ".py": true,
		".go": true, ".html": true, ".js": true, ".php": true, ".vba": true,
	}
	return soportadas[ext]
}

func archivoHaCambiado(relPath, hashActual string, manifest map[string]string) bool {
	return manifest[relPath] != hashActual
}
