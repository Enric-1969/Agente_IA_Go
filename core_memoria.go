package main

// NuevoGestorMemoria inicializa la memoria con un límite de mensajes retenidos.
func NuevoGestorMemoria(limite int) *GestorMemoria {
	if limite <= 0 {
		limite = 6 // Límite predeterminado: 6 mensajes (3 turnos completos)
	}
	return &GestorMemoria{
		mensajes: make([]Mensaje, 0, limite),
		limite:   limite,
	}
}

// Agregar registra un nuevo mensaje y aplica la ventana deslizante si supera el límite.
func (m *GestorMemoria) Agregar(rol, contenido string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.mensajes = append(m.mensajes, Mensaje{Role: rol, Content: contenido})

	if len(m.mensajes) > m.limite {
		m.mensajes = m.mensajes[len(m.mensajes)-m.limite:]
	}
}

// ObtenerHistorial retorna una copia segura del historial para evitar condiciones de carrera.
func (m *GestorMemoria) ObtenerHistorial() []Mensaje {
	m.mu.RLock()
	defer m.mu.RUnlock()

	copia := make([]Mensaje, len(m.mensajes))
	copy(copia, m.mensajes)
	return copia
}

// Limpiar vacía el historial de la sesión activa.
func (m *GestorMemoria) Limpiar() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.mensajes = make([]Mensaje, 0, m.limite)
}
