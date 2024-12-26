package enums

type League int

const (
	NBA League = iota
	MLBL
)

func (l League) FullGameTime() int {
	switch l {
	case NBA:
		return 48
	case MLBL:
		return 40
	default:
		return 0
	}

}

func (l League) QuarterDuration() int {
	switch l {
	case NBA:
		return 12
	case MLBL:
		return 10
	default:
		return 0
	}
}

func (l League) OvertimeDuration() int {
	switch l {
	case NBA:
		return 6
	case MLBL:
		return 5
	default:
		return 0
	}
}

func FromString(league string) League {
	switch league {
	case "NBA":
		return NBA
	case "MLBL":
		return MLBL
	default:
		return NBA
	}
}
