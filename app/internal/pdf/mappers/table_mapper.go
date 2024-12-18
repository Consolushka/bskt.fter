package mappers

type TableMapper interface {
	Headers() []string
	ToTable() []string
}
