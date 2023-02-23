package limiter_test

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ydmxcz/limiter"
)

func TestTokenBucket(t *testing.T) {
	tblr := limiter.NewTokenBucket(2, 6)
	time.Sleep(time.Second)
	fmt.Println(tblr.Allow())
	fmt.Println(tblr.Allow())
	fmt.Println(tblr.Allow())
	fmt.Println(tblr.Allow())
	time.Sleep(time.Second)
	fmt.Println(tblr.Allow())
	fmt.Println(tblr.Allow())
	fmt.Println(tblr.Allow())
	time.Sleep(time.Second)
	fmt.Println(tblr.Allow())
	fmt.Println(tblr.Allow())
	fmt.Println(tblr.Allow())
	time.Sleep(time.Second)
	fmt.Println(tblr.Allow())
	fmt.Println(tblr.Allow())
	fmt.Println(tblr.Allow())
	time.Sleep(time.Second)
	fmt.Println(tblr.Allow())
	fmt.Println(tblr.Allow())
	fmt.Println(tblr.Allow())
	time.Sleep(time.Second * 2)
	fmt.Println(tblr.Allow())
	fmt.Println(tblr.Allow())
	fmt.Println(tblr.Allow())
	fmt.Println(tblr.Allow())
	fmt.Println(tblr.Allow())
	fmt.Println(tblr.Allow())
	fmt.Println("=========")
	fmt.Println(tblr.Allow())

	time.Sleep(time.Second * 3)
	fmt.Println(tblr.Allow())
	fmt.Println(tblr.Allow())
	fmt.Println(tblr.Allow())
	fmt.Println(tblr.Allow())
	fmt.Println(tblr.Allow())
	fmt.Println(tblr.Allow())
	fmt.Println("=========")
	fmt.Println(tblr.Allow())
	// fmt.Println(tblr.Acquire())
	// fmt.Println(tblr.Acquire())
	// fmt.Println(tblr.Acquire())

}
func TestParallel(t *testing.T) {

	// tblr := limiter.NewTokenBucket(2, 6)
	tblr := limiter.NewLeakyBucket(2, 2)
	// time.Sleep(time.Second * 4)

	ctx, cf := context.WithCancel(context.Background())
	var num int64 = 0
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					wg.Done()
					return
				default:
					time.Sleep(time.Second)
					if tblr.Allow() {
						fmt.Println("AAAA")
						atomic.AddInt64(&num, 1)
					}
				}
			}
		}()
	}
	time.Sleep(time.Second * 10)
	cf()
	wg.Wait()
	fmt.Println(num)
}

func TestTime(t *testing.T) {
	n := time.Now()
	fmt.Println(n.UnixMilli())
	fmt.Println(n.UnixNano())
	na := n.UnixNano()
	mi := n.UnixMilli()
	if na > mi {
		fmt.Println("na")
	}
}
