package stak

import (
	"errors"
	"time"

	goc "github.com/pmylund/go-cache"
)

type GoCache struct {
	AllCharsIndex
	ram *goc.Cache
}

func NewGoCache(exprires time.Duration) GoCache {
	return GoCache{
		ram: goc.New(exprires, 5*time.Minute),
	}
}

func (g GoCache) Set(key string, data []byte, expireIn time.Duration) error {
	key = g.PrepareIndex(key)
	g.ram.Set(key, data, expireIn)
	return nil
}

func (g GoCache) Get(key string) ([]byte, error) {
	key = g.PrepareIndex(key)
	if i, ok := g.ram.Get(key); ok {
		return i.([]byte), nil
	}
	return nil, errors.New("key not found in cache: " + key)
}

func (g GoCache) Delete(key string) error {
	key = g.PrepareIndex(key)
	g.ram.Delete(key)
	return nil
}

func (g GoCache) Close() error {
	return nil
}
