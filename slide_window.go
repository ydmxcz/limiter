package limiter

import (
	"sync"
	"time"
)

type timeSlot struct {
	startTime time.Time
	count     int32
}

type SlideWindow struct {
	rate           int32
	windowDuration time.Duration
	slotDuration   time.Duration

	mu       sync.Mutex
	slotList []*timeSlot
}

type SlideWindowDuration struct {
	windowDuration time.Duration
	slotDuration   time.Duration
}

func NewSlideWindow(rate int32, slideWindowDuration SlideWindowDuration) *SlideWindow {
	return &SlideWindow{
		rate:           rate,
		windowDuration: slideWindowDuration.windowDuration,
		slotDuration:   slideWindowDuration.slotDuration,
		slotList:       make([]*timeSlot, 0, slideWindowDuration.windowDuration/slideWindowDuration.slotDuration),
	}
}

func (r *SlideWindow) Allow() bool {
	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	discardSlotIdx := -1
	// discard the slot before now
	for i := range r.slotList {
		slot := r.slotList[i]
		if slot.startTime.Add(r.slotDuration).After(now) {
			break
		}
		discardSlotIdx = i
	}
	if discardSlotIdx > -1 {
		r.slotList = r.slotList[discardSlotIdx+1:]
	}

	var reqCount int32 = 0
	for i := range r.slotList {
		reqCount += r.slotList[i].count
	}
	if reqCount >= r.rate {
		return false
	}

	if len(r.slotList) > 0 {
		r.slotList[len(r.slotList)-1].count++
	} else {
		r.slotList = append(r.slotList, r.newTimeSlot(now))
	}
	return true
}

func (r *SlideWindow) newTimeSlot(startTime time.Time) *timeSlot {
	return &timeSlot{startTime: startTime, count: 1}
}
