package models

import (
	"github.com/chuakid/cvwo-backend/db"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model  `json:"-"`
	ID          int      `json:"ID"`
	Description string   `json:"description"`
	Project     *Project `json:"-"`
	ProjectID   int      `json:"projectid"`
	Completed   bool     `json:"completed" gorm:"default:false"`
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
