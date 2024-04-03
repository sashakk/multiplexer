package limiter

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestNewRequestLimiter(t *testing.T) {
	t.Run("Test rate limiter", func(t *testing.T) {
		maxWant := int32(100)
		rt := NewRequestLimiter(maxWant)
		var maxCounter int32

		for i := 0; i < 1000; i++ {
			go func() {
				if rt.Allow() {
					numRequests := atomic.LoadInt32(&rt.concurrentRequests)
					if numRequests > atomic.LoadInt32(&maxCounter) {
						atomic.StoreInt32(&maxCounter, numRequests)
					}
					time.Sleep(time.Microsecond * 100)
					rt.Done()
				}
			}()
		}

		if maxWant < maxCounter {
			t.Errorf("maxCounter greate than maxWant: %d", maxCounter)
		}
	})
}
