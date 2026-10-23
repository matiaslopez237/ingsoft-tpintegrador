package domain

import (
	"errors"
	"strings"
)

// Errores de HU-03 (spec, sección 7).
var (
	ErrNombreIntegranteVacio  = errors.New("el nombre del integrante es obligatorio")
	ErrRolIntegranteVacio     = errors.New("el rol del integrante es obligatorio")
	ErrIntegranteDuplicado    = errors.New("ya hay un integrante con ese nombre en el proyecto")
	ErrIntegranteNoEncontrado = errors.New("integrante no encontrado en el proyecto")
	ErrRolDuplicado           = errors.New("el integrante ya tiene ese rol")
)

// Integrante es una persona que participa de un proyecto, con sus roles en el orden en que se
// agregaron. Implementa HU-03 (RN-03.5).
type Integrante struct {
	Nombre string
	Roles  []string
}

// AgregarIntegrante valida los datos y agrega un integrante al proyecto.
// La comparación de nombres duplicados usa ClaveNombre, igual que HU-01 (RN-01.3) hace para el nombre del proyecto.
// Implementa HU-03 (RN-03.1, RN-03.2, RN-03.3, RN-03.4).
func (p *Proyecto) AgregarIntegrante(nombre, rol string) error {
	nombre = strings.TrimSpace(nombre)
	rol = strings.TrimSpace(rol)
	switch {
	case nombre == "":
		return ErrNombreIntegranteVacio
	case rol == "":
		return ErrRolIntegranteVacio
	}

	clave := ClaveNombre(nombre)
	for _, existente := range p.Integrantes {
		if ClaveNombre(existente.Nombre) == clave {
			return ErrIntegranteDuplicado
		}
	}

	p.Integrantes = append(p.Integrantes, Integrante{Nombre: nombre, Roles: []string{rol}})
	return nil
}

// AgregarRol le suma un rol a un integrante que ya está en el proyecto.
// Los roles se comparan con ClaveNombre, igual que los nombres (RN-03.3).
// Implementa HU-03 (RN-03.4, RN-03.5, RN-03.6).
func (p *Proyecto) AgregarRol(nombre, rol string) error {
	rol = strings.TrimSpace(rol)
	if rol == "" {
		return ErrRolIntegranteVacio
	}

	clave := ClaveNombre(nombre)
	for i, existente := range p.Integrantes {
		if ClaveNombre(existente.Nombre) != clave {
			continue
		}
		for _, r := range existente.Roles {
			if ClaveNombre(r) == ClaveNombre(rol) {
				return ErrRolDuplicado
			}
		}
		p.Integrantes[i].Roles = append(existente.Roles, rol)
		return nil
	}
	return ErrIntegranteNoEncontrado
}

// QuitarIntegrante quita del proyecto al integrante con ese nombre.
// Implementa HU-03 (CA-03.3, CA-03.4).
func (p *Proyecto) QuitarIntegrante(nombre string) error {
	clave := ClaveNombre(nombre)
	for i, existente := range p.Integrantes {
		if ClaveNombre(existente.Nombre) == clave {
			p.Integrantes = append(p.Integrantes[:i], p.Integrantes[i+1:]...)
			return nil
		}
	}
	return ErrIntegranteNoEncontrado
}
