package crontab

import (
	"time"

	"github.com/ahaostudy/calendar_reminder/utils/crontab/cron_pool"
)

// Task timed task, which contains Work objects that can be used to interrupt a timed task if necessary
// yielding resources to an earlier task.
type Task struct {
	ID      string
	Time    time.Time
	Data    any
	Handler HandlerFunc

	work *cron_pool.Work
}

func NewTask(id string, t time.Time, data any, handler HandlerFunc) *Task {
	return &Task{ID: id, Time: t, Data: data, Handler: handler}
}

// Key generate a unique Key for the task, which is prefixed with a time string and can be sorted by time
func (t *Task) Key() string {
	return formatTime(t.Time) + t.ID
}

func (t *Task) Work(c *Crontab) *cron_pool.Work {
	if t.work == nil {
		t.work = cron_pool.NewWork(t.ID, t.Time, t.Data, func(id string, data any) {
			defer c.doneTask(t)
			t.Handler(t.ID, t.Data)
		})
	}
	return t.work
}

type HandlerFunc func(id string, data any)
