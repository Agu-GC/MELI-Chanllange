package classification

type ClassificationType struct {
	Name             string `json:"name" binding:"required"`
	Description      string `json:"descrption" binding:"required"`
	Pattern          string `json:"pattern" binding:"required"`
	Category         string `json:"category" binding:"required"`
	SensitivityLevel int    `json:"sensitivity_level" binding:"required"`
}
