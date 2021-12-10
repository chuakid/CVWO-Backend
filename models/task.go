package models

import (
	"github.com/chuakid/cvwo-backend/db"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ID          int      `json:"ID"`
	Description string   `json:"description"`
	Project     *Project `json:"project"`
	ProjectID   int      `json:"projectid"`
}

func (task *Task) CreateTask() error {
	result := db.DB.Create(task)
	return result.Error
}

func (task *Task) GetTask() error {
	result := db.DB.First(&task)
	return result.Error
}
func (task *Task) DeleteTask() error {
	result := db.DB.Delete(&task)
	return result.Error
}
