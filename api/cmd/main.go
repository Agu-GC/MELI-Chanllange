package main

import (
	"os"
	"strconv"

	"github.com/Agu-GC/MELI-Challenge/api/internal/handlers"
	"github.com/Agu-GC/MELI-Challenge/api/internal/infraestructure"
	"github.com/Agu-GC/MELI-Challenge/api/internal/repositories"
	"github.com/Agu-GC/MELI-Challenge/api/internal/usecases/cipher"
	"github.com/Agu-GC/MELI-Challenge/api/internal/usecases/classification"
	"github.com/Agu-GC/MELI-Challenge/api/internal/usecases/database"
	"github.com/Agu-GC/MELI-Challenge/api/internal/usecases/scan"
	"github.com/Agu-GC/MELI-Challenge/api/pkg"
	"github.com/gin-gonic/gin"
)

func getDBConnectionInfo(logger pkg.Logger) *database.DBConnectionInfo {
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		logger.Fatal("Database connection failed - Invalid Port", map[string]any{
			"error":   err.Error(),
			"DB_PORT": os.Getenv("DB_PORT"),
		})
	}

	return &database.DBConnectionInfo{
		Host:     os.Getenv("DB_HOST"),
		Port:     dbPort,
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		Dialect:  "mysql",
	}
}

func main() {
	logger := pkg.NewLogger("database-sacanner")
	router := gin.Default()
	router.Use(pkg.GinLogger(logger), pkg.GinRecovery(logger))

	authUsername := os.Getenv("AUTH_USERNAME")
	authPassword := os.Getenv("AUTH_PASSWORD")

	chiperkey := []byte(os.Getenv("CIPHER_KEY"))

	dbConnectionInfo := getDBConnectionInfo(logger)

	dbConnector := infraestructure.NewDatabaseConnector(logger)
	dbConnection, err := dbConnector.OpenConnection(dbConnectionInfo)
	if err != nil {
		logger.Fatal("Database connection failed", map[string]any{"error": err.Error()})
	}
	err = dbConnector.ExecuteMigrations(dbConnection)
	if err != nil {
		logger.Fatal("Database connection failed", map[string]any{"error": err.Error()})
	}

	classRepo := repositories.NewClassificationRepository(dbConnection, logger)
	databaseRepo := repositories.NewDatabaseRepository(dbConnection, logger)
	scanRepo := repositories.NewScannRepository(dbConnection, logger)

	classService := classification.NewClassificationService(classRepo, logger)
	cipherService := cipher.NewCipherService(logger)
	dbService := database.NewDatabaseService(databaseRepo, cipherService, chiperkey, logger)
	scanService := scan.NewScanService(dbConnector, scanRepo, classService, dbService, logger)

	scanHandler := handlers.NewScanHandler(scanService, logger)
	databaseHandler := handlers.NewDatabaseHandler(dbService, logger)
	classHandler := handlers.NewClassificationHandler(classService, logger)

	authorized := router.Group("/api/v1", gin.BasicAuth(gin.Accounts{
		authUsername: authPassword,
	}))

	authorized.POST("/database", databaseHandler.CreateDatabase)
	authorized.GET("/database/scan/:id", databaseHandler.GetDatabaseScanResult)

	authorized.POST("/database/scan/:id", scanHandler.CreateScan)

	authorized.POST("/classifications", classHandler.NewClassification)

	router.RunTLS(":443", "/app/ssl/server.crt", "/app/ssl/server.key")
}
