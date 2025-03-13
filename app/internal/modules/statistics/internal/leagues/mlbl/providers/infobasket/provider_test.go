package infobasket

import (
	"IMP/app/internal/infrastructure/infobasket"
	"IMP/app/internal/infrastructure/infobasket/mock_data"
	"IMP/app/internal/modules/statistics/models"
	"testing"
	"time"
)

// Setup a test Provider with mocked dependencies
func setupTestProvider(mockClient *infobasket.MockClient) *Provider {
	return &Provider{
		client: mockClient,
		mapper: newMapper(),
	}
}

// TestNewProvider tests the NewProvider function with dependencies (client and mapper)
func TestNewProvider(t *testing.T) {
	provider := NewProvider()

	// Check if provider is not nil
	if provider == nil {
		t.Error("Expected provider to be initialized, got nil")
	}

	// Check if client is initialized
	//goland:noinspection ALL
	if provider.client == nil {
		t.Error("Expected client to be initialized, got nil")
	}

	// Check if mapper is initialized
	if provider.mapper == nil {
		t.Error("Expected mapper to be initialized, got nil")
	}
}

// TestProvider_GameBoxScore test that we get exact game we requested by gameId
func TestProvider_GameBoxScore(t *testing.T) {
	// Setup mock data
	testGameId := "12345"
	mockBoxScore := mock_data.CreateMockGameBoxScoreResponse(101, 111, 1, 4, "23.12.2022", "13.00")

	// Mapping tested in mapper_test
	expectedGame := &models.GameBoxScoreDTO{
		Id: testGameId,
	}

	// Setup mock client
	mockClient := &infobasket.MockClient{
		BoxScoreFunc: func(gameId string) infobasket.GameBoxScoreResponse {
			return mockBoxScore
		},
	}

	provider := setupTestProvider(mockClient)

	// Execute test
	result, err := provider.GameBoxScore(testGameId)

	// Assert results
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result.Id != expectedGame.Id {
		t.Errorf("Expected game ID %s, got %s", expectedGame.Id, result.Id)
	}
}

// TestProvider_GamesByDate tests we get correct games by date (test how we filter games by date in provider)
func TestProvider_GamesByDate(t *testing.T) {
	format := "02.01.2006"

	searchingDate, _ := time.Parse(format, "08.03.2025")
	firstWrongDate := searchingDate.AddDate(0, 0, -1)
	secondWrongDate := searchingDate.AddDate(0, 0, 1)

	mockClient := &infobasket.MockClient{
		ScheduledGamesFunc: func(compId int) []infobasket.GameScheduleDto {
			if compId == 89960 {
				return []infobasket.GameScheduleDto{
					mock_data.CreateMockGameScheduleDto(1, searchingDate.Format(format), 1),
					mock_data.CreateMockGameScheduleDto(2, firstWrongDate.Format(format), -1),
					mock_data.CreateMockGameScheduleDto(3, secondWrongDate.Format(format), 0),
					mock_data.CreateMockGameScheduleDto(4, searchingDate.Format(format), 99999),
				}
			} else if compId == 89962 {
				return []infobasket.GameScheduleDto{
					mock_data.CreateMockGameScheduleDto(5, searchingDate.Format(format), 1),
					mock_data.CreateMockGameScheduleDto(6, searchingDate.Format(format), 0),
					mock_data.CreateMockGameScheduleDto(7, secondWrongDate.Format(format), 1),
					mock_data.CreateMockGameScheduleDto(8, firstWrongDate.Format(format), 1),
				}
			}
			panic("Unexpected compId")
		}}

	provider := setupTestProvider(mockClient)
	result, err := provider.GamesByDate(searchingDate)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result) != 4 {
		t.Fatalf("Expected 4 games, got %d", len(result))
	}

	if result[0] != "1" || result[1] != "4" || result[2] != "5" || result[3] != "6" {
		t.Fatalf("Expected game IDs 1 and 6, got %v", result)
	}
}

// TestProvider_GamesByTeam tests we get correct already finished games by team (finished == GameStatus == 1)
func TestProvider_GamesByTeam(t *testing.T) {
	mockClient := &infobasket.MockClient{
		TeamGamesFunc: func(teamId string) infobasket.TeamScheduleResponse {
			return infobasket.TeamScheduleResponse{Games: []infobasket.GameScheduleDto{
				mock_data.CreateMockGameScheduleDto(1, "01.01.2025", 1),
				mock_data.CreateMockGameScheduleDto(2, "02.01.2025", -1),
				mock_data.CreateMockGameScheduleDto(3, "03.01.2025", 0),
				mock_data.CreateMockGameScheduleDto(4, "04.01.2025", 99999),
			}}
		}}

	provider := setupTestProvider(mockClient)
	result, err := provider.GamesByTeam("123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Expected 3 games, got %d", len(result))
	}

	if result[0] != "1" {
		t.Fatalf("Expected game ID 1, got %s", result[0])
	}
}
