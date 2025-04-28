// situation/situation_routes.go

package situation

import "github.com/gin-gonic/gin"

func RegisterSituationRoutes(router *gin.Engine) {
	router.GET("/api/situation/actions/:index/:language", GetActionsByIndex)
	router.GET("/api/situation/actions/case/:slug/:language", GetActionsBySlug)
}
