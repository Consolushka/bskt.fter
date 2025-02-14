package games

import (
	"IMP/app/database"
	gamesDomain "IMP/app/internal/modules/games/domain"
	gamesModels "IMP/app/internal/modules/games/domain/models"
	"IMP/app/internal/modules/imp"
	"IMP/app/internal/modules/imp/domain/enums"
	impModels "IMP/app/internal/modules/imp/domain/models"
	leaguesDomain "IMP/app/internal/modules/leagues/domain"
	leaguesModels "IMP/app/internal/modules/leagues/domain/models"
	playersDomain "IMP/app/internal/modules/players/domain/models"
	teamModels "IMP/app/internal/modules/teams/domain/models"
	"IMP/app/internal/utils/array_utils"
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type Service struct {
	gamesRepository   *gamesDomain.Repository
	leaguesRepository *leaguesDomain.Repository
	dbConnection      *gorm.DB
}

func NewService() *Service {
	return &Service{
		gamesRepository:   gamesDomain.NewRepository(),
		leaguesRepository: leaguesDomain.NewRepository(),
		dbConnection:      database.GetDB(),
	}
}

// GetGame returns game by specific id
//
// Also preloads:
//   - League
//   - Teams
//   - Players stats
//   - Players impModels
func (s *Service) GetGame(id int) (*gamesModels.Game, error) {
	gameModel, err := s.retrieveGameModelById(id)

	if err != nil {
		return nil, err
	}

	return gameModel, nil
}

// GetGameMetrics returns game metrics by specific id
//
// Calculates IMP metrics for every player
func (s *Service) GetGameMetrics(id int, impPers []enums.ImpPERs) (*impModels.GameImpMetrics, error) {
	var leagueModel leaguesModels.League

	gameModel, err := s.retrieveGameModelById(id)
	if err != nil {
		return nil, err
	}

	tx := s.dbConnection.First(&leagueModel, leaguesModels.League{ID: gameModel.LeagueID})
	if tx.Error != nil {
		return nil, tx.Error
	}

	gameImpMetrics := s.mapGameModelToImpMetricsModel(gameModel, impPers, leagueModel)

	return gameImpMetrics, nil
}

// GetGames fetches all games for specific date and preloads all related impModels
func (s *Service) GetGames(date time.Time) ([]gamesModels.Game, error) {
	var gamesModel []gamesModels.Game

	tx := s.dbConnection.
		Model(&gamesModels.Game{}).
		Preload("League").
		Preload("HomeTeamStats").
		Preload("HomeTeamStats.Team").
		Preload("HomeTeamStats.PlayerGameStats").
		Preload("HomeTeamStats.PlayerGameStats.Player").
		Preload("AwayTeamStats").
		Preload("AwayTeamStats.Team").
		Preload("AwayTeamStats.PlayerGameStats").
		Preload("AwayTeamStats.PlayerGameStats.Player").
		Where("DATE(scheduled_at) = @date", sql.Named("date", date.Format("2006-01-02"))).
		Find(&gamesModel)

	return gamesModel, tx.Error
}

func (s *Service) retrieveGameModelById(id int) (*gamesModels.Game, error) {
	var gameModel gamesModels.Game

	tx := s.dbConnection.
		Preload("League").
		Preload("HomeTeamStats").
		Preload("HomeTeamStats.Team").
		Preload("HomeTeamStats.PlayerGameStats").
		Preload("HomeTeamStats.PlayerGameStats.Player").
		Preload("AwayTeamStats").
		Preload("AwayTeamStats.Team").
		Preload("AwayTeamStats.PlayerGameStats").
		Preload("AwayTeamStats.PlayerGameStats.Player").
		First(&gameModel, gamesModels.Game{ID: id})

	return &gameModel, tx.Error
}

func (s *Service) mapGameModelToImpMetricsModel(gameModel *gamesModels.Game, impPers []enums.ImpPERs, league leaguesModels.League) *impModels.GameImpMetrics {
	return &impModels.GameImpMetrics{
		Id:        gameModel.ID,
		Scheduled: &gameModel.ScheduledAt,
		Home: impModels.TeamImpMetrics{
			Alias:       gameModel.HomeTeamStats.Team.Alias,
			TotalPoints: gameModel.HomeTeamStats.Points,
			Players:     s.mapTeamPlayersMetrics(gameModel.HomeTeamStats, gameModel.AwayTeamStats, gameModel.PlayedMinutes, impPers, league),
		},
		Away: impModels.TeamImpMetrics{
			Alias:       gameModel.AwayTeamStats.Team.Alias,
			TotalPoints: gameModel.AwayTeamStats.Points,
			Players:     s.mapTeamPlayersMetrics(gameModel.AwayTeamStats, gameModel.HomeTeamStats, gameModel.PlayedMinutes, impPers, league),
		},
		FullGameTime: gameModel.PlayedMinutes,
	}
}

func (s *Service) mapTeamPlayersMetrics(currentTeam teamModels.TeamGameStats, opposingTeam teamModels.TeamGameStats, fullGameTime int, impPers []enums.ImpPERs, league leaguesModels.League) []impModels.PlayerImpMetrics {
	return array_utils.Map(currentTeam.PlayerGameStats, func(playerGameStats playersDomain.PlayerGameStats) impModels.PlayerImpMetrics {
		playerImpPers := make([]impModels.PlayerImpPersMetrics, len(impPers))

		cleanImp := imp.EvaluateClean(playerGameStats.PlayedSeconds, playerGameStats.PlsMin, currentTeam.Points-opposingTeam.Points, fullGameTime)

		for i, impPer := range impPers {

			playerImpPers[i] = impModels.PlayerImpPersMetrics{
				Per: impPer,
				IMP: imp.EvaluatePer(playerGameStats.PlayedSeconds, playerGameStats.PlsMin, currentTeam.Points-opposingTeam.Points, fullGameTime, impPer, league, &cleanImp),
			}
		}

		playerImpPers = array_utils.Sort(playerImpPers, func(a, b impModels.PlayerImpPersMetrics) bool {
			return a.Per.Order() < b.Per.Order()
		})

		return impModels.PlayerImpMetrics{
			FullNameLocal: playerGameStats.Player.FullNameLocal,
			FullNameEn:    playerGameStats.Player.FullNameEn,
			SecondsPlayed: playerGameStats.PlayedSeconds,
			PlsMin:        playerGameStats.PlsMin,
			ImpPers:       playerImpPers,
		}
	})
}
