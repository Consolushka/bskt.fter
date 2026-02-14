package service

import (
	"IMP/app/internal/core/tournaments"
	"IMP/app/internal/ports"
	"IMP/app/internal/service/providers"
	"IMP/app/pkg/logger"
	"encoding/json"
	"sort"
	"sync"
	"time"
)

type TournamentsOrchestrator struct {
	persistenceService PersistenceServiceInterface
	tournamentsRepo    ports.TournamentsRepo
	playersRepo        ports.PlayersRepo
	gamesRepo          ports.GamesRepo
}

func NewTournamentsOrchestrator(persistenceService PersistenceServiceInterface, tournamentsRepo ports.TournamentsRepo, playersRepo ports.PlayersRepo, gamesRepo ports.GamesRepo) *TournamentsOrchestrator {
	return &TournamentsOrchestrator{
		persistenceService: persistenceService,
		tournamentsRepo:    tournamentsRepo,
		playersRepo:        playersRepo,
		gamesRepo:          gamesRepo,
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
		"VTB",
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
		"count":   len(activeTournaments),
		"aliases": getTournamentsAliases(activeTournaments),
	})

	for _, tournament := range activeTournaments {
		go func(tournament tournaments.TournamentModel) {
			defer tournamentsGroup.Done()

			var params *map[string]interface{}

			if len(tournament.Provider.Params) == 0 {
				params = nil
			} else {
				err := json.Unmarshal(tournament.Provider.Params, &params)
				if err != nil {
					logger.Error("Error while unmarshalling provider params", map[string]interface{}{
						"error":      err,
						"tournament": tournament,
					})
					return
				}
			}

			statsProvider, err := providers.NewProvider(tournament.Provider.ProviderName, tournament.Provider.ExternalId, params)
			if err != nil {
				logger.Error("Error while creating stats provider", map[string]interface{}{
					"error":      err,
					"tournament": tournament,
				})
				return
			}

			processor := NewTournamentProcessor(statsProvider, t.persistenceService, t.playersRepo, t.gamesRepo, tournament.Id)

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

func getTournamentsAliases(activeTournaments []tournaments.TournamentModel) []string {
	uniqueAliases := make(map[string]struct{})

	for _, tournament := range activeTournaments {
		if tournament.League.Alias == "" {
			continue
		}
		uniqueAliases[tournament.League.Alias] = struct{}{}
	}

	aliases := make([]string, 0, len(uniqueAliases))
	for alias := range uniqueAliases {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)

	return aliases
}
