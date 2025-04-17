package handlers

import (
	"net/http"

	"github.com/Agu-GC/MELI-Challenge/api/internal/usecases/classification"
	"github.com/Agu-GC/MELI-Challenge/api/pkg"
	"github.com/gin-gonic/gin"
)

type ClassificationHandlerInterface interface {
	NewClassification(c *gin.Context)
}

type classificationHandler struct {
	classService classification.ClassificationServiceInterface
	logger       pkg.Logger
}

func NewClassificationHandler(classService classification.ClassificationServiceInterface, logger pkg.Logger) ClassificationHandlerInterface {
	return &classificationHandler{
		classService: classService,
		logger:       logger,
	}
}

func (ch *classificationHandler) NewClassification(c *gin.Context) {
	var newClassification classification.ClassificationType

	err := c.BindJSON(&newClassification)
	if err != nil {
		ch.logger.Error("bad request", map[string]any{"error": err})
		return
	}

	err = ch.classService.NewClassification(&newClassification)
	if err != nil {
		ch.logger.Error("error saving the new classfication", map[string]any{"error": err, "classification": newClassification})
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error saving the new classification"})
		return
	}
}
