// Command server levanta la aplicación web del TP Integrador
// "Software Metrics & Estimation".
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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
	repoBacklog := persistencia.NuevoRepositorioBacklog(db)
	mux := adapterhttp.NuevoRouter(adapterhttp.CasosDeUso{
		CrearProyecto:        app.NuevoCrearProyecto(repoProyectos),
		ModificarProyecto:    app.NuevoModificarProyecto(repoProyectos),
		AgregarIntegrante:    app.NuevoAgregarIntegrante(repoProyectos),
		QuitarIntegrante:     app.NuevoQuitarIntegrante(repoProyectos),
		AgregarRol:           app.NuevoAgregarRol(repoProyectos),
		ConsultarProyecto:    app.NuevoConsultarProyecto(repoProyectos),
		CrearElementoBacklog: app.NuevoCrearElementoBacklog(repoBacklog),
		// time.Now como reloj: "hoy" es el día actual en la zona del servidor (HU-18, RN-18.4).
		RegistrarEsfuerzo: app.NuevoRegistrarEsfuerzo(repoProyectos, repoBacklog,
			persistencia.NuevoRepositorioEsfuerzo(db), time.Now),
	})

	puerto := os.Getenv("PORT")
	if puerto == "" {
		puerto = "8080"
	}

	servidor := &http.Server{
		Addr:              ":" + puerto,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Apagado ordenado con Ctrl+C: termina los pedidos en curso y deja correr el defer db.Close().
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		<-ctx.Done()
		apagar, cancelar := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelar()
		if err := servidor.Shutdown(apagar); err != nil {
			log.Printf("apagar servidor: %v", err)
		}
	}()

	log.Printf("servidor escuchando en :%s", puerto)
	if err := servidor.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("error del servidor: %v", err)
	}
}
