package teams

import "IMP/app/internal/core/players"

type TeamStatEntity struct {
	TeamModel         TeamModel
	GameTeamStatModel GameTeamStatModel
	PlayerStats       []players.GameTeamPlayerStatModel
}
