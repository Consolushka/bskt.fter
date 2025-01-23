package games

import (
	"IMP/app/database"
	"IMP/app/internal/modules/calculations"
	"IMP/app/internal/modules/imp"
	"IMP/app/internal/modules/imp/enums"
	"IMP/app/internal/modules/imp/models"
	"IMP/app/internal/modules/leagues"
	"IMP/app/internal/modules/players"
	enums2 "IMP/app/internal/modules/statistics/enums"
	"IMP/app/internal/modules/teams"
	"IMP/app/internal/utils/array_utils"
	"gorm.io/gorm"
)

type Service struct {
	gamesRepository   *Repository
	leaguesRepository *leagues.Repository
	dbConnection      *gorm.DB
}

func NewService() *Service {
	return &Service{
		gamesRepository:   NewRepository(),
		leaguesRepository: leagues.NewRepository(),
		dbConnection:      database.GetDB().Debug(),
	}
}

// GetGame returns game by specific id
//
// Also preloads:
//   - League
//   - Teams
//   - Players stats
//   - Players models
func (s *Service) GetGame(id int) (*GameModel, error) {
	gameModel, err := s.retrieveGameModelById(id)

	if err != nil {
		return nil, err
	}

	return gameModel, nil
}

// GetGameMetrics returns game metrics by specific id
//
// Calculates IMP metrics for every player
func (s *Service) GetGameMetrics(id int, impPers []enums.ImpPERs) (*models.GameImpMetrics, error) {
	var leagueModel leagues.League

	gameModel, err := s.retrieveGameModelById(id)
	if err != nil {
		return nil, err
	}

	tx := s.dbConnection.Debug().
		First(&leagueModel, leagues.League{ID: gameModel.LeagueID})
	if tx.Error != nil {
		return nil, tx.Error
	}
	league := enums2.FromString(leagueModel.AliasEn)

	gameImpMetrics := s.mapGameModelToImpMetricsModel(gameModel, impPers, league)

	return gameImpMetrics, nil
}

func (s *Service) retrieveGameModelById(id int) (*GameModel, error) {
	var gameModel GameModel

	tx := s.dbConnection.Debug().
		Preload("League").
		Preload("HomeTeamStats").
		Preload("HomeTeamStats.Team").
		Preload("HomeTeamStats.PlayerGameStats").
		Preload("HomeTeamStats.PlayerGameStats.Player").
		Preload("AwayTeamStats").
		Preload("AwayTeamStats.Team").
		Preload("AwayTeamStats.PlayerGameStats").
		Preload("AwayTeamStats.PlayerGameStats.Player").
		First(&gameModel, GameModel{ID: id})

	return &gameModel, tx.Error
}

func (s *Service) mapGameModelToImpMetricsModel(gameModel *GameModel, impPers []enums.ImpPERs, league enums2.League) *models.GameImpMetrics {
	return &models.GameImpMetrics{
		Id:        gameModel.ID,
		Scheduled: &gameModel.ScheduledAt,
		Home: models.TeamImpMetrics{
			Alias:       gameModel.HomeTeamStats.Team.Alias,
			TotalPoints: gameModel.HomeTeamStats.Points,
			Players:     s.mapTeamPlayersMetrics(gameModel.HomeTeamStats, gameModel.AwayTeamStats, gameModel.PlayedMinutes, impPers, league),
		},
		Away: models.TeamImpMetrics{
			Alias:       gameModel.AwayTeamStats.Team.Alias,
			TotalPoints: gameModel.AwayTeamStats.Points,
			Players:     s.mapTeamPlayersMetrics(gameModel.AwayTeamStats, gameModel.HomeTeamStats, gameModel.PlayedMinutes, impPers, league),
		},
		FullGameTime: gameModel.PlayedMinutes,
	}
}

func (s *Service) mapTeamPlayersMetrics(currentTeam teams.TeamGameStats, oposingTeam teams.TeamGameStats, fullGameTime int, impPers []enums.ImpPERs, league enums2.League) []models.PlayerImpMetrics {
	return array_utils.Map(currentTeam.PlayerGameStats, func(playerGameStats players.PlayerGameStats) models.PlayerImpMetrics {
		playerImpPers := make([]models.PlayerImpPersMetrics, len(impPers))
		cleanImp := imp.CalculatePlayerImpPerMinute(float64(playerGameStats.PlayedSeconds)/60, playerGameStats.PlsMin, currentTeam.Points-oposingTeam.Points, fullGameTime)

		for i, impPer := range impPers {
			timeBase := enums.TimeBasesByLeagueAndPers(league, impPer)

			reliability := calculations.CalculateReliability(float64(playerGameStats.PlayedSeconds)/60, timeBase)
			pure := cleanImp * float64(timeBase.Minutes())

			playerImpPers[i] = models.PlayerImpPersMetrics{
				Per: impPer,
				IMP: pure * reliability,
			}
		}

		playerImpPers = array_utils.Sort(playerImpPers, func(a, b models.PlayerImpPersMetrics) bool {
			return a.Per.Order() < b.Per.Order()
		})

		return models.PlayerImpMetrics{
			FullNameLocal: playerGameStats.Player.FullNameLocal,
			FullNameEn:    playerGameStats.Player.FullNameEn,
			SecondsPlayed: playerGameStats.PlayedSeconds,
			PlsMin:        playerGameStats.PlsMin,
			ImpPers:       playerImpPers,
		}
	})
}
