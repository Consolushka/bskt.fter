package teams_repo

import (
	"IMP/app/internal/core/teams"
	"IMP/app/internal/ports"
	"fmt"
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
	existing, err := g.findByName(model)
	if err != nil {
		return model, false, err
	}
	if existing.Id != 0 {
		return existing, false, nil
	}

	if model.Alias == "" || utf8.RuneCountInString(model.Alias) > maxAliasLen {
		model.Alias = model.AutoGenerateAlias()
	}

	// Между поиском и вставкой ту же команду мог создать другой воркер.
	// Тогда вставка молча ничего не делает, и мы берём команду, которая уже лежит в базе.
	result := g.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model)
	if result.Error != nil {
		return model, false, result.Error
	}
	if result.RowsAffected == 0 {
		existing, err = g.findByName(model)
		if err != nil {
			return model, false, err
		}
		if existing.Id == 0 {
			return model, false, fmt.Errorf("team %q was neither created nor found", model.Name)
		}
		return existing, false, nil
	}

	return model, true, nil
}

// Пустой город в условие не входит: часть провайдеров его не отдаёт,
// и такая команда должна находиться по одному имени, как и раньше.
func (g Gorm) findByName(model teams.TeamModel) (teams.TeamModel, error) {
	query := g.db.Where("name = ?", model.Name)
	if model.HomeTown != "" {
		query = query.Where("home_town = ?", model.HomeTown)
	}

	var existing teams.TeamModel
	err := query.Order("id").Limit(1).Find(&existing).Error

	return existing, err
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
