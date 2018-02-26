package handler

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/zachvanuum/league-info-server/service"
)

func MakeHealthEndpoint(service service.HealthService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return service.Health()
	}
}

//DecodeHealthRequest stubbed for creating handler, works as just a simple passthrough
func DecodeHealthRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return nil, nil
}
