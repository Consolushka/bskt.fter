package infobasket

type GameScheduleDto struct {
	GameID                 int     `json:"GameID"`
	IsToday                bool    `json:"IsToday"`
	DaysFromToday          int     `json:"DaysFromToday"`
	GameNumber             string  `json:"GameNumber"`
	GameDateInt            int     `json:"GameDateInt"`
	GameDate               string  `json:"GameDate"`
	HasTime                bool    `json:"HasTime"`
	GameTime               string  `json:"GameTime"`
	GameTimeMsk            string  `json:"GameTimeMsk"`
	GameLocalDate          string  `json:"GameLocalDate"`
	GameDateTime           string  `json:"GameDateTime"`
	GameDateTimeMoscow     string  `json:"GameDateTimeMoscow"`
	DisplayDateTimeLocal   string  `json:"DisplayDateTimeLocal"`
	DisplayDateTimeMsk     string  `json:"DisplayDateTimeMsk"`
	DisplayDateTimeLocalEn string  `json:"DisplayDateTimeLocalEn"`
	DisplayDateTimeMskEn   string  `json:"DisplayDateTimeMskEn"`
	GameStatus             int     `json:"GameStatus"`
	TeamAid                int     `json:"TeamAid"`
	TeamBid                int     `json:"TeamBid"`
	ShortTeamNameAru       string  `json:"ShortTeamNameAru"`
	ShortTeamNameBru       string  `json:"ShortTeamNameBru"`
	TeamNameAru            string  `json:"TeamNameAru"`
	TeamNameBru            string  `json:"TeamNameBru"`
	CompTeamNameAru        string  `json:"CompTeamNameAru"`
	CompTeamNameBru        string  `json:"CompTeamNameBru"`
	RegionTeamNameAru      string  `json:"RegionTeamNameAru"`
	RegionTeamNameBru      string  `json:"RegionTeamNameBru"`
	ShortTeamNameAen       string  `json:"ShortTeamNameAen"`
	ShortTeamNameBen       string  `json:"ShortTeamNameBen"`
	TeamNameAen            string  `json:"TeamNameAen"`
	TeamNameBen            string  `json:"TeamNameBen"`
	CompTeamNameAen        string  `json:"CompTeamNameAen"`
	CompTeamNameBen        string  `json:"CompTeamNameBen"`
	RegionTeamNameAen      string  `json:"RegionTeamNameAen"`
	RegionTeamNameBen      string  `json:"RegionTeamNameBen"`
	ScoreA                 int     `json:"ScoreA"`
	ScoreB                 int     `json:"ScoreB"`
	Periods                int     `json:"Periods"`
	TeamLogoA              string  `json:"TeamLogoA"`
	TeamLogoB              string  `json:"TeamLogoB"`
	HasVideo               bool    `json:"HasVideo"`
	LiveID                 string  `json:"LiveID"`
	TvRu                   string  `json:"TvRu"`
	TvEn                   string  `json:"TvEn"`
	VideoID                *string `json:"VideoID"`
	RegionRu               string  `json:"RegionRu"`
	RegionEn               string  `json:"RegionEn"`
	RegionId               int     `json:"RegionId"`
	ArenaEn                string  `json:"ArenaEn"`
	ArenaRu                string  `json:"ArenaRu"`
	ArenaId                int     `json:"ArenaId"`
	GameAttendance         int     `json:"GameAttendance"`
	CompNameRu             string  `json:"CompNameRu"`
	CompNameEn             string  `json:"CompNameEn"`
	LeagueNameRu           string  `json:"LeagueNameRu"`
	LeagueNameEn           string  `json:"LeagueNameEn"`
	LeagueShortNameRu      string  `json:"LeagueShortNameRu"`
	LeagueShortNameEn      string  `json:"LeagueShortNameEn"`
	Gender                 int     `json:"Gender"`
}

type TeamBoxScoreDto struct {
	TeamNumber       int                 `json:"TeamNumber"`
	TeamID           int                 `json:"TeamID"`
	TeamName         TeamNameBoxScoreDto `json:"TeamName"`
	Score            int                 `json:"Score"`
	Points           int                 `json:"Points"`
	Shot1            int                 `json:"Shot1"`
	Goal1            int                 `json:"Goal1"`
	Shot2            int                 `json:"Shot2"`
	Goal2            int                 `json:"Goal2"`
	Shot3            int                 `json:"Shot3"`
	Goal3            int                 `json:"Goal3"`
	PaintShot        int                 `json:"PaintShot"`
	PaintGoal        int                 `json:"PaintGoal"`
	Shots1           string              `json:"Shots1"`
	Shot1Percent     string              `json:"Shot1Percent"`
	Shots2           string              `json:"Shots2"`
	Shot2Percent     string              `json:"Shot2Percent"`
	Shots3           string              `json:"Shots3"`
	Shot3Percent     string              `json:"Shot3Percent"`
	PaintShots       string              `json:"PaintShots"`
	PaintShotPercent string              `json:"PaintShotPercent"`
	Assist           int                 `json:"Assist"`
	Blocks           int                 `json:"Blocks"`
	DefRebound       int                 `json:"DefRebound"`
	OffRebound       int                 `json:"OffRebound"`
	Rebound          int                 `json:"Rebound"`
	Steal            int                 `json:"Steal"`
	Turnover         int                 `json:"Turnover"`
	TeamDefRebound   int                 `json:"TeamDefRebound"`
	TeamOffRebound   *int                `json:"TeamOffRebound"`
	TeamRebound      int                 `json:"TeamRebound"`
	TeamSteal        *int                `json:"TeamSteal"`
	TeamTurnover     int                 `json:"TeamTurnover"`
	Foul             int                 `json:"Foul"`
	OpponentFoul     int                 `json:"OpponentFoul"`
	Seconds          int                 `json:"Seconds"`
	PlayedTime       string              `json:"PlayedTime"`
	PlusMinus        *int                `json:"PlusMinus"`
	//Coach            map[string_utils]interface{} `json:"Coach"`
	Players []PlayerBoxScoreDto `json:"Players"`
	//Coaches          map[string_utils]interface{} `json:"Coaches"`
}

type TeamNameBoxScoreDto struct {
	CompTeamNameID       int           `json:"CompTeamNameID"`
	TeamID               int           `json:"TeamID"`
	TeamType             int           `json:"TeamType"`
	CompTeamShortNameRu  string        `json:"CompTeamShortNameRu"`
	CompTeamShortNameEn  string        `json:"CompTeamShortNameEn"`
	CompTeamNameRu       string        `json:"CompTeamNameRu"`
	CompTeamNameEn       string        `json:"CompTeamNameEn"`
	CompTeamRegionNameRu string        `json:"CompTeamRegionNameRu"`
	CompTeamRegionNameEn string        `json:"CompTeamRegionNameEn"`
	CompTeamAbcNameRu    string        `json:"CompTeamAbcNameRu"`
	CompTeamAbcNameEn    string        `json:"CompTeamAbcNameEn"`
	CompTeamNameChanged  interface{}   `json:"CompTeamNameChanged"`
	CompTeamNameDefault  bool          `json:"CompTeamNameDefault"`
	SysStatus            int           `json:"SysStatus"`
	SysLastChanged       string        `json:"SysLastChanged"`
	SysUser              interface{}   `json:"SysUser"`
	SysMyUser            interface{}   `json:"SysMyUser"`
	CompTeams            []interface{} `json:"CompTeams"`
	Team                 interface{}   `json:"Team"`
	IsRealTeam           bool          `json:"IsRealTeam"`
}

type PlayerBoxScoreDto struct {
	PersonID         int         `json:"PersonID"`
	TeamNumber       int         `json:"TeamNumber"`
	PlayerNumber     int         `json:"PlayerNumber"`
	DisplayNumber    string      `json:"DisplayNumber"`
	LastNameRu       string      `json:"LastNameRu"`
	LastNameEn       string      `json:"LastNameEn"`
	FirstNameRu      string      `json:"FirstNameRu"`
	FirstNameEn      string      `json:"FirstNameEn"`
	PersonNameRu     string      `json:"PersonNameRu"`
	PersonNameEn     string      `json:"PersonNameEn"`
	Capitan          int         `json:"Capitan"`
	PersonBirth      string      `json:"PersonBirth"`
	PosID            int         `json:"PosID"`
	CountryCodeIOC   string      `json:"CountryCodeIOC"`
	CountryNameRu    string      `json:"CountryNameRu"`
	CountryNameEn    string      `json:"CountryNameEn"`
	RankRu           string      `json:"RankRu"`
	RankEn           interface{} `json:"RankEn"`
	Height           int         `json:"Height"`
	Weight           int         `json:"Weight"`
	Points           int         `json:"Points"`
	Shot1            int         `json:"Shot1"`
	Goal1            int         `json:"Goal1"`
	Shots1           string      `json:"Shots1"`
	Shot1Percent     string      `json:"Shot1Percent"`
	Shot2            int         `json:"Shot2"`
	Goal2            int         `json:"Goal2"`
	Shots2           string      `json:"Shots2"`
	Shot2Percent     string      `json:"Shot2Percent"`
	PaintShot        int         `json:"PaintShot"`
	PaintGoal        int         `json:"PaintGoal"`
	PaintShots       string      `json:"PaintShots"`
	PaintShotPercent string      `json:"PaintShotPercent"`
	Shot3            int         `json:"Shot3"`
	Goal3            interface{} `json:"Goal3"`
	Shots3           string      `json:"Shots3"`
	Shot3Percent     string      `json:"Shot3Percent"`
	Assist           int         `json:"Assist"`
	Blocks           int         `json:"Blocks"`
	DefRebound       int         `json:"DefRebound"`
	OffRebound       int         `json:"OffRebound"`
	Rebound          int         `json:"Rebound"`
	Steal            int         `json:"Steal"`
	Turnover         int         `json:"Turnover"`
	Foul             int         `json:"Foul"`
	OpponentFoul     int         `json:"OpponentFoul"`
	PlusMinus        int         `json:"PlusMinus"`
	Seconds          int         `json:"Seconds"`
	PlayedTime       string      `json:"PlayedTime"`
	IsStart          bool        `json:"IsStart"`
	StartMark        string      `json:"StartMark"`
}

func CreateMockGameScheduleDto(gameId int, gameDate string, status int) GameScheduleDto {
	return GameScheduleDto{
		GameID:                 gameId,
		IsToday:                false,
		DaysFromToday:          0,
		GameNumber:             "G-001",
		GameDateInt:            20230101,
		GameDate:               gameDate,
		HasTime:                true,
		GameTime:               "19:00",
		GameTimeMsk:            "21:00",
		GameLocalDate:          gameDate,
		GameDateTime:           gameDate + " 19:00",
		GameDateTimeMoscow:     gameDate + " 21:00",
		DisplayDateTimeLocal:   gameDate + " 19:00",
		DisplayDateTimeMsk:     gameDate + " 21:00",
		DisplayDateTimeLocalEn: gameDate + " 7:00 PM",
		DisplayDateTimeMskEn:   gameDate + " 9:00 PM",
		GameStatus:             status,
		TeamAid:                101,
		TeamBid:                102,
		ShortTeamNameAru:       "Команда A",
		ShortTeamNameBru:       "Команда B",
		TeamNameAru:            "Полное имя команды A",
		TeamNameBru:            "Полное имя команды B",
		ShortTeamNameAen:       "Team A",
		ShortTeamNameBen:       "Team B",
		TeamNameAen:            "Full Team A Name",
		TeamNameBen:            "Full Team B Name",
		CompTeamNameAen:        "Team A",
		CompTeamNameBen:        "Team B",
		CompTeamNameAru:        "Команда A",
		CompTeamNameBru:        "Команда B",
		RegionTeamNameAru:      "Регион A",
		RegionTeamNameBru:      "Регион B",
		RegionTeamNameAen:      "Region A",
		RegionTeamNameBen:      "Region B",
		ScoreA:                 85,
		ScoreB:                 77,
		Periods:                4,
		TeamLogoA:              "https://example.com/logo_a.png",
		TeamLogoB:              "https://example.com/logo_b.png",
		HasVideo:               true,
		LiveID:                 "live123",
		TvRu:                   "Матч ТВ",
		TvEn:                   "Match TV",
		VideoID:                nil,
		RegionRu:               "Москва",
		RegionEn:               "Moscow",
		RegionId:               1,
		ArenaEn:                "Basketball Arena",
		ArenaRu:                "Баскетбольная Арена",
		ArenaId:                42,
		GameAttendance:         5000,
		CompNameRu:             "Чемпионат",
		CompNameEn:             "Championship",
		LeagueNameRu:           "Лига",
		LeagueNameEn:           "League",
		LeagueShortNameRu:      "Л",
		LeagueShortNameEn:      "L",
		Gender:                 1,
	}
}
