// language/language_service.go

package language

import (
	"context"
	"fmt"
	"encoding/json"
	
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/util"
)


type TranslatedMedicalInfo struct {
	BloodType  string `json:"BloodType"`
	Allergy    string `json:"Allergy"`
	Medication string `json:"Medication"`
	Height     string `json:"Height"`
	Weight     string `json:"Weight"`
	BirthDate  string `json:"BirthDate"`
	Notes      string `json:"Notes"`
}

type MedicalTranslationResponse struct {
	UserID     int    `json:"user_id"`
	BloodType  string `json:"blood_type"`
	Allergy    string `json:"allergy"`
	Medication string `json:"medication"`
	Height     string `json:"height"`
	Weight     string `json:"weight"`
	BirthDate  string `json:"birth_date"`
	Notes      string `json:"notes"`
}

type TranslationService interface {
	GetTranslatedMedicalInfo(ctx context.Context, userID int) (*MedicalTranslationResponse, error)
}

type translationService struct{}

func NewTranslationService() TranslationService {
	return &translationService{}
}


func (s *translationService) GetTranslatedMedicalInfo(ctx context.Context, userID int) (*MedicalTranslationResponse, error) {
	user, err := util.GetUserByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	medicalInfo, err := util.GetMedicalInfoByID(ctx, user.InfoID)
	if err != nil {
		return nil, err
	}

	langCode := util.CountryCodeToLangCode(user.CountryCode)

	prompt := fmt.Sprintf(`
		Translate the following medical information into %s:

		{
			"BloodType": "%s",
			"Allergy": "%s",
			"Medication": "%s",
			"Height": "%.2f cm",
			"Weight": "%.2f kg",
			"BirthDate": "%s",
			"Notes": "%s"
		}

		Return the result as JSON.
	`, langCode, medicalInfo.BloodType, medicalInfo.Allergy, medicalInfo.Medication,
	medicalInfo.Height, medicalInfo.Weight, medicalInfo.BirthDate, medicalInfo.Notes)

	geminiResp, err := callGeminiAPI(prompt)
	if err != nil {
		return nil, err
	}

	var translated TranslatedMedicalInfo
	err = json.Unmarshal([]byte(geminiResp), &translated)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Gemini API response: %v", err)
	}

	return &MedicalTranslationResponse{
		UserID:    userID,
		BloodType: translated.BloodType,
		Allergy:   translated.Allergy,
		Medication: translated.Medication,
		Height:    translated.Height,
		Weight:    translated.Weight,
		BirthDate: translated.BirthDate,
		Notes:     translated.Notes,
	}, nil
}