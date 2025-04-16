package handler

import (
	"net/http"
	"strconv"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/situation/service"
	"github.com/gin-gonic/gin"
)

// @Summary Get Actions by Situation Index
// @Description Retrieve the actions for a specific situation based on its index (numerical ID).
// @Tags situation
// @Accept json
// @Produce json
// @Param index path int true "Situation Index"
// @Param language path string true "Language" default(en)
// @Success 200 {object} model.Situation "Successfully retrieved situation and actions"
// @Failure 400 {object} map[string]string "error message"
// @Failure 404 {object} map[string]string "error message"
// @Router /situation/actions/{index}/{language} [get]
func GetActionsByIndex(c *gin.Context) {
	indexStr := c.Param("index")
	language := c.Param("language")

	// Convert index to integer
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid index format"})
		return
	}

	situation, err := service.GetSituationByIndex(index, language)
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, situation)
}

// @Summary Get Actions by Situation Slug
// @Description Retrieve the actions for a specific situation based on its slug (text-based ID).
// @Tags situation
// @Accept json
// @Produce json
// @Param slug path string true "Situation Slug"
// @Param language path string true "Language" default(en)
// @Success 200 {object} model.Situation "Successfully retrieved situation and actions"
// @Failure 404 {object} map[string]string "error message"
// @Router /situation/actions/case/{slug}/{language} [get]
func GetActionsBySlug(c *gin.Context) {
	slug := c.Param("slug")
	language := c.Param("language")

	situation, err := service.GetSituationBySlug(slug, language)
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, situation)
}
