package enums

type ImpPERs string

const (
	Clean    ImpPERs = "Clean"
	Bench    ImpPERs = "Bench"
	Start    ImpPERs = "Start"
	FullGame ImpPERs = "FullGame"
)

func (ip ImpPERs) Order() int {
	switch ip {
	case Clean:
		return 0
	case Bench:
		return 1
	case Start:
		return 2
	case FullGame:
		return 3
	}

	return -1
}

func (ip ImpPERs) ToString() string {
	return string(ip)
}
