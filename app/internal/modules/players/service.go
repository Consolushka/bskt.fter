package players

import (
	"IMP/app/internal/modules/players/domain"
	"IMP/app/internal/modules/players/domain/models"
)

type Service struct {
	repository *domain.Repository
}

func NewService() *Service {
	return &Service{
		repository: domain.NewRepository(),
	}
}

func (s *Service) GetPlayerByFullName(fullName string) ([]models.Player, error) {
	return s.repository.ListByFullName(fullName)
}
