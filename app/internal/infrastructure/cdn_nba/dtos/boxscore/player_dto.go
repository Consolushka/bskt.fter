package boxscore

const playedTimeFormat = "PT%mM%sS"

type PlayerDTO struct {
	Status     string              `json:"status"`
	Order      int                 `json:"order"`
	PersonId   int                 `json:"personId"`
	JerseyNum  string              `json:"jerseyNum"`
	Position   string              `json:"position"`
	Starter    string              `json:"starter"`
	Oncourt    string              `json:"oncourt"`
	Played     string              `json:"played"`
	Statistics PlayerEfficiencyDTO `json:"statistics"`
	Name       string              `json:"name"`
	NameI      string              `json:"nameI"`
	FirstName  string              `json:"firstName"`
	FamilyName string              `json:"familyName"`
}
