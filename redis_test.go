package stag

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var redis_host = "localhost"
var redis_port = 6379
var redis_db = 0
var redis_queue_name = fmt.Sprintf("name_%0.6f", rand.Float64())

func TestRedisGetSet(t *testing.T) {

	r := NewRedisCache(redis_host, redis_port, redis_db)

	// generate random key and value
	key := fmt.Sprintf("key:%0.6f", rand.Float64())
	val := fmt.Sprintf("val:%0.6f", rand.Float64())

	// unset key should give error
	_, e := r.Get(key)
	assert.NotNil(t, e)

	// set key should be easy to retrieve
	r.Set(key, []byte(val), 5*time.Second)
	v, e := r.Get(key)
	assert.Nil(t, e)
	assert.Equal(t, val, string(v))
}

func TestRedisQueuePushPop(t *testing.T) {
	q := NewRedisQueue(redis_host, redis_port, redis_db, redis_queue_name)

	// pop should give error on empty queue
	b, e := q.Pop()
	assert.Nil(t, b)
	assert.NotNil(t, e)

	// pop_wait should also give error on empty queue
	b, e = q.PopWait(time.Second * 1)
	assert.Nil(t, b)
	assert.NotNil(t, e)

	// when que is not empty, then pop works
	r1 := fmt.Sprintf("qval:%0.6f", rand.Float64())
	q.Push([]byte(r1))
	b, e = q.Pop()
	assert.Equal(t, r1, string(b))
	assert.Nil(t, e)
	len, _ := q.Len()
	assert.Zero(t, len)

	// when que is not empty, then popwait also works
	r2 := fmt.Sprintf("qval:%0.6f", rand.Float64())
	q.Push([]byte(r2))
	b, e = q.PopWait(time.Second * 1)
	assert.Equal(t, r2, string(b))
	assert.Nil(t, e)
	len, _ = q.Len()
	assert.Zero(t, len)

	// when 2 items are pushed, lenght must be 2
	r3 := fmt.Sprintf("qval:%0.6f", rand.Float64())
	r4 := fmt.Sprintf("qval:%0.6f", rand.Float64())
	q.Push([]byte(r3))
	q.Push([]byte(r4))
	len, e = q.Len()
	assert.Equal(t, 2, len)
	assert.Nil(t, e)
	// pop one value
	b, e = q.Pop()
	assert.Nil(t, e)
	assert.Equal(t, r3, string(b))
	len, e = q.Len()
	assert.Equal(t, 1, len)
	// pop_wait second value
	b, e = q.PopWait(time.Second * 1)
	assert.Nil(t, e)
	assert.Equal(t, r4, string(b))
	len, e = q.Len()
	assert.Equal(t, 0, len)
}
