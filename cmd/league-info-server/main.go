package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/zachvanuum/league-info-server"
	"github.com/zachvanuum/league-info-server/handler"
	"github.com/zachvanuum/league-info-server/service"
	"github.com/zachvanuum/league-lib"
)

type server struct {
	router       *mux.Router
	leagueClient *leaguelib.LeagueClient
	services     *services
	handlers     *handlers
	logger       kitlog.Logger
}

type services struct {
	HealthService    service.HealthService
	ChampionsService service.ChampionsService
}

type handlers struct {
	HealthHandler    *kithttp.Server
	ChampionsHandler *kithttp.Server
}

const (
	defaultPort = "8080"
	defaultHost = "localhost"
)

func main() {
	port := leagueinfoserver.EnvString("PORT", defaultPort)
	host := leagueinfoserver.EnvString("HOST", defaultHost)
	addr := host + ":" + port

	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stderr)
		logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
		logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
	}

	server := server{
		router:       mux.NewRouter(),
		leagueClient: leaguelib.NewLeagueClient("RGAPI-6a8118a1-3626-41c8-b9c3-239b06a718ad", leaguelib.NorthAmerica),
		logger:       logger,
	}

	server.createServices()
	server.createHandlers()
	server.createRoutes()

	srv := &http.Server{
		Handler:      server.router,
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	go func() {
		log.Printf("Starting server listening on port 8080\n")
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(ctx)
	log.Printf("Shutting down\n")
	os.Exit(0)
}

func (server *server) createServices() {
	championsService := service.NewChampionService(server.leagueClient)

	server.services = &services{
		ChampionsService: championsService,
	}
}

func (server *server) createHandlers() {
	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(server.logger),
	}

	health := handler.MakeHealthEndpoint(server.services.HealthService)
	health = leagueinfoserver.LoggingMiddleware(kitlog.With(server.logger, "method", "Health"))(health)
	healthHandler := kithttp.NewServer(
		health,
		handler.DecodeHealthRequest,
		encodeResponse,
		options...,
	)

	champions := handler.MakeChampionsEndpoint(server.services.ChampionsService)
	champions = leagueinfoserver.LoggingMiddleware(kitlog.With(server.logger, "method", "Champions"))(champions)
	championsHandler := kithttp.NewServer(
		champions,
		handler.DecodeChampionsRequest,
		encodeResponse,
		options...,
	)

	server.handlers = &handlers{
		HealthHandler:    healthHandler,
		ChampionsHandler: championsHandler,
	}
}

func (server *server) createRoutes() {
	server.router.Handle("/health", server.handlers.HealthHandler).Methods("GET")
	server.router.Handle("/champions", server.handlers.ChampionsHandler).Methods("GET")
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
