// user/user_handler.go

package user

import (
	"strconv"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/util"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *UserHandler {
	return &UserHandler{service: service}
}

type countryCodeRequest struct {
	CountryCode string `json:"country_code"`
}

// GET /api/users/me
// @Summary Get user info
// @Description Get user info by user ID
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} UserInfoResponse "User info retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/me [get]
func (h *UserHandler) GetUserInfo(c *gin.Context) {
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

	userInfoPtr, err := h.service.GetUserInfo(ctx, userID)
	if err != nil {
		util.RespondInternalError(c, err.Error())
		return
	}
	userInfoResponse := *userInfoPtr
	util.RespondSuccess(c, userInfoResponse)
}

// GET /api/users/info/:id
// @Summary Get user info by ID
// @Description Get user info by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserInfoResponse "User info retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/info/{id} [get]
func (h *UserHandler) GetUserInfoByID(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr := c.Param("member_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		util.RespondBadRequest(c, "Invalid user ID")
		return
	}

	userInfoPtr, err := h.service.GetUserInfo(ctx, userID)
	if err != nil {
		util.RespondInternalError(c, err.Error())
		return
	}
	userInfoResponse := *userInfoPtr
	util.RespondSuccess(c, userInfoResponse)
}

// PATCH /api/users/me/country
// @Summary Update user country
// @Description Update user country by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param country_code body int true "Country code"
// @Success 200 {object} UserInfoResponse "User country updated successfully"
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/me/country [patch]
func (h *UserHandler) UpdateCountry(c *gin.Context) {
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

	var req countryCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.RespondBadRequest(c, "Invalid input")
		return
	}

	userInfoPtr, err := h.service.UpdateCountry(ctx, userID, req.CountryCode)
	if err != nil {
		util.RespondInternalError(c, err.Error())
		return
	}
	userInfoResponse := *userInfoPtr
	util.RespondSuccess(c, userInfoResponse)
}
