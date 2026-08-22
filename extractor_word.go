package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// ExtraerTextoWord extrae tanto el texto de párrafos directos como el contenido de tablas.
func ExtraerTextoWord(ruta string) (string, error) {
	r, err := zip.OpenReader(ruta)
	if err != nil {
		return "", fmt.Errorf("error al abrir el archivo docx: %w", err)
	}
	defer r.Close()

	var documentXML *zip.File
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			documentXML = f
			break
		}
	}

	if documentXML == nil {
		return "", fmt.Errorf("no se encontro la estructura de texto dentro del docx")
	}

	docFile, err := documentXML.Open()
	if err != nil {
		return "", fmt.Errorf("error al leer el contenido XML de word: %w", err)
	}
	defer docFile.Close()

	byteValue, err := io.ReadAll(docFile)
	if err != nil {
		return "", fmt.Errorf("error al procesar los datos de word: %w", err)
	}

	var doc Document
	if err := xml.Unmarshal(byteValue, &doc); err != nil {
		return "", fmt.Errorf("error al decodificar el XML de word: %w", err)
	}

	var sb strings.Builder

	// 1. Extraer párrafos principales
	for _, p := range doc.Body.Paragraphs {
		var pText strings.Builder
		for _, t := range p.Text {
			pText.WriteString(t.Value)
		}
		linea := strings.TrimSpace(pText.String())
		if linea != "" {
			sb.WriteString(linea)
			sb.WriteString("\n")
		}
	}

	// 2. Extraer contenido dentro de tablas
	for _, tbl := range doc.Body.Tables {
		for _, row := range tbl.Rows {
			var rowText []string
			for _, cell := range row.Cells {
				var cellText strings.Builder
				for _, p := range cell.Paragraphs {
					for _, t := range p.Text {
						cellText.WriteString(t.Value)
					}
				}
				if txt := strings.TrimSpace(cellText.String()); txt != "" {
					rowText = append(rowText, txt)
				}
			}
			if len(rowText) > 0 {
				sb.WriteString(strings.Join(rowText, " | "))
				sb.WriteString("\n")
			}
		}
	}

	if sb.Len() == 0 {
		return "", fmt.Errorf("el archivo docx no contiene texto valido")
	}

	return sb.String(), nil
}
