package services

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ScanInteface interface {
	OpenConnection(dsn string) (*gorm.DB, error)
}

type ScanService struct {
}

func (s *ScanService) OpenConnection(dsn string) (*gorm.DB, error) {
	return gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{},
	)
}

func (s *ScanService) GetTableStructure(dsn string) (*gorm.DB, error) {
	return gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{},
	)
}

func (s *ScanService) BuildDSN(dataBaseInfo TargetDataBase)
