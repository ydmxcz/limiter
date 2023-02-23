package limiter

import (
	"sync"
	"time"
)

type FixedWindow struct {
	duration    time.Duration
	rate        int32
	mu          sync.Mutex
	count       int32
	lastGetTime time.Time
}

func NewFixedWindow(duration time.Duration, rate int32) *FixedWindow {
	return &FixedWindow{
		duration:    duration,
		rate:        rate,
		lastGetTime: time.Now(),
	}
}

func (r *FixedWindow) Allow() bool {
	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	if now.Sub(r.lastGetTime) > r.duration {
		r.reset(now)
	}
	if r.count >= r.rate {
		return false
	}
	r.count++
	return true
}

func (r *FixedWindow) reset(getTime time.Time) {
	r.lastGetTime = getTime
	r.count = 0
}
