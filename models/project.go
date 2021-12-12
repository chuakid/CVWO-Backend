package models

import (
	"strconv"

	"github.com/chuakid/cvwo-backend/db"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Users []*User `gorm:"many2many:UserProjects"`
	Tasks []*Task
}

type APIProjectSummary struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type APIProjectDetailed struct {
	ID    int        `json:"id"`
	Name  string     `json:"name"`
	Users []*APIUser `gorm:"many2many:UserProjects"`
	Tasks []*APITask
}

func (project *Project) GetProjectDetails() (*APIProjectDetailed, error) {
	var apiproject APIProjectDetailed
	result := db.DB.Model(&project).First(&apiproject, project.ID)
	if result.Error != nil {
		return nil, result.Error
	}
	//Manually do preloading, GORM doesn't work well with smart select + preloads
	err := db.DB.Model(&project).Association("Users").Find(&apiproject.Users)
	if err != nil {
		return nil, err
	}
	err = db.DB.Model(&project).Association("Tasks").Find(&apiproject.Tasks)
	if err != nil {
		return nil, err
	}

	return &apiproject, result.Error
}

func (project *Project) CreateProject(userid string) error {
	id, err := strconv.Atoi(userid)
	if err != nil {
		return err
	}
	user := User{
		ID: id,
	}
	project.Users = append(project.Users, &user)
	result := db.DB.Create(&project)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (project *Project) DeleteProject() error {
	return db.DB.Delete(&project).Error
}

func (project *Project) GetUsers() ([]User, error) {
	var users []User
	err := db.DB.Model(&project).Association("Users").Find(&users)
	return users, err
}

func (project *Project) RenameProject(name string) error {
	result := db.DB.Model(&project).Update("name", name)
	return result.Error
}

func (project *Project) AddUser(user *User) error {
	err := db.DB.Model(&project).Association("Users").Append(user)
	return err
}
