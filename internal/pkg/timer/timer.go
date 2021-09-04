package timer

import (
	"sync/atomic"
	"time"
)

type Timer struct {
	ticker        *time.Ticker
	expireTimes   int32
	maxRetryTimes int32
	expiredFunc   func(expireTimes int32)
	cancelFunc    func()
	stopped       bool
	done          chan bool
}

func NewTimer(d time.Duration, maxRetryTimes int32, expiredFunc func(expireTimes int32), cancelFunc func()) *Timer {
	return &Timer{
		ticker:        time.NewTicker(d),
		expireTimes:   0,
		maxRetryTimes: maxRetryTimes,
		expiredFunc:   expiredFunc,
		cancelFunc:    cancelFunc,
		stopped:       false,
		done:          make(chan bool, 1),
	}
}

func (t *Timer) ExpireTimes() int32 {
	return atomic.LoadInt32(&t.expireTimes)
}

func (t *Timer) MaxRetryTimes() int32 {
	return atomic.LoadInt32(&t.maxRetryTimes)
}

func (t *Timer) Stop() {
	t.done <- true
	close(t.done)
	t.stopped = true
}

func (t *Timer) Stopped() bool {
	return t.stopped
}
