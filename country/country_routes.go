// country/country_handler.go

package country

import (
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/auth"
	"github.com/gin-gonic/gin"
)

func RegisterCountryRoutes(r *gin.Engine, countryHandler *CountryHandler) {
	r.GET("/api/country/:country_code", auth.JWTAuthMiddleware(), countryHandler.GetCountryInfo)
}
