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
	s.db = dbtest.Setup(s.T(), &teams.TeamModel{}, &teams.GameTeamStatModel{}, &teams.TeamExternalIdModel{})
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
	res, created, err := s.repo.FirstOrCreate(input)

	// 3. Assert
	s.Require().NoError(err)
	s.False(created)
	s.Equal(lakers.Id, res.Id, "Should find the existing Lakers team")

	var count int64
	s.tx.Model(&teams.TeamModel{}).Count(&count)
	s.Equal(int64(2), count, "Total teams count should remain 2")
}

func (s *TeamRepoSuite) TestFirstOrCreate_ExistingKeepsAlias() {
	lakers := teams.TeamModel{Name: "Lakers", HomeTown: "Los Angeles", Alias: "LAL"}
	s.tx.Create(&lakers)

	res, created, err := s.repo.FirstOrCreate(teams.TeamModel{Name: "Lakers", HomeTown: "Los Angeles", Alias: "LAK"})

	s.Require().NoError(err)
	s.False(created)
	s.Equal("LAL", res.Alias)

	var stored teams.TeamModel
	s.tx.First(&stored, lakers.Id)
	s.Equal("LAL", stored.Alias)
}

func (s *TeamRepoSuite) TestFirstOrCreate_CreateNewGeneratesAlias() {
	res, created, err := s.repo.FirstOrCreate(teams.TeamModel{Name: `BK "Tambov"`, HomeTown: "Tambov"})

	s.Require().NoError(err)
	s.True(created)
	s.Equal("TAMBOV", res.Alias)
}

func (s *TeamRepoSuite) TestFirstOrCreate_CreateNew() {
	// 1. Seed one team
	lakers := teams.TeamModel{Name: "Lakers", HomeTown: "Los Angeles"}
	s.tx.Create(&lakers)

	// 2. Try to find/create Bulls
	input := teams.TeamModel{Name: "Bulls", HomeTown: "Chicago"}
	res, created, err := s.repo.FirstOrCreate(input)

	// 3. Assert
	s.Require().NoError(err)
	s.True(created)
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

func (s *TeamRepoSuite) TestFirstOrCreateStats_PersistsIsHome() {
	home, err := s.repo.FirstOrCreateStats(teams.GameTeamStatModel{GameId: 1, TeamId: 10, IsHome: true})
	s.Require().NoError(err)
	away, err := s.repo.FirstOrCreateStats(teams.GameTeamStatModel{GameId: 1, TeamId: 20, IsHome: false})
	s.Require().NoError(err)

	var storedHome, storedAway teams.GameTeamStatModel
	s.Require().NoError(s.tx.First(&storedHome, home.Id).Error)
	s.Require().NoError(s.tx.First(&storedAway, away.Id).Error)
	s.True(storedHome.IsHome)
	s.False(storedAway.IsHome)
}

func TestTeamRepoSuite(t *testing.T) {
	suite.Run(t, new(TeamRepoSuite))
}

func (s *TeamRepoSuite) TestFirstOrCreate_EmptyHomeTownMatchesByName() {
	existing := teams.TeamModel{Name: "Real Madrid", HomeTown: "Madrid", Alias: "RMA"}
	s.tx.Create(&existing)

	res, created, err := s.repo.FirstOrCreate(teams.TeamModel{Name: "Real Madrid"})

	s.Require().NoError(err)
	s.False(created)
	s.Equal(existing.Id, res.Id)
	s.Equal("RMA", res.Alias)

	var count int64
	s.tx.Model(&teams.TeamModel{}).Count(&count)
	s.Equal(int64(1), count)
}

func (s *TeamRepoSuite) TestFirstOrCreate_TooLongAliasIsRegenerated() {
	res, created, err := s.repo.FirstOrCreate(teams.TeamModel{Name: "Warriors", HomeTown: "Golden State", Alias: "ABCDEFGHIJK"})

	s.Require().NoError(err)
	s.True(created)
	s.Equal("WAR", res.Alias)
}

func (s *TeamRepoSuite) TestFindByExternalId_Found() {
	team := teams.TeamModel{Name: "Lakers", HomeTown: "Los Angeles", Alias: "LAL"}
	s.tx.Create(&team)
	s.tx.Create(&teams.TeamExternalIdModel{Scope: "ESPN_NBA", ExternalId: 9, TeamId: team.Id})

	res, found, err := s.repo.FindByExternalId("ESPN_NBA", 9)

	s.Require().NoError(err)
	s.True(found)
	s.Equal(team.Id, res.Id)
	s.Equal("Lakers", res.Name)
	s.Equal("LAL", res.Alias)
}

func (s *TeamRepoSuite) TestFindByExternalId_NotFound() {
	_, found, err := s.repo.FindByExternalId("ESPN_NBA", 9)

	s.Require().NoError(err)
	s.False(found)
}

func (s *TeamRepoSuite) TestFindByExternalId_ScopeIsPartOfKey() {
	team := teams.TeamModel{Name: "Lakers", HomeTown: "Los Angeles", Alias: "LAL"}
	s.tx.Create(&team)
	s.tx.Create(&teams.TeamExternalIdModel{Scope: "INFOBASKET:org", ExternalId: 9, TeamId: team.Id})

	_, found, err := s.repo.FindByExternalId("INFOBASKET:reg", 9)

	s.Require().NoError(err)
	s.False(found)
}

func (s *TeamRepoSuite) TestLinkExternalId_Creates() {
	team := teams.TeamModel{Name: "Lakers", HomeTown: "Los Angeles", Alias: "LAL"}
	s.tx.Create(&team)

	s.Require().NoError(s.repo.LinkExternalId("ESPN_NBA", 9, team.Id))

	res, found, err := s.repo.FindByExternalId("ESPN_NBA", 9)
	s.Require().NoError(err)
	s.True(found)
	s.Equal(team.Id, res.Id)
}

func (s *TeamRepoSuite) TestLinkExternalId_ExistingLinkIsKept() {
	first := teams.TeamModel{Name: "Lakers", HomeTown: "Los Angeles", Alias: "LAL"}
	second := teams.TeamModel{Name: "Bulls", HomeTown: "Chicago", Alias: "CHI"}
	s.tx.Create(&first)
	s.tx.Create(&second)

	s.Require().NoError(s.repo.LinkExternalId("ESPN_NBA", 9, first.Id))
	s.Require().NoError(s.repo.LinkExternalId("ESPN_NBA", 9, second.Id))

	res, found, err := s.repo.FindByExternalId("ESPN_NBA", 9)
	s.Require().NoError(err)
	s.True(found)
	s.Equal(first.Id, res.Id)

	var count int64
	s.tx.Model(&teams.TeamExternalIdModel{}).Count(&count)
	s.Equal(int64(1), count)
}
