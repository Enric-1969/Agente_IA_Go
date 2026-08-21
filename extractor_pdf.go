// (lectura y filtrado de .pdf)
package main

import (
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
)

// ExtraerTextoPDF procesa un archivo PDF página por página
// y recupera el texto plano delimitando el origen de cada página.
func ExtraerTextoPDF(ruta string) (string, error) {
	f, r, err := pdf.Open(ruta)
	if err != nil {
		return "", fmt.Errorf("error al abrir el archivo PDF: %w", err)
	}
	defer f.Close()

	var sb strings.Builder
	totalPaginas := r.NumPage()

	for i := 1; i <= totalPaginas; i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}

		content, err := p.GetPlainText(nil)
		if err != nil {
			continue
		}

		textoPagina := strings.TrimSpace(content)
		if textoPagina != "" {
			sb.WriteString(fmt.Sprintf("\n--- PAGINA %d ---\n", i))
			sb.WriteString(textoPagina)
			sb.WriteString("\n")
		}
	}

	if sb.Len() == 0 {
		return "", fmt.Errorf("el archivo PDF no contiene texto extraible")
	}

	return sb.String(), nil
}
