package service

import (
	"IMP/app/internal/core/leagues"
	"IMP/app/internal/core/tournaments"
	"IMP/app/internal/ports"
	"encoding/json"
	"fmt"
	"time"

	compositelogger "github.com/Consolushka/golang.composite_logger/pkg"
)

// ProviderFactory строит провайдера данных по его имени и конфигу.
// Инжектируется (а не вызывается напрямую), чтобы use-case можно было
// тестировать с подставным провайдером без реальных HTTP-клиентов.
type ProviderFactory func(providerName string, externalId *string, params *map[string]interface{}) (ports.StatsProvider, error)

// CreateTournamentInput — данные для регистрации турнира.
// StartAt/EndAt опциональны: если не заданы, даты пытаемся взять у провайдера.
type CreateTournamentInput struct {
	LeagueName  string
	LeagueAlias string
	Name        string

	ProviderName string
	ExternalId   *string
	Params       *map[string]interface{}
	// Season нужен провайдеру, чтобы выбрать нужный сезон при автозаполнении дат.
	Season string

	StartAt            *time.Time
	EndAt              *time.Time
	Tier               *int16
	RegulationDuration int
}

type TournamentRegistrar struct {
	leaguesRepo     ports.LeaguesRepo
	tournamentsRepo ports.TournamentsRepo
	newProvider     ProviderFactory
}

func NewTournamentRegistrar(
	leaguesRepo ports.LeaguesRepo,
	tournamentsRepo ports.TournamentsRepo,
	newProvider ProviderFactory,
) *TournamentRegistrar {
	return &TournamentRegistrar{
		leaguesRepo:     leaguesRepo,
		tournamentsRepo: tournamentsRepo,
		newProvider:     newProvider,
	}
}

// Create регистрирует турнир: проверяет конфиг провайдера (через фабрику),
// определяет даты сезона, гарантирует наличие лиги и атомарно сохраняет
// турнир с привязкой к провайдеру.
func (r *TournamentRegistrar) Create(input CreateTournamentInput) (tournaments.TournamentModel, error) {
	// Фабрика заодно валидирует, что для провайдера переданы нужные externalId/params.
	provider, err := r.newProvider(input.ProviderName, input.ExternalId, input.Params)
	if err != nil {
		return tournaments.TournamentModel{}, fmt.Errorf("create provider %q returned error: %w", input.ProviderName, err)
	}

	start, end := r.resolvePeriod(input, provider)

	league, err := r.leaguesRepo.FirstOrCreate(leagues.LeagueModel{
		Name:  input.LeagueName,
		Alias: input.LeagueAlias,
	})
	if err != nil {
		return tournaments.TournamentModel{}, fmt.Errorf("ensure league %q returned error: %w", input.LeagueName, err)
	}

	var paramsJSON []byte
	if input.Params != nil {
		paramsJSON, err = json.Marshal(*input.Params)
		if err != nil {
			return tournaments.TournamentModel{}, fmt.Errorf("marshal provider params returned error: %w", err)
		}
	}

	created, err := r.tournamentsRepo.Create(
		tournaments.TournamentModel{
			LeagueId:           league.Id,
			Name:               input.Name,
			Tier:               input.Tier,
			StartAt:            start,
			EndAt:              end,
			RegulationDuration: input.RegulationDuration,
		},
		tournaments.TournamentProvider{
			ProviderName: input.ProviderName,
			ExternalId:   input.ExternalId,
			Params:       paramsJSON,
		},
	)
	if err != nil {
		return tournaments.TournamentModel{}, fmt.Errorf("create tournament %q returned error: %w", input.Name, err)
	}

	return created, nil
}

// resolvePeriod определяет даты сезона по приоритету:
//  1. явные даты от пользователя;
//  2. даты от провайдера, если он умеет их сообщать;
//  3. иначе — предупреждение, даты остаются неизвестными (nil → NULL в БД).
//
// nil-даты сознательны: турнир с неизвестным окончанием считается активным
// (см. фильтр в TournamentsRepo.ListActive), а не "уже завершённым".
func (r *TournamentRegistrar) resolvePeriod(input CreateTournamentInput, provider ports.StatsProvider) (*time.Time, *time.Time) {
	if input.StartAt != nil && input.EndAt != nil {
		return input.StartAt, input.EndAt
	}

	periodProvider, ok := provider.(ports.TournamentPeriodProvider)
	if !ok {
		compositelogger.Warn("Provider can't report tournament period and no explicit dates given; dates left empty", map[string]interface{}{
			"provider": input.ProviderName,
		})
		return nil, nil
	}

	start, end, err := periodProvider.GetTournamentPeriod(input.Season)
	if err != nil {
		compositelogger.Warn("Couldn't fetch tournament period from provider; dates left empty", map[string]interface{}{
			"provider": input.ProviderName,
			"season":   input.Season,
			"error":    err,
		})
		return nil, nil
	}

	return &start, &end
}
