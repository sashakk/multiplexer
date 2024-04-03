package limiter

import (
	"sync"
)

type RequestLimiter struct {
	concurrentRequests int32
	maxConcurrentLimit int32
	l                  sync.Mutex
}

func (r *RequestLimiter) Allow() bool {
	r.l.Lock()
	defer r.l.Unlock()
	if r.concurrentRequests < r.maxConcurrentLimit {
		r.concurrentRequests++
		return true
	}

	return false
}

func (r *RequestLimiter) Done() {
	r.l.Lock()
	defer r.l.Unlock()
	r.concurrentRequests--
}

func NewRequestLimiter(maxConcurrentLimit int32) *RequestLimiter {
	return &RequestLimiter{
		maxConcurrentLimit: maxConcurrentLimit,
		concurrentRequests: 0,
	}
}
