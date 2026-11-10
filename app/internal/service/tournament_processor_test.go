package service

import (
	"IMP/app/internal/adapters/games_repo"
	"IMP/app/internal/adapters/players_repo"
	"IMP/app/internal/adapters/stats_provider"
	"IMP/app/internal/adapters/teams_repo"
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/core/teams"
	"IMP/app/internal/core/tournaments"
	"IMP/app/internal/ports"
	"errors"
	"strings"
	"testing"
	"time"

	composite_logger "github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTournamentProcessor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPersistence := NewMockPersistenceServiceInterface(ctrl)
	mockStatsProvider := stats_provider.NewMockStatsProvider(ctrl)
	mockPlayersRepo := players_repo.NewMockPlayersRepo(ctrl)
	mockTeamsRepo := teams_repo.NewMockTeamsRepo(ctrl)
	mockGamesRepo := games_repo.NewMockGamesRepo(ctrl)

	mockAggregator := ports.NewMockAggregator(ctrl)
	processor := NewTournamentProcessor(mockStatsProvider, mockPersistence, mockPlayersRepo, mockTeamsRepo, mockGamesRepo, mockAggregator, tournaments.TournamentModel{
		Id:       12,
		Provider: tournaments.TournamentProvider{ProviderName: "ESPN_NBA"},
	})

	assert.NotNil(t, processor)
	assert.Equal(t, uint(12), processor.tournament.Id)
	assert.Equal(t, mockPersistence, processor.persistenceService)
	assert.Equal(t, mockStatsProvider, processor.statsProvider)
	assert.Equal(t, mockPlayersRepo, processor.playersRepo)
	assert.Equal(t, mockTeamsRepo, processor.teamsRepo)
	assert.Equal(t, mockGamesRepo, processor.gamesRepo)
	assert.Equal(t, "ESPN_NBA", processor.teamIdScope)
}

func TestTournamentProcessor_ProcessByPeriod(t *testing.T) {
	from := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 2, 14, 23, 59, 59, 0, time.UTC)

	t.Run("returns error when stats provider fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockTeams := teams_repo.NewMockTeamsRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return(nil, errors.New("stats provider error"))

		mockAggregator := ports.NewMockAggregator(ctrl)
		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockTeams, mockGames, mockAggregator, tournaments.TournamentModel{Id: 1})
		count, err := processor.ProcessByPeriod(from, to)

		require.ErrorContains(t, err, "stats provider error")
		assert.Equal(t, 0, count)
	})

	t.Run("skips existing games without enrichment and saving", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockTeams := teams_repo.NewMockTeamsRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		gameEntities := []games.GameStatEntity{
			{GameModel: games.GameModel{Title: "LAL - BOS", ScheduledAt: from}},
		}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return(gameEntities, nil)
		expectTeamsFoundByName(mockTeams, "LAL", "BOS", 1)
		mockGames.EXPECT().Exists(gomock.Any()).Return(true, nil)

		mockAggregator := ports.NewMockAggregator(ctrl)
		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockTeams, mockGames, mockAggregator, tournaments.TournamentModel{Id: 99})
		count, err := processor.ProcessByPeriod(from, to)

		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("enriches and saves new game", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockTeams := teams_repo.NewMockTeamsRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		gameEntity := games.GameStatEntity{GameModel: games.GameModel{Title: "MIA - NYK", ScheduledAt: from}}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return([]games.GameStatEntity{gameEntity}, nil)
		expectTeamsFoundByName(mockTeams, "LAL", "BOS", 1)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, nil)
		mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(gameEntity, nil)
		mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil)
		mockAggregator := ports.NewMockAggregator(ctrl)
		mockAggregator.EXPECT().NotifyGameImported(uint(77)).Return(nil)

		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockTeams, mockGames, mockAggregator, tournaments.TournamentModel{Id: 77})
		count, err := processor.ProcessByPeriod(from, to)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("returns error when game existence check fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockTeams := teams_repo.NewMockTeamsRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		gameEntity := games.GameStatEntity{GameModel: games.GameModel{Title: "DAL - PHX", ScheduledAt: from}}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return([]games.GameStatEntity{gameEntity}, nil)
		expectTeamsFoundByName(mockTeams, "LAL", "BOS", 1)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, errors.New("db error"))

		mockAggregator := ports.NewMockAggregator(ctrl)
		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockTeams, mockGames, mockAggregator, tournaments.TournamentModel{Id: 11})
		count, err := processor.ProcessByPeriod(from, to)

		require.ErrorContains(t, err, "db error")
		assert.Equal(t, 0, count)
	})

	t.Run("continues when enrich game stats fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockTeams := teams_repo.NewMockTeamsRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		game1 := games.GameStatEntity{GameModel: games.GameModel{Title: "GAME-1", ScheduledAt: from}}
		game2 := games.GameStatEntity{GameModel: games.GameModel{Title: "GAME-2", ScheduledAt: from}}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return([]games.GameStatEntity{game1, game2}, nil)
		expectTeamsFoundByName(mockTeams, "LAL", "BOS", 2)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, nil).Times(2)
		gomock.InOrder(
			mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(games.GameStatEntity{}, errors.New("enrich error")),
			mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(game2, nil),
		)
		mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil).Times(1)
		mockAggregator := ports.NewMockAggregator(ctrl)
		mockAggregator.EXPECT().NotifyGameImported(uint(20)).Return(nil).Times(1)

		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockTeams, mockGames, mockAggregator, tournaments.TournamentModel{Id: 20})
		count, err := processor.ProcessByPeriod(from, to)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("continues when save game fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockTeams := teams_repo.NewMockTeamsRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		game1 := games.GameStatEntity{GameModel: games.GameModel{Title: "SAVE-1", ScheduledAt: from}}
		game2 := games.GameStatEntity{GameModel: games.GameModel{Title: "SAVE-2", ScheduledAt: from}}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return([]games.GameStatEntity{game1, game2}, nil)
		expectTeamsFoundByName(mockTeams, "LAL", "BOS", 2)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, nil).Times(2)
		mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(game1, nil)
		mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(game2, nil)
		gomock.InOrder(
			mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(errors.New("save error")),
			mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil),
		)
		mockAggregator := ports.NewMockAggregator(ctrl)
		mockAggregator.EXPECT().NotifyGameImported(uint(21)).Return(nil).Times(1)

		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockTeams, mockGames, mockAggregator, tournaments.TournamentModel{Id: 21})
		count, err := processor.ProcessByPeriod(from, to)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("continues when ListByFullName fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockTeams := teams_repo.NewMockTeamsRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		game := games.GameStatEntity{GameModel: games.GameModel{Title: "PLAYERS-ERR", ScheduledAt: from}}
		enriched := games.GameStatEntity{
			GameModel: game.GameModel,
			HomeTeamStat: teams.TeamStatEntity{
				PlayerStats: []players.PlayerStatisticEntity{
					{PlayerModel: players.PlayerModel{FullName: "John Doe"}},
				},
			},
		}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return([]games.GameStatEntity{game}, nil)
		expectTeamsFoundByName(mockTeams, "LAL", "BOS", 1)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, nil)
		mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(enriched, nil)
		mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil)
		mockPlayers.EXPECT().ListByFullName("John Doe").Return(nil, errors.New("players repo error"))

		mockAggregator := ports.NewMockAggregator(ctrl)
		mockAggregator.EXPECT().NotifyGameImported(uint(22)).Return(nil).Times(1)

		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockTeams, mockGames, mockAggregator, tournaments.TournamentModel{Id: 22})
		count, err := processor.ProcessByPeriod(from, to)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("continues when player not found and GetPlayerBio fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockTeams := teams_repo.NewMockTeamsRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		game := games.GameStatEntity{GameModel: games.GameModel{Title: "BIO-ERR", ScheduledAt: from}}
		enriched := games.GameStatEntity{
			GameModel: game.GameModel,
			HomeTeamStat: teams.TeamStatEntity{
				PlayerStats: []players.PlayerStatisticEntity{
					{
						PlayerExternalId: "123",
						PlayerModel:      players.PlayerModel{FullName: ""},
					},
				},
			},
		}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return([]games.GameStatEntity{game}, nil)
		expectTeamsFoundByName(mockTeams, "LAL", "BOS", 1)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, nil)
		mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(enriched, nil)
		mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil)
		mockPlayers.EXPECT().ListByFullName("").Return([]players.PlayerModel{}, nil)
		mockStats.EXPECT().GetPlayerBio("123").Return(players.PlayerBioEntity{}, errors.New("bio error"))

		mockAggregator := ports.NewMockAggregator(ctrl)
		mockAggregator.EXPECT().NotifyGameImported(uint(23)).Return(nil).Times(1)

		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockTeams, mockGames, mockAggregator, tournaments.TournamentModel{Id: 23})
		count, err := processor.ProcessByPeriod(from, to)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("continues when player not found and GetPlayerBio succeeds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockPlayers := players_repo.NewMockPlayersRepo(ctrl)
		mockTeams := teams_repo.NewMockTeamsRepo(ctrl)
		mockGames := games_repo.NewMockGamesRepo(ctrl)

		game := games.GameStatEntity{GameModel: games.GameModel{Title: "BIO-OK", ScheduledAt: from}}
		enriched := games.GameStatEntity{
			GameModel: game.GameModel,
			HomeTeamStat: teams.TeamStatEntity{
				PlayerStats: []players.PlayerStatisticEntity{
					{
						PlayerExternalId: "555",
						PlayerModel:      players.PlayerModel{FullName: "", BirthDate: time.Time{}},
					},
				},
			},
		}

		mockStats.EXPECT().GetGamesStatsByPeriod(from, to).Return([]games.GameStatEntity{game}, nil)
		expectTeamsFoundByName(mockTeams, "LAL", "BOS", 1)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, nil)
		mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(enriched, nil)
		mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil)
		mockPlayers.EXPECT().ListByFullName("").Return([]players.PlayerModel{}, nil)
		mockStats.EXPECT().GetPlayerBio("555").Return(players.PlayerBioEntity{
			FullName:  "Recovered Player",
			BirthDate: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		}, nil)

		mockAggregator := ports.NewMockAggregator(ctrl)
		mockAggregator.EXPECT().NotifyGameImported(uint(24)).Return(nil).Times(1)

		processor := NewTournamentProcessor(mockStats, mockPersistence, mockPlayers, mockTeams, mockGames, mockAggregator, tournaments.TournamentModel{Id: 24})
		count, err := processor.ProcessByPeriod(from, to)

		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}

type warningsRecorder struct {
	warnings  []map[string]interface{}
	byMessage map[string][]map[string]interface{}
}

func (r *warningsRecorder) InitLogger() composite_logger.Logger         { return r }
func (r *warningsRecorder) Info(msg string, ctx map[string]interface{}) { r.record(msg, ctx) }

// Считаем только предупреждение о датах: чужой warning не должен выдавать себя за него.
func (r *warningsRecorder) Warn(msg string, ctx map[string]interface{}) {
	r.record(msg, ctx)
	if strings.Contains(msg, "Imported game is scheduled outside tournament dates") {
		r.warnings = append(r.warnings, ctx)
	}
}
func (r *warningsRecorder) Error(msg string, ctx map[string]interface{}) { r.record(msg, ctx) }
func (r *warningsRecorder) Fatal(string, map[string]interface{})         {}

// Логгер дописывает перед сообщением уровень («[INFO] …»), а тесты ищут записи по самому тексту.
func (r *warningsRecorder) record(msg string, ctx map[string]interface{}) {
	if _, text, ok := strings.Cut(msg, "] "); ok {
		msg = text
	}
	r.byMessage[msg] = append(r.byMessage[msg], ctx)
}

func recordWarnings(t *testing.T) *warningsRecorder {
	recorder := &warningsRecorder{byMessage: map[string][]map[string]interface{}{}}
	composite_logger.Init(recorder)
	t.Cleanup(func() { composite_logger.Init() })
	return recorder
}

func TestTournamentProcessor_ProcessByPeriod_WarnsAboutGamesOutsideTournamentDates(t *testing.T) {
	startAt := time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC)
	endAt := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)
	insideDates := time.Date(2026, 3, 1, 19, 0, 0, 0, time.UTC)
	afterEnd := endAt.AddDate(0, 1, 15)
	beforeStart := startAt.AddDate(0, 0, -3)

	cases := []struct {
		name             string
		scheduledAt      []time.Time
		exists           bool
		saveErr          error
		expectedWarnings int
	}{
		{name: "game inside dates gives no warning", scheduledAt: []time.Time{insideDates}, expectedWarnings: 0},
		{name: "game after end gives warning", scheduledAt: []time.Time{afterEnd}, expectedWarnings: 1},
		{name: "game before start gives warning", scheduledAt: []time.Time{beforeStart}, expectedWarnings: 1},
		{name: "several games outside dates give one warning per poll", scheduledAt: []time.Time{beforeStart, insideDates, afterEnd, afterEnd}, expectedWarnings: 1},
		{name: "already imported game gives no warning", scheduledAt: []time.Time{afterEnd}, exists: true, expectedWarnings: 0},
		{name: "game that failed to save gives no warning", scheduledAt: []time.Time{afterEnd}, saveErr: errors.New("save error"), expectedWarnings: 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			gameEntities := make([]games.GameStatEntity, 0, len(tc.scheduledAt))
			for _, scheduledAt := range tc.scheduledAt {
				gameEntities = append(gameEntities, games.GameStatEntity{GameModel: games.GameModel{Title: "LAL - BOS", ScheduledAt: scheduledAt}})
			}

			mockStats := stats_provider.NewMockStatsProvider(ctrl)
			mockStats.EXPECT().GetGamesStatsByPeriod(gomock.Any(), gomock.Any()).Return(gameEntities, nil)
			mockTeams := teams_repo.NewMockTeamsRepo(ctrl)
			expectTeamsFoundByName(mockTeams, "LAL", "BOS", len(gameEntities))
			mockStats.EXPECT().EnrichGameStats(gomock.Any()).DoAndReturn(func(game games.GameStatEntity) (games.GameStatEntity, error) {
				return game, nil
			}).AnyTimes()
			mockGames := games_repo.NewMockGamesRepo(ctrl)
			mockGames.EXPECT().Exists(gomock.Any()).Return(tc.exists, nil).AnyTimes()
			mockPersistence := NewMockPersistenceServiceInterface(ctrl)
			mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(tc.saveErr).AnyTimes()
			mockAggregator := ports.NewMockAggregator(ctrl)
			mockAggregator.EXPECT().NotifyGameImported(uint(5)).Return(nil).AnyTimes()

			recorder := recordWarnings(t)
			tournament := tournaments.TournamentModel{Id: 5, Name: "NBA 2025-2026", StartAt: &startAt, EndAt: &endAt}
			processor := NewTournamentProcessor(mockStats, mockPersistence, players_repo.NewMockPlayersRepo(ctrl), mockTeams, mockGames, mockAggregator, tournament)

			_, err := processor.ProcessByPeriod(startAt, endAt)

			require.NoError(t, err)
			assert.Len(t, recorder.warnings, tc.expectedWarnings)
			assert.Equal(t, time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC), *tournament.StartAt)
			assert.Equal(t, time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC), *tournament.EndAt)
		})
	}

	t.Run("warning names tournament, its dates and game date", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		game := games.GameStatEntity{GameModel: games.GameModel{Title: "LAL - BOS", ScheduledAt: afterEnd}}

		mockStats := stats_provider.NewMockStatsProvider(ctrl)
		mockStats.EXPECT().GetGamesStatsByPeriod(gomock.Any(), gomock.Any()).Return([]games.GameStatEntity{game}, nil)
		mockTeams := teams_repo.NewMockTeamsRepo(ctrl)
		expectTeamsFoundByName(mockTeams, "LAL", "BOS", 1)
		mockStats.EXPECT().EnrichGameStats(gomock.Any()).Return(game, nil)
		mockGames := games_repo.NewMockGamesRepo(ctrl)
		mockGames.EXPECT().Exists(gomock.Any()).Return(false, nil)
		mockPersistence := NewMockPersistenceServiceInterface(ctrl)
		mockPersistence.EXPECT().SaveGame(gomock.Any()).Return(nil)
		mockAggregator := ports.NewMockAggregator(ctrl)
		mockAggregator.EXPECT().NotifyGameImported(uint(5)).Return(nil)

		recorder := recordWarnings(t)
		tournament := tournaments.TournamentModel{Id: 5, Name: "NBA 2025-2026", StartAt: &startAt, EndAt: &endAt}
		processor := NewTournamentProcessor(mockStats, mockPersistence, players_repo.NewMockPlayersRepo(ctrl), mockTeams, mockGames, mockAggregator, tournament)

		_, err := processor.ProcessByPeriod(startAt, endAt)

		require.NoError(t, err)
		require.Len(t, recorder.warnings, 1)
		assert.Equal(t, uint(5), recorder.warnings[0]["tournamentId"])
		assert.Equal(t, "NBA 2025-2026", recorder.warnings[0]["tournamentName"])
		assert.Equal(t, &startAt, recorder.warnings[0]["startAt"])
		assert.Equal(t, &endAt, recorder.warnings[0]["endAt"])
		assert.Equal(t, "LAL - BOS", recorder.warnings[0]["title"])
		assert.Equal(t, afterEnd, recorder.warnings[0]["scheduledAt"])
	})
}

// expectTeamsFoundByName ожидает поиск команд по имени для каждой игры: первой ищется команда хозяев, второй — гостей.
func expectTeamsFoundByName(mockTeams *teams_repo.MockTeamsRepo, homeAlias, awayAlias string, gamesCount int) {
	calls := 0
	mockTeams.EXPECT().FirstOrCreate(gomock.Any()).DoAndReturn(func(teams.TeamModel) (teams.TeamModel, bool, error) {
		calls++
		if calls%2 == 1 {
			return teams.TeamModel{Id: 1, Alias: homeAlias}, false, nil
		}
		return teams.TeamModel{Id: 2, Alias: awayAlias}, false, nil
	}).Times(2 * gamesCount)
}

const (
	homeExternalId = 31856
	awayExternalId = 3277
	teamsScope     = "INFOBASKET"
)

var teamsTournament = tournaments.TournamentModel{Id: 7, Provider: tournaments.TournamentProvider{ProviderName: teamsScope}}

// providerGame — игра в том виде, в каком её отдаёт провайдер: с обрезанным сокращением гостей в названии.
func providerGame(homeId, awayId int) games.GameStatEntity {
	return games.GameStatEntity{
		GameModel:          games.GameModel{Title: "UNIT - BK \"Ta", ScheduledAt: time.Date(2026, 2, 14, 19, 0, 0, 0, time.UTC)},
		HomeTeamExternalId: homeId,
		AwayTeamExternalId: awayId,
		HomeTeamStat:       teams.TeamStatEntity{TeamModel: teams.TeamModel{Name: "UNICS", HomeTown: "Kazan", Alias: "UNIT"}},
		AwayTeamStat:       teams.TeamStatEntity{TeamModel: teams.TeamModel{Name: "BK \"Tambov\"", HomeTown: "Tambov", Alias: "BK \"Ta"}},
	}
}

// savingRun — опрос одной новой игры, которая доходит до сохранения; запоминает, с каким названием
// проверялось существование игры и что ушло на сохранение.
type savingRun struct {
	processor    *TournamentProcessor
	teams        *teams_repo.MockTeamsRepo
	checkedTitle string
	saved        games.GameStatEntity
}

func newSavingRun(t *testing.T, game games.GameStatEntity) *savingRun {
	ctrl := gomock.NewController(t)
	run := &savingRun{teams: teams_repo.NewMockTeamsRepo(ctrl)}

	mockStats := stats_provider.NewMockStatsProvider(ctrl)
	mockStats.EXPECT().GetGamesStatsByPeriod(gomock.Any(), gomock.Any()).Return([]games.GameStatEntity{game}, nil)
	mockStats.EXPECT().EnrichGameStats(gomock.Any()).DoAndReturn(func(g games.GameStatEntity) (games.GameStatEntity, error) {
		return g, nil
	})
	mockGames := games_repo.NewMockGamesRepo(ctrl)
	mockGames.EXPECT().Exists(gomock.Any()).DoAndReturn(func(m games.GameModel) (bool, error) {
		run.checkedTitle = m.Title
		return false, nil
	})
	mockPersistence := NewMockPersistenceServiceInterface(ctrl)
	mockPersistence.EXPECT().SaveGame(gomock.Any()).DoAndReturn(func(g games.GameStatEntity) error {
		run.saved = g
		return nil
	})
	mockAggregator := ports.NewMockAggregator(ctrl)
	mockAggregator.EXPECT().NotifyGameImported(teamsTournament.Id).Return(nil)

	run.processor = NewTournamentProcessor(mockStats, mockPersistence, players_repo.NewMockPlayersRepo(ctrl), run.teams, mockGames, mockAggregator, teamsTournament)
	return run
}

func (r *savingRun) process(t *testing.T) {
	count, err := r.processor.ProcessByPeriod(time.Time{}, time.Time{})
	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func TestTournamentProcessor_ProcessByPeriod_BuildsTitleFromStoredAliases(t *testing.T) {
	run := newSavingRun(t, providerGame(homeExternalId, awayExternalId))
	run.teams.EXPECT().FindByExternalId(teamsScope, homeExternalId).Return(teams.TeamModel{Id: 1, Alias: "UNIT"}, true, nil)
	run.teams.EXPECT().FindByExternalId(teamsScope, awayExternalId).Return(teams.TeamModel{Id: 2, Alias: "TAMBOV"}, true, nil)

	run.process(t)

	assert.Equal(t, "UNIT - TAMBOV", run.checkedTitle)
	assert.Equal(t, "UNIT - TAMBOV", run.saved.GameModel.Title)
}

func TestTournamentProcessor_ProcessByPeriod_ResolvesTeamByExternalId(t *testing.T) {
	storedHome := teams.TeamModel{Id: 11, Name: "UNICS Kazan", HomeTown: "Kazan", Alias: "UNIT"}
	storedAway := teams.TeamModel{Id: 12, Name: "Tambov", HomeTown: "Tambov", Alias: "TAMBOV"}
	run := newSavingRun(t, providerGame(homeExternalId, awayExternalId))
	run.teams.EXPECT().FindByExternalId(teamsScope, homeExternalId).Return(storedHome, true, nil)
	run.teams.EXPECT().FindByExternalId(teamsScope, awayExternalId).Return(storedAway, true, nil)

	run.process(t)

	assert.Equal(t, storedHome, run.saved.HomeTeamStat.TeamModel)
	assert.Equal(t, storedAway, run.saved.AwayTeamStat.TeamModel)
}

// expectLinkByName ожидает, что команду не нашли по id провайдера, нашли или создали по имени
// и связали с id; повторное чтение связи возвращает linked.
func expectLinkByName(mockTeams *teams_repo.MockTeamsRepo, externalId int, candidate, byName teams.TeamModel, created bool, linked teams.TeamModel) {
	gomock.InOrder(
		mockTeams.EXPECT().FindByExternalId(teamsScope, externalId).Return(teams.TeamModel{}, false, nil),
		mockTeams.EXPECT().FirstOrCreate(candidate).Return(byName, created, nil),
		mockTeams.EXPECT().LinkExternalId(teamsScope, externalId, byName.Id).Return(nil),
		mockTeams.EXPECT().FindByExternalId(teamsScope, externalId).Return(linked, true, nil),
	)
}

func TestTournamentProcessor_ProcessByPeriod_LinksExistingTeamFoundByName(t *testing.T) {
	game := providerGame(homeExternalId, awayExternalId)
	home := teams.TeamModel{Id: 1, Name: "UNICS", Alias: "UNIT"}
	away := teams.TeamModel{Id: 2, Name: "BK \"Tambov\"", Alias: "TAMBOV"}
	run := newSavingRun(t, game)
	expectLinkByName(run.teams, homeExternalId, game.HomeTeamStat.TeamModel, home, false, home)
	expectLinkByName(run.teams, awayExternalId, game.AwayTeamStat.TeamModel, away, false, away)
	recorder := recordWarnings(t)

	run.process(t)

	assert.Equal(t, "UNIT - TAMBOV", run.saved.GameModel.Title)
	assert.Empty(t, recorder.byMessage["New team created"])
	linkedLogs := recorder.byMessage["Team linked to provider id"]
	require.Len(t, linkedLogs, 2)
	assert.Equal(t, map[string]interface{}{"teamId": uint(1), "scope": teamsScope, "externalId": homeExternalId}, linkedLogs[0])
	assert.Equal(t, map[string]interface{}{"teamId": uint(2), "scope": teamsScope, "externalId": awayExternalId}, linkedLogs[1])
}

func TestTournamentProcessor_ProcessByPeriod_CreatesAndLinksNewTeam(t *testing.T) {
	game := providerGame(homeExternalId, awayExternalId)
	home := teams.TeamModel{Id: 41, Name: "UNICS", HomeTown: "Kazan", Alias: "UNI"}
	away := teams.TeamModel{Id: 42, Name: "BK \"Tambov\"", HomeTown: "Tambov", Alias: "TAM"}
	run := newSavingRun(t, game)
	expectLinkByName(run.teams, homeExternalId, game.HomeTeamStat.TeamModel, home, true, home)
	expectLinkByName(run.teams, awayExternalId, game.AwayTeamStat.TeamModel, away, true, away)
	recorder := recordWarnings(t)

	run.process(t)

	assert.Equal(t, "UNI - TAM", run.saved.GameModel.Title)
	assert.Empty(t, recorder.byMessage["Team linked to provider id"])
	createdLogs := recorder.byMessage["New team created"]
	require.Len(t, createdLogs, 2)
	assert.Equal(t, map[string]interface{}{
		"teamId": uint(41), "name": "UNICS", "homeTown": "Kazan", "alias": "UNI",
		"scope": teamsScope, "externalId": homeExternalId, "tournamentId": teamsTournament.Id,
	}, createdLogs[0])
	assert.Equal(t, map[string]interface{}{
		"teamId": uint(42), "name": "BK \"Tambov\"", "homeTown": "Tambov", "alias": "TAM",
		"scope": teamsScope, "externalId": awayExternalId, "tournamentId": teamsTournament.Id,
	}, createdLogs[1])
}

func TestTournamentProcessor_ProcessByPeriod_UsesTeamStoredInLinkAfterRace(t *testing.T) {
	game := providerGame(homeExternalId, awayExternalId)
	byName := teams.TeamModel{Id: 50, Name: "UNICS", Alias: "UNICS"}
	linkedByOtherWorker := teams.TeamModel{Id: 1, Name: "UNICS Kazan", Alias: "UNIT"}
	run := newSavingRun(t, game)
	expectLinkByName(run.teams, homeExternalId, game.HomeTeamStat.TeamModel, byName, false, linkedByOtherWorker)
	run.teams.EXPECT().FindByExternalId(teamsScope, awayExternalId).Return(teams.TeamModel{Id: 2, Alias: "TAMBOV"}, true, nil)

	run.process(t)

	assert.Equal(t, "UNIT - TAMBOV", run.checkedTitle)
	assert.Equal(t, linkedByOtherWorker, run.saved.HomeTeamStat.TeamModel)
}

func TestTournamentProcessor_ProcessByPeriod_SkipsLinkWithoutExternalId(t *testing.T) {
	game := providerGame(0, 0)
	run := newSavingRun(t, game)
	gomock.InOrder(
		run.teams.EXPECT().FirstOrCreate(game.HomeTeamStat.TeamModel).Return(teams.TeamModel{Id: 1, Alias: "UNIT"}, false, nil),
		run.teams.EXPECT().FirstOrCreate(game.AwayTeamStat.TeamModel).Return(teams.TeamModel{Id: 2, Alias: "TAMBOV"}, false, nil),
	)

	run.process(t)

	assert.Equal(t, "UNIT - TAMBOV", run.saved.GameModel.Title)
}

func TestTournamentProcessor_ProcessByPeriod_ReturnsErrorWhenTeamLookupFails(t *testing.T) {
	repoErr := errors.New("teams repo error")
	game := providerGame(homeExternalId, awayExternalId)
	home := teams.TeamModel{Id: 1, Name: "UNICS", Alias: "UNIT"}

	cases := []struct {
		name        string
		setup       func(*teams_repo.MockTeamsRepo)
		checkLogged bool
	}{
		{
			name: "home team lookup by external id fails",
			setup: func(m *teams_repo.MockTeamsRepo) {
				m.EXPECT().FindByExternalId(teamsScope, homeExternalId).Return(teams.TeamModel{}, false, repoErr)
			},
			checkLogged: true,
		},
		{
			name: "away team lookup by external id fails",
			setup: func(m *teams_repo.MockTeamsRepo) {
				m.EXPECT().FindByExternalId(teamsScope, homeExternalId).Return(home, true, nil)
				m.EXPECT().FindByExternalId(teamsScope, awayExternalId).Return(teams.TeamModel{}, false, repoErr)
			},
		},
		{
			name: "lookup by name fails",
			setup: func(m *teams_repo.MockTeamsRepo) {
				m.EXPECT().FindByExternalId(teamsScope, homeExternalId).Return(teams.TeamModel{}, false, nil)
				m.EXPECT().FirstOrCreate(gomock.Any()).Return(teams.TeamModel{}, false, repoErr)
			},
		},
		{
			name: "link fails",
			setup: func(m *teams_repo.MockTeamsRepo) {
				m.EXPECT().FindByExternalId(teamsScope, homeExternalId).Return(teams.TeamModel{}, false, nil)
				m.EXPECT().FirstOrCreate(gomock.Any()).Return(home, false, nil)
				m.EXPECT().LinkExternalId(teamsScope, homeExternalId, home.Id).Return(repoErr)
			},
		},
		{
			name: "reading link after insert fails",
			setup: func(m *teams_repo.MockTeamsRepo) {
				gomock.InOrder(
					m.EXPECT().FindByExternalId(teamsScope, homeExternalId).Return(teams.TeamModel{}, false, nil),
					m.EXPECT().FirstOrCreate(gomock.Any()).Return(home, false, nil),
					m.EXPECT().LinkExternalId(teamsScope, homeExternalId, home.Id).Return(nil),
					m.EXPECT().FindByExternalId(teamsScope, homeExternalId).Return(teams.TeamModel{}, false, repoErr),
				)
			},
		},
		{
			name: "link is not found after insert",
			setup: func(m *teams_repo.MockTeamsRepo) {
				gomock.InOrder(
					m.EXPECT().FindByExternalId(teamsScope, homeExternalId).Return(teams.TeamModel{}, false, nil),
					m.EXPECT().FirstOrCreate(gomock.Any()).Return(home, false, nil),
					m.EXPECT().LinkExternalId(teamsScope, homeExternalId, home.Id).Return(nil),
					m.EXPECT().FindByExternalId(teamsScope, homeExternalId).Return(teams.TeamModel{}, false, nil),
				)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockStats := stats_provider.NewMockStatsProvider(ctrl)
			mockStats.EXPECT().GetGamesStatsByPeriod(gomock.Any(), gomock.Any()).Return([]games.GameStatEntity{game}, nil)
			mockTeams := teams_repo.NewMockTeamsRepo(ctrl)
			tc.setup(mockTeams)
			recorder := recordWarnings(t)

			// Репозиторий игр, сохранение и уведомление без ожиданий: любой их вызов уронит тест.
			processor := NewTournamentProcessor(mockStats, NewMockPersistenceServiceInterface(ctrl), players_repo.NewMockPlayersRepo(ctrl), mockTeams, games_repo.NewMockGamesRepo(ctrl), ports.NewMockAggregator(ctrl), teamsTournament)
			count, err := processor.ProcessByPeriod(time.Time{}, time.Time{})

			require.ErrorContains(t, err, "resolve team")
			assert.Equal(t, 0, count)
			require.Len(t, recorder.byMessage["Failed to resolve team"], 1)
			if !tc.checkLogged {
				return
			}
			require.ErrorIs(t, err, repoErr)
			logged := recorder.byMessage["Failed to resolve team"][0]
			// К записи об ошибке логгер сам добавляет стек вызовов, поэтому сверяем только наши поля.
			assert.Equal(t, teamsTournament.Id, logged["tournamentId"])
			assert.Equal(t, teamsScope, logged["scope"])
			assert.Equal(t, homeExternalId, logged["externalId"])
			assert.Equal(t, "UNICS", logged["name"])
			assert.Equal(t, "Kazan", logged["homeTown"])
			loggedErr, ok := logged["error"].(error)
			require.True(t, ok)
			require.ErrorIs(t, loggedErr, repoErr)
		})
	}
}

// scopedStatsProvider — провайдер, который дополнительно сообщает свою базу id команд.
type scopedStatsProvider struct {
	*stats_provider.MockStatsProvider
	scope string
}

func (p scopedStatsProvider) TeamIdScope() string { return p.scope }

func TestNewTournamentProcessor_TakesTeamIdScopeFromProvider(t *testing.T) {
	cases := []struct {
		name          string
		provider      func(*stats_provider.MockStatsProvider) ports.StatsProvider
		expectedScope string
	}{
		{
			name: "provider with own team id scope",
			provider: func(m *stats_provider.MockStatsProvider) ports.StatsProvider {
				return scopedStatsProvider{MockStatsProvider: m, scope: "INFOBASKET:org"}
			},
			expectedScope: "INFOBASKET:org",
		},
		{
			name:          "provider without own team id scope",
			provider:      func(m *stats_provider.MockStatsProvider) ports.StatsProvider { return m },
			expectedScope: teamsScope,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockStats := stats_provider.NewMockStatsProvider(ctrl)
			mockStats.EXPECT().GetGamesStatsByPeriod(gomock.Any(), gomock.Any()).Return([]games.GameStatEntity{providerGame(homeExternalId, awayExternalId)}, nil)
			mockTeams := teams_repo.NewMockTeamsRepo(ctrl)
			// Опрос обрывается на первой команде: здесь важно только, с какой базой id ушёл поиск.
			mockTeams.EXPECT().FindByExternalId(tc.expectedScope, homeExternalId).Return(teams.TeamModel{}, false, errors.New("stop"))
			recordWarnings(t)

			processor := NewTournamentProcessor(tc.provider(mockStats), NewMockPersistenceServiceInterface(ctrl), players_repo.NewMockPlayersRepo(ctrl), mockTeams, games_repo.NewMockGamesRepo(ctrl), ports.NewMockAggregator(ctrl), teamsTournament)
			_, err := processor.ProcessByPeriod(time.Time{}, time.Time{})

			require.Error(t, err)
			assert.Equal(t, tc.expectedScope, processor.teamIdScope)
		})
	}
}
