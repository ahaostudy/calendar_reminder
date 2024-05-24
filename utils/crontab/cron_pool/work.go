package cron_pool

import (
	"sync"
	"time"
)

type Work struct {
	sync.Mutex

	id      string
	time    time.Time
	data    any
	handler HandlerFunc
	status  int

	worker *Worker
}

func (w *Work) ID() string {
	return w.id
}

const (
	WorkStatusReady = iota
	WorkStatusWaiting
	WorkStatusSuccess
	WorkStatusStopped
)

type HandlerFunc func(string, any)

func NewWork(id string, t time.Time, data any, handler HandlerFunc) *Work {
	return &Work{
		id:      id,
		time:    t,
		data:    data,
		handler: handler,
		status:  WorkStatusReady,
	}
}

func (w *Work) Wait() (*time.Timer, bool) {
	if w.status != WorkStatusReady {
		return nil, false
	}
	w.status = WorkStatusWaiting
	timer := time.NewTimer(time.Until(w.time))
	return timer, true
}

func (w *Work) Interrupt() bool {
	if w.status == WorkStatusSuccess || w.status == WorkStatusStopped {
		return false
	}
	if w.status == WorkStatusWaiting {
		w.worker.InterruptChan <- struct{}{}
	}
	w.status = WorkStatusStopped
	return true
}

func (w *Work) Run() {
	if w.status == WorkStatusWaiting {
		w.handler(w.id, w.data)
		w.status = WorkStatusSuccess
	}
}

func (w *Work) Reset() {
	w.status = WorkStatusReady
}
