package domain

type TableMapper interface {
	Headers() []string
	ToTable() []string
}
