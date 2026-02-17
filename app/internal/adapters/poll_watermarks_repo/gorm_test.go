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

func (s *PollWatermarkRepoSuite) TestFirstOrCreate_FindExisting() {
	// 1. Seed two watermarks
	now := time.Now().UTC().Round(time.Second)
	w1 := poll_watermarks.PollWatermarkModel{TournamentId: 1, LastSuccessfulPollAt: now}
	w2 := poll_watermarks.PollWatermarkModel{TournamentId: 2, LastSuccessfulPollAt: now.Add(time.Hour)}
	s.tx.Create(&w1)
	s.tx.Create(&w2)

	// 2. Try to find watermark for Tournament 1 with DIFFERENT time
	input := poll_watermarks.PollWatermarkModel{
		TournamentId:         1,
		LastSuccessfulPollAt: now.Add(24 * time.Hour),
	}
	res, err := s.repo.FirstOrCreate(input)

	// 3. Assert
	s.Require().NoError(err)
	s.Equal(uint(1), res.TournamentId)
	s.True(now.Equal(res.LastSuccessfulPollAt), "Should return EXISTING time, not the one from input")

	var count int64
	s.tx.Model(&poll_watermarks.PollWatermarkModel{}).Count(&count)
	s.Equal(int64(2), count)
}

func (s *PollWatermarkRepoSuite) TestFirstOrCreate_CreateNew() {
	// 1. Seed one watermark
	now := time.Now().UTC().Round(time.Second)
	s.tx.Create(&poll_watermarks.PollWatermarkModel{TournamentId: 1, LastSuccessfulPollAt: now})

	// 2. Create watermark for Tournament 2
	input := poll_watermarks.PollWatermarkModel{
		TournamentId:         2,
		LastSuccessfulPollAt: now.Add(time.Hour),
	}
	res, err := s.repo.FirstOrCreate(input)

	// 3. Assert
	s.Require().NoError(err)
	s.Equal(uint(2), res.TournamentId)
	s.True(now.Add(time.Hour).Equal(res.LastSuccessfulPollAt))

	var count int64
	s.tx.Model(&poll_watermarks.PollWatermarkModel{}).Count(&count)
	s.Equal(int64(2), count)
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
