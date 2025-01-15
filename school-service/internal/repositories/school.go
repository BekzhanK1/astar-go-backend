package repositories

import (
	"school-service/internal/models"

	"gorm.io/gorm"
)

// SchoolRepository will handle all CRUD operations for School
type SchoolRepository struct {
	db *gorm.DB
}

// NewSchoolRepository is a constructor for SchoolRepository
func NewSchoolRepository(db *gorm.DB) *SchoolRepository {
	return &SchoolRepository{db}
}

// Create will create a new School
func (r *SchoolRepository) Create(school *models.School) error {
	return r.db.Create(school).Error
}

// GetByID will fetch a School by its ID
func (r *SchoolRepository) GetByID(id uint) (*models.School, error) {
	var school models.School
	if err := r.db.First(&school, id).Error; err != nil {
		return nil, err
	}
	return &school, nil
}

// Update will update a School
func (r *SchoolRepository) Update(school *models.School) error {
	return r.db.Save(school).Error
}

// Delete will delete a School
func (r *SchoolRepository) Delete(school *models.School) error {
	return r.db.Delete(school).Error
}
