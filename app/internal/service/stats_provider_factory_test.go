package service

import (
	"IMP/app/internal/adapters/stats_provider"
	"IMP/app/internal/ports"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewStatsProvider tests the NewStatsProvider function
// Verify that when NBA league name is provided while creating stats provider - returns CdnNbaStatsProviderAdapter and no error
// Verify that when MLBL league name is provided while creating stats provider - returns InfobasketStatsProviderAdapter and no error
// Verify that when unknown league name is provided while creating stats provider - returns nil and specific error message
func TestNewStatsProvider(t *testing.T) {
	cases := []struct {
		name             string
		data             string
		expectedProvider ports.StatsProvider
		expectedError    error
	}{
		{
			name:             "successfully creates NBA stats provider",
			data:             "NBA",
			expectedProvider: stats_provider.ApiNbaStatsProviderAdapter{},
			expectedError:    nil,
		},
		{
			name:             "successfully creates MLBL stats provider",
			data:             "MLBL",
			expectedProvider: stats_provider.InfobasketStatsProviderAdapter{},
			expectedError:    nil,
		},
		{
			name:             "returns error for unknown league name - differs from supported leagues",
			data:             "UNKNOWN",
			expectedProvider: nil,
			expectedError:    errors.New("unknown league: UNKNOWN"),
		},
		{
			name:             "returns error for empty league name - differs from supported leagues",
			data:             "",
			expectedProvider: nil,
			expectedError:    errors.New("unknown league: "),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			result, err := NewStatsProvider(tc.data)

			if tc.expectedError != nil {
				assert.Equal(t, err, tc.expectedError)
				return
			}

			assert.NoError(t, err)
			assert.IsType(t, tc.expectedProvider, result)
		})
	}
}
