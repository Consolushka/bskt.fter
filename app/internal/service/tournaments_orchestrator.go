package service

import (
	"IMP/app/internal/core/tournaments"
	"IMP/app/internal/ports"
	"IMP/app/pkg/logger"
	"sync"
	"time"
)

// todo: refactor to create table with schedule time etc.

type TournamentsOrchestrator struct {
	persistenceService PersistenceServiceInterface
	tournamentsRepo    ports.TournamentsRepo
}

func NewTournamentsOrchestrator(persistenceService PersistenceServiceInterface, tournamentsRepo ports.TournamentsRepo) *TournamentsOrchestrator {
	return &TournamentsOrchestrator{
		persistenceService: persistenceService,
		tournamentsRepo:    tournamentsRepo,
	}
}

// ProcessUrgentEuropeanTournaments
// Fetches and processing tournaments, which need to be processed as soon as possible
func (t TournamentsOrchestrator) ProcessUrgentEuropeanTournaments() error {
	leaguesAliases := []string{
		"MLBL",
		"FONBETSL",
	}

	processingTournaments, err := t.tournamentsRepo.ListTournamentsByLeagueAliases(leaguesAliases)
	if err != nil {
		return err
	}

	t.processTournamentsByDate(processingTournaments, time.Now())

	return nil
}

// ProcessAmericanTournaments
// Fetches and processing tournaments, which played at night by MSK
func (t TournamentsOrchestrator) ProcessAmericanTournaments() error {
	leaguesAliases := []string{
		"NBA",
	}

	processingTournaments, err := t.tournamentsRepo.ListTournamentsByLeagueAliases(leaguesAliases)
	if err != nil {
		return err
	}

	t.processTournamentsByDate(processingTournaments, time.Now().Add(-time.Hour*24))

	return nil
}

// ProcessNotUrgentEuropeanTournaments
// Fetches and processing tournaments, which need to process, but not urgent
func (t TournamentsOrchestrator) ProcessNotUrgentEuropeanTournaments() error {
	leaguesAliases := []string{
		"UBA",
	}

	processingTournaments, err := t.tournamentsRepo.ListTournamentsByLeagueAliases(leaguesAliases)
	if err != nil {
		return err
	}

	t.processTournamentsByDate(processingTournaments, time.Now().Add(-time.Hour*24))

	return nil
}

// ProcessAllTournamentsToday
// Fetches all active tournaments from repository and processes today games
func (t TournamentsOrchestrator) ProcessAllTournamentsToday() error {
	activeTournaments, err := t.tournamentsRepo.ListActiveTournaments()

	if err != nil {
		return err
	}

	t.processTournamentsByDate(activeTournaments, time.Now().Add(-time.Hour*24))

	return nil
}

func (t TournamentsOrchestrator) processTournamentsByDate(activeTournaments []tournaments.TournamentModel, date time.Time) {
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

			processor := NewTournamentProcessor(statsProvider, t.persistenceService, tournament.Id)

			err = processor.Process(date)
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
