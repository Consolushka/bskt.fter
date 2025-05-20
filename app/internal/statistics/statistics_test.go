package statistics

import (
	"IMP/app/internal/domain"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

// TestNewLeagueProvider tests the NewLeagueProvider function under various conditions:
// - Verify that when NBA league is provided - an NBA provider is returned
// - Verify that when MLBL league is provided - an MLBL provider is returned
// - Verify that when an unsupported league is provided - an error is returned
// - Verify that the function is case-insensitive when comparing league aliases
// - Verify that the error message is correct when an unsupported league is provided
func TestNewLeagueProvider(t *testing.T) {
	cases := []struct {
		name           string
		league         domain.League
		expectedType   interface{}
		expectedErrMsg string
	}{
		{
			name: "NBA league returns NBA provider",
			league: domain.League{
				AliasEn: strings.ToUpper(domain.NBAAlias),
			},
			expectedType:   &nbaProvider{},
			expectedErrMsg: "",
		},
		{
			name: "MLBL league returns MLBL provider",
			league: domain.League{
				AliasEn: strings.ToUpper(domain.MLBLAlias),
			},
			expectedType:   &mlblProvider{},
			expectedErrMsg: "",
		},
		{
			name: "Lowercase NBA alias still returns NBA provider",
			league: domain.League{
				AliasEn: strings.ToLower(domain.NBAAlias),
			},
			expectedType:   &nbaProvider{},
			expectedErrMsg: "",
		},
		{
			name: "Lowercase MLBL alias still returns MLBL provider",
			league: domain.League{
				AliasEn: strings.ToLower(domain.MLBLAlias),
			},
			expectedType:   &mlblProvider{},
			expectedErrMsg: "",
		},
		{
			name: "Mixed case NBA alias still returns NBA provider",
			league: domain.League{
				AliasEn: "NbA",
			},
			expectedType:   &nbaProvider{},
			expectedErrMsg: "",
		},
		{
			name: "Mixed case MLBL alias still returns MLBL provider",
			league: domain.League{
				AliasEn: "MlBl",
			},
			expectedType:   &mlblProvider{},
			expectedErrMsg: "",
		},
		{
			name: "Unsupported league returns error",
			league: domain.League{
				AliasEn: "UNSUPPORTED",
			},
			expectedType:   nil,
			expectedErrMsg: "There is no provider for league: " + strings.ToUpper(domain.MLBLAlias),
		},
		{
			name: "Empty league alias returns error",
			league: domain.League{
				AliasEn: "",
			},
			expectedType:   nil,
			expectedErrMsg: "There is no provider for league: " + strings.ToUpper(domain.MLBLAlias),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			provider, err := NewLeagueProvider(&tc.league)

			if tc.expectedErrMsg != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedErrMsg)
				assert.Nil(t, provider)

				// Verify that the error message contains the MLBL alias, not the NBA alias
				if tc.league.AliasEn != strings.ToUpper(domain.MLBLAlias) &&
					tc.league.AliasEn != strings.ToLower(domain.MLBLAlias) {
					assert.Contains(t, err.Error(), strings.ToUpper(domain.MLBLAlias))
					assert.NotContains(t, err.Error(), strings.ToUpper(domain.NBAAlias))
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, provider)

				// Check the concrete type of the provider
				expectedType := reflect.TypeOf(tc.expectedType)
				actualType := reflect.TypeOf(provider)
				assert.Equal(t, expectedType, actualType,
					"Expected provider of type %v but got %v", expectedType, actualType)

				// Verify that the provider implements the StatsProvider interface
				_, ok := provider.(StatsProvider)
				assert.True(t, ok, "Provider should implement StatsProvider interface")
			}
		})
	}
}
