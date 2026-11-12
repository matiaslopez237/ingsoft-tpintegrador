package domain

// NumerosDeSprint devuelve el número de cada sprint dentro de su proyecto, por su ID. El ID es único
// en toda la base (RN-08.9) y no dice nada al usuario; las pantallas muestran el número (RN-08.7).
// Corrige DEF-05.
func NumerosDeSprint(sprints []Sprint) map[int64]int {
	return nil
}

// ContarSprints cuenta los sprints Abiertos y los Cerrados. Corrige DEF-05 (spec HU-26, sección 3:
// el Resumen del reporte de un proyecto los muestra).
func ContarSprints(sprints []Sprint) (abiertos, cerrados int) {
	return 0, 0
}
