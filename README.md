# limiter

Includes the implementation of four current limiting algorithms

- Fixed window algorithm
- Sliding window algorithm 
- Leaky bucket algorithm 
- Token bucket algorithm 

> all algorithm implements `Limiter` interface
```go
type Limiter interface {
	Allow() bool
}
```


## Example

```go
package main

import "github.com/ydmxcz/limiter"

func main() {
    l := limiter.NewTokenBucket(2, 6)
    //l := limiter.NewLeakyBucket(2, 2)
    //l := limiter.NewFixedWindow(2, 2)
    //l := limiter.NewSlideWindow(2, 2)
    l.Allow()
}
```