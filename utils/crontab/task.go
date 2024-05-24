package crontab

import (
	"github.com/ahaostudy/calendar_reminder/utils/crontab/cron_pool"
	"github.com/google/uuid"
	"time"
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

func NewTask(t time.Time, d any, h HandlerFunc) *Task {
	return &Task{ID: uuid.NewString(), Time: t, Data: d, Handler: h}
}

// Key generate a unique Key for the task, which is prefixed with a time string and can be sorted by time
func (t *Task) Key() string {
	return formatTime(t.Time) + t.ID
}

func (t *Task) Work(c *Crontab) *cron_pool.Work {
	if t.work == nil {
		t.work = cron_pool.NewWork(t.ID, t.Time, t.Data, func(id string, data any) {
			defer c.doneTask(t)
			t.Handler(t.Data)
		})
	}
	return t.work
}

type HandlerFunc func(data any)
