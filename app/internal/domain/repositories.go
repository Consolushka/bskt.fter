package domain

import (
	"IMP/app/database"
	"errors"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

type GamesRepository struct {
	dbConnection *gorm.DB
}

func NewGamesRepository() *GamesRepository {
	return &GamesRepository{
		dbConnection: database.GetDB(),
	}
}

func (r *GamesRepository) FirstOrCreate(game Game) (Game, error) {
	var result Game

	tx := r.dbConnection.
		Attrs(Game{
			PlayedMinutes: game.PlayedMinutes,
			OfficialId:    game.OfficialId,
		}).
		FirstOrCreate(&result, Game{
			HomeTeamID:  game.HomeTeamID,
			AwayTeamID:  game.AwayTeamID,
			LeagueID:    game.LeagueID,
			ScheduledAt: game.ScheduledAt,
		})

	return result, tx.Error
}

// Exists checks if game exists in db. Can check by id or official_id
func (r *GamesRepository) Exists(game Game) (bool, error) {
	var exists bool
	var condition string

	if game.ID != 0 {
		condition = "id = " + strconv.Itoa(game.ID)
	} else {
		if game.OfficialId != "" {
			condition = "official_id = '" + game.OfficialId + "'"
		}
	}

	err := r.dbConnection.
		Model(&Game{}).
		Select("count(*) > 0").
		Where(condition).
		Find(&exists).
		Error

	if err != nil {
		return false, err
	}

	return exists, nil
}

type LeaguesRepository struct {
	db *gorm.DB
}

func NewLeaguesRepository() *LeaguesRepository {
	return &LeaguesRepository{
		db: database.GetDB(),
	}
}

func (r *LeaguesRepository) List() ([]League, error) {
	var result []League

	tx := r.db.Find(&result)

	return result, tx.Error
}

func (r *LeaguesRepository) FirstByAliasEn(aliasEn string) (*League, error) {
	var result League
	tx := r.db.First(&result, League{AliasEn: strings.ToUpper(aliasEn)})
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &result, nil
}

type PlayersRepository struct {
	dbConnection *gorm.DB
}

func NewPlayersRepository() *PlayersRepository {
	return &PlayersRepository{
		dbConnection: database.GetDB(),
	}
}

func (r *PlayersRepository) FirstByOfficialId(id string) (*Player, error) {
	var result Player

	tx := r.dbConnection.
		First(
			&result,
			Player{
				OfficialId: id,
			})

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &result, tx.Error
}

func (r *PlayersRepository) FirstOrCreate(player Player) (*Player, error) {
	var result Player

	tx := r.dbConnection.
		Attrs(Player{
			FullNameLocal: player.FullNameLocal,
			FullNameEn:    player.FullNameEn,
			BirthDate:     player.BirthDate,
		}).
		FirstOrCreate(
			&result,
			Player{
				OfficialId: player.OfficialId,
			})

	return &result, tx.Error
}

func (r *PlayersRepository) FirstOrCreatePlayerGameStats(stats PlayerGameStats) error {
	tx := r.dbConnection.Attrs(
		PlayerGameStats{
			PlayedSeconds: stats.PlayedSeconds,
			PlsMin:        stats.PlsMin,
			IsBench:       stats.IsBench,
			IMPClean:      stats.IMPClean,
		}).
		FirstOrCreate(
			&PlayerGameStats{},
			PlayerGameStats{
				PlayerID:   stats.PlayerID,
				TeamGameId: stats.TeamGameId,
			})

	return tx.Error
}

type TeamsRepository struct {
	dbConnection *gorm.DB
}

func NewTeamsRepository() *TeamsRepository {
	return &TeamsRepository{
		dbConnection: database.GetDB(),
	}
}

func (r *TeamsRepository) FirstOrCreate(team Team) (Team, error) {
	var result Team

	tx := r.dbConnection.
		Attrs(
			Team{
				Name:       team.Name,
				OfficialId: team.OfficialId,
			}).
		FirstOrCreate(&result,
			Team{
				Alias:    team.Alias,
				LeagueID: team.LeagueID,
			},
		)

	return result, tx.Error
}

func (r *TeamsRepository) FirstOrCreateTeamGameStats(stats TeamGameStats) (TeamGameStats, error) {
	var result TeamGameStats

	tx := r.dbConnection.Attrs(
		TeamGameStats{
			Points: stats.Points,
		}).
		FirstOrCreate(
			&result,
			TeamGameStats{
				TeamId: stats.TeamId,
				GameId: stats.GameId,
			})

	return result, tx.Error
}
