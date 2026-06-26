package leagues_repo

import (
	"IMP/app/internal/core/leagues"
	"IMP/app/pkg/dbtest"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type LeaguesRepoSuite struct {
	suite.Suite
	db   *gorm.DB
	tx   *gorm.DB
	repo Gorm
}

func (s *LeaguesRepoSuite) SetupSuite() {
	s.db = dbtest.Setup(s.T(), &leagues.LeagueModel{})
}

func (s *LeaguesRepoSuite) SetupTest() {
	s.tx = s.db.Begin()
	s.repo = NewGormRepo(s.tx)
}

func (s *LeaguesRepoSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *LeaguesRepoSuite) TestFirstOrCreate_CreatesNew() {
	result, err := s.repo.FirstOrCreate(leagues.LeagueModel{Name: "Euroleague", Alias: "euroleague"})

	s.Require().NoError(err)
	s.NotZero(result.Id)
	s.Equal("Euroleague", result.Name)
	s.Equal("euroleague", result.Alias)
}

func (s *LeaguesRepoSuite) TestFirstOrCreate_ReturnsExisting() {
	existing := leagues.LeagueModel{Name: "Euroleague", Alias: "euroleague"}
	s.Require().NoError(s.tx.Create(&existing).Error)

	result, err := s.repo.FirstOrCreate(leagues.LeagueModel{Name: "Euroleague", Alias: "euroleague"})

	s.Require().NoError(err)
	// тот же id — новой строки не создалось
	s.Equal(existing.Id, result.Id)

	var count int64
	s.Require().NoError(s.tx.Model(&leagues.LeagueModel{}).Count(&count).Error)
	s.Equal(int64(1), count)
}

func TestLeaguesRepoSuite(t *testing.T) {
	suite.Run(t, new(LeaguesRepoSuite))
}
