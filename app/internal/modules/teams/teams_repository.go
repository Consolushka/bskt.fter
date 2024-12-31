package teams

import (
	"IMP/app/database"
)

func FirstOrCreate(team TeamModel) (TeamModel, error) {
	var result TeamModel
	dbConnection := database.Connect()

	tx := dbConnection.
		Attrs(TeamModel{
			Name: team.Name,
		}).
		FirstOrCreate(&result, TeamModel{Alias: team.Alias, LeagueID: team.LeagueID})

	return result, tx.Error
}
