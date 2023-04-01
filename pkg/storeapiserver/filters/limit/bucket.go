package limit

import (
	"k8s.io/klog/v2"
	"sync"
	"time"
)

type Bucket struct {
	cap      int64
	tokens   int64
	lock     sync.Mutex
	rate     int64
	lastTime int64
}


const (
	DefaultCap  = 5
	DefaultRate = 1
)

func NewBucket(cap, rate int64) *Bucket {

	if cap < 0 || rate < 0 {
		klog.Error("cap and rate param is wrong")
		cap = DefaultCap
		rate = DefaultRate
	}

	b := &Bucket{
		cap:    cap,
		tokens: cap,
		rate:   rate,
	}

	return b

}

// IsAccept 是否接受请求
func (b *Bucket) IsAccept() bool {
	b.lock.Lock()
	defer b.lock.Unlock()

	now := time.Now().Unix()
	b.tokens = b.tokens + (now - b.lastTime) * b.rate

	if b.tokens >= b.cap {
		b.tokens = b.cap
	}

	b.lastTime = now
	if b.tokens > 0 {
		b.tokens--
		return true
	}

	return false

}
