package service

import (
	"IMP/app/internal/core/games"
	"IMP/app/internal/core/players"
	"IMP/app/internal/core/teams"
	"IMP/app/internal/core/tournaments"
	"IMP/app/internal/ports"
	"fmt"
	"reflect"
	"strconv"
	"time"

	composite_logger "github.com/Consolushka/golang.composite_logger/pkg"
)

type TournamentProcessorInterface interface {
	ProcessByPeriod(from, to time.Time) (int, error)
}

type TournamentProcessor struct {
	tournament         tournaments.TournamentModel
	statsProvider      ports.StatsProvider
	persistenceService PersistenceServiceInterface
	playersRepo        ports.PlayersRepo
	teamsRepo          ports.TeamsRepo
	gamesRepo          ports.GamesRepo
	aggregator         ports.Aggregator
	teamIdScope        string
}

func NewTournamentProcessor(statsProvider ports.StatsProvider, serviceInterface PersistenceServiceInterface, playersRepo ports.PlayersRepo, teamsRepo ports.TeamsRepo, gamesRepo ports.GamesRepo, aggregator ports.Aggregator, tournament tournaments.TournamentModel) *TournamentProcessor {
	// Обычно id команд уникальны в пределах провайдера; провайдер с несколькими базами сообщает свою сам.
	teamIdScope := tournament.Provider.ProviderName
	if scoped, ok := statsProvider.(ports.TeamIdScopeProvider); ok {
		teamIdScope = scoped.TeamIdScope()
	}

	return &TournamentProcessor{
		tournament:         tournament,
		statsProvider:      statsProvider,
		persistenceService: serviceInterface,
		playersRepo:        playersRepo,
		teamsRepo:          teamsRepo,
		gamesRepo:          gamesRepo,
		aggregator:         aggregator,
		teamIdScope:        teamIdScope,
	}
}

func (t TournamentProcessor) ProcessByPeriod(from, to time.Time) (int, error) {
	gameEntities, err := t.statsProvider.GetGamesStatsByPeriod(from, to)
	if err != nil {
		return 0, fmt.Errorf("GetGamesStatsByPeriod with %v, %v from %s returned error: %w", from, to, reflect.TypeOf(t.statsProvider), err)
	}
	savedGames := make([]string, 0, len(gameEntities))
	outOfPeriodReported := false

	if len(gameEntities) > 0 {
		composite_logger.Info("Start processing tournament games", map[string]interface{}{
			"tournamentId": t.tournament.Id,
		})
	}

	for _, gameEntity := range gameEntities {
		gameEntity.GameModel.TournamentId = t.tournament.Id

		// Команды ищем до проверки существования игры: название игры — часть её ключа
		// и собирается из алиасов команд, лежащих в базе.
		home, err := t.resolveTeam(gameEntity.HomeTeamExternalId, gameEntity.HomeTeamStat.TeamModel)
		if err != nil {
			return 0, err
		}
		away, err := t.resolveTeam(gameEntity.AwayTeamExternalId, gameEntity.AwayTeamStat.TeamModel)
		if err != nil {
			return 0, err
		}
		gameEntity.AssignTeams(home, away)

		isExists, err := t.gamesRepo.Exists(gameEntity.GameModel)
		if err != nil {
			composite_logger.Error("Failed to check whether game exists", map[string]interface{}{
				"tournamentId": gameEntity.GameModel.TournamentId,
				"title":        gameEntity.GameModel.Title,
				"scheduledAt":  gameEntity.GameModel.ScheduledAt,
				"error":        err,
			})
			return 0, fmt.Errorf("exists with %v from %s returned error: %w", gameEntity.GameModel, reflect.TypeOf(t.gamesRepo), err)
		}
		if isExists {
			composite_logger.Info("Game already exists. Skip game processing", map[string]interface{}{
				"gameModel": gameEntity.GameModel,
			})
			continue
		}

		gameEntity, err = t.statsProvider.EnrichGameStats(gameEntity)
		if err != nil {
			composite_logger.Warn("Couldn't enrich game stats", map[string]interface{}{
				"gameModel": gameEntity.GameModel,
				"error":     err,
			})
			continue
		}

		// DISCOVERY & INGESTION for players
		t.discoverAndIngestPlayers(&gameEntity)

		err = t.persistenceService.SaveGame(gameEntity)
		if err != nil {
			composite_logger.Error("t.persistenceService.SaveGame returned error", map[string]interface{}{
				"error":      err,
				"gameEntity": gameEntity,
			})
			continue
		}

		// Проверяем только что сохранённые игры: опрос каждый раз пересматривает прошлые сутки,
		// и проверка уже известных игр повторяла бы одно и то же предупреждение.
		if !outOfPeriodReported && !t.tournament.IsWithinPeriod(gameEntity.GameModel.ScheduledAt) {
			composite_logger.Warn("Imported game is scheduled outside tournament dates", map[string]interface{}{
				"tournamentId":   t.tournament.Id,
				"tournamentName": t.tournament.Name,
				"startAt":        t.tournament.StartAt,
				"endAt":          t.tournament.EndAt,
				"title":          gameEntity.GameModel.Title,
				"scheduledAt":    gameEntity.GameModel.ScheduledAt,
			})
			outOfPeriodReported = true
		}

		err = t.aggregator.NotifyGameImported(t.tournament.Id)
		if err != nil {
			// We only log here, not failing the whole process
			composite_logger.Error("t.aggregator.NotifyGameImported returned error", map[string]interface{}{
				"error":        err,
				"tournamentId": t.tournament.Id,
			})
		}

		savedGames = append(savedGames, formatSavedGameLog(gameEntity))
	}

	if len(gameEntities) > 0 {
		composite_logger.Info("Finished processing tournament games", map[string]interface{}{
			"tournamentId": t.tournament.Id,
			"savedCount":   len(savedGames),
			"savedGames":   savedGames,
		})
	}

	return len(savedGames), nil
}

// resolveTeam находит команду в базе: сначала по id провайдера, затем по имени.
// Найденную по имени команду связывает с id, чтобы следующее переименование у провайдера её не раздвоило.
func (t TournamentProcessor) resolveTeam(externalId int, candidate teams.TeamModel) (teams.TeamModel, error) {
	team, err := t.lookupTeam(externalId, candidate)
	if err != nil {
		composite_logger.Error("Failed to resolve team", map[string]interface{}{
			"tournamentId": t.tournament.Id,
			"scope":        t.teamIdScope,
			"externalId":   externalId,
			"name":         candidate.Name,
			"homeTown":     candidate.HomeTown,
			"error":        err,
		})
		return teams.TeamModel{}, fmt.Errorf("resolve team %q (scope %s, external id %d) returned error: %w", candidate.Name, t.teamIdScope, externalId, err)
	}
	return team, nil
}

func (t TournamentProcessor) lookupTeam(externalId int, candidate teams.TeamModel) (teams.TeamModel, error) {
	if externalId != 0 {
		team, found, err := t.teamsRepo.FindByExternalId(t.teamIdScope, externalId)
		if err != nil {
			return teams.TeamModel{}, fmt.Errorf("FindByExternalId returned error: %w", err)
		}
		if found {
			return team, nil
		}
	}

	team, created, err := t.teamsRepo.FirstOrCreate(candidate)
	if err != nil {
		return teams.TeamModel{}, fmt.Errorf("FirstOrCreate returned error: %w", err)
	}
	if created {
		composite_logger.Warn("New team created", map[string]interface{}{
			"teamId":       team.Id,
			"name":         team.Name,
			"homeTown":     team.HomeTown,
			"alias":        team.Alias,
			"scope":        t.teamIdScope,
			"externalId":   externalId,
			"tournamentId": t.tournament.Id,
		})
	}

	// Провайдер не дал id — связывать нечего.
	if externalId == 0 {
		return team, nil
	}

	err = t.teamsRepo.LinkExternalId(t.teamIdScope, externalId, team.Id)
	if err != nil {
		return teams.TeamModel{}, fmt.Errorf("LinkExternalId with team %d returned error: %w", team.Id, err)
	}

	// Перечитываем: связь мог первым записать другой воркер, и источник истины — то, что лежит в базе.
	linked, found, err := t.teamsRepo.FindByExternalId(t.teamIdScope, externalId)
	if err != nil {
		return teams.TeamModel{}, fmt.Errorf("FindByExternalId after link returned error: %w", err)
	}
	if !found {
		return teams.TeamModel{}, fmt.Errorf("team link (scope %s, external id %d) not found after insert", t.teamIdScope, externalId)
	}
	if !created {
		composite_logger.Info("Team linked to provider id", map[string]interface{}{
			"teamId":     linked.Id,
			"scope":      t.teamIdScope,
			"externalId": externalId,
		})
	}

	return linked, nil
}

func (t TournamentProcessor) discoverAndIngestPlayers(gameEntity *games.GameStatEntity) {
	homePlayers := gameEntity.HomeTeamStat.PlayerStats
	awayPlayers := gameEntity.AwayTeamStat.PlayerStats

	for i := range homePlayers {
		t.ensurePlayerBio(&homePlayers[i])
	}

	for i := range awayPlayers {
		t.ensurePlayerBio(&awayPlayers[i])
	}
}

func (t TournamentProcessor) ensurePlayerBio(playerStat *players.PlayerStatisticEntity) {
	// DISCOVERY: Check if player exists in our DB
	playersByFullName, err := t.playersRepo.ListByFullName(playerStat.PlayerModel.FullName)
	if err != nil {
		composite_logger.Error("Failed to search players by full name", map[string]interface{}{
			"playerFullName": playerStat.PlayerModel.FullName,
			"error":          err,
		})
		return
	}

	// If player exists, we already have their bio (likely)
	if len(playersByFullName) == 1 {
		return
	}

	// INGESTION: If player is unknown or missing critical data, fetch from provider
	if playerStat.PlayerModel.FullName == "" || time.Time.IsZero(playerStat.PlayerModel.BirthDate) {
		playerBio, err := t.statsProvider.GetPlayerBio(playerStat.PlayerExternalId)
		if err != nil {
			composite_logger.Warn("error while fetching player bio", map[string]interface{}{
				"playerId": playerStat.PlayerExternalId,
				"err":      err,
			})
		} else {
			playerStat.PlayerModel.FullName = playerBio.FullName
			playerStat.PlayerModel.BirthDate = playerBio.BirthDate
		}
	}
}

func formatSavedGameLog(gameEntity games.GameStatEntity) string {
	return gameEntity.GameModel.Title + " " +
		formatScore(gameEntity.HomeTeamStat.GameTeamStatModel.Score) + ":" +
		formatScore(gameEntity.AwayTeamStat.GameTeamStatModel.Score)
}

func formatScore(score int) string {
	return strconv.Itoa(score)
}
