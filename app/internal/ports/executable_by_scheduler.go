package ports

import "time"

type ExecutableByScheduler interface {
	GetName() string
	GetPeriodicity() time.Duration
	Execute() error
}
