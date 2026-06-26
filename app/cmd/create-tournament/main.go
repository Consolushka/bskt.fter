package main

import (
	"IMP/app/database"
	"IMP/app/internal/adapters/leagues_repo"
	"IMP/app/internal/adapters/tournaments_repo"
	"IMP/app/internal/infra/config"
	"IMP/app/internal/infra/logger"
	"IMP/app/internal/ports"
	"IMP/app/internal/service"
	"IMP/app/internal/service/providers"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	compositelogger "github.com/Consolushka/golang.composite_logger/pkg"
	"github.com/joho/godotenv"
)

// Утилита для регистрации турнира из командной строки.
// Пример:
//
//	go run ./app/cmd/create-tournament \
//	  -name "Euroleague 2025-2026" -league-name Euroleague -league-alias euroleague \
//	  -provider API_BASKETBALL -external-id 120 -season 2025
//
// Если -start/-end не заданы, даты пытаемся взять у провайдера (если он умеет);
// иначе турнир создаётся без дат (см. предупреждение в логе).
func main() {
	time.Local = time.UTC

	name := flag.String("name", "", "tournament name (required)")
	leagueName := flag.String("league-name", "", "league name (required)")
	leagueAlias := flag.String("league-alias", "", "league alias (required)")
	provider := flag.String("provider", "", "provider: API_NBA | API_BASKETBALL | INFOBASKET | SPORTOTEKA (required)")
	externalID := flag.String("external-id", "", "provider external id (e.g. basketball league id, infobasket comp id)")
	paramsJSON := flag.String("params", "", `provider params as JSON, e.g. '{"leadHost":"reg","year":2024}'`)
	season := flag.String("season", "", "season the provider uses to look up dates, e.g. 2025")
	startStr := flag.String("start", "", "explicit start date YYYY-MM-DD (optional, overrides provider)")
	endStr := flag.String("end", "", "explicit end date YYYY-MM-DD (optional, overrides provider)")
	tier := flag.Int("tier", -1, "tournament tier (optional)")
	regulationDuration := flag.Int("regulation-duration", 0, "regulation duration in minutes")
	flag.Parse()

	if *name == "" || *leagueName == "" || *leagueAlias == "" || *provider == "" {
		fmt.Fprintln(os.Stderr, "name, league-name, league-alias and provider are required")
		flag.Usage()
		os.Exit(2)
	}

	//nolint:errcheck
	godotenv.Load()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	compositelogger.Init(logger.BuildSettings(cfg.Logger)...)

	db := database.OpenDbConnection(cfg.Database)

	factory := func(providerName string, externalId *string, params *map[string]interface{}) (ports.StatsProvider, error) {
		return providers.NewProvider(providerName, externalId, params, cfg.Providers)
	}

	registrar := service.NewTournamentRegistrar(
		leagues_repo.NewGormRepo(db),
		tournaments_repo.NewGormRepo(db),
		factory,
	)

	created, err := registrar.Create(service.CreateTournamentInput{
		LeagueName:         *leagueName,
		LeagueAlias:        *leagueAlias,
		Name:               *name,
		ProviderName:       *provider,
		ExternalId:         optionalString(*externalID),
		Params:             mustParseParams(*paramsJSON),
		Season:             *season,
		StartAt:            mustParseDate("start", *startStr),
		EndAt:              mustParseDate("end", *endStr),
		Tier:               optionalTier(*tier),
		RegulationDuration: *regulationDuration,
	})
	if err != nil {
		compositelogger.Error("Failed to create tournament", map[string]interface{}{"error": err})
		os.Exit(1)
	}

	fmt.Printf("Created tournament id=%d name=%q leagueId=%d start=%s end=%s\n",
		created.Id, created.Name, created.LeagueId, formatDate(created.StartAt), formatDate(created.EndAt))
}

func optionalString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func optionalTier(v int) *int16 {
	if v < 0 {
		return nil
	}
	t := int16(v)
	return &t
}

func mustParseParams(raw string) *map[string]interface{} {
	if raw == "" {
		return nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		fmt.Fprintf(os.Stderr, "invalid -params JSON: %v\n", err)
		os.Exit(2)
	}
	return &m
}

func mustParseDate(flagName, raw string) *time.Time {
	if raw == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid -%s date (want YYYY-MM-DD): %v\n", flagName, err)
		os.Exit(2)
	}
	return &t
}

func formatDate(t *time.Time) string {
	if t == nil {
		return "<none>"
	}
	return t.Format("2006-01-02")
}
