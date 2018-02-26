package service

import (
	"context"
	"net/url"

	"github.com/zachvanuum/league-lib"
	leaguemodel "github.com/zachvanuum/league-lib/model"
)

type ChampionServicer interface {
	Champions(context.Context, url.Values) (leaguemodel.Champions, error)
}

type ChampionService struct {
	leagueClient *leaguelib.LeagueClient
}

func NewChampionService(leagueClient *leaguelib.LeagueClient) ChampionService {
	return ChampionService{
		leagueClient: leagueClient,
	}
}

func (championService *ChampionService) Champions(_ context.Context, values url.Values) (leaguemodel.Champions, error) {
	champions := championService.leagueClient.GetChampions(nil)

	return champions, nil
}
