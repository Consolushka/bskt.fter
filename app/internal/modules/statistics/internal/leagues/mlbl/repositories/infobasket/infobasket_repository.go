package infobasket

import (
	"IMP/app/internal/infrastructure/infobasket"
	models2 "IMP/app/internal/modules/statistics/models"
)

type Repository struct {
	client             *infobasket.Client
	persistenceService *persistenceService
}

func (i Repository) GameBoxScore(gameId string) (*models2.GameBoxScoreDTO, error) {
	panic("implement me")
	//var gameDto boxscore2.GameInfo
	//
	//boxscoreJson := i.client.BoxScore(gameId)
	//boxscoreRaw, _ := json.Marshal(boxscoreJson)
	//
	//err := json.Unmarshal(boxscoreRaw, &gameDto)
	//if err != nil {
	//	return nil, err
	//}
	//
	//i.persistenceService.saveGame(gameDto)
	//return gameDto.ToImpModel(), nil
}

func (i Repository) TodayGames() (string, []string, error) {
	panic("implement me")
}

func NewRepository() *Repository {
	return &Repository{
		client:             infobasket.NewInfobasketClient(),
		persistenceService: newPersistenceService(),
	}
}
