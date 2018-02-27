package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/zachvanuum/league-info-server/service"
)

func MakeChampionsEndpoint(svc service.ChampionsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(service.ChampionsRequest)
		champions, err := svc.Champions(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("Failed to retrieve champions, %v", err)
		}

		return champions, nil
	}
}

func DecodeChampionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	role := r.URL.Query().Get("role")

	request := service.ChampionsRequest{
		Role: role,
	}

	return request, nil
}
