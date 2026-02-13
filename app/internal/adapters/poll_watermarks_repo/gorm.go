package poll_watermarks_repo

import (
	"IMP/app/internal/core/poll_watermarks"
	"IMP/app/internal/ports"
	"errors"

	"gorm.io/gorm"
)

type Gorm struct {
	db *gorm.DB
}

func (g Gorm) FirstOrCreate(model poll_watermarks.PollWatermarkModel) (poll_watermarks.PollWatermarkModel, error) {
	var foundModel poll_watermarks.PollWatermarkModel

	tx := g.db.Where("task_type = ?", model.TaskType).First(&foundModel)
	if tx.Error == nil {
		return foundModel, nil
	}
	if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return poll_watermarks.PollWatermarkModel{}, tx.Error
	}

	model.LastSuccessfulPollAt = model.LastSuccessfulPollAt.UTC()

	tx = g.db.Create(&model)
	return model, tx.Error
}

func (g Gorm) Update(model poll_watermarks.PollWatermarkModel) (poll_watermarks.PollWatermarkModel, error) {
	model.LastSuccessfulPollAt = model.LastSuccessfulPollAt.UTC()

	tx := g.db.Save(&model)

	return model, tx.Error
}

func NewGormRepo(db *gorm.DB) ports.PollWatermarkRepo {
	return Gorm{db: db}
}
