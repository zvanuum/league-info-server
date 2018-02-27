package service

import (
	"context"
	"net/url"

	"github.com/zachvanuum/league-info-server"
	"github.com/zachvanuum/league-lib"
	leaguemodel "github.com/zachvanuum/league-lib/model"
)

type ChampionsRequest struct {
	Role string `json:"role,omitempty"`
}

type ChampionsServicer interface {
	Champions(context.Context, url.Values) (leaguemodel.Champions, error)
}

type ChampionsService struct {
	leagueClient *leaguelib.LeagueClient
}

func NewChampionService(leagueClient *leaguelib.LeagueClient) ChampionsService {
	return ChampionsService{
		leagueClient: leagueClient,
	}
}

func (championsService *ChampionsService) Champions(_ context.Context, req ChampionsRequest) (leaguemodel.Champions, error) {
	query := createChampionsQuery(req)

	champions := championsService.leagueClient.GetChampions(query)
	champions = applyChampionsRequestFilters(champions, req)

	return champions, nil
}

func createChampionsQuery(req ChampionsRequest) url.Values {
	query := url.Values{}
	if len(req.Role) != 0 {
		query.Add("champListData", "tags")
	}

	return query
}

func applyChampionsRequestFilters(champions leaguemodel.Champions, req ChampionsRequest) leaguemodel.Champions {
	if len(req.Role) != 0 {
		champions = applyRoleFilter(champions, req.Role)
	}

	return champions
}

func applyRoleFilter(champions leaguemodel.Champions, role string) leaguemodel.Champions {
	filteredChampions := make(map[string]leaguemodel.Champion)

	for key, champion := range champions.Data {
		switch role {
		case "adc", "marksman":
			if leagueinfoserver.SliceContains("Marksman", champion.Tags) {
				filteredChampions[key] = champion
			}
		default:
		}
	}

	champions.Data = filteredChampions
	return champions
}
