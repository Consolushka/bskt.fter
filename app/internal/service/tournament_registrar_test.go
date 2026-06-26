package service

import (
	"IMP/app/internal/adapters/leagues_repo"
	"IMP/app/internal/adapters/tournaments_repo"
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/leagues"
	"IMP/app/internal/core/players"
	"IMP/app/internal/core/tournaments"
	"IMP/app/internal/ports"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeStatsProvider реализует только ports.StatsProvider — даты сообщать не умеет.
type fakeStatsProvider struct{}

func (fakeStatsProvider) GetGamesStatsByPeriod(_, _ time.Time) ([]games.GameStatEntity, error) {
	return nil, nil
}
func (fakeStatsProvider) EnrichGameStats(g games.GameStatEntity) (games.GameStatEntity, error) {
	return g, nil
}
func (fakeStatsProvider) GetPlayerBio(_ string) (players.PlayerBioEntity, error) {
	return players.PlayerBioEntity{}, nil
}

// fakePeriodProvider дополнительно реализует ports.TournamentPeriodProvider.
type fakePeriodProvider struct {
	fakeStatsProvider
	start time.Time
	end   time.Time
	err   error
}

func (f fakePeriodProvider) GetTournamentPeriod(_ string) (time.Time, time.Time, error) {
	return f.start, f.end, f.err
}

func newRegistrar(ctrl *gomock.Controller, provider ports.StatsProvider) (*TournamentRegistrar, *tournaments_repo.MockTournamentsRepo) {
	leaguesRepo := leagues_repo.NewMockLeaguesRepo(ctrl)
	leaguesRepo.EXPECT().FirstOrCreate(gomock.Any()).Return(leagues.LeagueModel{Id: 7}, nil).AnyTimes()

	tournamentsRepo := tournaments_repo.NewMockTournamentsRepo(ctrl)

	factory := func(_ string, _ *string, _ *map[string]interface{}) (ports.StatsProvider, error) {
		return provider, nil
	}

	return NewTournamentRegistrar(leaguesRepo, tournamentsRepo, factory), tournamentsRepo
}

// captureCreate настраивает мок Create так, чтобы вернуть переданный турнир и
// сохранить его для проверок.
func captureCreate(repo *tournaments_repo.MockTournamentsRepo, captured *tournaments.TournamentModel) {
	repo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
		func(t tournaments.TournamentModel, _ tournaments.TournamentProvider) (tournaments.TournamentModel, error) {
			*captured = t
			return t, nil
		},
	)
}

func TestTournamentRegistrar_Create(t *testing.T) {
	t.Run("explicit dates take precedence over provider", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)

		// провайдер умеет отдавать даты, но они должны быть проигнорированы
		provider := fakePeriodProvider{
			start: time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC),
			end:   time.Date(2099, 6, 1, 0, 0, 0, 0, time.UTC),
		}
		registrar, repo := newRegistrar(ctrl, provider)

		var captured tournaments.TournamentModel
		captureCreate(repo, &captured)

		_, err := registrar.Create(CreateTournamentInput{
			Name: "T", ProviderName: "API_BASKETBALL", Season: "2025",
			StartAt: &start, EndAt: &end,
		})

		require.NoError(t, err)
		assert.Equal(t, start, captured.StartAt)
		assert.Equal(t, end, captured.EndAt)
	})

	t.Run("uses provider dates when no explicit dates given", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pStart := time.Date(2025, 9, 30, 0, 0, 0, 0, time.UTC)
		pEnd := time.Date(2026, 5, 24, 0, 0, 0, 0, time.UTC)
		provider := fakePeriodProvider{start: pStart, end: pEnd}
		registrar, repo := newRegistrar(ctrl, provider)

		var captured tournaments.TournamentModel
		captureCreate(repo, &captured)

		_, err := registrar.Create(CreateTournamentInput{
			Name: "T", ProviderName: "API_BASKETBALL", Season: "2025",
		})

		require.NoError(t, err)
		assert.Equal(t, pStart, captured.StartAt)
		assert.Equal(t, pEnd, captured.EndAt)
	})

	t.Run("leaves dates empty when provider can't report period", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		registrar, repo := newRegistrar(ctrl, fakeStatsProvider{})

		var captured tournaments.TournamentModel
		captureCreate(repo, &captured)

		_, err := registrar.Create(CreateTournamentInput{
			Name: "T", ProviderName: "INFOBASKET", Season: "2025",
		})

		require.NoError(t, err)
		assert.True(t, captured.StartAt.IsZero())
		assert.True(t, captured.EndAt.IsZero())
	})

	t.Run("leaves dates empty when provider returns an error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		provider := fakePeriodProvider{err: errors.New("api down")}
		registrar, repo := newRegistrar(ctrl, provider)

		var captured tournaments.TournamentModel
		captureCreate(repo, &captured)

		_, err := registrar.Create(CreateTournamentInput{
			Name: "T", ProviderName: "API_BASKETBALL", Season: "2025",
		})

		require.NoError(t, err)
		assert.True(t, captured.StartAt.IsZero())
		assert.True(t, captured.EndAt.IsZero())
	})

	t.Run("returns error and skips persistence when provider config is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// репозитории без ожиданий: любой их вызов провалит тест
		leaguesRepo := leagues_repo.NewMockLeaguesRepo(ctrl)
		tournamentsRepo := tournaments_repo.NewMockTournamentsRepo(ctrl)
		factory := func(_ string, _ *string, _ *map[string]interface{}) (ports.StatsProvider, error) {
			return nil, errors.New("external id must be set")
		}
		registrar := NewTournamentRegistrar(leaguesRepo, tournamentsRepo, factory)

		_, err := registrar.Create(CreateTournamentInput{Name: "T", ProviderName: "API_BASKETBALL"})

		require.Error(t, err)
	})
}
