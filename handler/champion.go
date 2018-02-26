package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/endpoint"
	"github.com/zachvanuum/league-info-server/service"
)

type championsRequest struct {
	Values url.Values `json:"values,omitempty"`
}

func MakeChampionsEndpoint(service service.ChampionService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(championsRequest)
		champions, err := service.Champions(ctx, req.Values)
		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve champions, %v", err)
		}

		return champions, nil
	}
}

func DecodeChampionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	query := r.URL.Query()
	if len(query) == 0 {
		query = nil
	}
	request := championsRequest{
		Values: query,
	}

	return request, nil
}
