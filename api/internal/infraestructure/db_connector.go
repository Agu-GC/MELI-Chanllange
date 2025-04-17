package infraestructure

import (
	"strings"
	"time"

	"fmt"

	"github.com/Agu-GC/MELI-Challenge/api/internal/domain"
	"github.com/Agu-GC/MELI-Challenge/api/internal/usecases/database"
	"github.com/Agu-GC/MELI-Challenge/api/pkg"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

const MAX_RETRIES = 3
const DELAY_SECONDS = 5

type DBConnectorInterface interface {
	OpenConnection(dbConn *database.DBConnectionInfo) (*gorm.DB, error)
	CloseConnection(db *gorm.DB) error
	ExecuteMigrations(db *gorm.DB) error
}

type dbConnector struct {
	logger pkg.Logger
}

func NewDatabaseConnector(logger pkg.Logger) DBConnectorInterface {
	return &dbConnector{logger: logger}
}

func (cb *dbConnector) OpenConnection(dbConn *database.DBConnectionInfo) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch strings.ToLower(dbConn.Dialect) {
	case "sqlserver":
		dsn := fmt.Sprintf(
			"sqlserver://%s:%s@%s:%d?database=%s",
			dbConn.Username, dbConn.Password, dbConn.Host, dbConn.Port, dbConn.Name,
		)
		dialector = sqlserver.Open(dsn)

	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			dbConn.Host, dbConn.Port, dbConn.Username, dbConn.Password, dbConn.Name,
		)
		dialector = postgres.Open(dsn)

	default: //Asumimos por defecto mysql
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbConn.Username, dbConn.Password, dbConn.Host, dbConn.Port, dbConn.Name,
		)
		dialector = mysql.Open(dsn)
	}

	var err error
	for i := 0; i < MAX_RETRIES; i++ {
		db, err := gorm.Open(dialector, &gorm.Config{})
		if err == nil {
			return db, nil
		}

		cb.logger.Info("Connection failed. Retrying...", map[string]any{"attempt": i})
		time.Sleep(DELAY_SECONDS * time.Second)
	}

	cb.logger.Error("Failed to establish connection", map[string]any{"error": err})
	return nil, fmt.Errorf("failed to establish connection: %v", err)
}

func (cb *dbConnector) CloseConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		cb.logger.Error("Closing the connection", map[string]any{"error": err, "database": db})
		return fmt.Errorf("closing the connection: %v", err)
	}
	return sqlDB.Close()
}

func (cb *dbConnector) ExecuteMigrations(db *gorm.DB) error {
	time.Sleep(DELAY_SECONDS * time.Second)
	return db.AutoMigrate(
		&domain.Database{},
		&domain.Scan{},
		&domain.Table{},
		&domain.Column{},
		&domain.Classification{},
	)
}
