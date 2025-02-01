package teams

import (
	"IMP/app/log"
	"github.com/sirupsen/logrus"
)

type Service struct {
	repository *Repository

	logger *logrus.Logger
}

func NewService() *Service {
	return &Service{
		repository: NewRepository(),
		logger:     log.GetLogger(),
	}
}

func (s *Service) GetTeams() ([]Team, error) {
	var teams []Team

	tx := s.repository.dbConnection.Model(&Team{}).
		Find(&teams)

	return teams, tx.Error
}

func (s *Service) GetTeamById(id int) (Team, error) {
	var team Team

	tx := s.repository.dbConnection.Model(&Team{}).
		Where("id = ?", id).
		First(&team)

	return team, tx.Error
}
