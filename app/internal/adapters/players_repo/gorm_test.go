package players_repo

import (
	"IMP/app/internal/core/players"
	"IMP/app/pkg/dbtest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type PlayerRepoSuite struct {
	suite.Suite
	db   *gorm.DB
	tx   *gorm.DB
	repo Gorm
}

func (s *PlayerRepoSuite) SetupSuite() {
	s.db = dbtest.Setup(s.T(), &players.PlayerModel{}, &players.GameTeamPlayerStatModel{})
}

func (s *PlayerRepoSuite) SetupTest() {
	s.tx = s.db.Begin()
	s.repo = NewGormRepo(s.tx)
}

func (s *PlayerRepoSuite) TearDownTest() {
	s.tx.Rollback()
}

func (s *PlayerRepoSuite) TestFirstOrCreate_FindExisting() {
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	p1 := players.PlayerModel{FullName: "LeBron James", BirthDate: birthDate}
	p2 := players.PlayerModel{FullName: "Kevin Durant", BirthDate: birthDate}
	s.tx.Create(&p1)
	s.tx.Create(&p2)

	input := players.PlayerModel{FullName: "LeBron James", BirthDate: birthDate}
	res, err := s.repo.FirstOrCreate(input)

	s.Require().NoError(err)
	s.Equal(p1.Id, res.Id)

	var count int64
	s.tx.Model(&players.PlayerModel{}).Count(&count)
	s.Equal(int64(2), count)
}

func (s *PlayerRepoSuite) TestFirstOrCreate_CreateNew() {
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	s.tx.Create(&players.PlayerModel{FullName: "LeBron James", BirthDate: birthDate})

	input := players.PlayerModel{FullName: "Stephen Curry", BirthDate: birthDate}
	res, err := s.repo.FirstOrCreate(input)

	s.Require().NoError(err)
	s.NotZero(res.Id)
	s.Equal("Stephen Curry", res.FullName)

	var count int64
	s.tx.Model(&players.PlayerModel{}).Count(&count)
	s.Equal(int64(2), count)
}

func (s *PlayerRepoSuite) TestListByFullName() {
	s.tx.Create(&players.PlayerModel{FullName: "John Doe"})
	s.tx.Create(&players.PlayerModel{FullName: "John Doe"}) // Same name, different ID
	s.tx.Create(&players.PlayerModel{FullName: "Jane Doe"})

	res, err := s.repo.ListByFullName("John Doe")

	s.Require().NoError(err)
	s.Len(res, 2)
	s.Equal("John Doe", res[0].FullName)
	s.Equal("John Doe", res[1].FullName)
}

func (s *PlayerRepoSuite) TestFirstOrCreateStat_FindExisting() {
	stat1 := players.GameTeamPlayerStatModel{GameId: 1, TeamId: 10, PlayerId: 100, PlsMin: 5}
	stat2 := players.GameTeamPlayerStatModel{GameId: 1, TeamId: 10, PlayerId: 101, PlsMin: -2}
	s.tx.Create(&stat1)
	s.tx.Create(&stat2)

	input := players.GameTeamPlayerStatModel{GameId: 1, TeamId: 10, PlayerId: 100, PlsMin: 10}
	res, err := s.repo.FirstOrCreateStat(input)

	s.Require().NoError(err)
	s.Equal(stat1.Id, res.Id)
	s.Equal(int8(5), res.PlsMin, "Should return existing value")

	var count int64
	s.tx.Model(&players.GameTeamPlayerStatModel{}).Count(&count)
	s.Equal(int64(2), count)
}

func (s *PlayerRepoSuite) TestFirstOrCreateStat_CreateNew() {
	stat1 := players.GameTeamPlayerStatModel{GameId: 1, TeamId: 10, PlayerId: 100, PlsMin: 5}
	s.tx.Create(&stat1)

	input := players.GameTeamPlayerStatModel{GameId: 2, TeamId: 10, PlayerId: 100, PlsMin: 8}
	res, err := s.repo.FirstOrCreateStat(input)

	s.Require().NoError(err)
	s.NotZero(res.Id)
	s.Equal(int8(8), res.PlsMin)

	var count int64
	s.tx.Model(&players.GameTeamPlayerStatModel{}).Count(&count)
	s.Equal(int64(2), count)
}

func TestPlayerRepoSuite(t *testing.T) {
	suite.Run(t, new(PlayerRepoSuite))
}
