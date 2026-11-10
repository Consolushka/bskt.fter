package teams_repo

import (
	"IMP/app/internal/core/teams"
	"IMP/app/internal/ports"
	"unicode/utf8"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ ports.TeamsRepo = (*Gorm)(nil)

type Gorm struct {
	db *gorm.DB
}

// Длина колонки teams.alias: сокращение провайдера длиннее неё дало бы ошибку базы на каждом опросе.
const maxAliasLen = 10

// Алиас задаётся только при создании команды: иначе провайдеры с разными сокращениями
// переписывали бы его друг за другом, а ручные исправления откатывались бы.
func (g Gorm) FirstOrCreate(model teams.TeamModel) (teams.TeamModel, bool, error) {
	// Пустой город в условие не входит: часть провайдеров его не отдаёт,
	// и такая команда должна находиться по одному имени, как и раньше.
	query := g.db.Where("name = ?", model.Name)
	if model.HomeTown != "" {
		query = query.Where("home_town = ?", model.HomeTown)
	}

	var existing teams.TeamModel
	if err := query.Order("id").Limit(1).Find(&existing).Error; err != nil {
		return model, false, err
	}
	if existing.Id != 0 {
		return existing, false, nil
	}

	if model.Alias == "" || utf8.RuneCountInString(model.Alias) > maxAliasLen {
		model.Alias = model.AutoGenerateAlias()
	}

	err := g.db.Create(&model).Error

	return model, err == nil, err
}

func (g Gorm) FindByExternalId(scope string, externalId int) (teams.TeamModel, bool, error) {
	var team teams.TeamModel
	err := g.db.
		Joins("JOIN team_external_ids ON team_external_ids.team_id = teams.id").
		Where("team_external_ids.scope = ? AND team_external_ids.external_id = ?", scope, externalId).
		Limit(1).
		Find(&team).Error

	return team, team.Id != 0, err
}

func (g Gorm) LinkExternalId(scope string, externalId int, teamId uint) error {
	return g.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&teams.TeamExternalIdModel{
		Scope:      scope,
		ExternalId: externalId,
		TeamId:     teamId,
	}).Error
}

func (g Gorm) FirstOrCreateStats(model teams.GameTeamStatModel) (teams.GameTeamStatModel, error) {
	tx := g.db.FirstOrCreate(&model, teams.GameTeamStatModel{
		GameId: model.GameId,
		TeamId: model.TeamId,
	})

	return model, tx.Error
}

func NewGormRepo(db *gorm.DB) Gorm {
	return Gorm{db: db}
}
