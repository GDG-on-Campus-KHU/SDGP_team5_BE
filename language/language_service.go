// language/language_service.go

package language

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	unitUtil "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/language/util"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/util"
)

type TranslationService interface {
	GetTranslatedMedicalInfo(ctx context.Context, userID int) (*MedicalTranslationResponse, error)
}

type translationService struct{}

func NewTranslationService() TranslationService {
	return &translationService{}
}

func callGeminiForMedicalInfo(apiKey string, langCode string, medicalInfo map[string]string) (string, error) {
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=" + apiKey

	infoJSON, err := json.MarshalIndent(medicalInfo, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal medical info: %v", err)
	}

	// request prompt
	prompt := fmt.Sprintf(`
	Translate the following medical information into %s.
	
	- The fields "Allergy", "Medication", and "Notes" should be translated carefully with proper medical terminology and precision.
	- The field "Name" refers to a person's full name and should be transliterated or adapted to the target language based on cultural or linguistic conventions, rather than translated literally.

	Return the translated result in JSON format with the same keys as the original:
	%s
	`, langCode, string(infoJSON))

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// log.Println("Gemini API Response:", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Println(string(bodyBytes))
		return "", fmt.Errorf("Gemini API error: %s", string(bodyBytes))
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return "", fmt.Errorf("failed to parse Gemini response: %v", err)
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no translation returned")
	}

	rawText := result.Candidates[0].Content.Parts[0].Text
	cleaned := cleanGeminiJSON(rawText)

	return cleaned, nil
}

func (s *translationService) GetTranslatedMedicalInfo(ctx context.Context, userID int) (*MedicalTranslationResponse, error) {

	// 요청한 사용자의 기본 정보 조회
	user, err := util.GetUserByIDFromRequest(ctx, userID)
	if err != nil {
		log.Println("Failed to find user:", err)
		return nil, fmt.Errorf("failed to find user: %v", err)
	}

	// 해당 사용자의 medical info 조회
	medicalInfo, err := util.GetMedicalInfoByID(ctx, int32(user.UserID))
	if err != nil {
		log.Println("Error fetching medical info:", err)
		return nil, fmt.Errorf("failed to find medical info: %v", err)
	}

	if medicalInfo == nil {
		log.Println("No medical info found for user ID:", userID)
		return nil, fmt.Errorf("no medical info found")
	}

	log.Println("Medical info found:", medicalInfo)

	// 사용자의 'country_code'에 따라 번역 언어 결정
	langCode := util.CountryCodeToLangCode(user.CountryCode)

	targetHeightUnit := unitUtil.DefaultHeightUnit(user.CountryCode)
	targetWeightUnit := unitUtil.DefaultWeightUnit(user.CountryCode)

	convertedHeight, convertedHeightUnit := unitUtil.ConvertHeight(medicalInfo.Height, medicalInfo.HeightUnit, targetHeightUnit)
	convertedWeight, convertedWeightUnit := unitUtil.ConvertWeight(medicalInfo.Weight, medicalInfo.WeightUnit, targetWeightUnit)
	log.Println(medicalInfo.Height, medicalInfo.HeightUnit, targetHeightUnit, medicalInfo.Weight, medicalInfo.WeightUnit, targetWeightUnit)

	// 번역할 필드 (allergy, medication, notes, name)
	medicalInfoMap := map[string]string{
		"Allergy":    medicalInfo.Allergy,
		"Medication": medicalInfo.Medication,
		"Notes":      medicalInfo.Notes,
		"Name":       user.Name,
	}

	// Gemini API 호출
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("Gemini API key is not set")
		return nil, fmt.Errorf("Gemini API key is not set")
	}

	cleanedJSON, err := callGeminiForMedicalInfo(apiKey, langCode, medicalInfoMap)
	if err != nil {
		log.Println("Error calling Gemini API:", err)
		return nil, err
	}

	// 6. Gemini response 처리
	var translated TranslatedMedicalInfo
	err = json.Unmarshal([]byte(cleanedJSON), &translated)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Gemini response: %v", err)
	}

	// API response 객체 생성
	response := &MedicalTranslationResponse{
		UserID:     userID,
		Name:       translated.Name,
		Allergy:    translated.Allergy,
		Medication: translated.Medication,
		Notes:      translated.Notes,
		BloodType:  string(medicalInfo.BloodType),
		Height:     convertedHeight,
		HeightUnit: string(convertedHeightUnit),
		Weight:     convertedWeight,
		WeightUnit: string(convertedWeightUnit),
		BirthDate:  medicalInfo.BirthDate,
	}

	responseJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Println("Error marshalling response:", err)
	} else {
		log.Println("Sending response:", string(responseJSON))
	}

	return response, nil
}

// Gemini response 추가 처리
func cleanGeminiJSON(raw string) string {
	re := regexp.MustCompile("(?s)```json\\s*(.*?)\\s*```")
	matches := re.FindStringSubmatch(raw)
	if len(matches) >= 2 {
		return matches[1]
	}

	return strings.TrimSpace(raw)
}
