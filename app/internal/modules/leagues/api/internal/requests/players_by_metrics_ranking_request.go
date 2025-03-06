package requests

import (
	"IMP/app/internal/abstract/custom_request"
	"IMP/app/internal/base/components/request_components"
	"strconv"
)

type PlayersByMetricsRankingRequest struct {
	custom_request.BaseRequest
	request_components.HasIdPathParam
	request_components.HasPersQueryParam

	limit             *int
	minGamesPlayed    int
	minMinutesPerGame int
	avgMinutesPerGame int
}

func (s *PlayersByMetricsRankingRequest) Validators() []func(storage custom_request.CustomRequestStorage) error {
	var parentValidators []func(storage custom_request.CustomRequestStorage) error

	currenValidators := []func(storage custom_request.CustomRequestStorage) error{
		s.parseLimit,
		s.parseMinGamesPlayed,
		s.parseMinMinutesPerGame,
		s.HasIdPathParam.Validators()[0],
	}

	for _, validator := range s.HasIdPathParam.Validators() {
		parentValidators = append(parentValidators, validator)
	}

	for _, validator := range s.HasPersQueryParam.Validators() {
		parentValidators = append(parentValidators, validator)
	}

	for _, validator := range currenValidators {
		parentValidators = append(parentValidators, validator)
	}

	return parentValidators
}

func (s *PlayersByMetricsRankingRequest) parseLimit(storage custom_request.CustomRequestStorage) error {
	limit := storage.GetQueryParam("limit")

	if limit == "" {
		s.limit = nil
		return nil
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return err
	}
	s.limit = &limitInt

	return nil
}

func (s *PlayersByMetricsRankingRequest) Limit() *int {
	return s.limit
}

func (s *PlayersByMetricsRankingRequest) parseMinGamesPlayed(storage custom_request.CustomRequestStorage) error {
	minGamesPlayed := storage.GetQueryParam("min_games")

	if minGamesPlayed == "" {
		s.minGamesPlayed = 0
		return nil
	}

	minGamesPlayedInt, err := strconv.Atoi(minGamesPlayed)
	if err != nil {
		return err
	}
	s.minGamesPlayed = minGamesPlayedInt

	return err
}

func (s *PlayersByMetricsRankingRequest) MinGamesPlayed() int {
	return s.minGamesPlayed
}

func (s *PlayersByMetricsRankingRequest) parseMinMinutesPerGame(storage custom_request.CustomRequestStorage) error {
	minMinutesPerGame := storage.GetQueryParam("min_minutes")

	if minMinutesPerGame == "" {
		s.minMinutesPerGame = 0
		return nil
	}

	minMinutesPerGameInt, err := strconv.Atoi(minMinutesPerGame)
	if err != nil {
		return err
	}
	s.minMinutesPerGame = minMinutesPerGameInt

	return err
}

func (s *PlayersByMetricsRankingRequest) MinMinutesPerGame() int {
	return s.minMinutesPerGame
}

func (s *PlayersByMetricsRankingRequest) parseAvgMinutesPerGame(storage custom_request.CustomRequestStorage) error {
	avgMinutesPerGame := storage.GetQueryParam("avg_minutes")

	if avgMinutesPerGame == "" {
		s.avgMinutesPerGame = 0
		return nil
	}

	avgMinutesPerGameInt, err := strconv.Atoi(avgMinutesPerGame)
	if err != nil {
		return err
	}
	s.avgMinutesPerGame = avgMinutesPerGameInt
	return nil
}

func (s *PlayersByMetricsRankingRequest) AvgMinutesPerGame() int {
	return s.avgMinutesPerGame
}
