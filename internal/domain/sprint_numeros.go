package domain

// NumerosDeSprint devuelve el número de cada sprint dentro de su proyecto, por su ID. El ID es único
// en toda la base (RN-08.9) y no dice nada al usuario; las pantallas muestran el número (RN-08.7).
// Corrige DEF-05.
func NumerosDeSprint(sprints []Sprint) map[int64]int {
	numeros := make(map[int64]int, len(sprints))
	for _, s := range sprints {
		numeros[s.ID] = s.Numero
	}
	return numeros
}

// ContarSprints cuenta los sprints Abiertos y los Cerrados. Corrige DEF-05 (spec HU-26, sección 3:
// el Resumen del reporte de un proyecto los muestra).
func ContarSprints(sprints []Sprint) (abiertos, cerrados int) {
	for _, s := range sprints {
		switch s.Estado {
		case SprintAbierto:
			abiertos++
		case SprintCerrado:
			cerrados++
		}
	}
	return abiertos, cerrados
}
