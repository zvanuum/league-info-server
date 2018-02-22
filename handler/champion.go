package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	"github.com/zachvanuum/league-info-server/service"
)

type championsRequest struct {
	Values url.Values `json:"values"`
}

func MakeChampionsEndpoint(service service.ChampionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(championsRequest)
		champions := service.LeagueClient.GetChampions(req.Values)
		return champions, nil
	}
}

func DecodeChampionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request championsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
