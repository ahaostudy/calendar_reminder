package crontab

import (
	"github.com/ahaostudy/calendar_reminder/utils/crontab/cron_pool"
	"github.com/google/uuid"
	"github.com/tidwall/btree"
	//"log"
	"sync"
	"time"
)

// Task timed task
type Task struct {
	ID      string
	Time    time.Time
	Data    any
	Handler HandlerFunc
}

func NewTask(t time.Time, d any, h HandlerFunc) *Task {
	return &Task{ID: uuid.NewString(), Time: t, Data: d, Handler: h}
}

// Key generate a unique Key for the task, which is prefixed with a time string and can be sorted by time
func (t *Task) Key() string {
	return formatTime(t.Time) + t.ID
}

// RunningTask contains the Work object, which can be used to interrupt the timed task when necessary
// and give up resources to earlier tasks
type RunningTask struct {
	Task
	Work *cron_pool.Work
}

type HandlerFunc func(data any)

// Crontab a timed task scheduler that automatically reorders tasks based on their execution time
type Crontab struct {
	cron         *cron_pool.CronPoll
	RunningTasks btree.Map[string, *RunningTask]
	WaitingTasks btree.Map[string, *Task]
	RMX          sync.Mutex
	WMX          sync.Mutex
}

func NewCrontab(gs, mws int) *Crontab {
	return &Crontab{cron: cron_pool.NewCronPoll(gs, mws)}
}

func (crontab *Crontab) Run() {
	crontab.cron.Run()
}

func (crontab *Crontab) AddTask(task *Task) {
	// when the number of running tasks is less than the number of goroutines, it can be run directly
	if crontab.RunningTasks.Len() < crontab.cron.Gos {
		crontab.runTask(task)
		return
	}
	// when the execution time of the task is earlier than the execution time of the last one in the running queue
	// its resources are seized
	_, lastRunningTask, ok := crontab.RunningTasks.Max()
	if ok && task.Time.Before(lastRunningTask.Time) {
		crontab.PopRunningTasks(&lastRunningTask.Task)
		lastRunningTask.Work.Interrupt()
		crontab.PushWaitingTasks(&lastRunningTask.Task)
		crontab.runTask(task)
		return
	}

	// otherwise push the waiting queue
	crontab.PushWaitingTasks(task)
}

// run the specified task without focusing on the scheduling logic
func (crontab *Crontab) runTask(task *Task) {
	work := crontab.work(task)
	crontab.PushRunningTasks(&RunningTask{Task: *task, Work: work})
	crontab.cron.AddWork(work)
}

func (crontab *Crontab) work(task *Task) (work *cron_pool.Work) {
	d := task.Time.Sub(time.Now())
	if d < 0 {
		d = 0
	}
	work = cron_pool.NewWork(task.ID, d, task.Data, func(id string, data any) {
		defer crontab.doneTask(task)
		task.Handler(task.Data)
	})
	return work
}

// clear from queue after task done
func (crontab *Crontab) doneTask(task *Task) {
	crontab.PopRunningTasks(task)
	crontab.moveTask()
}

func (crontab *Crontab) moveTask() {
	task, ok := crontab.PopWaitingTasks()
	if !ok {
		return
	}
	crontab.AddTask(task)
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func (crontab *Crontab) PushRunningTasks(task *RunningTask) {
	crontab.RMX.Lock()
	crontab.RunningTasks.Set(task.Key(), task)
	crontab.RMX.Unlock()
}

func (crontab *Crontab) PopRunningTasks(task *Task) {
	crontab.RMX.Lock()
	if task == nil {
		crontab.RunningTasks.PopMax()
	} else {
		crontab.RunningTasks.Delete(task.Key())
	}
	crontab.RMX.Unlock()
}

func (crontab *Crontab) PushWaitingTasks(task *Task) {
	crontab.WMX.Lock()
	crontab.WaitingTasks.Set(task.Key(), task)
	crontab.WMX.Unlock()
}

func (crontab *Crontab) PopWaitingTasks() (*Task, bool) {
	crontab.WMX.Lock()
	defer crontab.WMX.Unlock()
	_, task, ok := crontab.WaitingTasks.PopMin()
	return task, ok
}
