package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// GenerarHashMD5 crea un hash rápido a partir de una cadena de texto.
func GenerarHashMD5(texto string) string {
	hash := md5.Sum([]byte(texto))
	return hex.EncodeToString(hash[:])
}

// ObtenerHashArchivo genera una firma usando nombre, fecha de modificación y tamaño.
// Es rapidísimo (O(1)) porque evita leer el contenido del archivo.
func ObtenerHashArchivo(ruta string) (string, error) {
	info, err := os.Stat(ruta)
	if err != nil {
		return "", err
	}
	clave := fmt.Sprintf("%s_%d_%d", info.Name(), info.ModTime().UnixNano(), info.Size())
	return GenerarHashMD5(clave), nil
}

// CargarManifest lee el JSON y devuelve el mapa de archivos procesados.
// Devuelve mapa vacío si el archivo no existe, pero lanza error si el JSON está corrupto.
func CargarManifest(ruta string) (map[string]string, error) {
	manifest := make(map[string]string)

	data, err := os.ReadFile(ruta)
	if err != nil {
		// Si el error es simplemente que el archivo no existe, devolvemos el mapa vacío (es normal la primera vez)
		if os.IsNotExist(err) {
			return manifest, nil
		}
		// Si es un error de permisos o disco, lo devolvemos
		return nil, err
	}

	// NO usamos el guion bajo (_). Si el JSON está mal formado, queremos saberlo.
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		return nil, fmt.Errorf("error al leer el JSON de manifest (posible corrupción): %w", err)
	}

	return manifest, nil
}

// GuardarManifest sobrescribe el archivo JSON con el mapa actualizado en formato legible (indentado).
func GuardarManifest(ruta string, manifest map[string]string) error {
	err := os.MkdirAll(filepath.Dir(ruta), 0755)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ruta, data, 0644)
}

// --- FILTROS Y VALIDACIÓN DE ARCHIVOS ---

// esExtensionSoportada verifica si el archivo puede ser procesado por los extractores.
func esExtensionSoportada(ext string) bool {
	soportadas := map[string]bool{
		".pdf": true, ".txt": true, ".md": true, ".py": true,
		".go": true, ".html": true, ".js": true, ".php": true, ".vba": true,
	}
	return soportadas[ext]
}

// archivoHaCambiado compara la firma del disco con la firma guardada en el manifest.
func archivoHaCambiado(relPath, hashActual string, manifest map[string]string) bool {
	return manifest[relPath] != hashActual
}
