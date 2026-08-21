// (lectura de .docx)
package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type Document struct {
	XMLName xml.Name `xml:"document"`
	Body    Body     `xml:"body"`
}

type Body struct {
	Paragraphs []Paragraph `xml:"p"`
}

type Paragraph struct {
	Runs []Run `xml:"r"`
}

type Run struct {
	Text string `xml:"t"`
}

// ExtraerTextoWord procesa archivos .docx extrayendo el contenido de word/document.xml
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
	for _, p := range doc.Body.Paragraphs {
		var pText strings.Builder
		for _, r := range p.Runs {
			pText.WriteString(r.Text)
		}
		linea := strings.TrimSpace(pText.String())
		if linea != "" {
			sb.WriteString(linea)
			sb.WriteString("\n")
		}
	}

	if sb.Len() == 0 {
		return "", fmt.Errorf("el archivo docx no contiene texto valido")
	}

	return sb.String(), nil
}
