package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// NuevoGestorModelos genera el gestor e inicializa la lista inicial de modelos probados.
func NuevoGestorModelos(apiKey string) *GestorModelos {
	gm := &GestorModelos{apiKey: apiKey}
	gm.ModelosValidos = ObtenerModelosValidos(apiKey)
	return gm
}

// ObtenerModeloActivo devuelve el modelo prioritario (el primero de la lista).
func (gm *GestorModelos) ObtenerModeloActivo() string {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	if len(gm.ModelosValidos) == 0 {
		return ""
	}
	return gm.ModelosValidos[0]
}

// ReportarFalloYRotar elimina el modelo que ha fallado, conmuta al siguiente
// y dispara la búsqueda en segundo plano si queda solo 1 modelo en la reserva.
func (gm *GestorModelos) ReportarFalloYRotar(modeloFallido string) string {
	gm.mu.Lock()

	var listaFiltrada []string
	for _, m := range gm.ModelosValidos {
		if m != modeloFallido {
			listaFiltrada = append(listaFiltrada, m)
		}
	}
	gm.ModelosValidos = listaFiltrada

	restantes := len(gm.ModelosValidos)
	debeBuscarNuevos := restantes <= 1 && !gm.recargando

	if debeBuscarNuevos {
		gm.recargando = true
	}

	siguienteModelo := ""
	if restantes > 0 {
		siguienteModelo = gm.ModelosValidos[0]
		fmt.Printf("🔄 Modelo '%s' descartado. Conmutando al siguiente: '%s'\n", modeloFallido, siguienteModelo)
	}

	gm.mu.Unlock()

	if debeBuscarNuevos {
		fmt.Println("⚠️ Queda solo 1 modelo en la reserva. Generando nueva lista en segundo plano...")
		go func() {
			nuevosModelos := ObtenerModelosValidos(gm.apiKey)

			gm.mu.Lock()
			defer gm.mu.Unlock()

			if len(nuevosModelos) > 0 {
				gm.ModelosValidos = nuevosModelos
				fmt.Printf("✅ Lista de modelos actualizada en segundo plano (%d modelos activos)\n", len(nuevosModelos))
			}
			gm.recargando = false
		}()
	}

	return siguienteModelo
}

// ObtenerModelosValidos realiza el escaneo directo contra Groq en paralelo.
func ObtenerModelosValidos(apiKey string) []string {
	req, err := http.NewRequest("GET", "https://api.groq.com/openai/v1/models", nil)
	if err != nil {
		return nil
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil
	}
	defer resp.Body.Close()

	var catalogo CatalogoModelos
	if err := json.NewDecoder(resp.Body).Decode(&catalogo); err != nil {
		return nil
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var modelosValidos []string
	descartes := []string{"whisper", "guard", "vision"}

	for _, m := range catalogo.Data {
		modelIDLower := strings.ToLower(m.ID)

		ignorar := false
		for _, exc := range descartes {
			if strings.Contains(modelIDLower, exc) {
				ignorar = true
				break
			}
		}
		if ignorar {
			continue
		}

		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			if ProbarSoporteToolCalling(apiKey, id) {
				mu.Lock()
				modelosValidos = append(modelosValidos, id)
				mu.Unlock()
			}
		}(m.ID)
	}

	wg.Wait()
	return modelosValidos
}
