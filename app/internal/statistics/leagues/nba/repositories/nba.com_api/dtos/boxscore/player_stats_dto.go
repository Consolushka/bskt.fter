package boxscore

type PlayerEfficiencyDTO struct {
	Assists                 int     `json:"assists"`
	Blocks                  int     `json:"blocks"`
	BlocksReceived          int     `json:"blocksReceived"`
	FieldGoalsAttempted     int     `json:"fieldGoalsAttempted"`
	FieldGoalsMade          int     `json:"fieldGoalsMade"`
	FieldGoalsPercentage    float64 `json:"fieldGoalsPercentage"`
	FoulsOffensive          int     `json:"foulsOffensive"`
	FoulsDrawn              int     `json:"foulsDrawn"`
	FoulsPersonal           int     `json:"foulsPersonal"`
	FoulsTechnical          int     `json:"foulsTechnical"`
	FreeThrowsAttempted     int     `json:"freeThrowsAttempted"`
	FreeThrowsMade          int     `json:"freeThrowsMade"`
	FreeThrowsPercentage    float64 `json:"freeThrowsPercentage"`
	Minus                   int     `json:"minus"`
	Minutes                 string  `json:"minutes"`
	MinutesCalculated       string  `json:"minutesCalculated"`
	Plus                    int     `json:"plus"`
	PlusMinusPoints         int     `json:"plusMinusPoints"`
	Points                  int     `json:"points"`
	PointsFastBreak         int     `json:"pointsFastBreak"`
	PointsInThePaint        int     `json:"pointsInThePaint"`
	PointsSecondChance      int     `json:"pointsSecondChance"`
	ReboundsDefensive       int     `json:"reboundsDefensive"`
	ReboundsOffensive       int     `json:"reboundsOffensive"`
	ReboundsTotal           int     `json:"reboundsTotal"`
	Steals                  int     `json:"steals"`
	ThreePointersAttempted  int     `json:"threePointersAttempted"`
	ThreePointersMade       int     `json:"threePointersMade"`
	ThreePointersPercentage float64 `json:"threePointersPercentage"`
	Turnovers               int     `json:"turnovers"`
	TwoPointersAttempted    int     `json:"twoPointersAttempted"`
	TwoPointersMade         int     `json:"twoPointersMade"`
	TwoPointersPercentage   float64 `json:"twoPointersPercentage"`
}
