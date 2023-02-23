package limiter

import (
	"math"
	"sync"
	"time"
)

type LeakBucket struct {
	rate     float64 //每秒速度
	capacity float64

	mutex    sync.Mutex
	lastTime time.Time
	water    float64
}

func NewLeakBucket(rate float64, capacity float64) *LeakBucket {
	return &LeakBucket{
		rate:     rate,
		capacity: capacity,
		lastTime: time.Now(),
		water:    0,
	}
}

func (r *LeakBucket) Allow() bool {
	now := time.Now()

	r.mutex.Lock()
	defer r.mutex.Unlock()

	remainWaterCap := math.Max(0, r.water-(now.Sub(r.lastTime).Seconds()*r.rate))
	r.lastTime = now
	if remainWaterCap+1 < r.capacity {
		r.water = remainWaterCap + 1
		return true
	}
	return false
}

type LeakyBucket struct {
	rate       float64 //固定每秒出水速率
	capacity   float64 //桶的容量
	water      float64 //桶中当前水量
	lastLeakMs int64   //桶上次漏水时间戳 ms

	lock sync.Mutex
}

func (l *LeakyBucket) Allow() bool {
	l.lock.Lock()
	defer l.lock.Unlock()

	now := time.Now().UnixNano() / 1e6
	eclipse := float64((now - l.lastLeakMs)) * l.rate / 1000 //先执行漏水
	l.water = l.water - eclipse                              //计算剩余水量
	l.water = math.Max(0, l.water)                           //桶干了
	l.lastLeakMs = now
	if (l.water + 1) < l.capacity {
		// 尝试加水,并且水还未满
		l.water++
		return true
	} else {
		// 水满，拒绝加水
		return false
	}
}

func NewLeakyBucket(rate float64, capacity float64) *LeakyBucket {
	return &LeakyBucket{
		rate:       rate,
		capacity:   capacity,
		lastLeakMs: time.Now().UnixNano() / 1e6,
		water:      0,
	}
}

func (l *LeakyBucket) Set(r, c float64) {
	l.rate = r
	l.capacity = c
	l.water = 0
	l.lastLeakMs = time.Now().UnixNano() / 1e6
}
