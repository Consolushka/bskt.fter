package enums

type League int

const (
	NBA League = iota
	MLBL
)

func (l League) FullGameDuration() int {
	switch l {
	case NBA:
		return 48
	case MLBL:
		return 40
	default:
		return 0
	}
}
