package timer

import (
	"context"
	"sync"
	"sync/atomic"
)

type TimerController struct {
	timers          []*Timer
	timersNum       int32
	stoppedTimerNum int32

	ctx    context.Context
	cancel context.CancelFunc
	mutex  sync.Mutex
}

func NewTimerController(timers ...*Timer) *TimerController {
	tc := &TimerController{
		timers:          timers,
		timersNum:       int32(len(timers)),
		stoppedTimerNum: int32(0),
	}
	return tc
}

func (p *TimerController) Start() {
	go func(p *TimerController) {
		for {
			select {
			case <-p.ctx.Done():
				return
			default:
				if p.stoppedTimerNum != p.timersNum {
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
		}
	}(p)
}

func (p *TimerController) AddTimer(timer *Timer) {
	p.mutex.Lock()
	p.timers = append(p.timers, timer)
	p.mutex.Unlock()
}

func (p *TimerController) Stop() {
	p.mutex.Lock()
	for _, t := range p.timers {
		if !t.Stopped() {
			t.Stop()
		}
	}
	p.mutex.Unlock()
	p.cancel()
}

func (p *TimerController) TimersNum() int32 {
	return atomic.LoadInt32(&p.timersNum)
}

func (p *TimerController) StoppedTimersNum() int32 {
	return atomic.LoadInt32(&p.stoppedTimerNum)
}
