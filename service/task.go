package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ahaostudy/calendar_reminder/job/reminder"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ahaostudy/calendar_reminder/dal/mysql"
	"github.com/ahaostudy/calendar_reminder/model"
)

func GetTask(ctx context.Context, id string, userId uint) (*model.Task, error) {
	task, err := model.GetTaskById(mysql.DB, ctx, id, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("task does not exists")
		}
		return nil, fmt.Errorf("get task by id error: %v", err)
	}
	return task, nil
}

func GetTaskList(ctx context.Context, userId uint) ([]*model.Task, error) {
	tasks, err := model.GetTaskList(mysql.DB, ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("get task list error: %v", err)
	}
	return tasks, nil
}

func CreateTask(ctx context.Context, userId uint, title string, t int64) (*model.Task, error) {
	task := &model.Task{
		ID:     uuid.New().String(),
		UserId: userId,
		Title:  title,
		Time:   t,
	}
	err := model.CreateTask(mysql.DB, ctx, task)
	if err != nil {
		return nil, fmt.Errorf("create task failed: %v" + err.Error())
	}
	CheckAndPushCrontab(ctx, task)
	return task, nil
}

func UpdateTask(ctx context.Context, id string, userId uint, title string, t int64) (*model.Task, error) {
	task := &model.Task{
		Title: title,
		Time:  t,
	}
	err := model.UpdateTask(mysql.DB, ctx, id, userId, task)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("task does not exists")
		}
		return nil, fmt.Errorf("update task failed: %v" + err.Error())
	}
	CheckAndPushCrontab(ctx, task)
	return task, nil
}

func DeleteTask(ctx context.Context, id string, userId uint) error {
	err := model.DeleteTask(mysql.DB, ctx, id, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("task does not exists")
		}
		return err
	}
	// try to delete the task from Crontab to prevent error reminders
	reminder.DeleteTask(&model.Task{ID: id})
	return nil
}

func GetTaskListByDate(ctx context.Context, userId uint, date time.Time) ([]*model.Task, error) {
	start, end := GetDateStartAndEnd(date)
	tasks, err := model.GetUserTaskListByTimeRange(mysql.DB, ctx, userId, start.Unix(), end.Unix())
	if err != nil {
		return nil, fmt.Errorf("get task list by date error: %v", err)
	}
	return tasks, nil
}

// CheckAndPushCrontab check if the task missed being pulled by the ReminderService and manually add to it
func CheckAndPushCrontab(ctx context.Context, task *model.Task) {
	taskTime := time.Unix(task.Time, 0)
	// if the task's time is in the current round time, it needs to be manually added to the wait queue
	if reminder.InCurrentRoundTime(taskTime) {
		var err error
		task.User, err = GetUser(ctx, task.UserId)
		if err != nil {
			logrus.Errorf("get user err: %v", err.Error())
			return
		}
		reminder.AddTask(task)
	}
}

func ParseDate(date string) (time.Time, error) {
	const layout = "2006-01-02"
	return time.Parse(layout, date)
}

// GetDateStartAndEnd get the start and last second of the day
func GetDateStartAndEnd(date time.Time) (time.Time, time.Time) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.AddDate(0, 0, 1).Add(time.Second * -1)
	return start, end
}
