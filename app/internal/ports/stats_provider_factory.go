package ports

type StatsProviderFactory interface {
	Create() (StatsProvider, error)
}
