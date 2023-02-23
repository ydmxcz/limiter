package limiter

import (
	"math"
	"sync"
	"time"
)

type TokenBucket struct {
	rate     float64
	capacity float64
	mutex    sync.Mutex
	token    float64
	lastTime time.Time
}

func NewTokenBucket(rate float64, capacity float64) *TokenBucket {
	return &TokenBucket{
		rate:     rate,
		capacity: capacity,
		token:    0,
		lastTime: time.Now(),
	}
}

func (r *TokenBucket) Allow() bool {
	now := time.Now()

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// add some tokens as a constant rate every second
	addTokenNum := now.Sub(r.lastTime).Seconds() * r.rate
	// if the bucket are full,drop
	r.token = math.Min(r.capacity, r.token+addTokenNum)
	r.lastTime = now

	if int64(r.token) > 0 {
		r.token--
		return true
	}
	// no token ,refuse request
	return false
}
