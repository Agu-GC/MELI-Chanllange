package database

import (
	"fmt"

	"github.com/Agu-GC/MELI-Challenge/api/internal/domain"
	"github.com/Agu-GC/MELI-Challenge/api/internal/repositories"
	"github.com/Agu-GC/MELI-Challenge/api/pkg"

	"github.com/Agu-GC/MELI-Challenge/api/internal/usecases/cipher"
)

type DatabaseServiceInteface interface {
	CreateDatabase(dbInfo *DBConnectionInfo) (uint, error)
	GetDatabaseByID(dbID uint) (*DBConnectionInfo, error)
	GetDatabaseScanResult(dbID uint) (*DBScanResult, error)
}

type databaseService struct {
	dbRepo        repositories.DatabaseRepositoryInterface
	cipherService cipher.CipherServiceInterface
	cipherKey     []byte
	logger        pkg.Logger
}

func NewDatabaseService(dbRepo repositories.DatabaseRepositoryInterface, cipherService cipher.CipherServiceInterface, cipherKey []byte, logger pkg.Logger) DatabaseServiceInteface {
	return &databaseService{
		dbRepo:        dbRepo,
		cipherService: cipherService,
		cipherKey:     cipherKey,
		logger:        logger,
	}
}

func (ds *databaseService) CreateDatabase(dbInfo *DBConnectionInfo) (uint, error) {
	cipherPass, err := ds.cipherService.Encrypt(dbInfo.Password, ds.cipherKey)
	if err != nil {
		return 0, err
	}

	newDb := domain.Database{
		Host:              dbInfo.Host,
		Port:              dbInfo.Port,
		Username:          dbInfo.Username,
		EncryptedPassword: cipherPass,
		Dialect:           dbInfo.Dialect,
		Name:              dbInfo.Name,
	}
	err = ds.dbRepo.Create(&newDb)
	if err != nil {
		return 0, fmt.Errorf("saving db error: %w", err)
	}
	return newDb.ID, nil
}

func (ds *databaseService) GetDatabaseByID(dbID uint) (*DBConnectionInfo, error) {
	db, err := ds.dbRepo.GetByID(dbID)
	if err != nil {
		ds.logger.Error("error getting database connection info", map[string]any{"error": err, "database_id": dbID})
		return nil, fmt.Errorf("retriving db error: %w", err)
	}
	pass, err := ds.cipherService.Decrypt(db.EncryptedPassword, ds.cipherKey)
	if err != nil {
		ds.logger.Error("error decrypting the password", map[string]any{"error": err, "database_id": dbID})
		return nil, fmt.Errorf("decrypting error: %w", err)
	}
	return &DBConnectionInfo{Host: db.Host, Port: db.Port, Username: db.Username, Password: pass, Dialect: db.Dialect, Name: db.Name, ID: db.ID}, nil
}

func (ds *databaseService) GetDatabaseScanResult(dbID uint) (*DBScanResult, error) {
	db, err := ds.dbRepo.GetWithLastScanInfo(dbID)
	if err != nil {
		return nil, fmt.Errorf("retriving db error: %w", err)
	}
	if len(db.Scans) < 1 {
		ds.logger.Error("no scans found", map[string]any{"error": err, "database_id": dbID})
		return nil, fmt.Errorf("no scans found")
	}

	ds.logger.Info("scan restored", map[string]any{"database_id": dbID, "database": db})
	var tableScannedList []TableScanned
	for _, table := range db.Scans[0].Tables {
		tableScanned := TableScanned{TableName: table.TableName}

		for _, column := range table.Columns {
			tableScanned.Columns = append(tableScanned.Columns, ColumnScanned{ColumnName: column.ColumnName, DataType: column.Classification.Name})
		}
		tableScannedList = append(tableScannedList, tableScanned)
	}

	scanResult := DBScanResult{
		DatabaseID:   db.ID,
		DatabaseName: db.Name,
		Host:         db.Host,
		Port:         db.Port,
		Tables:       tableScannedList,
	}
	return &scanResult, nil
}
