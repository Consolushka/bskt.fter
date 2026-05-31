package ports

type Aggregator interface {
	NotifyGameImported(tournamentId uint) error
}
