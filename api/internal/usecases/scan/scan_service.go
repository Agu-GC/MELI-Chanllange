package scan

import (
	"fmt"
	"time"

	"github.com/Agu-GC/MELI-Challenge/api/internal/domain"
	"github.com/Agu-GC/MELI-Challenge/api/internal/infraestructure"
	"github.com/Agu-GC/MELI-Challenge/api/internal/repositories"
	"github.com/Agu-GC/MELI-Challenge/api/internal/usecases/classification"
	"github.com/Agu-GC/MELI-Challenge/api/internal/usecases/database"
	"github.com/Agu-GC/MELI-Challenge/api/pkg"
	"gorm.io/gorm"
)

type ScanServiceInteface interface {
	ScanDatabase(uint) error
	GetDatabaseSchema(*database.DBConnectionInfo) (*database.ExternalDBSchema, error)
}

type scanService struct {
	dbConn       infraestructure.DBConnectorInterface
	classService classification.ClassificationServiceInterface
	dbService    database.DatabaseServiceInteface
	scanRepo     repositories.ScanRepositoryInterface
	logger       pkg.Logger
}

func NewScanService(dbConn infraestructure.DBConnectorInterface, scanRepo repositories.ScanRepositoryInterface, classService classification.ClassificationServiceInterface, dbService database.DatabaseServiceInteface, logger pkg.Logger) ScanServiceInteface {
	return &scanService{
		dbConn:       dbConn,
		classService: classService,
		scanRepo:     scanRepo,
		dbService:    dbService,
		logger:       logger,
	}
}

func (ss *scanService) ScanDatabase(dbId uint) (err error) {
	dbInfo, err := ss.dbService.GetDatabaseByID(dbId)
	if err != nil {
		return err
	}
	return ss.initScan(dbInfo)
}

func (ss *scanService) initScan(dbInfo *database.DBConnectionInfo) (err error) {
	scan := domain.Scan{
		DatabaseID: dbInfo.ID,
		Status:     "running",
		StartedAt:  time.Now(),
	}

	defer func() {
		scan.FinishedAt = time.Now()
		scan.Status = "completed"
		if err != nil {
			scan.Status = "failed"
			scan.ErrorMessage = err.Error()
		}
		ss.saveScan(&scan)
	}()

	extDb, err := ss.GetDatabaseSchema(dbInfo)
	if err != nil {
		return err
	}

	lclTables, err := ss.classifySchema(extDb)
	if err != nil {
		return err
	}
	scan.Tables = lclTables

	return nil
}

func (ss *scanService) saveScan(scan *domain.Scan) {
	if err := ss.scanRepo.Create(scan); err != nil {
		ss.logger.Fatal("Failed to save scan", map[string]any{"error": err, "scan": scan})
	}
}

func (ss *scanService) classifySchema(extDb *database.ExternalDBSchema) ([]domain.Table, error) {
	var lclTableList []domain.Table
	for extTable, extColumnList := range extDb.Tables {

		lclTable := domain.Table{
			TableName:  extTable,
			SchemaName: extDb.DatabaseName,
		}

		for _, extColumn := range extColumnList {
			datatype, err := ss.classService.ClassifyData(extColumn.Name())
			if err != nil {
				ss.logger.Error("error classifing the column", map[string]any{"error": err, "database": extDb, "column": extColumn})
				return nil, fmt.Errorf("classification error: %w", err)
			}

			lclColumn := domain.Column{
				ColumnName: extColumn.Name(),
				DataType:   extColumn.DatabaseTypeName(),
			}

			if datatype != nil {
				lclColumn.ClassificationID = &datatype.ID
				lclColumn.Classification = *datatype
			}

			lclTable.Columns = append(lclTable.Columns, lclColumn)
		}
		lclTableList = append(lclTableList, lclTable)
	}

	return lclTableList, nil
}

func (ss *scanService) GetDatabaseSchema(dbInfo *database.DBConnectionInfo) (*database.ExternalDBSchema, error) {
	dbConn, err := ss.dbConn.OpenConnection(dbInfo)
	if err != nil {
		ss.logger.Error("error opening the database connection", map[string]any{"error": err, "database": dbInfo})
		return nil, fmt.Errorf("connection error: %w", err)
	}
	defer ss.dbConn.CloseConnection(dbConn)

	migrator := dbConn.Migrator()
	extTableList, err := migrator.GetTables()
	if err != nil {
		ss.logger.Error("error getting tables", map[string]any{"error": err, "database": dbInfo})
		return nil, fmt.Errorf("tables error: %w", err)
	}

	dbSchema := database.ExternalDBSchema{DatabaseName: dbInfo.Name, Tables: make(map[string][]gorm.ColumnType)}
	for _, extTable := range extTableList {
		extColumnList, err := migrator.ColumnTypes(extTable)
		if err != nil {
			ss.logger.Error("getting the table columns", map[string]any{"error": err, "database": dbInfo, "table": extTable})
			return nil, fmt.Errorf("columns error: %w", err)
		}
		dbSchema.Tables[extTable] = extColumnList
	}
	return &dbSchema, nil
}
