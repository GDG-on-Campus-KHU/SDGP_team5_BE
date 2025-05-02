// favorite/favorite_handler.go

package favorite

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/util"
)

type FavoriteHandler struct {
	service *FavoriteService
}

func NewFavoriteHandler(service *FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{service: service}
}

// POST /api/favorites/:id
// @Summary Add a new favorite
// @Description Add a favorite using its ID
// @Tags favorites
// @Accept json
// @Produce json
// @Param id path string true "Situation ID"
// @Success 200 {object} map[string]string "Favorite added successfully"
// @Failure 400 {object} map[string]string "Invalid situation ID or user ID"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 403 {object} map[string]string "Already in favorites"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/favorites/{id} [post]
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	ctx := c.Request.Context()

	situationIndexStr := c.Param("index")
	situationIndex, err := strconv.Atoi(situationIndexStr)
	if err != nil {
		util.RespondBadRequest(c, "Invalid situation ID")
		return
	}

	userIDStr, err := util.GetUserIDFromContext(c)
	if err != nil {
		util.RespondUnauthorized(c, "Unauthorized access")
		return
	}

	userId, err := strconv.Atoi(userIDStr)
	if err != nil {
		util.RespondBadRequest(c, "Invalid user ID")
		return
	}

	if err := h.service.AddFavorite(ctx, userId, situationIndex); err != nil {
		if err.Error() == "already in favorites" {
			util.RespondBadRequest(c, "Already in favorites")
			return
		}
		util.RespondInternalError(c, err.Error())
		return
	}

	util.RespondSuccess(c, "Favorite added successfully")
}

// GET /api/favorites
// @Summary Get a list of favorites
// @Description Get a list of favorites for the authenticated user
// @Tags favorites
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string "Invalid user ID"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/favorites [get]
func (h *FavoriteHandler) GetFavorites(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr, err := util.GetUserIDFromContext(c)
	if err != nil {
		util.RespondUnauthorized(c, "Unauthorized access")
		return
	}

	userId, err := strconv.Atoi(userIDStr)
	if err != nil {
		util.RespondBadRequest(c, "Invalid user ID")
		return
	}

	favorites, err := h.service.GetFavorites(ctx, userId)
	if err != nil {
		util.RespondInternalError(c, err.Error())
		return
	}

	util.RespondSuccess(c, favorites)
}
