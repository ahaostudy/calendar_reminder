package model

import (
	"context"

	"gorm.io/gorm"
)

type Task struct {
	BaseModel
	// ID generate 36-bit string using uuid
	ID     string `gorm:"primary_key;auto_increment;type:varchar(36)" json:"id"`
	UserId uint   `gorm:"index:idx_user_date,priority:1" json:"user_id"`
	Title  string `gorm:"type:varchar(127)" json:"title"`
	// Time uses int64 to store timestamp, which can be easily indexed and easily used for JSON serialization
	// Create a union index using UserId and Time to save index space and improve union query efficiency
	Time int64 `gorm:"index:idx_user_date,priority:2" json:"time"` // timestamp
}

func GetTaskById(db *gorm.DB, ctx context.Context, id string, userId uint) (task *Task, err error) {
	err = db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userId).First(&task).Error
	return
}

func GetTaskList(db *gorm.DB, ctx context.Context, userId uint) (tasks []*Task, err error) {
	err = db.WithContext(ctx).Where("user_id = ?", userId).Find(&tasks).Error
	return
}

func CreateTask(db *gorm.DB, ctx context.Context, task *Task) error {
	return db.WithContext(ctx).Create(task).Error
}

func UpdateTask(db *gorm.DB, ctx context.Context, id string, userId uint, task *Task) error {
	result := db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userId).Updates(task).Scan(task)
	// `RowsAffected == 0` indicates that the record does not exist
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func DeleteTask(db *gorm.DB, ctx context.Context, id string, userId uint) error {
	result := db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userId).Delete(&Task{})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func GetTaskListByTimeRange(db *gorm.DB, ctx context.Context, userId uint, start int64, end int64) (tasks []*Task, err error) {
	err = db.WithContext(ctx).Where("user_id = ? AND time BETWEEN ? AND ?", userId, start, end).Find(&tasks).Error
	return
}
