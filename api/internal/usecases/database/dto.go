package database

import "gorm.io/gorm"

type ExternalDBSchema struct {
	DatabaseName string
	Tables       map[string][]gorm.ColumnType
}

type DBConnectionInfo struct {
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required,min=1,max=65535"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Dialect  string `json:"dialect"`
	Name     string `json:"name" binding:"required,alphanum"`
	ID       uint   `json:"-"`
}

type DBScanResult struct { //TODO: agregar datos del escaneo en si
	DatabaseID   uint           `json:"id"`
	DatabaseName string         `json:"name"`
	Host         string         `json:"host"`
	Port         int            `json:"port"`
	Tables       []TableScanned `json:"tables"`
}

type TableScanned struct {
	TableName string          `json:"name"`
	Columns   []ColumnScanned `json:"columns"`
}

type ColumnScanned struct {
	ColumnName string `json:"name"`
	DataType   string `json:"type"`
}
