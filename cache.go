package stag

import (
	"strings"
	"time"

	"github.com/rightjoin/fig"
)

type Cache interface {
	IndexFormatter
	Set(key string, data []byte, expireIn time.Duration) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	Close() error
}

func NewCache(container ...string) (out Cache) {

	parent := strings.Join(container, ".")

	// find engine
	engn := fig.String(parent, "engine")

	switch engn {
	case "redis":
		out = NewRedisConfig(container...)

	case "go-cache", "gocache":
		out = NewGoCache(15 * time.Minute) // default expiry time is 15 min

	case "logger":
		out = cacheLogger{
			Inner: NewCache(parent, "inner"),
		}

	default:
		panic("unknown cache engine:" + engn)
	}

	return
}
