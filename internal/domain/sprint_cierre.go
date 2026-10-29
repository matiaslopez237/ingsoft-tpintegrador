package domain

// ResumenCierre es lo que deja el cierre de un sprint: los story points que registra y las historias
// que vuelven al Product Backlog. Implementa HU-11 (RN-11.3, RN-11.6).
type ResumenCierre struct {
	PuntosPlanificados int     // todas las historias asignadas al cerrar, terminadas o no
	PuntosCompletados  int     // solo las Hecho
	DevueltasAlBacklog []int64 // IDs de las historias sin terminar, en el orden recibido
}

// ResumirCierre aplica RN-11.3 y RN-11.6 a las historias asignadas al sprint en el momento del
// cierre. Una historia sin estimar suma 0 (P-12); las Hecho quedan en el sprint (P-13). No modifica
// las historias: quitarlas del sprint es tarea de la persistencia.
func ResumirCierre(historias []ElementoBacklog) ResumenCierre {
	var r ResumenCierre
	for _, h := range historias {
		sp := 0
		if h.StoryPoints != nil {
			sp = *h.StoryPoints
		}
		r.PuntosPlanificados += sp
		if h.Estado == EstadoHecho {
			r.PuntosCompletados += sp
		} else {
			r.DevueltasAlBacklog = append(r.DevueltasAlBacklog, h.ID)
		}
	}
	return r
}
