// country/country_handler.go

package country

import (
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/util"
	"github.com/gin-gonic/gin"
)

type CountryHandler struct {
	service *CountryService
}

func NewCountryHandler(service *CountryService) *CountryHandler {
	return &CountryHandler{service: service}
}

type countryCodeRequest struct {
	CountryCode string `json:"country_code"`
}

// GET /api/users/:country_code
// @Summary Get country info
// @Description Get country info by country code
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} CountryInfoResponse "Country info retrieved successfully"
// @Failure 400 {object} map[string]string "Invalid country code"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/users/:country_code [get]
func (h *CountryHandler) GetCountryInfo(c *gin.Context) {
	ctx := c.Request.Context()

	countryCode := c.Param("country_code")
	if countryCode == "" {
		util.RespondBadRequest(c, "Country code is required")
		return
	}

	countryInfoPtr, err := h.service.GetCountryInfo(ctx, countryCode)
	if err != nil {
		util.RespondInternalError(c, err.Error())
		return
	}
	countryInfoResponse := *countryInfoPtr
	util.RespondSuccess(c, countryInfoResponse)
}
