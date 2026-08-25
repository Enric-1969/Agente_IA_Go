// EstimarTokensLocal calcula los tokens de entrada aproximados en microsegundos (~3.5 caracteres por token).
func EstimarTokensLocal(texto string) int {
	if len(texto) == 0 {
		return 0
	}
	return len(texto) * 10 / 35
}