package teams_repo

import (
	"IMP/app/internal/core/teams"
	"IMP/app/pkg/dbtest"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TeamRepoSuite struct {
	suite.Suite
	db   *gorm.DB
	tx   *gorm.DB
	repo Gorm
}

func (s *TeamRepoSuite) SetupSuite() {
	s.db = dbtest.Setup(s.T(), &teams.TeamModel{}, &teams.GameTeamStatModel{})
}

func (s *TeamRepoSuite) SetupTest() {
	s.tx = s.db.Begin()
	s.repo = NewGormRepo(s.tx)
}

func (s *TeamRepoSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *TeamRepoSuite) TestFirstOrCreate_FindExisting() {
	// 1. Seed two teams
	lakers := teams.TeamModel{Name: "Lakers", HomeTown: "Los Angeles"}
	bulls := teams.TeamModel{Name: "Bulls", HomeTown: "Chicago"}
	s.tx.Create(&lakers)
	s.tx.Create(&bulls)

	// 2. Try to find Lakers
	input := teams.TeamModel{Name: "Lakers", HomeTown: "Los Angeles"}
	res, err := s.repo.FirstOrCreate(input)

	// 3. Assert
	s.Require().NoError(err)
	s.Equal(lakers.Id, res.Id, "Should find the existing Lakers team")

	var count int64
	s.tx.Model(&teams.TeamModel{}).Count(&count)
	s.Equal(int64(2), count, "Total teams count should remain 2")
}

func (s *TeamRepoSuite) TestFirstOrCreate_CreateNew() {
	// 1. Seed one team
	lakers := teams.TeamModel{Name: "Lakers", HomeTown: "Los Angeles"}
	s.tx.Create(&lakers)

	// 2. Try to find/create Bulls
	input := teams.TeamModel{Name: "Bulls", HomeTown: "Chicago"}
	res, err := s.repo.FirstOrCreate(input)

	// 3. Assert
	s.Require().NoError(err)
	s.NotZero(res.Id)
	s.NotEqual(lakers.Id, res.Id)
	s.Equal("Bulls", res.Name)

	var count int64
	s.tx.Model(&teams.TeamModel{}).Count(&count)
	s.Equal(int64(2), count, "Total teams count should be 2 after creation")
}

func (s *TeamRepoSuite) TestFirstOrCreateStats_FindExisting() {
	// 1. Seed two stats
	stat1 := teams.GameTeamStatModel{GameId: 1, TeamId: 10, Score: 100}
	stat2 := teams.GameTeamStatModel{GameId: 2, TeamId: 10, Score: 90}
	s.tx.Create(&stat1)
	s.tx.Create(&stat2)

	// 2. Try to find stat for Game 1
	input := teams.GameTeamStatModel{GameId: 1, TeamId: 10, Score: 120} // Score is different in input
	res, err := s.repo.FirstOrCreateStats(input)

	// 3. Assert
	s.Require().NoError(err)
	s.Equal(stat1.Id, res.Id)
	s.Equal(100, res.Score, "Should return existing score, not the one from input")

	var count int64
	s.tx.Model(&teams.GameTeamStatModel{}).Count(&count)
	s.Equal(int64(2), count)
}

func (s *TeamRepoSuite) TestFirstOrCreateStats_CreateNew() {
	// 1. Seed one stat
	stat1 := teams.GameTeamStatModel{GameId: 1, TeamId: 10, Score: 100}
	s.tx.Create(&stat1)

	// 2. Try to find/create stat for Game 2
	input := teams.GameTeamStatModel{GameId: 2, TeamId: 10, Score: 95}
	res, err := s.repo.FirstOrCreateStats(input)

	// 3. Assert
	s.Require().NoError(err)
	s.NotZero(res.Id)
	s.NotEqual(stat1.Id, res.Id)
	s.Equal(95, res.Score)

	var count int64
	s.tx.Model(&teams.GameTeamStatModel{}).Count(&count)
	s.Equal(int64(2), count)
}

func TestTeamRepoSuite(t *testing.T) {
	suite.Run(t, new(TeamRepoSuite))
}
