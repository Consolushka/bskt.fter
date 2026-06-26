package tournaments_repo

import (
	"IMP/app/internal/core/tournaments"
	"IMP/app/internal/ports"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var _ ports.TournamentsRepo = (*Gorm)(nil)

type Gorm struct {
	db *gorm.DB
}

// Create пишет турнир и его привязку к провайдеру в одной транзакции:
// сначала турнир (чтобы получить сгенерированный id), затем provider-маппинг с этим id.
func (g Gorm) Create(tournament tournaments.TournamentModel, provider tournaments.TournamentProvider) (tournaments.TournamentModel, error) {
	err := g.db.Transaction(func(tx *gorm.DB) error {
		// Omit ассоциаций: лигу и провайдера сохраняем сами, не даём GORM трогать их автоматически.
		if err := tx.Omit("League", "Provider").Create(&tournament).Error; err != nil {
			return fmt.Errorf("create tournament returned error: %w", err)
		}

		provider.TournamentId = tournament.Id
		if err := tx.Create(&provider).Error; err != nil {
			return fmt.Errorf("create tournament provider returned error: %w", err)
		}

		return nil
	})
	if err != nil {
		return tournaments.TournamentModel{}, fmt.Errorf("create tournament %q transaction returned error: %w", tournament.Name, err)
	}

	tournament.Provider = provider
	return tournament, nil
}

func (g Gorm) ListByLeagueAliases(aliases []string) ([]tournaments.TournamentModel, error) {
	var models []tournaments.TournamentModel
	if len(aliases) == 0 {
		return models, nil
	}

	err := g.db.Model(&tournaments.TournamentModel{}).
		Joins("JOIN leagues ON tournaments.league_id = leagues.id").
		Preload("League").
		Preload("Provider").
		Where("leagues.alias IN ?", aliases).
		Find(&models).Error

	return models, err
}

func (g Gorm) Get(id uint) (tournaments.TournamentModel, error) {
	var model tournaments.TournamentModel
	err := g.db.Preload("League").Preload("Provider").First(&model, id).Error
	return model, err
}

func NewGormRepo(db *gorm.DB) Gorm {
	return Gorm{db: db}
}

func (g Gorm) ListActive() ([]tournaments.TournamentModel, error) {
	var models []tournaments.TournamentModel

	// Турниры, дата окончания которых уже прошла, опрашивать незачем:
	// новых игр там не будет, а лимиты внешних API расходовать на них не хочется.
	err := g.db.
		Preload("League").
		Preload("Provider").
		Where("end_at >= ?", time.Now().UTC()).
		Find(&models).Error

	return models, err
}
