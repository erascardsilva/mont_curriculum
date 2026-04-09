package parser

import "mont_curriculum/internal/models"

// ExtractedData contém todos os dados que conseguimos "pescar" do arquivo importado
type ExtractedData struct {
	Profile    models.Profile
	Education  []models.Education
	Experience []models.Experience
	Projects   []models.Project
}
