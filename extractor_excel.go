// (lectura de .xlsx)
package main

import (
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ExtraerTextoExcel procesa libros .xlsx y .xlsm iterando por todas sus hojas.
// Formatea las filas con separadores (|) para preservar la estructura tabular.
func ExtraerTextoExcel(ruta string) (string, error) {
	f, err := excelize.OpenFile(ruta)
	if err != nil {
		return "", fmt.Errorf("error al abrir el archivo excel: %w", err)
	}
	defer f.Close()

	var sb strings.Builder

	hojas := f.GetSheetList()
	for _, hoja := range hojas {
		rows, err := f.GetRows(hoja)
		if err != nil || len(rows) == 0 {
			continue
		}

		sb.WriteString(fmt.Sprintf("\n--- HOJA: %s ---\n", hoja))

		for _, row := range rows {
			linea := strings.Join(row, " | ")
			if strings.TrimSpace(linea) != "" {
				sb.WriteString(linea)
				sb.WriteString("\n")
			}
		}
	}

	if sb.Len() == 0 {
		return "", fmt.Errorf("el archivo excel no contiene datos de texto validos")
	}

	return sb.String(), nil
}
