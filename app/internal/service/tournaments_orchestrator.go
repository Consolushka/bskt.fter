package service

import (
	"IMP/app/internal/core/tournaments"
	"IMP/app/internal/ports"
	"IMP/app/pkg/logger"
	"sync"
	"time"
)

type TournamentsOrchestrator struct {
	persistenceService PersistenceServiceInterface
	tournamentsRepo    ports.TournamentsRepo
	playersRepo        ports.PlayersRepo
}

func NewTournamentsOrchestrator(persistenceService PersistenceServiceInterface, tournamentsRepo ports.TournamentsRepo, playersRepo ports.PlayersRepo) *TournamentsOrchestrator {
	return &TournamentsOrchestrator{
		persistenceService: persistenceService,
		tournamentsRepo:    tournamentsRepo,
		playersRepo:        playersRepo,
	}
}

// ProcessUrgentEuropeanTournaments
// Fetches and processing tournaments, which need to be processed as soon as possible
func (t TournamentsOrchestrator) ProcessUrgentEuropeanTournaments(from, to time.Time) error {
	leaguesAliases := []string{
		"MLBL",
		"FONBETSL",
	}

	processingTournaments, err := t.tournamentsRepo.ListTournamentsByLeagueAliases(leaguesAliases)
	if err != nil {
		return err
	}

	t.processTournamentsByPeriod(processingTournaments, from, to)

	return nil
}

// ProcessAmericanTournaments
// Fetches and processing tournaments, which played at night by MSK
func (t TournamentsOrchestrator) ProcessAmericanTournaments(from, to time.Time) error {
	leaguesAliases := []string{
		"NBA",
	}

	processingTournaments, err := t.tournamentsRepo.ListTournamentsByLeagueAliases(leaguesAliases)
	if err != nil {
		return err
	}

	t.processTournamentsByPeriod(processingTournaments, from, to)

	return nil
}

// ProcessNotUrgentEuropeanTournaments
// Fetches and processing tournaments, which need to process, but not urgent
func (t TournamentsOrchestrator) ProcessNotUrgentEuropeanTournaments(from, to time.Time) error {
	leaguesAliases := []string{
		"UBA",
	}

	processingTournaments, err := t.tournamentsRepo.ListTournamentsByLeagueAliases(leaguesAliases)
	if err != nil {
		return err
	}

	t.processTournamentsByPeriod(processingTournaments, from, to)

	return nil
}

// ProcessAllTournamentsToday
// Fetches all active tournaments from repository and processes today games
func (t TournamentsOrchestrator) ProcessAllTournamentsToday(from, to time.Time) error {
	activeTournaments, err := t.tournamentsRepo.ListActiveTournaments()

	if err != nil {
		return err
	}

	t.processTournamentsByPeriod(activeTournaments, from, to)

	return nil
}

func (t TournamentsOrchestrator) processTournamentsByPeriod(activeTournaments []tournaments.TournamentModel, from, to time.Time) {
	var tournamentsGroup sync.WaitGroup
	tournamentsGroup.Add(len(activeTournaments))

	logger.Info("Start processing tournaments", map[string]interface{}{
		"tournaments": activeTournaments,
	})

	for _, tournament := range activeTournaments {
		go func(tournament tournaments.TournamentModel) {
			defer tournamentsGroup.Done()

			statsProvider, err := NewTournamentStatsProvider(tournament)
			if err != nil {
				logger.Error("Error while creating stats provider", map[string]interface{}{
					"error":      err,
					"tournament": tournament,
				})
				return
			}

			processor := NewTournamentProcessor(statsProvider, t.persistenceService, t.playersRepo, tournament.Id)

			err = processor.ProcessByPeriod(from, to)
			if err != nil {
				logger.Error("Error while processing tournament games", map[string]interface{}{
					"error":      err,
					"tournament": tournament,
					"processor":  processor,
				})
				return
			}
		}(tournament)
	}

	tournamentsGroup.Wait()
}
