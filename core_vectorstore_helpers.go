package main

import (
	"math"
)

// SimilitudCoseno calcula el ángulo de coincidencia entre dos vectores de igual dimensión.
func SimilitudCoseno(v1, v2 []float32) float32 {
	if len(v1) != len(v2) || len(v1) == 0 {
		return 0.0
	}

	var productoEscalar, normaA, normaB float64
	for i := range v1 {
		productoEscalar += float64(v1[i]) * float64(v2[i])
		normaA += float64(v1[i]) * float64(v1[i])
		normaB += float64(v2[i]) * float64(v2[i])
	}

	if normaA == 0 || normaB == 0 {
		return 0.0
	}

	return float32(productoEscalar / (math.Sqrt(normaA) * math.Sqrt(normaB)))
}
