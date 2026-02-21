package service

// todo: tests
import (
	"IMP/app/internal/core/tournament_poll_logs"
	"IMP/app/internal/core/tournaments"
	"IMP/app/internal/ports"
	"IMP/app/internal/service/providers"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"sync"
	"time"

	compositelogger "github.com/Consolushka/golang.composite_logger/pkg"
)

type TournamentsOrchestrator struct {
	persistenceService PersistenceServiceInterface
	tournamentsRepo    ports.TournamentsRepo
	playersRepo        ports.PlayersRepo
	gamesRepo          ports.GamesRepo
	pollLogRepo        ports.TournamentPollLogsRepo
}

func NewTournamentsOrchestrator(
	persistenceService PersistenceServiceInterface,
	tournamentsRepo ports.TournamentsRepo,
	playersRepo ports.PlayersRepo,
	gamesRepo ports.GamesRepo,
	pollLogRepo ports.TournamentPollLogsRepo,
) *TournamentsOrchestrator {
	return &TournamentsOrchestrator{
		persistenceService: persistenceService,
		tournamentsRepo:    tournamentsRepo,
		playersRepo:        playersRepo,
		gamesRepo:          gamesRepo,
		pollLogRepo:        pollLogRepo,
	}
}

// ProcessAll
// Fetches all active tournaments from repository and processes games for the given period
func (t TournamentsOrchestrator) ProcessAll(from, to time.Time) error {
	activeTournaments, err := t.tournamentsRepo.ListActive()

	if err != nil {
		return fmt.Errorf("ListActive from %s returned error: %w", reflect.TypeOf(t.tournamentsRepo), err)
	}

	t.processTournamentsByPeriod(activeTournaments, from, to)

	return nil
}

// ProcessTournament
// Processes a single tournament for a given period and records the result in poll logs
func (t TournamentsOrchestrator) ProcessTournament(tournament tournaments.TournamentModel, from, to time.Time) error {
	pollStartAt := time.Now()

	var params *map[string]interface{}
	if len(tournament.Provider.Params) == 0 {
		params = nil
	} else {
		err := json.Unmarshal(tournament.Provider.Params, &params)
		if err != nil {
			return fmt.Errorf("error while unmarshalling provider params for tournament %d: %w", tournament.Id, err)
		}
	}

	statsProvider, err := providers.NewProvider(tournament.Provider.ProviderName, tournament.Provider.ExternalId, params)
	if err != nil {
		return fmt.Errorf("error while creating stats provider for tournament %d: %w", tournament.Id, err)
	}

	processor := NewTournamentProcessor(statsProvider, t.persistenceService, t.playersRepo, t.gamesRepo, tournament.Id)

	savedGamesCount, processErr := processor.ProcessByPeriod(from, to)

	// LOGGING: Record results to database
	pollEndAt := time.Now()
	status := tournament_poll_logs.StatusSuccess
	var errMsg *string
	if processErr != nil {
		status = tournament_poll_logs.StatusError
		s := processErr.Error()
		errMsg = &s
	}

	_, logErr := t.pollLogRepo.Create(tournament_poll_logs.TournamentPollLogModel{
		TournamentId:    tournament.Id,
		PollStartAt:     pollStartAt,
		PollEndAt:       &pollEndAt,
		IntervalStart:   from,
		IntervalEnd:     to,
		SavedGamesCount: savedGamesCount,
		Status:          status,
		ErrorMessage:    errMsg,
	})

	if logErr != nil {
		compositelogger.Error("Couldn't create tournament poll log", map[string]interface{}{
			"tournamentId": tournament.Id,
			"error":        logErr,
		})
	}

	return processErr
}

func (t TournamentsOrchestrator) processTournamentsByPeriod(activeTournaments []tournaments.TournamentModel, from, to time.Time) {
	var tournamentsGroup sync.WaitGroup
	tournamentsGroup.Add(len(activeTournaments))

	compositelogger.Info("Start processing tournaments", map[string]interface{}{
		"count":   len(activeTournaments),
		"aliases": getTournamentsAliases(activeTournaments),
	})

	for _, tournament := range activeTournaments {
		go func(tournament tournaments.TournamentModel) {
			defer tournamentsGroup.Done()

			err := t.ProcessTournament(tournament, from, to)
			if err != nil {
				compositelogger.Error("Error while processing tournament", map[string]interface{}{
					"error":      err,
					"tournament": tournament,
				})
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
