package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/zachvanuum/league-info-server/handler"
	"github.com/zachvanuum/league-info-server/service"
	"github.com/zachvanuum/league-lib"
)

type server struct {
	router       *mux.Router
	leagueClient *leaguelib.LeagueClient
	services     *services
	handlers     *handlers
}

type services struct {
	HealthService    service.HealthService
	ChampionsService service.ChampionService
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
	port := envString("PORT", defaultPort)
	host := envString("HOST", defaultHost)
	addr := host + ":" + port

	server := server{
		router:       mux.NewRouter(),
		leagueClient: leaguelib.NewLeagueClient("RGAPI-59a66e58-77d9-4fb4-993f-2e1542ad8241", leaguelib.NorthAmerica),
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
	healthHandler := kithttp.NewServer(
		handler.MakeHealthEndpoint(server.services.HealthService),
		handler.DecodeHealthRequest,
		encodeResponse,
	)

	championsHandler := kithttp.NewServer(
		handler.MakeChampionsEndpoint(server.services.ChampionsService),
		handler.DecodeChampionsRequest,
		encodeResponse,
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

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
