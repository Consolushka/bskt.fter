package games

import (
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
	"time"
)

type Service struct {
	gamesRepository   *gamesDomain.Repository
	leaguesRepository *leaguesDomain.Repository
}

func NewService() *Service {
	return &Service{
		gamesRepository:   gamesDomain.NewRepository(),
		leaguesRepository: leaguesDomain.NewRepository(),
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
	gameModel, err := s.gamesRepository.FirstGameStatsById(id)

	if err != nil {
		return nil, err
	}

	return gameModel, nil
}

// GetGameMetrics returns game metrics by specific id
//
// Calculates IMP metrics for every player
func (s *Service) GetGameMetrics(id int, impPers []enums.ImpPERs) (*impModels.GameImpMetrics, error) {
	gameModel, err := s.gamesRepository.FirstGameStatsById(id)
	if err != nil {
		return nil, err
	}

	leagueModel, err := s.leaguesRepository.FirstById(gameModel.LeagueID)
	if err != nil {
		return nil, err
	}

	gameImpMetrics := s.mapGameModelToImpMetricsModel(gameModel, impPers, leagueModel)

	return gameImpMetrics, nil
}

// GetGames fetches all games for specific date and preloads all related impModels
func (s *Service) GetGames(date time.Time) ([]gamesModels.Game, error) {
	return s.gamesRepository.GamesStatsByDateList(date)
}

func (s *Service) mapGameModelToImpMetricsModel(gameModel *gamesModels.Game, impPers []enums.ImpPERs, league *leaguesModels.League) *impModels.GameImpMetrics {
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

func (s *Service) mapTeamPlayersMetrics(currentTeam teamModels.TeamGameStats, opposingTeam teamModels.TeamGameStats, fullGameTime int, impPers []enums.ImpPERs, league *leaguesModels.League) []impModels.PlayerImpMetrics {
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
