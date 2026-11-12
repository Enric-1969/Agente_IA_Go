package main

import (
	"testing"
	"time"
)

func TestGestorModelos_ObtenerYRotar(t *testing.T) {
	// 1. Preparar (Arrange): Creamos un gestor en memoria con 3 modelos simulados
	gm := &GestorModelos{
		APIKey:         "clave-falsa-para-test",
		ModelosValidos: []string{"modelo-A", "modelo-B", "modelo-C"},
	}

	// Comprobamos que el modelo activo inicial sea el primero ("modelo-A")
	if activo := gm.ObtenerModeloActivo(); activo != "modelo-A" {
		t.Fatalf("Se esperaba 'modelo-A', se obtuvo '%s'", activo)
	}

	// 2. Actuar (Act): Simulamos que 'modelo-A' falla (quedan 2 modelos: B y C)
	siguiente := gm.ReportarFalloYRotar("modelo-A")

	// 3. Comprobar (Assert): Debe conmutar a 'modelo-B'
	if siguiente != "modelo-B" {
		t.Errorf("Se esperaba conmutar a 'modelo-B', se obtuvo '%s'", siguiente)
	}
	if activo := gm.ObtenerModeloActivo(); activo != "modelo-B" {
		t.Errorf("El modelo activo tras rotar debería ser 'modelo-B', es '%s'", activo)
	}

	// Simulamos un segundo fallo ('modelo-B'): queda solo 1 modelo ('modelo-C')
	// Esto dispara la goroutine de recarga en segundo plano (que fallará rápido por clave falsa sin romper nada)
	ultimo := gm.ReportarFalloYRotar("modelo-B")
	if ultimo != "modelo-C" {
		t.Errorf("Se esperaba conmutar a 'modelo-C', se obtuvo '%s'", ultimo)
	}

	// Breve pausa para permitir que la goroutine en segundo plano termine limpia antes de cerrar el test
	time.Sleep(150 * time.Millisecond)

	// Si falla el último modelo, debe devolver cadena vacía sin entrar en pánico (panic)
	vacio := gm.ReportarFalloYRotar("modelo-C")
	if vacio != "" {
		t.Errorf("Se esperaba cadena vacía al agotar modelos, se obtuvo '%s'", vacio)
	}
}
