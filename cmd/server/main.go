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
	"syscall"
	"time"

	adapterhttp "github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/http"
	"github.com/matiaslopez237/ingsoft-tpintegrador/internal/adapters/pdf"
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
	repoSprints := persistencia.NuevoRepositorioSprints(db)
	repoRondas := persistencia.NuevoRepositorioRondas(db)     // HU-14
	repoEsfuerzo := persistencia.NuevoRepositorioEsfuerzo(db) // HU-18, HU-19 y HU-23
	repoDefectos := persistencia.NuevoRepositorioDefectos(db) // HU-20/HU-21
	armarReporteProyecto := app.NuevoArmarReporteProyecto(repoProyectos, repoBacklog, repoSprints, repoEsfuerzo, repoDefectos)
	armarReporteSprint := app.NuevoArmarReporteSprint(repoProyectos, repoSprints, repoBacklog, repoEsfuerzo, repoDefectos) // HU-26
	// HU-27: la fuente convierte el reporte de HU-26 a DocumentoPDF (spec HU-27, P-11).
	fuenteReporte := adapterhttp.NuevaFuenteDocumentoReporte(armarReporteProyecto, armarReporteSprint)
	generadorPDF := pdf.NuevoGenerador()
	mux := adapterhttp.NuevoRouter(adapterhttp.CasosDeUso{
		CrearProyecto:            app.NuevoCrearProyecto(repoProyectos),
		ModificarProyecto:        app.NuevoModificarProyecto(repoProyectos),
		AgregarIntegrante:        app.NuevoAgregarIntegrante(repoProyectos),
		QuitarIntegrante:         app.NuevoQuitarIntegrante(repoProyectos),
		AgregarRol:               app.NuevoAgregarRol(repoProyectos),
		ConsultarProyecto:        app.NuevoConsultarProyecto(repoProyectos),
		CrearElementoBacklog:     app.NuevoCrearElementoBacklog(repoBacklog),
		EditarElementoBacklog:    app.NuevoEditarElementoBacklog(repoBacklog),    // HU-06
		EliminarElementoBacklog:  app.NuevoEliminarElementoBacklog(repoBacklog),  // HU-06
		ConsultarElementoBacklog: app.NuevoConsultarElementoBacklog(repoBacklog), // HU-06
		ListarProyectos:          app.NuevoListarProyectos(repoProyectos),        // HU-30
		// time.Now como reloj: "hoy" es el día actual en la zona del servidor (HU-18, RN-18.4).
		RegistrarEsfuerzo:  app.NuevoRegistrarEsfuerzo(repoProyectos, repoBacklog, repoEsfuerzo, time.Now),
		CrearSprint:        app.NuevoCrearSprint(repoProyectos, repoSprints),
		ListarSprints:      app.NuevoListarSprints(repoProyectos, repoSprints),
		ListarBacklog:      app.NuevoListarBacklog(repoProyectos, repoBacklog),
		AsignarHistoria:    app.NuevoAsignarHistoria(repoBacklog),
		DesasignarHistoria: app.NuevoDesasignarHistoria(repoBacklog),
		CerrarSprint:       app.NuevoCerrarSprint(repoProyectos, repoSprints),      // HU-11
		CompletarHistoria:  app.NuevoCompletarHistoria(repoProyectos, repoBacklog), // HU-10
		// time.Now como reloj: "hoy" es el día actual en la zona del servidor (HU-20, RN-20.5).
		RegistrarDefecto:        app.NuevoRegistrarDefecto(repoProyectos, repoBacklog, repoSprints, repoDefectos, time.Now),
		IniciarRonda:            app.NuevoIniciarRonda(repoProyectos, repoBacklog, repoRondas),           // HU-14
		Votar:                   app.NuevoVotar(repoProyectos, repoRondas),                               // HU-14
		ConsultarRonda:          app.NuevoConsultarRonda(repoProyectos, repoBacklog, repoRondas),         // HU-14
		ListarSprintsAnteriores: app.NuevoListarSprintsAnteriores(repoProyectos, repoSprints),            // HU-12
		ConsultarVelocidad:      app.NuevoConsultarVelocidad(repoProyectos, repoSprints),                 // HU-22
		ListarRondas:            app.NuevoListarRondas(repoProyectos, repoBacklog, repoRondas),           // HU-16
		RegistrarEstimacion:     app.NuevoRegistrarEstimacion(repoProyectos, repoBacklog, repoRondas),    // HU-17
		CompararEsfuerzo:        app.NuevoCompararEsfuerzo(repoProyectos, repoBacklog, repoEsfuerzo),     // HU-19
		ConsultarHorasProyecto:  app.NuevoConsultarHorasProyecto(repoProyectos, repoEsfuerzo),            // HU-23
		ConsultarHorasSprint:    app.NuevoConsultarHorasSprint(repoProyectos, repoSprints, repoEsfuerzo), // HU-23
		// time.Now como reloj: "hoy" es el día actual en la zona del servidor (HU-21, RN-21.5).
		ActualizarEstadoDefecto: app.NuevoActualizarEstadoDefecto(repoProyectos, repoSprints, repoDefectos, time.Now),
		ListarDefectos:          app.NuevoListarDefectos(repoProyectos, repoDefectos),                       // HU-21
		ConsultarCalidadSprints: app.NuevoConsultarCalidadSprints(repoProyectos, repoSprints, repoDefectos), // HU-24
		ArmarReporteProyecto:    armarReporteProyecto,
		ArmarReporteSprint:      armarReporteSprint, // HU-26
		// time.Now como reloj: la emisión del PDF es la hora del pedido, en UTC-3 (HU-27, RN-27.10).
		ExportarReporteProyecto: app.NuevoExportarReporteProyecto(repoProyectos, fuenteReporte, generadorPDF, time.Now),
		ExportarReporteSprint:   app.NuevoExportarReporteSprint(repoProyectos, repoSprints, fuenteReporte, generadorPDF, time.Now), // HU-27
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

	// Apagado ordenado con Ctrl+C o SIGTERM (la señal con que Render detiene el servidor en cada
	// despliegue): termina los pedidos en curso y deja correr el defer db.Close().
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
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
