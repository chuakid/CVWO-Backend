package models

import (
	"log"
	"strconv"

	"github.com/chuakid/cvwo-backend/db"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ID          int      `json:"ID"`
	Description string   `json:"description"`
	Project     *Project `json:"projects"`
	ProjectID   int
}

func (task *Task) CreateTask(projectid string) (int, error) {
	//Check if user is allowed to access this project
	var projectUsers []User
	projectidInt, err := strconv.Atoi(projectid)
	if err != nil {
		log.Print("Error getting project for task creation:", err)
		return -1, err
	}
	db.DB.Model(Project{ID: projectidInt}).Association("Users").Find(&projectUsers)
	log.Print(projectUsers)

	return 0, nil
}
