package models

import (
	"time"
)

type Staff struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SchoolID  uint      `gorm:"not null" json:"school_id"`
	School    School    `gorm:"foreignKey:SchoolID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"school"`
	FirstName string    `gorm:"type:varchar(255);not null" json:"first_name"`
	LastName  string    `gorm:"type:varchar(255);not null" json:"last_name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
