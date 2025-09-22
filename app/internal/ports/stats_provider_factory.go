package ports

type StatsProviderFactory interface {
	Create() (StatsProvider, error)
	ProviderName() string
}
