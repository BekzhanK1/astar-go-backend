package user

import "gorm.io/gorm"

// User model represents a user in the system.
type User struct {
	ID        uint           `gorm:"primaryKey"`
	Email     string         `gorm:"unique;not null"`
	FirstName string         `gorm:"size:100;not null"`
	LastName  string         `gorm:"size:100;not null"`
	Role      string         `gorm:"size:50;not null"`
	Password  string         `gorm:"size:255;not null"`
	CreatedAt int64          `gorm:"autoCreateTime"`
	UpdatedAt int64          `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

const (
	RoleSuperAdmin  = "superadmin"
	RoleSchoolAdmin = "school_admin"
	RoleAdvisor     = "advisor"
	RoleTeacher     = "teacher"
)
