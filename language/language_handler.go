// language/language_handler.go

package language

import (
	"github.com/gin-gonic/gin"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/util"
)


type TranslationHandler struct {
	service TranslationService
}

func NewTranslationHandler(service TranslationService) *TranslationHandler {
	return &TranslationHandler{service: service}
}


type MedicalTranslationRequest struct {
	UserID int `json:"user_id" binding:"required"`
}

type TranslatedMedicalInfo struct {
	Name       string `json:"name"`
	Allergy    string `json:"allergy"`
	Medication string `json:"medication"`
	Notes      string `json:"notes"`
}



type MedicalTranslationResponse struct {
	UserID     int		`json:"user_id"`
	Name       string	`json:"name"`
	Allergy    string	`json:"allergy"`
	Medication string	`json:"medication"`
	Notes      string 	`json:"notes"`
	BloodType  string 	`json:"blood_type"`
	Height     float64	`json:"height"`
	HeightUnit string  	`json:"height_unit"`
	Weight     float64 	`json:"weight"`
	WeightUnit string  	`json:"weight_unit"`
	BirthDate  string  	`json:"birth_date"`
	InfoTitles []string `json:"info_titles,omitempty"`
}


func (h *TranslationHandler) TranslateMedicalInfoHandler(c *gin.Context) {
	var req MedicalTranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.RespondBadRequest(c, "Invalid request body")
		return
	}

	result, err := h.service.GetTranslatedMedicalInfo(c.Request.Context(), req.UserID)
	if err != nil {
		util.RespondInternalError(c, err.Error())
		return
	}

	util.RespondSuccess(c, result)
}