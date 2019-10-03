package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	log "github.com/Financial-Times/go-logger"
	"github.com/Financial-Times/http-handlers-go/httphandlers"
	status "github.com/Financial-Times/service-status-go/httphandlers"
	"github.com/gorilla/mux"
	cli "github.com/jawher/mow.cli"
	"github.com/rcrowley/go-metrics"
)

const appDescription = "UPP Golang Microservice Template short description - please amend"

func main() {
	app := cli.App("graph-query-api", appDescription)

	appSystemCode := app.String(cli.StringOpt{
		Name:   "app-system-code",
		Value:  "graph-query-api",
		Desc:   "System Code of the application",
		EnvVar: "APP_SYSTEM_CODE",
	})

	appName := app.String(cli.StringOpt{
		Name:   "app-name",
		Value:  "graph-query-api",
		Desc:   "Application name",
		EnvVar: "APP_NAME",
	})

	port := app.String(cli.StringOpt{
		Name:   "port",
		Value:  "8080",
		Desc:   "Port to listen on",
		EnvVar: "APP_PORT",
	})
	neoURL := app.String(cli.StringOpt{
		Name:   "neo-url",
		Value:  "bolt://localhost:7687/",
		Desc:   "neo4j endpoint URL",
		EnvVar: "NEO_URL"})

	logLevel := app.String(cli.StringOpt{
		Name:   "logLevel",
		Value:  "INFO",
		Desc:   "Logging level (DEBUG, INFO, WARN, ERROR)",
		EnvVar: "LOG_LEVEL",
	})

	log.InitLogger(*appSystemCode, *logLevel)
	log.Infof("[Startup] graph-query-api is starting ")

	app.Action = func() {
		log.Infof("System code: %s, App Name: %s, Port: %s", *appSystemCode, *appName, *port)

		neoClient, err := NewNeoClient(*neoURL)
		if err != nil {
			log.Fatal(err)
		}

		handler := requestHandler{neoClient}

		go func() {
			serveEndpoints(*appSystemCode, *appName, *port, handler)
		}()

		waitForSignal()

		err = neoClient.Close()
		if err != nil {
			log.Error(err)
		}
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Errorf("App could not start, error=[%s]\n", err)
		return
	}
}

func serveEndpoints(appSystemCode string, appName string, port string, requestHandler requestHandler) {
	healthService := newHealthService(appSystemCode, appName, appDescription)

	serveMux := http.NewServeMux()

	serveMux.HandleFunc(healthPath, http.HandlerFunc(fthealth.Handler(healthService.Health())))
	serveMux.HandleFunc(status.GTGPath, status.NewGoodToGoHandler(healthService.GTG))
	serveMux.HandleFunc(status.BuildInfoPath, status.BuildInfoHandler)

	servicesRouter := mux.NewRouter()
	servicesRouter.HandleFunc("/search", requestHandler.searchEndpoint).Methods("GET")

	var monitoringRouter http.Handler = servicesRouter
	monitoringRouter = httphandlers.TransactionAwareRequestLoggingHandler(log.Logger(), monitoringRouter)
	monitoringRouter = httphandlers.HTTPMetricsHandler(metrics.DefaultRegistry, monitoringRouter)

	serveMux.Handle("/", monitoringRouter)

	server := &http.Server{Addr: ":" + port, Handler: serveMux}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Infof("HTTP server closing with message: %v", err)
		}
	}()

	waitForSignal()
	log.Infof("[Shutdown] graph-query-api is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("Unable to stop http server: %v", err)
	}
}

func waitForSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
