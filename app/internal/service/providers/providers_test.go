package providers

import (
	"IMP/app/internal/adapters/stats_provider"
	"IMP/app/internal/infra/config"
	"strings"
	"testing"
)

func TestNewProvider(t *testing.T) {
	extID := "12345"
	invalidExtID := "abc"
	sportotekaTag := "vbl"

	infobasketParams := map[string]interface{}{
		"leadHost": "tambov",
	}
	sportotekaParams := map[string]interface{}{
		"year": float64(2025),
	}
	emptyParams := map[string]interface{}{}

	cfg := config.ProvidersConfig{
		ApiSportApiKey:         "test-token",
		ApiNbaRateLimit:        10,
		ApiBasketballRateLimit: 10,
		InfobasketRateLimit:    25,
		SportotekaRateLimit:    25,
	}

	tests := []struct {
		name         string
		providerName string
		externalID   *string
		params       *map[string]interface{}
		wantErr      string
		expectedType string
	}{
		{
			name:         "returns API_NBA provider",
			providerName: ApiNba,
			wantErr:      "",
			expectedType: "api_nba",
		},
		{
			name:         "returns error for CDN_NBA provider",
			providerName: CdnNba,
			wantErr:      "not implemented",
		},
		{
			name:         "returns error for unknown provider",
			providerName: "UNKNOWN",
			wantErr:      "unknown provider: UNKNOWN",
		},
		{
			name:         "returns error for INFOBASKET with nil external id",
			providerName: Infobasket,
			params:       &infobasketParams,
			wantErr:      "external id must be set",
		},
		{
			name:         "returns error for INFOBASKET with invalid external id",
			providerName: Infobasket,
			externalID:   &invalidExtID,
			params:       &infobasketParams,
			wantErr:      "invalid syntax",
		},
		{
			name:         "returns error for INFOBASKET with nil params",
			providerName: Infobasket,
			externalID:   &extID,
			params:       nil,
			wantErr:      "params must be set",
		},
		{
			name:         "returns error for INFOBASKET without leadHost",
			providerName: Infobasket,
			externalID:   &extID,
			params:       &emptyParams,
			wantErr:      "leadHost must be set",
		},
		{
			name:         "returns INFOBASKET provider",
			providerName: Infobasket,
			externalID:   &extID,
			params:       &infobasketParams,
			expectedType: "infobasket",
		},
		{
			name:         "returns error for SPORTOTEKA with nil external id",
			providerName: Sportoteka,
			params:       &sportotekaParams,
			wantErr:      "external id must be set",
		},
		{
			name:         "returns error for SPORTOTEKA with nil params",
			providerName: Sportoteka,
			externalID:   &sportotekaTag,
			params:       nil,
			wantErr:      "params must be set",
		},
		{
			name:         "returns error for SPORTOTEKA without year",
			providerName: Sportoteka,
			externalID:   &sportotekaTag,
			params:       &emptyParams,
			wantErr:      "year must be set",
		},
		{
			name:         "returns SPORTOTEKA provider",
			providerName: Sportoteka,
			externalID:   &sportotekaTag,
			params:       &sportotekaParams,
			expectedType: "sportoteka",
		},
		{
			name:         "returns API_BASKETBALL provider",
			providerName: ApiBasketball,
			externalID:   &extID,
			params:       &emptyParams,
			expectedType: "api_basketball",
		},
		{
			name:         "returns error for API_BASKETBALL with nil external id",
			providerName: ApiBasketball,
			params:       &emptyParams,
			wantErr:      "external id (league id) must be set",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			provider, err := NewProvider(tc.providerName, tc.externalID, tc.params, cfg)

			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error %q, got nil", tc.wantErr)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error to contain %q, got %q", tc.wantErr, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if provider == nil {
				t.Fatal("expected non-nil provider")
			}

			switch tc.expectedType {
			case "api_nba":
				if _, ok := provider.(stats_provider.ApiNbaStatsProviderAdapter); !ok {
					t.Fatalf("expected ApiNbaStatsProviderAdapter, got %T", provider)
				}
			case "api_basketball":
				if _, ok := provider.(stats_provider.ApiBasketballStatsProviderAdapter); !ok {
					t.Fatalf("expected ApiBasketballStatsProviderAdapter, got %T", provider)
				}
			case "infobasket":
				if _, ok := provider.(stats_provider.InfobasketStatsProviderAdapter); !ok {
					t.Fatalf("expected InfobasketStatsProviderAdapter, got %T", provider)
				}
			case "sportoteka":
				if _, ok := provider.(stats_provider.SportotekaStatsProviderAdapter); !ok {
					t.Fatalf("expected SportotekaStatsProviderAdapter, got %T", provider)
				}
			}
		})
	}
}
