package boxscore

import (
	"FTER/app/internal/models"
	"FTER/app/internal/translator"
	"FTER/app/internal/utils/string_utils"
)

type PlayerBoxscore struct {
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

func (p *PlayerBoxscore) ToFterModel() models.PlayerModel {
	personName := p.PersonNameEn
	if string_utils.HasNonLanguageChars(personName, string_utils.Latin) {
		personName = translator.Translate(personName, nil, "en")
	}

	return models.PlayerModel{
		FullName:      personName,
		SecondsPlayed: p.Seconds,
		PlsMin:        p.PlusMinus,
	}
}
