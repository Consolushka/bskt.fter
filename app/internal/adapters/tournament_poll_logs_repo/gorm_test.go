package tournament_poll_logs_repo

import (
	"IMP/app/internal/core/tournament_poll_logs"
	"IMP/app/pkg/dbtest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type PollLogRepoSuite struct {
	suite.Suite
	db   *gorm.DB
	tx   *gorm.DB
	repo Gorm
}

func (s *PollLogRepoSuite) SetupSuite() {
	s.db = dbtest.Setup(s.T(), &tournament_poll_logs.TournamentPollLogModel{})
}

func (s *PollLogRepoSuite) SetupTest() {
	s.tx = s.db.Begin()
	s.repo = Gorm{db: s.tx}
}

func (s *PollLogRepoSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *PollLogRepoSuite) TestCreate() {
	now := time.Now().UTC().Round(time.Second)
	log := tournament_poll_logs.TournamentPollLogModel{
		TournamentId:    1,
		PollStartAt:     now.Add(-time.Minute),
		IntervalStart:   now.Add(-time.Hour),
		IntervalEnd:     now,
		SavedGamesCount: 5,
		Status:          tournament_poll_logs.StatusSuccess,
	}

	res, err := s.repo.Create(log)

	s.Require().NoError(err)
	s.NotZero(res.Id)
	s.Equal(uint(1), res.TournamentId)
	s.Equal(5, res.SavedGamesCount)

	var dbLog tournament_poll_logs.TournamentPollLogModel
	s.tx.First(&dbLog, res.Id)
	s.Equal(tournament_poll_logs.StatusSuccess, dbLog.Status)
}

func (s *PollLogRepoSuite) TestGetLatestSuccess() {
	now := time.Now().UTC().Round(time.Second)

	// Seed some logs
	s.tx.Create(&tournament_poll_logs.TournamentPollLogModel{
		TournamentId: 1,
		Status:       tournament_poll_logs.StatusSuccess,
		IntervalEnd:  now.Add(-time.Hour),
	})
	s.tx.Create(&tournament_poll_logs.TournamentPollLogModel{
		TournamentId: 1,
		Status:       tournament_poll_logs.StatusSuccess,
		IntervalEnd:  now, // This should be the latest
	})
	s.tx.Create(&tournament_poll_logs.TournamentPollLogModel{
		TournamentId: 1,
		Status:       tournament_poll_logs.StatusError,
		IntervalEnd:  now.Add(time.Hour), // Error, should be ignored
	})
	s.tx.Create(&tournament_poll_logs.TournamentPollLogModel{
		TournamentId: 2, // Different tournament
		Status:       tournament_poll_logs.StatusSuccess,
		IntervalEnd:  now.Add(2 * time.Hour),
	})

	latest, err := s.repo.GetLatestSuccess(1)

	s.Require().NoError(err)
	s.Equal(uint(1), latest.TournamentId)
	s.True(now.Equal(latest.IntervalEnd), "Should return log with latest interval_end")
}

func (s *PollLogRepoSuite) TestGetLatestSuccess_NotFound() {
	_, err := s.repo.GetLatestSuccess(999)
	s.Require().Error(err)
	s.Equal(gorm.ErrRecordNotFound, err)
}

func TestPollLogRepoSuite(t *testing.T) {
	suite.Run(t, new(PollLogRepoSuite))
}
