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
	Tasks []*Task
}

func (project *Project) GetProjectDetails() (*APIProjectDetailed, error) {
	var apiproject APIProjectDetailed
	var summary APIProjectSummary
	result := db.DB.Model(&project).First(&summary, project.ID)

	if result.Error != nil {
		return nil, result.Error
	}
	//Do this or GORM will throw errors about relations due to custom structs
	apiproject.ID = summary.ID
	apiproject.Name = summary.Name

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

	//set role
	userproject := UserProject{
		UserID:    id,
		ProjectID: project.ID,
	}
	err = userproject.ChangeRole(1)
	if err != nil {
		return err
	}

	return nil
}

func (project *Project) DeleteProject() error {
	return db.DB.Select("Tasks").Delete(&project).Error
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

func (project *Project) GetUsersWithRoles() ([]UserRole, error) {
	rows, err := db.DB.Raw(`SELECT users.username, user_projects.role 
						FROM users
						JOIN user_projects ON user_projects.user_id = users.id 
						WHERE user_projects.project_id = ?`, project.ID).Rows()
	var userroles []UserRole
	for rows.Next() {
		userrole := UserRole{}
		db.DB.ScanRows(rows, &userrole)
		userroles = append(userroles, userrole)
	}

	return userroles, err
}
