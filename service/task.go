package service

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	return nil
}

func GetTaskListByDate(ctx context.Context, userId uint, date time.Time) ([]*model.Task, error) {
	start, end := GetDateStartAndEnd(date)
	tasks, err := model.GetTaskListByTimeRange(mysql.DB, ctx, userId, start.Unix(), end.Unix())
	if err != nil {
		return nil, fmt.Errorf("get task list by date error: %v", err)
	}
	return tasks, nil
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
