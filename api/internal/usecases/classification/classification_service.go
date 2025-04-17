package classification

import (
	"fmt"
	"regexp"

	"github.com/Agu-GC/MELI-Challenge/api/internal/domain"
	"github.com/Agu-GC/MELI-Challenge/api/internal/repositories"
	"github.com/Agu-GC/MELI-Challenge/api/pkg"
)

type ClassificationServiceInterface interface {
	ClassifyData(inputWord string) (*domain.Classification, error)
	NewClassification(newClassification *ClassificationType) error
}

type classificationService struct {
	classRepo repositories.ClassificationRepositoryInterface
	logger    pkg.Logger
}

func NewClassificationService(classRepo repositories.ClassificationRepositoryInterface, logger pkg.Logger) ClassificationServiceInterface {
	return &classificationService{classRepo: classRepo, logger: logger}
}

func (cs *classificationService) NewClassification(newClassification *ClassificationType) error {
	classToSave := domain.Classification{
		Name:             newClassification.Name,
		Description:      newClassification.Description,
		Pattern:          newClassification.Pattern,
		Category:         newClassification.Category,
		SensitivityLevel: newClassification.SensitivityLevel,
	}
	err := cs.classRepo.Create(&classToSave)
	if err != nil {
		cs.logger.Error("Error saving the new classification", map[string]any{"error": err, "classification": newClassification})
		return fmt.Errorf("saving classification error: %w", err)
	}
	return nil
}

func (cs *classificationService) ClassifyData(inputWord string) (*domain.Classification, error) {
	dataTypes, err := cs.classRepo.GetAll()
	if err != nil {
		cs.logger.Error("error getting the classification patterns", map[string]any{"error": err, "inputWord": inputWord})
		return nil, fmt.Errorf("retriving classification error: %w", err)
	}
	for _, dataType := range dataTypes {
		regexp.Compile(dataType.Pattern)
		matched, err := regexp.Match(dataType.Pattern, []byte(inputWord))
		if err != nil {
			cs.logger.Error("error matching the pattern", map[string]any{"error": err, "inputWord": inputWord, "pattern": dataType.Pattern})
			return nil, fmt.Errorf("classifing error: %w", err)
		}
		if matched {
			return dataType, nil
		}
	}
	return nil, nil
}
