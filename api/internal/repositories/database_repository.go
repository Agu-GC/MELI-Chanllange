package repositories

import (
	"github.com/Agu-GC/MELI-Challenge/api/internal/domain"
	"github.com/Agu-GC/MELI-Challenge/api/pkg"
	"gorm.io/gorm"
)

type DatabaseRepositoryInterface interface {
	Create(dbConn *domain.Database) error
	GetByID(id uint) (*domain.Database, error)
	GetWithLastScanInfo(id uint) (*domain.Database, error)
}

type databaseRepository struct {
	db     *gorm.DB
	logger pkg.Logger
}

func NewDatabaseRepository(dbConnection *gorm.DB, logger pkg.Logger) DatabaseRepositoryInterface {
	return &databaseRepository{db: dbConnection, logger: logger}
}

func (r *databaseRepository) Create(database *domain.Database) error {
	return r.db.Create(database).Error
}

func (r *databaseRepository) GetByID(id uint) (*domain.Database, error) {
	var dbConn domain.Database
	err := r.db.First(&dbConn, id).Error
	return &dbConn, err
}

func (r *databaseRepository) GetWithLastScanInfo(id uint) (*domain.Database, error) {
	var dbConn domain.Database

	err := r.db.Preload("Scans", func(db *gorm.DB) *gorm.DB {
		return db.Order("id DESC").Limit(1).Preload("Tables.Columns.Classification")
	}).First(&dbConn, id).Error

	if err != nil {
		r.logger.Error("error fetching database with latest scan", map[string]any{"error": err, "database_id": id})
		return nil, err
	}

	return &dbConn, nil
}
