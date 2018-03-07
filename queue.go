package stag

import (
	"strings"
	"time"

	"github.com/rightjoin/fig"
)

type Queue interface {
	Push(data []byte) error
	Pop() ([]byte, error)
	PopWait(time.Duration) ([]byte, error)
	Len() (int, error)
	Close() error
}

func NewQueue(container ...string) (out Queue) {
	parent := strings.Join(container, ".")

	// find engine
	engn := fig.String(parent, "engine")

	switch engn {
	case "redis":
		redis := NewRedisConfig(container...)
		if redis.name == "" {
			panic("redis queue name missing")
		}
		out = redis

	default:
		panic("unknown queue engine:" + engn)
	}

	return
}
