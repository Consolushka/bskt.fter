package ports

import "time"

type ExecutableByScheduler interface {
	GetName() string
	ShouldBeExecutedAt() time.Time
	Execute() error
}
