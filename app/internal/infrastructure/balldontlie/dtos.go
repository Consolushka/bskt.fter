package balldontlie

type Team struct {
	ID           int    `json:"id"`
	Abbreviation string `json:"abbreviation"`
	City         string `json:"city"`
	Conference   string `json:"conference"`
	Division     string `json:"division"`
	FullName     string `json:"full_name"`
	Name         string `json:"name"`
}

type Player struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Position     string `json:"position"`
	Height       string `json:"height"`
	Weight       string `json:"weight"`
	College      string `json:"college"`
	Country      string `json:"country"`
	DraftYear    int    `json:"draft_year"`
	DraftRound   int    `json:"draft_round"`
	DraftNumber  int    `json:"draft_number"`
	Team         Team   `json:"team"`
	JerseyNumber string `json:"jersey_number"`
}

type Meta struct {
	PerPage int `json:"per_page"`
}
