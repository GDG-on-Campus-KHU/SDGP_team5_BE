// group/group_handler.go

package group

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/util"
)


type GroupHandler struct {
	service *GroupService
}

func NewGroupHandler(service *GroupService) *GroupHandler {
	return &GroupHandler{service: service}
}


type CreateGroupRequest struct {
	GroupName string `json:"group_name" binding:"required"`
}


// POST /api/groups
// @Summary Create a new group
// @Description Create a group with a specified name
// @Tags groups
// @Accept json
// @Produce json
// @Param request body CreateGroupRequest true "Group name"
// @Success 200 {object} model.Group "Group created successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/groups [post]
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		util.RespondUnauthorized(c, "Unauthorized access")
		return
	}

	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.RespondBadRequest(c, "Invalid request body")
		return
	}

	group, err := h.service.CreateGroup(ctx, req.GroupName, userID)
	if err != nil {
		util.RespondInternalError(c, err.Error())
		return
	}
	util.RespondSuccess(c, group)
}


// GET /api/groups/{id}
// @Summary Get a group by ID
// @Description Retrieve a group's information using its ID
// @Tags groups
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} model.Group "Group found"
// @Failure 400 {object} map[string]string "Invalid group ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/groups/{id} [get]
func (h *GroupHandler) GetGroupByGroupID(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		util.RespondBadRequest(c, "Invalid group ID")
		return
	}

	group, err := h.service.GetGroupByGroupID(ctx, objID)
	if err != nil {
		util.RespondInternalError(c, err.Error())
		return
	}
	util.RespondSuccess(c, group)
}


// PUT /api/groups/{id}
// @Summary Update a group
// @Description Update a group's name by its ID
// @Tags groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body CreateGroupRequest true "Updated group name"
// @Success 200 {object} map[string]string "Group updated"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/groups/{id} [put]
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		util.RespondBadRequest(c, "Invalid group ID")
		return
	}

	var updateReq CreateGroupRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		util.RespondBadRequest(c, "Invalid request body")
		return
	}

	group := &model.Group{
		ID:        objID,
		GroupName: updateReq.GroupName,
	}

	if err := h.service.UpdateGroup(ctx, group); err != nil {
		util.RespondInternalError(c, err.Error())
		return
	}
	util.RespondSuccess(c, gin.H{"message": "Group updated"})
}


// DELETE /api/groups/{id}
// @Summary Delete a group
// @Description Delete a group by its ID
// @Tags groups
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} map[string]string "Group deleted"
// @Failure 400 {object} map[string]string "Invalid group ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/groups/{id} [delete]
func (h *GroupHandler) DeleteGroup(c *gin.Context) {

	ctx := c.Request.Context()

	idStr := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		util.RespondBadRequest(c, "Invalid group ID")
		return
	}

	if err := h.service.DeleteGroup(ctx, objID); err != nil {
		util.RespondInternalError(c, err.Error())
		return
	}
	util.RespondSuccess(c, gin.H{"message": "Group deleted"})
}


// GET /api/groups
// @Summary List all groups
// @Description Get a list of all groups
// @Tags groups
// @Produce json
// @Success 200 {array} model.Group "List of groups"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/groups [get]
func (h *GroupHandler) ListGroups(c *gin.Context) {
	ctx := c.Request.Context()

	groups, err := h.service.ListGroups(ctx)
	if err != nil {
		util.RespondInternalError(c, err.Error())
		return
	}
	util.RespondSuccess(c, groups)
}