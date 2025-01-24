package enums

type ImpPERs string

const (
	Clean    ImpPERs = "Clean"
	Bench    ImpPERs = "Bench"
	Starter  ImpPERs = "Starter"
	FullGame ImpPERs = "FullGame"
)

func (ip ImpPERs) Order() int {
	switch ip {
	case Clean:
		return 0
	case Bench:
		return 1
	case Starter:
		return 2
	case FullGame:
		return 3
	}

	return -1
}
