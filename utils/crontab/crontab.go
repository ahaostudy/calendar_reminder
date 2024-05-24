package crontab

import (
	"github.com/ahaostudy/calendar_reminder/utils/crontab/cron_pool"
	"github.com/tidwall/btree"
	"sync"
	"time"
)

// Crontab a timed task scheduler that automatically reorders tasks based on their execution time
type Crontab struct {
	sync.Mutex
	cron         *cron_pool.CronPoll
	RunningTasks btree.Map[string, *Task]
	WaitingTasks btree.Map[string, *Task]
	RMX          sync.RWMutex
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
	if crontab.RunningTasksLen() < crontab.cron.Gos {
		crontab.runTask(task)
		return
	}
	// when the execution time of the task is earlier than the execution time of the last one in the running queue
	// its resources are seized
	lastRunningTask, ok := crontab.RunningTasksMax()
	if ok && task.Time.Before(lastRunningTask.Time) {
		if lastRunningTask.Work(crontab).Interrupt() {
			crontab.PopRunningTasks(lastRunningTask)
			crontab.PushWaitingTasks(lastRunningTask)
			crontab.runTask(task)
			return
		}
	}

	// otherwise push the waiting queue
	crontab.PushWaitingTasks(task)
}

// run the specified task without focusing on the scheduling logic
func (crontab *Crontab) runTask(task *Task) {
	crontab.PushRunningTasks(task)
	work := task.Work(crontab)
	work.Lock()
	defer work.Unlock()
	work.Reset()
	crontab.cron.AddWork(work)
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

func (crontab *Crontab) PushRunningTasks(task *Task) {
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

func (crontab *Crontab) RunningTasksMax() (*Task, bool) {
	crontab.RMX.RLock()
	_, task, ok := crontab.RunningTasks.Max()
	crontab.RMX.RUnlock()
	return task, ok
}

func (crontab *Crontab) RunningTasksLen() int {
	crontab.RMX.RLock()
	length := crontab.RunningTasks.Len()
	crontab.RMX.RUnlock()
	return length
}

func (crontab *Crontab) PushWaitingTasks(task *Task) {
	crontab.WMX.Lock()
	crontab.WaitingTasks.Set(task.Key(), task)
	crontab.WMX.Unlock()
}

func (crontab *Crontab) PopWaitingTasks() (*Task, bool) {
	crontab.WMX.Lock()
	_, task, ok := crontab.WaitingTasks.PopMin()
	crontab.WMX.Unlock()
	return task, ok
}
