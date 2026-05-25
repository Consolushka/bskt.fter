package api_nba

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_DynamicURLBuilding(t *testing.T) {
	var capturedURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedURL = r.URL.String()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := &Client{
		baseUrl:     server.URL,
		baseHeaders: make(map[string]string),
		limiter:     nil, // We don't need limiter for this test if we call methods directly
	}

	t.Run("Games with multiple parameters", func(t *testing.T) {
		_, _ = client.Games(1, "2023-01-01", "1", "2022", "10", "UTC")
		assert.Contains(t, capturedURL, "/games")
		assert.Contains(t, capturedURL, "id=1")
		assert.Contains(t, capturedURL, "date=2023-01-01")
		assert.Contains(t, capturedURL, "league=1")
		assert.Contains(t, capturedURL, "season=2022")
		assert.Contains(t, capturedURL, "team=10")
		assert.Contains(t, capturedURL, "timezone=UTC")
	})

	t.Run("PlayersStatistics with zero/empty values", func(t *testing.T) {
		_, _ = client.PlayersStatistics(0, 123, 0, "")
		assert.Equal(t, "/players/statistics?game=123", capturedURL)
	})

	t.Run("PlayerInfo with multiple parameters and some zeros", func(t *testing.T) {
		_, _ = client.PlayerInfo(101, "James", 0, 2023, "USA", "")
		assert.Contains(t, capturedURL, "/players")
		assert.Contains(t, capturedURL, "id=101")
		assert.Contains(t, capturedURL, "name=James")
		assert.NotContains(t, capturedURL, "team=")
		assert.Contains(t, capturedURL, "season=2023")
		assert.Contains(t, capturedURL, "country=USA")
		assert.NotContains(t, capturedURL, "search=")
	})
}
