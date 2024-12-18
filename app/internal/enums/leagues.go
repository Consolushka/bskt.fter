package enums

type League int

const (
	NBA League = iota
	MLBL
)

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
