package reminder

import (
	"context"
	"github.com/ahaostudy/calendar_reminder/dal/mysql"
	"github.com/ahaostudy/calendar_reminder/job/email"
	"github.com/ahaostudy/calendar_reminder/model"
	"github.com/ahaostudy/calendar_reminder/utils/crontab"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	// time interval of each round
	ti = 10 * time.Second

	// start time of the every round
	// new tasks is loaded every ti, the time range is: [st, st + ti), then st += ti
	st = time.Now()

	gs  = 100  // goroutines
	mws = 1000 // max work channel size

	// global crontab object
	cron *crontab.Crontab

	tasksKey = make(map[string]string)
)

func init() {
	cron = crontab.NewCrontab(gs, mws)
}

// InCurrentRoundTime check if t is within the time of the current round
func InCurrentRoundTime(t time.Time) bool {
	return t.After(st.Add(-1*ti)) && t.Before(st)
}

func RunReminderService() {
	cron.Run()
	// execute once at startup to move back the st pointer
	timedPolling()
	st = time.Now().Add(ti)
	// timed polling to get new tasks
	ticker := time.NewTicker(ti)
	for {
		go timedPolling()
		<-ticker.C
	}
}

// execution entry for each polling
func timedPolling() {
	logrus.Info("timed reminder service running")
	ctx, cancel := context.WithTimeout(context.Background(), ti/2)
	LoadNewTasks(ctx)
	cancel()
	st = st.Add(ti)
}

// LoadNewTasks loads new timed tasks
func LoadNewTasks(ctx context.Context) {
	ed := st.Add(ti)
	tasks, err := model.GetTaskListByTimeRange(mysql.DB, ctx, st.Unix(), ed.Unix()-1)
	logrus.Info("load new task lists:", len(tasks))
	if err != nil && ctx.Err() == nil {
		time.Sleep(time.Millisecond * 500)
		LoadNewTasks(ctx)
	}
	for _, task := range tasks {
		AddTask(task)
	}
}

// AddTask inserts a new task into the scheduler and cleans up if the task has already been added to the scheduler
func AddTask(task *model.Task) {
	logrus.Info("add task to crontab:", task)
	// delete old timed task
	DeleteTask(task)
	taskTime := time.Unix(task.Time, 0)
	cronTask := crontab.NewTask(task.ID, taskTime, task, TaskReminderHandler)
	tasksKey[task.ID] = cronTask.Key()
	cron.AddTask(cronTask)
}

func DeleteTask(task *model.Task) {
	if key, ok := tasksKey[task.ID]; ok {
		cron.DeleteTask(key)
	}
}

// TaskReminderHandler timed task callbacks, which will be called when the time is up
func TaskReminderHandler(id string, data any) {
	logrus.Info("task reminder:", id)
	task, ok := data.(*model.Task)
	if !ok {
		logrus.Error("task reminder handle error: data type is not Task")
		return
	}
	msg := email.Message{Task: task}
	msg.Send()
}
