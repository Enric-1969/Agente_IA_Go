package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// EstimarTokensLocal calcula los tokens de entrada aproximados en microsegundos (~3.5 caracteres por token)
func EstimarTokensLocal(texto string) int {
	if len(texto) == 0 {
		return 0
	}
	return len(texto) * 10 / 35
}

// obtenerEnvInt lee un entero del .env, aplica valor por defecto si falta y valida un piso mínimo de seguridad.
func obtenerEnvInt(clave string, valorPorDefecto int, minimoPermitido int) int {
	if val := os.Getenv(clave); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			if parsed < minimoPermitido {
				fmt.Printf("⚠️ [%s=%d] es inferior al mínimo seguro (%d). Elevando automáticamente a %d.\n", clave, parsed, minimoPermitido, minimoPermitido)
				return minimoPermitido
			}
			return parsed
		}
	}
	return valorPorDefecto
}

// obtenerEnvFloat lee un decimal del .env y aplica un valor por defecto si falta o es inválido.
func obtenerEnvFloat(clave string, valorPorDefecto float64) float64 {
	if val := os.Getenv(clave); val != "" {
		if parsed, err := strconv.ParseFloat(val, 64); err == nil && parsed > 0 {
			return parsed
		}
	}
	return valorPorDefecto
}

// cargarEnv lee el archivo .env e inyecta las variables en la memoria del proceso
func cargarEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		linea := strings.TrimSpace(scanner.Text())
		if linea == "" || strings.HasPrefix(linea, "#") {
			continue
		}
		partes := strings.SplitN(linea, "=", 2)
		if len(partes) == 2 {
			clave := strings.TrimSpace(partes[0])
			valor := strings.TrimSpace(partes[1])
			valor = strings.Trim(valor, `"'`)
			os.Setenv(clave, valor)
		}
	}
	if err := scanner.Err(); err != nil {
		return
	}
}
