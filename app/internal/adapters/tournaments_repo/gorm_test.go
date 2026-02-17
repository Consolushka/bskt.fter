package tournaments_repo

import (
	"IMP/app/internal/core/leagues"
	"IMP/app/internal/core/tournaments"
	"IMP/app/pkg/dbtest"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type TournamentRepoSuite struct {
	suite.Suite
	db   *gorm.DB
	tx   *gorm.DB
	repo Gorm
}

// SetupSuite выполняется один раз перед всеми тестами в наборе
func (s *TournamentRepoSuite) SetupSuite() {
	// Создаем базу и мигрируем модели один раз
	s.db = dbtest.Setup(s.T(),
		&leagues.LeagueModel{},
		&tournaments.TournamentModel{},
		&tournaments.TournamentProvider{},
	)
}

// SetupTest выполняется перед каждым тестом
func (s *TournamentRepoSuite) SetupTest() {
	// Начинаем транзакцию
	s.tx = s.db.Begin()
	// Инициализируем репозиторий с транзакционным объектом БД
	s.repo = NewGormRepo(s.tx)
}

// TearDownTest выполняется после каждого теста
func (s *TournamentRepoSuite) TearDownTest() {
	// Откатываем все изменения, сделанные в тесте
	s.tx.Rollback()
}

func (s *TournamentRepoSuite) TestListByLeagueAliases() {
	// Seed (внутри транзакции)
	targetLeague := leagues.LeagueModel{Name: "Target", Alias: "target"}
	s.tx.Create(&targetLeague)

	otherLeague := leagues.LeagueModel{Name: "Other", Alias: "other"}
	s.tx.Create(&otherLeague)

	t1 := tournaments.TournamentModel{Name: "T1", LeagueId: targetLeague.Id}
	t2 := tournaments.TournamentModel{Name: "T2", LeagueId: targetLeague.Id}
	s.tx.Create(&t1)
	s.tx.Create(&t2)

	tOther := tournaments.TournamentModel{Name: "Other T", LeagueId: otherLeague.Id}
	s.tx.Create(&tOther)

	externalId := "123"
	s.tx.Create(&tournaments.TournamentProvider{
		TournamentId: t1.Id,
		ProviderName: "api_nba",
		ExternalId:   &externalId,
	})

	// Execute
	results, err := s.repo.ListByLeagueAliases([]string{"target"})

	// Assert
	s.Require().NoError(err)
	s.Len(results, 2)

	var foundT1 bool
	for _, r := range results {
		if r.Id == t1.Id {
			foundT1 = true
			s.Equal("T1", r.Name)
			s.NotNil(r.League)
			s.Equal("target", r.League.Alias)
			s.Equal("api_nba", r.Provider.ProviderName)
		}
	}
	s.True(foundT1)
}

func (s *TournamentRepoSuite) TestGet() {
	// Seed (внутри транзакции)
	league := leagues.LeagueModel{Name: "Target", Alias: "target"}
	s.tx.Create(&league)

	t1 := tournaments.TournamentModel{Name: "Target T", LeagueId: league.Id}
	s.tx.Create(&t1)

	externalId := "ext-123"
	s.tx.Create(&tournaments.TournamentProvider{
		TournamentId: t1.Id,
		ProviderName: "api_nba",
		ExternalId:   &externalId,
	})

	// Execute
	res, err := s.repo.Get(t1.Id)

	// Assert
	s.Require().NoError(err)
	s.Equal(t1.Id, res.Id)
	s.Equal("Target T", res.Name)
	s.Equal("target", res.League.Alias)
	s.Equal("api_nba", res.Provider.ProviderName)
}

func TestTournamentRepoSuite(t *testing.T) {
	suite.Run(t, new(TournamentRepoSuite))
}
