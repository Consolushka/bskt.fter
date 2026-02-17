package poll_watermarks_repo

import (
	"IMP/app/internal/core/poll_watermarks"
	"IMP/app/pkg/dbtest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type PollWatermarkRepoSuite struct {
	suite.Suite
	db   *gorm.DB
	tx   *gorm.DB
	repo Gorm
}

func (s *PollWatermarkRepoSuite) SetupSuite() {
	s.db = dbtest.Setup(s.T(), &poll_watermarks.PollWatermarkModel{})
}

func (s *PollWatermarkRepoSuite) SetupTest() {
	s.tx = s.db.Begin()
	s.repo = Gorm{db: s.tx}
}

func (s *PollWatermarkRepoSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *PollWatermarkRepoSuite) TestFirstOrCreate() {
	now := time.Now().Round(time.Second)
	watermark := poll_watermarks.PollWatermarkModel{
		TournamentId:         1,
		LastSuccessfulPollAt: now,
	}

	// 1. Сценарий: Создание новой записи
	res, err := s.repo.FirstOrCreate(watermark)
	s.Require().NoError(err)
	s.Equal(uint(1), res.TournamentId)
	s.True(now.Equal(res.LastSuccessfulPollAt))

	// 2. Сценарий: Возврат существующей записи
	newTime := now.Add(time.Hour)
	watermark2 := poll_watermarks.PollWatermarkModel{
		TournamentId:         1,
		LastSuccessfulPollAt: newTime,
	}
	res2, err := s.repo.FirstOrCreate(watermark2)
	s.Require().NoError(err)
	s.Equal(uint(1), res2.TournamentId)
	s.True(now.Equal(res2.LastSuccessfulPollAt), "Должен вернуться старый watermark, а не новый")
}

func (s *PollWatermarkRepoSuite) TestUpdate() {
	now := time.Now().UTC().Round(time.Second)
	watermark := poll_watermarks.PollWatermarkModel{
		TournamentId:         1,
		LastSuccessfulPollAt: now.Add(-time.Hour),
	}
	s.tx.Create(&watermark)

	// Update
	updatedTime := now
	watermark.LastSuccessfulPollAt = updatedTime
	res, err := s.repo.Update(watermark)

	s.Require().NoError(err)
	s.True(updatedTime.Equal(res.LastSuccessfulPollAt))

	// Verify in DB
	var dbModel poll_watermarks.PollWatermarkModel
	s.tx.First(&dbModel, 1)
	s.True(updatedTime.Equal(dbModel.LastSuccessfulPollAt))
}

func TestPollWatermarkRepoSuite(t *testing.T) {
	suite.Run(t, new(PollWatermarkRepoSuite))
}
