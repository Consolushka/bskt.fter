package teams

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTeamModel_TableName(t *testing.T) {
	model := TeamModel{}
	assert.Equal(t, "teams", model.TableName())
}

func TestTeamModel_AutoGenerateAlias(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty name", input: "", expected: ""},
		{name: "long name takes first three letters", input: "Warriors", expected: "WAR"},
		{name: "short name kept whole", input: "OKC", expected: "OKC"},
		{name: "six-char name kept whole", input: "CSKA-2", expected: "CSKA-2"},
		{name: "lowercase is uppercased", input: "cs", expected: "CS"},
		{name: "multibyte long name", input: "Динамо-Москва", expected: "ДИН"},
		{name: "multibyte short name", input: "МГУ", expected: "МГУ"},
		{name: "space skipped in short alias", input: "Le Mans Sarthe", expected: "LEM"},
		{name: "punctuation skipped in short alias", input: "St. Petersburg", expected: "STP"},
		{name: "quotes and latin club prefix dropped", input: `BK "Tambov"`, expected: "TAMBOV"},
		{name: "guillemets and BC prefix dropped", input: "BC «Derzhava»", expected: "DER"},
		{name: "low-high quotes dropped", input: "БК „Тамбов“", expected: "ТАМБОВ"},
		{name: "single curly quotes dropped", input: "BC ‘Avtodor’", expected: "AVT"},
		{name: "cyrillic club prefix dropped", input: "БК Московский", expected: "МОС"},
		{name: "club prefix alone is kept", input: "БК", expected: "БК"},
		{name: "prefix glued to word is not dropped", input: "БК10", expected: "БК10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := TeamModel{Name: tt.input}
			assert.Equal(t, tt.expected, model.AutoGenerateAlias())
		})
	}
}

func TestProviderAlias(t *testing.T) {
	tests := []struct {
		name     string
		abbr     string
		fullName string
		expected string
	}{
		{name: "real abbreviation kept", abbr: "LAL", fullName: "Lakers", expected: "LAL"},
		{name: "truncated name dropped", abbr: `BK "Ta`, fullName: `BK "Tambov"`, expected: ""},
		{name: "truncated name dropped case-insensitively", abbr: "SOYUZ-", fullName: "Soyuz-12", expected: ""},
		{name: "truncated multibyte name dropped", abbr: "Союз-1", fullName: "Союз-12", expected: ""},
		{name: "short name prefix kept as real abbreviation", abbr: "CSKA", fullName: "CSKA Moscow", expected: "CSKA"},
		{name: "abbreviation equal to name kept", abbr: "МГУ", fullName: "МГУ", expected: "МГУ"},
		{name: "six-char abbreviation equal to name dropped", abbr: "BK Ufa", fullName: "BK Ufa", expected: ""},
		{name: "abbreviation uppercased", abbr: "Lak", fullName: "Lakers", expected: "LAK"},
		{name: "empty abbreviation stays empty", abbr: "", fullName: "Lakers", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ProviderAlias(tt.abbr, tt.fullName))
		})
	}
}

func TestTeamExternalIdModel_TableName(t *testing.T) {
	model := TeamExternalIdModel{}
	assert.Equal(t, "team_external_ids", model.TableName())
}
