package stak

import (
	"encoding/json"
	"time"

	log "github.com/rightjoin/log15"
)

type cacheLogger struct {
	AllCharsIndex
	Inner Cache
}

func (c cacheLogger) Set(key string, data []byte, expireIn time.Duration) error {
	k := c.PrepareIndex(key)

	// use json format, if data is convertible to it
	var obj interface{}
	err := json.Unmarshal(data, &obj)
	if err != nil {
		obj = data
	}
	log.Info("setting cache", "module", "cache", "key", k, "key-final", c.Inner.PrepareIndex(key), "data", obj)

	return c.Inner.Set(key, data, expireIn)
}

func (c cacheLogger) Get(key string) ([]byte, error) {
	data, err := c.Inner.Get(key)
	k := c.PrepareIndex(key)

	if err != nil {
		log.Info("get cache - not found", "module", "cache", "key", k, "key-final", c.Inner.PrepareIndex(key))
		return nil, err
	}

	// use json format, if data is convertible to it
	var obj interface{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		obj = data
	}
	log.Info("get cache - success", "module", "cache", "key", k, "key-final", c.Inner.PrepareIndex(key), "data", obj)
	return data, nil
}

func (c cacheLogger) Delete(key string) error {
	k := c.PrepareIndex(key)
	log.Info("delete cache", "module", "cache", "key", k, "key-final", c.Inner.PrepareIndex(key))
	return c.Inner.Delete(key)
}

func (c cacheLogger) Close() error {
	log.Info("close cache", "module", "cache")
	return c.Inner.Close()
}
