package aggregator

import (
	"bytes"
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

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		composite_logger.Error("Failed to send game-imported webhook", map[string]interface{}{
			"url":          url,
			"tournamentId": tournamentId,
			"error":        err,
		})
		return fmt.Errorf("failed to send post request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		composite_logger.Warn("Aggregator returned non-OK status", map[string]interface{}{
			"url":          url,
			"tournamentId": tournamentId,
			"status":       resp.Status,
		})
	}

	return nil
}
