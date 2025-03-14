package internal

import (
	"IMP/app/internal/abstract"
	"IMP/app/internal/base/components/request_components"
	"IMP/app/internal/modules/imp"
	"IMP/app/internal/modules/leagues"
	"IMP/app/internal/modules/leagues/api/internal/requests"
	"IMP/app/internal/modules/leagues/api/internal/responses"
	"IMP/app/internal/modules/leagues/domain/models"
	"IMP/app/internal/modules/leagues/domain/resources"
	teamsModels "IMP/app/internal/modules/teams/domain/models"
	"IMP/app/internal/utils/array_utils"
	"IMP/app/internal/utils/time_utils"
	"net/http"
)

type Controller struct {
	abstract.BaseController

	service *leagues.Service
}

func NewController() *Controller {
	return &Controller{
		service: leagues.NewService(),
	}
}

func (c *Controller) GetLeagues(w http.ResponseWriter, r *requests.GetLeaguesRequest) {
	var response []resources.LeagueResource

	leagueModels, err := c.service.GetAllLeagues()
	if err != nil {
		c.InternalServerError(w, err)
		return
	}

	for _, league := range leagueModels {
		response = append(response, resources.NewLeagueResponse(&league))
	}

	c.Ok(w, response)
}

func (c *Controller) GetGamesByLeagueAndDate(w http.ResponseWriter, r *requests.GetGamesByLeagueAndDate) {
	gamesModel, err := c.service.GetGamesByLeagueAndDate(r.Id(), *r.Date())
	if err != nil {
		c.InternalServerError(w, err)
		return
	}

	response := responses.NewGamesByDateResponse(*r.Date(), gamesModel)

	c.Ok(w, response)
}

func (c *Controller) GetTeamsByLeague(w http.ResponseWriter, r *request_components.HasIdPathParam) {
	leagueModel, err := c.service.GetLeague(r.Id())
	teamsModel, err := c.service.GetTeamsByLeague(r.Id())
	if err != nil {
		c.InternalServerError(w, err)
		return
	}

	response := responses.TeamsInLeagueResponse{League: resources.NewLeagueResponse(leagueModel),
		Teams: array_utils.Map(teamsModel, func(team teamsModels.Team) responses.TeamInLeagueResponse {
			return responses.TeamInLeagueResponse{
				Id:    team.ID,
				Name:  team.Name,
				Alias: team.Alias,
			}
		}),
	}

	c.Ok(w, response)
}

func (c *Controller) PlayersRanking(w http.ResponseWriter, r *requests.PlayersByMetricsRankingRequest) {
	limit := r.Limit()

	leagueModel, err := c.service.GetLeague(r.Id())
	var rankingResources []resources.PlayerMetricRank

	ranking, err := c.service.GetPlayersRanking(r.Id(), limit, r.MinMinutesPerGame(), r.AvgMinutesPerGame(), r.MinGamesPlayed())
	if err != nil {
		c.InternalServerError(w, err)
		return
	}

	rankingResources = array_utils.Map(*ranking, func(playerMetricRank models.PlayerImpRanking) resources.PlayerMetricRank {
		return resources.PlayerMetricRank{
			Rank:             playerMetricRank.Rank,
			FullName:         playerMetricRank.FullNameLocal,
			AvgMinutesPlayed: time_utils.SecondsToFormat(int(playerMetricRank.AvgPlayedSeconds), time_utils.PlayedTimeFormat),
			GamesPlayed:      playerMetricRank.GamesPlayed,
			ImpPer:           imp.EvaluatePer(int(playerMetricRank.AvgPlayedSeconds), nil, nil, nil, r.Per(), leagueModel, &playerMetricRank.AvgImpClean),
		}
	})

	rankingResources = array_utils.Sort(rankingResources, func(a, b resources.PlayerMetricRank) bool {
		return a.ImpPer > b.ImpPer
	})

	if limit != nil {
		rankingResources = rankingResources[:*limit]
	}

	for i := range rankingResources {
		rankingResources[i].Rank = i + 1
	}

	response := responses.RankingResponse{
		Ranking: rankingResources,
	}

	c.Ok(w, response)
}
