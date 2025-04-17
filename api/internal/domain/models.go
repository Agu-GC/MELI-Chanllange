package domain

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Database struct {
	Base
	Name              string `gorm:"uniqueIndex;size:100"`
	Dialect           string `gorm:"size:50"`
	Host              string `gorm:"size:100"`
	Port              int
	Username          string `gorm:"size:50"`
	EncryptedPassword string `gorm:"size:255"`
	Scans             []Scan
}

type Scan struct {
	Base
	DatabaseID   uint
	Status       string `gorm:"size:20"` // pending, running, completed, failed
	StartedAt    time.Time
	FinishedAt   time.Time
	ErrorMessage string `gorm:"type:text"`
	Tables       []Table
}

type Table struct {
	Base
	ScanID     uint
	TableName  string `gorm:"size:100;index"`
	SchemaName string `gorm:"size:50;index"`
	Columns    []Column
}

type Column struct {
	Base
	TableID          uint
	ColumnName       string `gorm:"size:100;index"`
	DataType         string `gorm:"size:50"` // Tipo original
	IsNullable       bool
	ClassificationID *uint `gorm:"index"`
	Classification   Classification
}

type Classification struct {
	Base
	Name             string `gorm:"uniqueIndex;size:50"`
	Description      string `gorm:"type:text"`
	Pattern          string `gorm:"type:text"`
	Category         string `gorm:"size:50"` // PII, Financial, etc.
	SensitivityLevel int
}
