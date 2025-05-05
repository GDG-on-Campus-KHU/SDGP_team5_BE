// medical_info/medical_info_handler.go

package medical_info

import (
	"strconv"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/util"
	"github.com/gin-gonic/gin"
)

type MedicalInfoHandler struct {
	service *MedicalInfoService
}

func NewMedicalInfoHandler(service *MedicalInfoService) *MedicalInfoHandler {
	return &MedicalInfoHandler{service: service}
}

type MedicalInfoRequest struct {
	BloodType  string  `json:"blood_type" binding:"required"`
	Allergy    string  `json:"allergy" binding:"required"`
	Medication string  `json:"medication" binding:"required"`
	Height     float64 `json:"height" binding:"required"`
	HeightUnit string  `json:"height_unit" binding:"required"`
	Weight     float64 `json:"weight" binding:"required"`
	WeightUnit string  `json:"weight_unit" binding:"required"`
	BirthDate  string  `json:"birth_date" binding:"required"`
	Notes      string  `json:"notes" binding:"required"`
}

// POST /api/medical-info
// @Summary Create medical info
// @Description Create medical info
// @Tags medical_info
// @Accept json
// @Produce json
// @Param medical_info body MedicalInfo true "Medical info"
// @Success 200 {object} map[string]string "Medical info created successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/medical-info [post]
func (h *MedicalInfoHandler) CreateMedicalInfo(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr, err := util.GetUserIDFromContext(c)
	if err != nil {
		util.RespondUnauthorized(c, "Unauthorized access")
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		util.RespondBadRequest(c, "Invalid user ID")
		return
	}

	var req MedicalInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.RespondBadRequest(c, "Invalid request body")
		return
	}

	medical_info, err := h.service.CreateMedicalInfo(ctx, userID, req)
	if err != nil {
		util.RespondInternalError(c, err.Error())
		return
	}

	util.RespondSuccess(c, medical_info)
}

// GET /api/medical-info/me
// @Summary Get medical info
// @Description Get medical info
// @Tags medical_info
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Medical info retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/medical-info/me [get]
func (h *MedicalInfoHandler) GetMedicalInfo(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr, err := util.GetUserIDFromContext(c)
	if err != nil {
		util.RespondUnauthorized(c, "Unauthorized access")
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		util.RespondBadRequest(c, "Invalid user ID")
		return
	}

	medical_info, err := h.service.GetMedicalInfo(ctx, userID)
	if err != nil {
		util.RespondInternalError(c, err.Error())
		return
	}

	util.RespondSuccess(c, medical_info)
}

// GET /api/medical-info/:id
// @Summary Get medical info by ID
// @Description Get medical info by user ID
// @Tags medical_info
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Medical info retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/medical-info/:id [get]
func (h *MedicalInfoHandler) GetMedicalInfoByID(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		util.RespondBadRequest(c, "Invalid user ID")
		return
	}

	medical_info, err := h.service.GetMedicalInfo(ctx, userID)
	if err != nil {
		util.RespondInternalError(c, err.Error())
		return
	}

	util.RespondSuccess(c, medical_info)
}

// PUT /api/medical-info
// @Summary Update medical info
// @Description Update medical info
// @Tags medical_info
// @Accept json
// @Produce json
// @Param medical_info body MedicalInfo true "Medical info"
// @Success 200 {object} map[string]string "Medical info updated successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/medical-info [put]
func (h *MedicalInfoHandler) UpdateMedicalInfo(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr, err := util.GetUserIDFromContext(c)
	if err != nil {
		util.RespondUnauthorized(c, "Unauthorized access")
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		util.RespondBadRequest(c, "Invalid user ID")
		return
	}

	var req MedicalInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.RespondBadRequest(c, "Invalid request body")
		return
	}

	medical_info, err := h.service.UpdateMedicalInfo(ctx, userID, req)
	if err != nil {
		util.RespondInternalError(c, err.Error())
		return
	}

	util.RespondSuccess(c, medical_info)
}
