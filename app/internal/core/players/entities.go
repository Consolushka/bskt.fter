package players

import "time"

type PlayerStatisticEntity struct {
	PlayerModel             PlayerModel
	PlayerExternalId        string
	GameTeamPlayerStatModel GameTeamPlayerStatModel
}

type PlayerBioEntity struct {
	FullName  string
	BirthDate time.Time
}
