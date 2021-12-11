package models

type UserProjects struct {
	UserID    int `gorm:"primaryKey"`
	ProjectID int `gorm:"primaryKey"`
	Role      int `json:"role" gorm:"default:2"`
}
