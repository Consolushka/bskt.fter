package statistics

import (
	"IMP/app/internal/domain"
	"strings"
	"time"
)

type StatsProvider interface {
	// GameBoxScore returns boxscore data from stats provider
	GameBoxScore(gameId string) (*GameBoxScoreDTO, error)
	// GamesByDate returns list of games for given date
	GamesByDate(date time.Time) ([]string, error)
	// GamesByTeam returns list of already played games for given team
	GamesByTeam(teamId string) ([]string, error)
}

func NewLeagueProvider(league *domain.League) StatsProvider {
	switch league.AliasEn {
	case strings.ToUpper(domain.NBAAlias):
		return newNbaProvider(league)
	case strings.ToUpper(domain.MLBLAlias):
		return newMlblProvider(league)
	default:
		return nil
	}
}
