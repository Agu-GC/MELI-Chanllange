package repositories

import (
	"github.com/Agu-GC/MELI-Challenge/api/internal/domain"
	"github.com/Agu-GC/MELI-Challenge/api/pkg"

	"gorm.io/gorm"
)

type ClassificationRepositoryInterface interface {
	Create(classification *domain.Classification) error
	GetAll() ([]*domain.Classification, error)
}

type classificationRepository struct {
	db     *gorm.DB
	logger pkg.Logger
}

func NewClassificationRepository(dbConnection *gorm.DB, logger pkg.Logger) ClassificationRepositoryInterface {
	return &classificationRepository{db: dbConnection, logger: logger}
}

func (r *classificationRepository) Create(classification *domain.Classification) error {
	return r.db.Create(classification).Error
}

func (r *classificationRepository) GetAll() ([]*domain.Classification, error) {
	var classifications []*domain.Classification
	result := r.db.Find(&classifications)
	return classifications, result.Error
}
