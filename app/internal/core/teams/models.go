package teams

import (
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type TeamModel struct {
	Id        uint      `gorm:"column:id"`
	Name      string    `gorm:"column:name"`
	HomeTown  string    `gorm:"column:home_town"`
	Alias     string    `gorm:"column:alias"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (TeamModel) TableName() string {
	return "teams"
}

// Слова, которыми клубы начинают название («БК „Тамбов“»): в алиасе они бесполезны,
// у половины команд алиас начинался бы одинаково.
var clubPrefixes = map[string]struct{}{"BK": {}, "BC": {}, "БК": {}}

const maxWholeNameAlias = 6

func (t *TeamModel) AutoGenerateAlias() string {
	unquoted := strings.Map(func(r rune) rune {
		if unicode.Is(unicode.Quotation_Mark, r) {
			return -1
		}
		return r
	}, t.Name)
	words := strings.Fields(unquoted)
	if len(words) > 1 {
		if _, ok := clubPrefixes[strings.ToUpper(words[0])]; ok {
			words = words[1:]
		}
	}

	runes := []rune(strings.Join(words, " "))
	if len(runes) > maxWholeNameAlias {
		runes = runes[:3]
	}

	return strings.ToUpper(string(runes))
}

// Infobasket и Sportoteka вместо сокращения отдают первые 6 символов имени: `BK "Ta` для `BK "Tambov"`.
const providerTruncatedLen = 6

// ProviderAlias отбрасывает сокращение провайдера, если это просто обрезанное имя.
// Пустой алиас потом сгенерируется из имени. Более короткое начало имени («CSKA» для «CSKA Moscow»)
// считаем настоящим сокращением.
func ProviderAlias(abbr, fullName string) string {
	if utf8.RuneCountInString(abbr) == providerTruncatedLen &&
		strings.HasPrefix(strings.ToUpper(fullName), strings.ToUpper(abbr)) {
		return ""
	}

	return strings.ToUpper(abbr)
}

type GameTeamStatModel struct {
	Id        uint      `gorm:"column:id"`
	GameId    uint      `gorm:"column:game_id"`
	TeamId    uint      `gorm:"column:team_id"`
	Score     int       `gorm:"column:score"`
	FinalDiff int       `gorm:"column:final_differential"`
	IsHome    bool      `gorm:"column:is_home"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (GameTeamStatModel) TableName() string {
	return "game_team_stats"
}
