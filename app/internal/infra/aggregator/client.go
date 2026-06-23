package aggregator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	composite_logger "github.com/Consolushka/golang.composite_logger/pkg"
)

type ApiClient struct {
	host string
}

func NewApiClient(host string) *ApiClient {
	return &ApiClient{
		host: host,
	}
}

func (c *ApiClient) NotifyGameImported(tournamentId uint) error {
	if c.host == "" {
		return nil
	}

	url := fmt.Sprintf("%s/webhooks/aggregator/game-imported", c.host)
	payload := map[string]interface{}{
		"tournament_id": tournamentId,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Уведомление агрегатора — самостоятельная операция "отправил и забыл";
	// её не нужно отменять вместе с обработкой турнира, поэтому отдельный фоновый контекст.
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		composite_logger.Error("Failed to send game-imported webhook", map[string]interface{}{
			"url":          url,
			"tournamentId": tournamentId,
			"error":        err,
		})
		return fmt.Errorf("failed to send post request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		composite_logger.Warn("Aggregator returned non-OK status", map[string]interface{}{
			"url":          url,
			"tournamentId": tournamentId,
			"status":       resp.Status,
		})
	}

	return nil
}
