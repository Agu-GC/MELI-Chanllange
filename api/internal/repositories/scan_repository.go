package repositories

import (
	"github.com/Agu-GC/MELI-Challenge/api/internal/domain"
	"github.com/Agu-GC/MELI-Challenge/api/pkg"
	"gorm.io/gorm"
)

type ScanRepositoryInterface interface {
	Create(*domain.Scan) error
	GetByID(id uint) (*domain.Scan, error)
}

type scanRepository struct {
	dbConn *gorm.DB
	logger pkg.Logger
}

func NewScannRepository(dbConnection *gorm.DB, logger pkg.Logger) ScanRepositoryInterface {
	return &scanRepository{dbConn: dbConnection, logger: logger}
}

func (sr *scanRepository) Create(scan *domain.Scan) error {
	return sr.dbConn.Create(scan).Error
}

func (sr *scanRepository) GetByID(id uint) (*domain.Scan, error) {
	var scan domain.Scan
	err := sr.dbConn.Preload("Tables.Columns.Classification").Preload("Database").First(&scan, id).Error
	return &scan, err
}
