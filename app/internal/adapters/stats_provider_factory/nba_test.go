package stats_provider_factory

import (
	"IMP/app/internal/adapters/stats_provider"
	"IMP/app/internal/ports"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNbaStatsProviderFactory_Create tests the Create method
// Verify that when Create is called while creating stats provider - returns CdnNbaStatsProviderAdapter instance and no error
func TestNbaStatsProviderFactory_Create(t *testing.T) {
	cases := []struct {
		name     string
		data     struct{}
		expected ports.StatsProvider
		errorMsg string
	}{
		{
			name:     "successfully creates CdnNbaStatsProviderAdapter",
			data:     struct{}{},
			expected: stats_provider.ApiNbaStatsProviderAdapter{},
			errorMsg: "",
		},
	}

	factory := NbaStatsProviderFactory{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			result, err := factory.Create()

			if tc.errorMsg != "" {
				assert.EqualError(t, err, tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expected, result)
		})
	}
}
