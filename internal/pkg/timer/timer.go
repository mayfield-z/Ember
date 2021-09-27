package timer

import (
	"sync"
	"sync/atomic"
	"time"
)

var timerPool = sync.Pool{
	New: func() interface{} { return new(Timer) },
}

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
	t := timerPool.Get().(*Timer)
	t.maxRetryTimes = maxRetryTimes
	t.expiredFunc = expiredFunc
	t.cancelFunc = cancelFunc
	t.done = make(chan bool, 1)
	t.stopped = false
	t.ticker = time.NewTicker(d)
	return t
}

func (t *Timer) ExpireTimes() int32 {
	return atomic.LoadInt32(&t.expireTimes)
}

func (t *Timer) MaxRetryTimes() int32 {
	return atomic.LoadInt32(&t.maxRetryTimes)
}

func (t *Timer) Stop() {
	t.stopped = true
	t.done <- true
	close(t.done)
	timerPool.Put(t)
}

func (t *Timer) Stopped() bool {
	return t.stopped
}
