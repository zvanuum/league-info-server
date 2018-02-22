package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/zachvanuum/league-info-server/handler"
	"github.com/zachvanuum/league-info-server/service"
	"github.com/zachvanuum/league-lib"
)

// type Server struct {
// 	http.Server
// 	Client   *leaguelib.LeagueClient
// 	ReqCount uint32
// }

func main() {
	leagueClient := leaguelib.NewLeagueClient("RGAPI-676bb774-b3c1-44ec-b695-e3703e31125f", leaguelib.NorthAmerica)
	championsService := service.ChampionService{
		LeagueClient: leagueClient,
	}

	championsHandler := httptransport.NewServer(
		handler.MakeChampionsEndpoint(championsService),
		handler.DecodeChampionsRequest,
		encodeResponse,
	)

	http.Handle("/champions", championsHandler)

	done := make(chan bool)
	go func() {
		log.Printf("Starting server listening on port 8080\n")
		log.Printf("%v\n", http.ListenAndServe(":8080", nil))
		// err := server.ListenAndServe()
		// if err != nil {
		// 	log.Printf("Error with HTTP server: %v", err)
		// }
		done <- true
	}()

	// server.waitShutdown()

	<-done
	log.Printf("Finished")
}

// func (server *Server) waitShutdown() {
// 	irqSig := make(chan os.Signal, 1)
// 	signal.Notify(irqSig, syscall.SIGINT, syscall.SIGTERM)

// 	select {
// 	case sig := <-irqSig:
// 		log.Printf("Shutdown request (signal: %v)", sig)
// 	}

// 	log.Printf("Shutting down HTTP server...")

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	err := server.Shutdown(ctx)
// 	if err != nil {
// 		log.Printf("Shutdown request error: %v", err)
// 	}
// }

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
