package middleware

import (
	"pethub_api/models"
)

// FetchAdoptionApplicationAndQuestionnaire fetches both the adoption application and questionnaire for an adopter
func FetchAdoptionApplicationAndQuestionnaire(adopterID int) (*models.AdoptionApplication, *models.Questionnaires, error) {
	var app models.AdoptionApplication
	if err := DBConn.Where("adopter_id = ?", adopterID).First(&app).Error; err != nil {
		return nil, nil, err
	}

	var questionnaire models.Questionnaires
	if err := DBConn.Where("application_id = ?", app.ApplicationID).First(&questionnaire).Error; err != nil {
		return &app, nil, err
	}

	return &app, &questionnaire, nil
}
