package main

import "testing"

func TestEvaluarIntencionHeuristica(t *testing.T) {
	// Definimos una tabla de casos de prueba (Table-Driven Test)
	casos := []struct {
		nombre        string
		consulta      string
		rutaEsperada  TipoRuta
		resueltoEsper bool
	}{
		{
			nombre:        "Detección de refactorización / desarrollo",
			consulta:      "  Refactoriza esta FUNCION en Go por favor  ",
			rutaEsperada:  RutaDesarrollo,
			resueltoEsper: true,
		},
		{
			nombre:        "Detección de consulta RAG por prefijo",
			consulta:      "Como funciona el sistema vectorial?",
			rutaEsperada:  RutaRAG,
			resueltoEsper: true,
		},
		{
			nombre:        "Detección de consulta RAG por palabra clave",
			consulta:      "Revisa la documentacion del proyecto",
			rutaEsperada:  RutaRAG,
			resueltoEsper: true,
		},
		{
			nombre:        "Detección de herramienta local",
			consulta:      "Puedes listar archivos de la carpeta?",
			rutaEsperada:  RutaHerramienta,
			resueltoEsper: true,
		},
		{
			nombre:        "Consulta general sin coincidencia heurística (delega al LLM)",
			consulta:      "Hola, buenos dias, que tal estas hoy?",
			rutaEsperada:  RutaGeneral,
			resueltoEsper: false,
		},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			ruta, resuelto := EvaluarIntencionHeuristica(tc.consulta)

			if ruta != tc.rutaEsperada || resuelto != tc.resueltoEsper {
				t.Errorf("Para '%s': se esperaba (%v, %v), pero se obtuvo (%v, %v)",
					tc.consulta, tc.rutaEsperada, tc.resueltoEsper, ruta, resuelto)
			}
		})
	}
}
