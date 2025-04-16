package handler

import "github.com/gin-gonic/gin"

func RegisterSituationRoutes(router *gin.Engine) {
	router.GET("/situation/actions/:index/:language", GetActionsByIndex)
	router.GET("/situation/actions/case/:slug/:language", GetActionsBySlug)
}
