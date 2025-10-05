package ports

import "time"

type CachedRemoteResource interface {
	LocalFileName() string
	GetLifeTime() time.Duration
	Load() (any, error)
}
