package main

import (
	"math"
	"testing"
)

func TestSimilitudCoseno(t *testing.T) {
	// Tolerancia mínima para comparar números decimales (float32)
	const epsilon = 1e-5

	casiIgual := func(a, b float32) bool {
		return math.Abs(float64(a-b)) < epsilon
	}

	// 1. Vectores idénticos -> Similitud máxima (1.0)
	v1 := []float32{1.0, 2.0, 3.0}
	if res := SimilitudCoseno(v1, v1); !casiIgual(res, 1.0) {
		t.Errorf("Vectores idénticos deberían dar 1.0, se obtuvo %f", res)
	}

	// 2. Vectores perpendiculares (sin relación) -> Similitud 0.0
	vA := []float32{1.0, 0.0}
	vB := []float32{0.0, 1.0}
	if res := SimilitudCoseno(vA, vB); !casiIgual(res, 0.0) {
		t.Errorf("Vectores perpendiculares deberían dar 0.0, se obtuvo %f", res)
	}

	// 3. Vectores exactamente opuestos -> Similitud -1.0
	vPos := []float32{1.0, 1.0}
	vNeg := []float32{-1.0, -1.0}
	if res := SimilitudCoseno(vPos, vNeg); !casiIgual(res, -1.0) {
		t.Errorf("Vectores opuestos deberían dar -1.0, se obtuvo %f", res)
	}

	// 4. Protección contra división por cero (vector nulo) o dimensiones distintas
	vCero := []float32{0.0, 0.0}
	if res := SimilitudCoseno(vPos, vCero); res != 0.0 {
		t.Errorf("Un vector nulo debe devolver 0.0 sin NaN, se obtuvo %f", res)
	}
	if res := SimilitudCoseno([]float32{1.0}, []float32{1.0, 2.0}); res != 0.0 {
		t.Errorf("Dimensiones distintas deben devolver 0.0, se obtuvo %f", res)
	}
}
