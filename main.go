package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	// "github.com/zachvanuum/league-lib"
)

type Server struct {
	http.Server
	shutdownReq chan bool
	reqCount    uint32
}

func main() {
	// leagueClient := leaguelib.NewLeagueClient("RGAPI-9ce99ca9-6f83-4af2-a30b-e50c5de8ea05", leaguelib.NorthAmerica)

	server := newServer()

	done := make(chan bool)
	go func() {
		log.Printf("Starting server listening on %s\n", server.Server.Addr)
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Error with HTTP server: %v", err)
		}
		done <- true
	}()

	server.waitShutdown()

	<-done
	log.Printf("Finished")
}

func newServer() *Server {
	server := &Server{
		Server: http.Server{
			Addr:         ":8080",
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		shutdownReq: make(chan bool),
	}

	server.Handler = initRouter()
	return server
}

func initRouter() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/health", HealthHandler).Methods("GET")

	return router
}

func (server *Server) waitShutdown() {
	irqSig := make(chan os.Signal, 1)
	signal.Notify(irqSig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-irqSig:
		log.Printf("Shutdown request (signal: %v)", sig)
	}

	log.Printf("Shutting down HTTP server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Printf("Shutdown request error: %v", err)
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[HealthHandler] Received request")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{ "alive": true }`)
}
