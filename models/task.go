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

type APITask struct {
	ID          int    `json:"ID"`
	Description string `json:"description"`
	ProjectID   int    `json:"-"`
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
func (task *Task) EditTask(description string) error {
	result := db.DB.Model(&task).Update("description", description)
	return result.Error
}
