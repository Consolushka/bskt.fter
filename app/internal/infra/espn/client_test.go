package espn

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_URLBuilding(t *testing.T) {
	var siteURL, coreURL, userAgent string
	status := http.StatusOK

	handler := func(captured *string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			*captured = r.URL.String()
			userAgent = r.Header.Get("User-Agent")
			w.WriteHeader(status)
			_, _ = w.Write([]byte(`{}`))
		}
	}
	site := httptest.NewServer(handler(&siteURL))
	defer site.Close()
	core := httptest.NewServer(handler(&coreURL))
	defer core.Close()

	c := NewClient(0).(*Client)
	c.siteBaseUrl = site.URL
	c.coreBaseUrl = core.URL

	t.Run("Scoreboard builds dates query", func(t *testing.T) {
		_, err := c.Scoreboard("20261028")
		require.NoError(t, err)
		assert.Equal(t, "/scoreboard?dates=20261028", siteURL)
	})

	t.Run("Summary builds event query", func(t *testing.T) {
		_, err := c.Summary("401909888")
		require.NoError(t, err)
		assert.Equal(t, "/summary?event=401909888", siteURL)
	})

	t.Run("Athlete uses core base url", func(t *testing.T) {
		_, err := c.Athlete("5175737")
		require.NoError(t, err)
		assert.Equal(t, "/athletes/5175737", coreURL)
	})

	t.Run("sends User-Agent header", func(t *testing.T) {
		_, err := c.Scoreboard("20261028")
		require.NoError(t, err)
		assert.Equal(t, "Argus/1.0 (+bskt.imp)", userAgent)
	})

	t.Run("returns error on non-2xx", func(t *testing.T) {
		status = http.StatusNotFound
		_, err := c.Scoreboard("20261028")
		require.ErrorContains(t, err, "status 404")
	})
}
