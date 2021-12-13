package models

import "github.com/chuakid/cvwo-backend/db"

type UserProject struct {
	UserID    int `gorm:"primaryKey"`
	ProjectID int `gorm:"primaryKey"`
	Role      int `json:"role" gorm:"default:2"`
}

func (userproject UserProject) ChangeRole(role int) error {
	result := db.DB.Model(&userproject).Update("role", role)
	return result.Error
}
