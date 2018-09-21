package stak

import (
	"math/rand"
	"time"
)

func init() {
	// initialize the random number seed
	rand.Seed(time.Now().UTC().UnixNano())
}
