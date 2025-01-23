package enums

import (
	"IMP/app/internal/modules/statistics/enums"
	"strconv"
)

type TimeBasedImpCoefficient int

const (
	// Per20 Per24 Minutes for role players
	Per20 TimeBasedImpCoefficient = 20
	Per24 TimeBasedImpCoefficient = 24
	// Per30 Per38 Minutes for starters/all stars
	Per30 TimeBasedImpCoefficient = 30
	Per38 TimeBasedImpCoefficient = 38
	// Per40 Per48 Full game
	Per40 TimeBasedImpCoefficient = 40
	Per48 TimeBasedImpCoefficient = 48
)

func TimeBasesByLeagueAndPers(league enums.League, impPer ImpPERs) TimeBasedImpCoefficient {
	switch league {
	case enums.NBA:
		switch impPer {
		case Bench:
			return Per24
		case Starter:
			return Per38
		case FullGame:
			return Per48
		}
	case enums.MLBL:
		switch impPer {
		case Bench:
			return Per20
		case Starter:
			return Per30
		case FullGame:
			return Per40
		}
	}

	return 0
}

func TimeBasesByLeague(league enums.League) []TimeBasedImpCoefficient {
	switch league {
	case enums.NBA:
		return []TimeBasedImpCoefficient{
			Per24,
			Per38,
			Per48,
		}
	case enums.MLBL:
		return []TimeBasedImpCoefficient{
			Per20,
			Per30,
			Per40,
		}
	default:
		return []TimeBasedImpCoefficient{}
	}
}

func (t TimeBasedImpCoefficient) Title() string {
	return "Per" + strconv.Itoa(int(t))
}

func (t TimeBasedImpCoefficient) Minutes() int {
	switch t {
	case Per20:
		return 20
	case Per24:
		return 24
	case Per30:
		return 30
	case Per38:
		return 38
	case Per40:
		return 40
	case Per48:
		return 48
	default:
		return 0
	}
}

func (t TimeBasedImpCoefficient) Seconds() int {
	return t.Minutes() * 60
}

func (t TimeBasedImpCoefficient) InsufficientDistanceCoef() float64 {
	switch t {
	case Per20:
		return 0.05
	case Per24:
		return 0.0003
	case Per30:
		return 0.00024
	case Per38:
		return 0.00012
	case Per40:
		return 0.00008
	case Per48:
		return 0.000034
	default:
		return 0
	}
}

func (t TimeBasedImpCoefficient) InsufficientDistancePower() float64 {
	switch t {
	case Per20:
		return 1
	case Per24:
		return 2.8
	case Per30:
		return 3
	case Per38:
		return 3
	case Per40:
		return 3
	case Per48:
		return 3.2
	default:
		return 0
	}
}

func (t TimeBasedImpCoefficient) SufficientDistanceOffset() float64 {
	switch t {
	case Per20:
		return 0.4
	case Per24:
		return 0.189
	case Per30:
		return 0.252
	case Per38:
		return 0.405
	case Per40:
		return 0.327
	case Per48:
		return 0.42
	default:
		return 0
	}

}

func (t TimeBasedImpCoefficient) SufficientDistancePower() float64 {
	switch t {
	case Per20:
		return 0.6
	case Per24:
		return 0.6
	case Per30:
		return 0.6
	case Per38:
		return 0.6
	case Per40:
		return 0.8
	case Per48:
		return 0.4
	default:
		return 0
	}

}

func (t TimeBasedImpCoefficient) OverSufficientDistanceUpperPower() float64 {
	switch t {
	case Per20:
		return 2
	case Per24:
		return 2
	case Per30:
		return 2
	case Per38:
		return 2
	case Per40:
		return 2
	case Per48:
		return 2
	default:
		return 0
	}

}

func (t TimeBasedImpCoefficient) OverSufficientDistanceLowerPower() float64 {
	switch t {
	case Per20:
		return 3.1
	case Per24:
		return 2.8
	case Per30:
		return 2.3
	case Per38:
		return 2.2
	case Per40:
		return 2
	case Per48:
		return 1.9
	default:
		return 0
	}

}
