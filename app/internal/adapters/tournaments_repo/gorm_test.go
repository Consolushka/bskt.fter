package tournaments_repo

import (
	"IMP/app/internal/core/leagues"
	"IMP/app/internal/core/tournaments"
	"IMP/app/pkg/dbtest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TournamentRepoSuite struct {
	suite.Suite
	repo Gorm
}

func (s *TournamentRepoSuite) SetupTest() {
	db := dbtest.Setup(s.T(),
		&leagues.LeagueModel{},
		&tournaments.TournamentModel{},
		&tournaments.TournamentProvider{},
	)
	s.repo = NewGormRepo(db)
}

func (s *TournamentRepoSuite) TestListByLeagueAliases() {
	// Seed
	league := leagues.LeagueModel{Name: "Test", Alias: "test-alias"}
	s.repo.db.Create(&league)

	tournament := tournaments.TournamentModel{Name: "Tournament 1", LeagueId: league.Id}
	s.repo.db.Create(&tournament)

	// Execute
	results, err := s.repo.ListByLeagueAliases([]string{"test-alias"})

	// Assert
	s.Require().NoError(err)
	s.Len(results, 1)
	s.Equal("Tournament 1", results[0].Name)
}

func (s *TournamentRepoSuite) TestGet() {
	tournament := tournaments.TournamentModel{Name: "Target"}
	s.repo.db.Create(&tournament)

	res, err := s.repo.Get(tournament.Id)

	s.Require().NoError(err)
	s.Equal(tournament.Id, res.Id)
}

func TestTournamentRepoSuite(t *testing.T) {
	suite.Run(t, new(TournamentRepoSuite))
}
