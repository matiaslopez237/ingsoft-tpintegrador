// Command server levanta la aplicación web del TP Integrador
// "Software Metrics & Estimation".
package main

import (
	"log"
	"net/http"
	"os"

	adapterhttp "github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/http"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/persistencia"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/app"
)

func main() {
	rutaDB := os.Getenv("DB_PATH")
	if rutaDB == "" {
		rutaDB = "metrics.db"
	}

	db, err := persistencia.AbrirDB(rutaDB)
	if err != nil {
		log.Fatalf("no se pudo inicializar la base de datos: %v", err)
	}
	defer db.Close()
	log.Printf("conectado a la base de datos (%s)", rutaDB)

	repoProyectos := persistencia.NuevoRepositorioProyectos(db)
	mux := adapterhttp.NuevoRouter(app.NuevoCrearProyecto(repoProyectos))

	puerto := os.Getenv("PORT")
	if puerto == "" {
		puerto = "8080"
	}

	log.Printf("servidor escuchando en :%s", puerto)
	if err := http.ListenAndServe(":"+puerto, mux); err != nil {
		log.Fatalf("error del servidor: %v", err)
	}
}
