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
	UserID int `json:"user_id"`
}


func (h *TranslationHandler) TranslateMedicalInfoHandler(c *gin.Context) {

	var req MedicalTranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.RespondBadRequest(c, "Invalid request body")
		return
	}

	userID := req.UserID

	result, err := h.service.GetTranslatedMedicalInfo(c.Request.Context(), userID)
	if err != nil {
		util.RespondInternalError(c, err.Error())
		return
	}

	util.RespondSuccess(c, result)
}