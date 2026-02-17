package games_repo

import (
	"IMP/app/internal/core/games"
	"IMP/app/pkg/dbtest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type GameRepoSuite struct {
	suite.Suite
	db   *gorm.DB
	tx   *gorm.DB
	repo Gorm
}

func (s *GameRepoSuite) SetupSuite() {
	s.db = dbtest.Setup(s.T(), &games.GameModel{})
}

func (s *GameRepoSuite) SetupTest() {
	s.tx = s.db.Begin()
	s.repo = NewGormRepo(s.tx)
}

func (s *GameRepoSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *GameRepoSuite) TestExists() {
	now := time.Now().UTC().Truncate(time.Second)

	// Seed one game
	existingGame := games.GameModel{
		TournamentId: 1,
		Title:        "Lakers vs Bulls",
		ScheduledAt:  now,
	}
	s.tx.Create(&existingGame)

	// 1. Case: Something exists but NOT what we are looking for (Different Title)
	exists, err := s.repo.Exists(games.GameModel{
		TournamentId: 1,
		Title:        "Warriors vs Celtics",
		ScheduledAt:  now,
	})
	s.Require().NoError(err)
	s.False(exists, "Should return false for different title")

	// 1.1 Case: Different TournamentId
	exists, err = s.repo.Exists(games.GameModel{
		TournamentId: 2,
		Title:        "Lakers vs Bulls",
		ScheduledAt:  now,
	})
	s.Require().NoError(err)
	s.False(exists, "Should return false for different tournament id")

	// 2. Case: What we are looking for exists
	exists, err = s.repo.Exists(games.GameModel{
		TournamentId: 1,
		Title:        "Lakers vs Bulls",
		ScheduledAt:  now,
	})
	s.Require().NoError(err)
	s.True(exists, "Should return true for existing game")
}

func (s *GameRepoSuite) TestFirstOrCreate_FindExisting() {
	now := time.Now().UTC().Truncate(time.Second)

	// Seed
	g1 := games.GameModel{TournamentId: 1, Title: "Match 1", ScheduledAt: now, Duration: 48}
	g2 := games.GameModel{TournamentId: 1, Title: "Match 2", ScheduledAt: now, Duration: 48}
	s.tx.Create(&g1)
	s.tx.Create(&g2)

	// Try to find Match 1
	input := games.GameModel{TournamentId: 1, Title: "Match 1", ScheduledAt: now, Duration: 60} // Different duration
	res, err := s.repo.FirstOrCreate(input)

	s.Require().NoError(err)
	s.Equal(g1.Id, res.Id)
	s.Equal(48, res.Duration, "Should return existing duration, not the input one")

	var count int64
	s.tx.Model(&games.GameModel{}).Count(&count)
	s.Equal(int64(2), count)
}

func (s *GameRepoSuite) TestFirstOrCreate_CreateNew() {
	now := time.Now().UTC().Truncate(time.Second)
	s.tx.Create(&games.GameModel{TournamentId: 1, Title: "Match 1", ScheduledAt: now})

	// Try to create Match 2
	input := games.GameModel{TournamentId: 1, Title: "Match 2", ScheduledAt: now, Duration: 48}
	res, err := s.repo.FirstOrCreate(input)

	s.Require().NoError(err)
	s.NotZero(res.Id)
	s.Equal("Match 2", res.Title)

	var count int64
	s.tx.Model(&games.GameModel{}).Count(&count)
	s.Equal(int64(2), count)
}

func TestGameRepoSuite(t *testing.T) {
	suite.Run(t, new(GameRepoSuite))
}
