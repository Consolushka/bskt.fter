package teams

import (
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type TeamModel struct {
	Id        uint      `gorm:"column:id"`
	Name      string    `gorm:"column:name;uniqueIndex:teams_name_home_town_unique_idx"`
	HomeTown  string    `gorm:"column:home_town;uniqueIndex:teams_name_home_town_unique_idx"`
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
		// Только буквы и цифры: иначе «Le Mans» дал бы алиас с пробелом на конце.
		short := make([]rune, 0, 3)
		for _, r := range runes {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				if short = append(short, r); len(short) == 3 {
					break
				}
			}
		}
		runes = short
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

// TeamExternalIdModel связывает id команды у провайдера с командой в нашей базе.
// Scope — чьи это id: одна база id у провайдера (у Infobasket она своя на каждом хосте).
type TeamExternalIdModel struct {
	Id         uint      `gorm:"column:id"`
	Scope      string    `gorm:"column:scope;uniqueIndex:team_external_ids_scope_external_id_unique_idx"`
	ExternalId int       `gorm:"column:external_id;uniqueIndex:team_external_ids_scope_external_id_unique_idx"`
	TeamId     uint      `gorm:"column:team_id"`
	CreatedAt  time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
}

func (TeamExternalIdModel) TableName() string {
	return "team_external_ids"
}
