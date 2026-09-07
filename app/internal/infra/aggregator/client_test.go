package aggregator

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type capturedRequest struct {
	method string
	path   string
	token  string
	body   string
}

func startServer(t *testing.T) (*httptest.Server, *capturedRequest) {
	t.Helper()
	captured := &capturedRequest{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		*captured = capturedRequest{
			method: r.Method,
			path:   r.URL.Path,
			token:  r.Header.Get("X-AGGREGATOR-TOKEN"),
			body:   string(body),
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	t.Cleanup(srv.Close)
	return srv, captured
}

func TestNotifyGameImported_SendsTokenHeader(t *testing.T) {
	srv, captured := startServer(t)

	err := NewApiClient(srv.URL, "secret").NotifyGameImported(42)

	require.NoError(t, err)
	require.Equal(t, http.MethodPost, captured.method)
	require.Equal(t, "/webhooks/aggregator/game-imported", captured.path)
	require.Equal(t, "secret", captured.token)
	require.JSONEq(t, `{"tournament_id":42}`, captured.body)
}

func TestNotifyGameImported_EmptyTokenOmitsHeader(t *testing.T) {
	srv, captured := startServer(t)

	err := NewApiClient(srv.URL, "").NotifyGameImported(42)

	require.NoError(t, err)
	require.Equal(t, http.MethodPost, captured.method)
	require.Empty(t, captured.token)
}

func TestNotifyGameImported_EmptyHostSkipsRequest(t *testing.T) {
	_, captured := startServer(t)

	err := NewApiClient("", "secret").NotifyGameImported(42)

	require.NoError(t, err)
	require.Empty(t, captured.method, "request must not be sent when host is empty")
}
