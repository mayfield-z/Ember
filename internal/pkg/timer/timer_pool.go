package timer

import (
	"sync"
	"sync/atomic"
)

type TimerPool struct {
	timers          []*Timer
	timersNum       int32
	stoppedTimerNum int32

	mutex sync.Mutex
}

func NewTimerPool(timers ...*Timer) *TimerPool {
	th := &TimerPool{
		timers:          timers,
		timersNum:       int32(len(timers)),
		stoppedTimerNum: int32(0),
	}
	return th
}

func (p *TimerPool) Start() {
	go func(p *TimerPool) {
		if p.TimersNum() != p.StoppedTimersNum() {
			for {
				p.mutex.Lock()
				for i, t := range p.timers {
					if !t.stopped {
						select {
						case <-t.ticker.C:
							atomic.AddInt32(&t.expireTimes, 1)
							//TODO: must atomic?
							if t.ExpireTimes() > t.MaxRetryTimes() {
								t.cancelFunc()
								t.Stop()
								p.timers = append(p.timers[:i], p.timers[i+1:]...)
							} else {
								t.expiredFunc(t.ExpireTimes())
							}
						}
					}
				}
				p.mutex.Unlock()
			}
		}
	}(p)
}

func (p *TimerPool) AddTimer(timer *Timer) {
	p.mutex.Lock()
	p.timers = append(p.timers, timer)
	p.mutex.Unlock()
}

func (p *TimerPool) StopAll() {
	for _, t := range p.timers {
		if !t.Stopped() {
			t.Stop()
		}
	}
}

func (p *TimerPool) TimersNum() int32 {
	return atomic.LoadInt32(&p.timersNum)
}

func (p *TimerPool) StoppedTimersNum() int32 {
	return atomic.LoadInt32(&p.stoppedTimerNum)
}
