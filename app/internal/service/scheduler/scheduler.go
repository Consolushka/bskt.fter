package scheduler

import (
	"IMP/app/internal/adapters/games_repo"
	"IMP/app/internal/adapters/players_repo"
	"IMP/app/internal/adapters/teams_repo"
	"IMP/app/internal/adapters/tournament_poll_logs_repo"
	"IMP/app/internal/adapters/tournaments_repo"
	"IMP/app/internal/core/tournaments"
	"IMP/app/internal/infra/config"
	"IMP/app/internal/ports"
	"IMP/app/internal/service"
	"context"
	"sync"
	"time"

	composite_logger "github.com/Consolushka/golang.composite_logger/pkg"
	"gorm.io/gorm"
)

type Scheduler struct {
	pollInterval    time.Duration
	staggerInterval time.Duration
	refreshInterval time.Duration

	workers map[uint]context.CancelFunc
	mu      sync.Mutex

	orchestrator    *service.TournamentsOrchestrator
	pollLogRepo     ports.TournamentPollLogsRepo
	tournamentsRepo ports.TournamentsRepo
}

func NewScheduler(db *gorm.DB, cfg *config.Config) *Scheduler {
	tournamentsRepo := tournaments_repo.NewGormRepo(db)
	gamesRepo := games_repo.NewGormRepo(db)
	teamsRepo := teams_repo.NewGormRepo(db)
	playersRepo := players_repo.NewGormRepo(db)
	pollLogRepo := tournament_poll_logs_repo.NewGormRepo(db)

	persistenceService := service.NewPersistenceService(gamesRepo, teamsRepo, playersRepo)
	orchestrator := service.NewTournamentsOrchestrator(
		persistenceService,
		tournamentsRepo,
		playersRepo,
		gamesRepo,
		pollLogRepo,
		cfg.Providers,
	)

	return &Scheduler{
		pollInterval:    time.Duration(cfg.Scheduler.PollInterval) * time.Minute,
		staggerInterval: time.Duration(cfg.Scheduler.StaggerInterval) * time.Minute,
		refreshInterval: time.Duration(cfg.Scheduler.RefreshInterval) * time.Minute,
		workers:         make(map[uint]context.CancelFunc),
		orchestrator:    orchestrator,
		pollLogRepo:     pollLogRepo,
		tournamentsRepo: tournamentsRepo,
	}
}

func Handle(db *gorm.DB, cfg *config.Config) {
	s := NewScheduler(db, cfg)
	s.Run()
}

func (s *Scheduler) Run() {
	composite_logger.Info("Scheduler starting", map[string]interface{}{
		"pollInterval":    s.pollInterval.String(),
		"staggerInterval": s.staggerInterval.String(),
		"refreshInterval": s.refreshInterval.String(),
	})

	// Initial refresh
	s.refreshTournaments()

	// Periodic refresh
	ticker := time.NewTicker(s.refreshInterval)
	defer ticker.Stop()

	for range ticker.C {
		s.refreshTournaments()
	}
}

func (s *Scheduler) refreshTournaments() {
	activeTournaments, err := s.tournamentsRepo.ListActive()
	if err != nil {
		composite_logger.Error("Couldn't fetch active tournaments", map[string]interface{}{
			"error": err,
		})
		return
	}

	activeIds := make(map[uint]struct{})
	for _, t := range activeTournaments {
		activeIds[t.Id] = struct{}{}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. Stop workers for tournaments that are no longer active
	for id, cancel := range s.workers {
		if _, ok := activeIds[id]; !ok {
			composite_logger.Info("Stopping worker for inactive tournament", map[string]interface{}{
				"tournamentId": id,
			})
			cancel()
			delete(s.workers, id)
		}
	}

	// 2. Start workers for new tournaments
	newWorkersCount := 0
	for _, t := range activeTournaments {
		if _, ok := s.workers[t.Id]; !ok {
			ctx, cancel := context.WithCancel(context.Background())
			s.workers[t.Id] = cancel

			delay := time.Duration(newWorkersCount) * s.staggerInterval
			go s.runTournamentWorker(ctx, t, delay)
			newWorkersCount++
		}
	}

	if newWorkersCount > 0 {
		composite_logger.Info("Queued new workers", map[string]interface{}{
			"count":           newWorkersCount,
			"staggerInterval": s.staggerInterval.String(),
		})
	}
}

func (s *Scheduler) runTournamentWorker(ctx context.Context, tournament tournaments.TournamentModel, initialDelay time.Duration) {
	defer composite_logger.Recover(map[string]interface{}{
		"tournamentId": tournament.Id,
	})

	if initialDelay > 0 {
		select {
		case <-time.After(initialDelay):
		case <-ctx.Done():
			return
		}
	}

	composite_logger.Info("Worker started", map[string]interface{}{
		"tournamentId": tournament.Id,
		"interval":     s.pollInterval.String(),
	})

	// Immediate first run
	s.processTournament(tournament)

	// Periodic runs
	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.processTournament(tournament)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Scheduler) processTournament(tournament tournaments.TournamentModel) {
	now := time.Now().UTC()

	// 1. DISCOVERY: Get latest successful poll interval end
	latestLog, err := s.pollLogRepo.GetLatestSuccess(tournament.Id)
	var intervalStart time.Time
	if err != nil {
		// If no logs found, start from today
		intervalStart = toStartOfUTCDay(now)
	} else {
		intervalStart = latestLog.IntervalEnd
	}

	// 2. INGESTION: Run orchestration (it will handle internal poll logging)
	nextPollAt := now.Add(s.pollInterval)
	if err := s.orchestrator.ProcessTournament(tournament, intervalStart, now, &nextPollAt); err != nil {
		composite_logger.Error("Error while processing tournament games", map[string]interface{}{
			"tournamentId": tournament.Id,
			"error":        err,
		})
		return
	}

	composite_logger.Info("Tournament worker cycle finished", map[string]interface{}{
		"tournamentId": tournament.Id,
		"from":         intervalStart,
		"to":           now,
	})
}

func toStartOfUTCDay(value time.Time) time.Time {
	value = value.UTC()
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
}
