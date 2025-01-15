package models

import "time"

type Group struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	FlowID    uint      `gorm:"not null" json:"flow_id"`
	Flow      Flow      `gorm:"foreignKey:FlowID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"flow"`
	LevelName string    `gorm:"not null" json:"level_name"`
	Level     Level     `gorm:"foreignKey:LevelName;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"level"`
	TeacherID uint      `gorm:"not null" json:"teacher_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
