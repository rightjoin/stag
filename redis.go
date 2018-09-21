package stak

import (
	"fmt"
	"strings"
	"time"

	"github.com/rightjoin/fig"

	redis "gopkg.in/redis.v5"
)

type Redis struct {
	AllCharsIndex
	r    *redis.Client
	name string
}

func NewRedisCache(host string, port int, db int) Redis {
	serv := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", host, port),
		DB:   db,
	})

	_, err := serv.Ping().Result()
	if err != nil {
		panic(err)
	}

	return Redis{
		r: serv,
	}
}

func NewRedisQueue(host string, port int, db int, name string) Redis {
	serv := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", host, port),
		DB:   db,
	})

	_, err := serv.Ping().Result()
	if err != nil {
		panic(err)
	}

	return Redis{
		r:    serv,
		name: name,
	}
}

func NewRedisConfig(container ...string) Redis {

	parent := strings.Join(container, ".")

	type Config struct {
		Host string
		Port int
		Db   int
		Name string `fig:"optional"`
	}
	var c Config
	fig.Struct(&c, parent)
	if c.Name == "" {
		return NewRedisCache(c.Host, c.Port, c.Db)
	}
	return NewRedisQueue(c.Host, c.Port, c.Db, c.Name)
}

/* Cache Methods */

func (rd Redis) Set(key string, data []byte, expireIn time.Duration) error {
	key = rd.PrepareIndex(key)
	cmd := rd.r.Set(key, data, expireIn)
	if err := cmd.Err(); err != nil {
		return err
	}
	return nil
}

func (rd Redis) Get(key string) ([]byte, error) {
	key = rd.PrepareIndex(key)
	data, err := rd.r.Get(key).Bytes()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (rd Redis) Delete(key string) error {
	key = rd.PrepareIndex(key)
	return rd.r.Del(key).Err()
}

func (rd Redis) Close() error {
	return rd.r.Close()
}

/* End Cache Methods */

/* Queue Methods */

func (rd Redis) Push(data []byte) error {
	return rd.r.LPush(rd.name, string(data)).Err()
}

func (rd Redis) Pop() ([]byte, error) {
	cmd := rd.r.RPop(rd.name)
	if err := cmd.Err(); err != nil {
		return nil, err
	}

	return cmd.Bytes()
}

func (rd Redis) PopWait(dur time.Duration) ([]byte, error) {

	cmd := rd.r.BRPop(dur, rd.name)

	// check for error
	if err := cmd.Err(); err != nil {
		return nil, err
	}

	rslt, err := cmd.Result()
	if err != nil {
		return nil, err
	}

	return []byte(rslt[1]), nil
}

func (rd Redis) Len() (int, error) {
	i, e := rd.r.LLen(rd.name).Result()
	return int(i), e
}

/* End Queue Methods */
