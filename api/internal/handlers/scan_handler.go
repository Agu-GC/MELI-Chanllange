package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Agu-GC/MELI-Challenge/api/internal/usecases/scan"
	"github.com/Agu-GC/MELI-Challenge/api/pkg"
	"github.com/gin-gonic/gin"
)

type ScanHandlerInterface interface {
	CreateScan(c *gin.Context)
}

type scanHandler struct {
	scanService scan.ScanServiceInteface
	logger      pkg.Logger
}

func NewScanHandler(ss scan.ScanServiceInteface, logger pkg.Logger) ScanHandlerInterface {
	return &scanHandler{scanService: ss, logger: logger}
}

func (sh *scanHandler) CreateScan(c *gin.Context) {
	stringDbID := c.Params.ByName("id")
	dbID, err := strconv.Atoi(stringDbID)
	if err != nil || dbID < 0 {
		sh.logger.Info("bad resquest - invalid database-id", map[string]any{"database_id": dbID})
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid database id"})
		return
	}

	sh.logger.Info("creating a database scan", map[string]any{"database_id": dbID})
	err = sh.scanService.ScanDatabase(uint(dbID))
	if err != nil {
		sh.logger.Error("error scanning the database", map[string]any{"error": err, "database_id": dbID})
		statusCode := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
		}
		c.IndentedJSON(statusCode, gin.H{"message": err.Error()})
		return
	}
	sh.logger.Info("database scaned", map[string]any{"database_id": dbID})
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "scan executed"})
}
