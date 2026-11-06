package ejemplo

import "time"

// contenidoTP es el proyecto de ejemplo: este mismo TP, "Software Metrics & Estimation" (sección 3
// de la spec). El backlog, los sprints, el esfuerzo, las rondas y los defectos se suman en la
// tarea 31.5.
func contenidoTP() contenido {
	return contenido{
		proyecto: "Software Metrics & Estimation",
		inicio:   time.Date(2026, time.October, 22, 0, 0, 0, 0, time.Local), // P-1
		integrantes: []integrante{
			{"Matías López", "Agile Enabler"},
			{"Ezequiel Assenza", "Product Builder"},
			{"Nicolás Andino", "Product Builder"},
			{"Lucas Orellana", "Product Builder"},
		},
	}
}
