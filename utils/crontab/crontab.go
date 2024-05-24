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
			// push to waiting tasks and reset the task
			work := lastRunningTask.Work(crontab)
			crontab.Lock()
			crontab.PopRunningTasks(lastRunningTask.Key())
			crontab.PushWaitingTasks(lastRunningTask)
			work.Reset()
			crontab.Unlock()
			// run the current, earlier task
			crontab.runTask(task)
			return
		}
	}

	// otherwise push the waiting queue
	crontab.PushWaitingTasks(task)
}

func (crontab *Crontab) DeleteTask(key string) bool {
	if key == "" {
		return false
	}
	task, ok := crontab.GetRunningTask(key)
	if ok {
		crontab.Lock()
		defer crontab.Unlock()
		work := task.Work(crontab)
		crontab.PopRunningTasks(key)
		// interrupt the execution of the task
		return work.Interrupt()
	}
	task, ok = crontab.GetWaitingTask(key)
	if ok {
		crontab.Lock()
		defer crontab.Unlock()
		work := task.Work(crontab)
		crontab.PopWaitingTasks(key)
		return work.Interrupt()
	}
	return false
}

// run the specified task without focusing on the scheduling logic
func (crontab *Crontab) runTask(task *Task) {
	crontab.Lock()
	defer crontab.Unlock()
	crontab.PushRunningTasks(task)
	crontab.cron.AddWork(task.Work(crontab))
}

// clear from queue after task done
func (crontab *Crontab) doneTask(task *Task) {
	crontab.PopRunningTasks(task.Key())
	crontab.moveTask()
}

func (crontab *Crontab) moveTask() {
	task, ok := crontab.PopWaitingTasks("")
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

func (crontab *Crontab) PopRunningTasks(key string) (task *Task, ok bool) {
	crontab.RMX.Lock()
	if key == "" {
		_, task, ok = crontab.RunningTasks.PopMax()
	} else {
		task, ok = crontab.RunningTasks.Delete(key)
	}
	crontab.RMX.Unlock()
	return task, ok
}

func (crontab *Crontab) GetRunningTask(key string) (task *Task, ok bool) {
	crontab.RMX.Lock()
	task, ok = crontab.RunningTasks.Get(key)
	crontab.RMX.Unlock()
	return
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

func (crontab *Crontab) GetWaitingTask(key string) (task *Task, ok bool) {
	crontab.WMX.Lock()
	task, ok = crontab.WaitingTasks.Get(key)
	crontab.WMX.Unlock()
	return
}

func (crontab *Crontab) PopWaitingTasks(key string) (task *Task, ok bool) {
	crontab.WMX.Lock()
	if key == "" {
		_, task, ok = crontab.WaitingTasks.PopMin()
	} else {
		task, ok = crontab.WaitingTasks.Delete(key)
	}
	crontab.WMX.Unlock()
	return task, ok
}
