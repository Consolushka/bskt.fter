package resources

import (
	"IMP/app/internal/modules/games/domain/models"
	teamsResources "IMP/app/internal/modules/teams/domain/resources"
	"strconv"
)

type Game struct {
	GameId    string              `json:"game_id"`
	Scheduled string              `json:"scheduled"`
	Home      teamsResources.Team `json:"home"`
	Away      teamsResources.Team `json:"away"`
}

func NewGameResource(model models.Game) Game {
	return Game{
		GameId:    strconv.Itoa(model.ID),
		Scheduled: model.ScheduledAt.Format("02.01.2006 15:04"),
		Home:      teamsResources.NewTeamResource(model.HomeTeamStats),
		Away:      teamsResources.NewTeamResource(model.AwayTeamStats),
	}
}
