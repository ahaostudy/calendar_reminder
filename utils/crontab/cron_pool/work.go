package cron_pool

import (
	"time"
)

type Work struct {
	id       string
	duration time.Duration
	data     any
	handler  HandlerFunc
	status   string

	worker *Worker
}

func (w *Work) ID() string {
	return w.id
}

const (
	WorkStatusReady   = "ready"
	WorkStatusWaiting = "waiting"
	WorkStatusSuccess = "success"
	WorkStatusStopped = "stopped"
)

type HandlerFunc func(string, any)

func NewWork(id string, duration time.Duration, data any, handler HandlerFunc) *Work {
	return &Work{
		id:       id,
		duration: duration,
		data:     data,
		handler:  handler,
		status:   WorkStatusReady,
	}
}

func (w *Work) Wait() *time.Timer {
	w.status = WorkStatusWaiting
	timer := time.NewTimer(w.duration)
	return timer
}

func (w *Work) Interrupt() {
	if w.status == WorkStatusWaiting {
		if w.worker != nil && w.worker.InterruptChan != nil {
			w.worker.InterruptChan <- struct{}{}
		}
	}
	w.status = WorkStatusStopped
}

func (w *Work) Run() {
	if w.status == WorkStatusWaiting {
		w.handler(w.id, w.data)
		w.status = WorkStatusSuccess
	}
}
