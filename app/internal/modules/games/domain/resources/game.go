package resources

import (
	"IMP/app/internal/modules/games/domain/models"
	"IMP/app/internal/modules/teams/resources"
	"strconv"
)

type Game struct {
	GameId    string         `json:"game_id"`
	Scheduled string         `json:"scheduled"`
	Home      resources.Team `json:"home"`
	Away      resources.Team `json:"away"`
}

func NewGameResource(model models.Game) Game {
	return Game{
		GameId:    strconv.Itoa(model.ID),
		Scheduled: model.ScheduledAt.Format("02.01.2006 15:04"),
		Home:      resources.NewTeamResource(model.HomeTeamStats),
		Away:      resources.NewTeamResource(model.AwayTeamStats),
	}
}
