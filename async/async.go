package async

import (
	"fmt"
	"github.com/trichner/elog"
	"log"
	"sync"
	"sync/atomic"
)

func NewAsyncLogger(e elog.EventLogger) *AsyncEventLogger {
	return &AsyncEventLogger{
		el: e,
	}
}

type AsyncEventLogger struct {
	el      elog.EventLogger
	events  chan *elog.Event
	wg      sync.WaitGroup
	running int32
}

func (a *AsyncEventLogger) Start() {
	wasStarted := atomic.CompareAndSwapInt32(&a.running, 0, 1)
	if !wasStarted {
		panic("already running")
	}

	a.events = make(chan *elog.Event, 1024)

	a.wg.Add(1)

	go func() {
		defer a.wg.Done()
		for {
			select {
			case event, ok := <-a.events:
				if !ok {
					return
				}

				err := a.el.Log(event)
				if err != nil {
					log.Printf("error sending async event: %s", err)
				}
			}
		}
	}()
}

func (a *AsyncEventLogger) Shutdown() {
	wasRunning := atomic.CompareAndSwapInt32(&a.running, 1, 0)
	if !wasRunning {
		return
	}

	close(a.events)
	a.wg.Wait()
}

func (a *AsyncEventLogger) IsRunning() {
	close(a.events)
	a.wg.Wait()
}

func (a *AsyncEventLogger) Log(e *elog.Event) error {
	if atomic.LoadInt32(&a.running) == 0 {
		return fmt.Errorf("async event logger not running")
	}
	a.events <- e
	return nil
}
