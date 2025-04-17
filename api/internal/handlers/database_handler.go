package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Agu-GC/MELI-Challenge/api/internal/usecases/database"
	"github.com/Agu-GC/MELI-Challenge/api/pkg"
	"github.com/gin-gonic/gin"
)

type DatabaseHandlerInterface interface {
	CreateDatabase(c *gin.Context)
	GetDatabaseScanResult(c *gin.Context)
}

type databaseHandler struct {
	databaseService database.DatabaseServiceInteface
	logger          pkg.Logger
}

func NewDatabaseHandler(ds database.DatabaseServiceInteface, logger pkg.Logger) DatabaseHandlerInterface {
	return &databaseHandler{databaseService: ds, logger: logger}
}

func (dh *databaseHandler) CreateDatabase(c *gin.Context) {
	var databaseInfo database.DBConnectionInfo
	if err := c.BindJSON(&databaseInfo); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	dh.logger.Info("saving database", map[string]any{"database": databaseInfo})
	dbID, err := dh.databaseService.CreateDatabase(&databaseInfo)
	if err != nil {
		dh.logger.Error("error saving database", map[string]any{"error": err, "database": databaseInfo})
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		c.IndentedJSON(statusCode, gin.H{"message": err.Error()}) //TODO: mejorar respuesta de error		return
	}
	dh.logger.Info("database saved", map[string]any{"database_id": dbID})
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "created", "database_id": dbID})
}

func (dh *databaseHandler) GetDatabaseScanResult(c *gin.Context) {
	stringDbID := c.Params.ByName("id")
	dbID, err := strconv.Atoi(stringDbID)
	if err != nil || dbID < 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid database id"})
		return
	}

	dh.logger.Info("getting database scan", map[string]any{"database_id": dbID})
	dbScanResult, err := dh.databaseService.GetDatabaseScanResult(uint(dbID))
	if err != nil {
		dh.logger.Error("getting database scan", map[string]any{"error": err, "database_id": dbID})
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		c.IndentedJSON(statusCode, gin.H{"message": err.Error()}) //TODO: mejorar respuesta de error
		return
	}
	dh.logger.Info("database scan retrived", map[string]any{"database_id": dbID})
	c.IndentedJSON(http.StatusOK, dbScanResult)
}
