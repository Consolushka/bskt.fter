package cdn_nba

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestNewCdnNbaClient verifies that the NewCdnNbaClient function
// correctly initializes a new Client with the expected base URL
func TestNewCdnNbaClient(t *testing.T) {
	// Call the function being tested
	client := NewCdnNbaClient()

	// Verify the client was initialized with the correct base URL
	assert.Equal(t, "https://cdn.nba.com/static/json", client.baseUrl)
}

// TestClient_BoxScore verifies the behavior of the BoxScore method
// in the Client struct:
// - Verify that when a valid game ID is provided - returns correct BoxScoreDto with all expected fields
// - Verify that when a valid game with overtime ID is provided - returns correct BoxScoreDto with all expected fields
func TestClient_BoxScore(t *testing.T) {
	cases := []struct {
		name                 string
		gameId               string
		expectedJsonFileName string // Оставляю пустым для заполнения
	}{
		{
			name:                 "Valid game ID returns correct box score - Lakers vs Warriors",
			gameId:               "0022201063", // Lakers vs Warriors game
			expectedJsonFileName: "mock_cdn_nba_0022201063.json",
		},
		{
			name:                 "Valid game with overtime ID returns correct box score - Lakers vs Warriors",
			gameId:               "0042400223", // Lakers vs Warriors game
			expectedJsonFileName: "mock_cdn_nba_0042400223.json",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client := NewCdnNbaClient()

			result := client.BoxScore(tc.gameId)

			var expectedDto GameBoxScoreResponse
			err := loadTestData(tc.expectedJsonFileName, &expectedDto)

			assert.NoError(t, err, "Failed to load test data")

			assert.Equal(t, expectedDto.Game, result)
		})
	}
}

// TestClient_ScheduleSeason verifies the behavior of the ScheduleSeason method
// in the Client struct:
// - Verify that the method correctly fetches and returns the full season schedule
// - Verify that the schedule contains valid game data for multiple dates
func TestClient_ScheduleSeason(t *testing.T) {
	// Create client
	client := NewCdnNbaClient()

	// Call the method being tested
	result := client.ScheduleSeason()

	// Verify that the response contains games
	assert.NotEmpty(t, result.Games, "Schedule should contain games")
	assert.True(t, len(result.Games) > 10, "Schedule should contain multiple game dates")

	// Check structure of game dates and games
	for _, gameDate := range result.Games {
		// Verify each game date has a valid date string
		assert.NotEmpty(t, gameDate.GameDate, "Game date should not be empty")

		// Verify each game date has games
		assert.NotEmpty(t, gameDate.Games, "Game date should contain games")

		// Check structure of individual games
		for _, game := range gameDate.Games {
			// Verify game has an ID
			assert.NotEmpty(t, game.GameId, "Game should have an ID")

			// Verify game has a status
			assert.NotZero(t, game.GameStatus, "Game should have a status")

			// Skip games that are not yet played
			if game.GameStatus == 1 {
				continue
			}

			// Verify game has team information
			assert.NotEmpty(t, game.AwayTeam.TeamName, "Away team should have a name")
			assert.NotEmpty(t, game.HomeTeam.TeamName, "Home team should have a name")
			assert.NotZero(t, game.AwayTeam.TeamId, "Away team should have an ID")
			assert.NotZero(t, game.HomeTeam.TeamId, "Home team should have an ID")
		}
	}

	// Verify that the schedule contains games from different months
	months := make(map[string]bool)
	for _, gameDate := range result.Games {
		// Extract month from date string (assuming format like "MM/DD/YYYY HH:MM:SS")
		if len(gameDate.GameDate) >= 2 {
			month := gameDate.GameDate[0:2]
			months[month] = true
		}
	}
	assert.True(t, len(months) >= 3, "Schedule should contain games from at least 3 different months")

	// Verify that the schedule contains both regular season and playoff games
	// This is a simple heuristic - regular season games typically have IDs starting with "002"
	// while playoff games have IDs starting with "004"
	hasRegularSeasonGames := false

	for _, gameDate := range result.Games {
		for _, game := range gameDate.Games {
			if len(game.GameId) >= 3 {
				if game.GameId[0:3] == "002" {
					hasRegularSeasonGames = true
					break
				}
			}
		}
	}

	assert.True(t, hasRegularSeasonGames, "Schedule should contain regular season games")
	// Playoff games might not be available depending on the time of year
	// so we don't assert this must be true
}

func loadTestData(fileName string, target interface{}) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(target)
}
