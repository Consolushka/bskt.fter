package statistics

import (
	"IMP/app/internal/domain"
	"errors"
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

func NewLeagueProvider(league *domain.League) (StatsProvider, error) {
	switch strings.ToUpper(league.AliasEn) {
	case strings.ToUpper(domain.NBAAlias):
		return newNbaProvider(league), nil
	case strings.ToUpper(domain.MLBLAlias):
		return newMlblProvider(league), nil
	default:
		return nil, errors.New("There is no provider for league: " + strings.ToUpper(domain.MLBLAlias))
	}
}
