// group/group_handler.go

package group

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

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

type UpdateGroupRequest struct {
	GroupName string `json:"group_name"`
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


// PATCH /api/groups/{id}
// @Summary Update a group
// @Description Update a group's name by its ID (Partial update)
// @Tags groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body UpdateGroupRequest true "Updated group name"
// @Success 200 {object} map[string]string "Group updated"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/groups/{id} [patch]
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

	group, err := h.service.GetGroupByGroupID(ctx, objID)
	if err != nil {
		util.RespondInternalError(c, "Error fetching group")
		return
	}
	if group == nil {
		util.RespondNotFound(c, "Group not found")
		return
	}

	if updateReq.GroupName != "" {
		group.GroupName = updateReq.GroupName
	}

	if err := h.service.UpdateGroup(ctx, group); err != nil {
		util.RespondInternalError(c, err.Error())
		return
	}
	util.RespondSuccess(c, gin.H{"message": "Group name updated"})
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


// GET /api/groups/me
// @Summary Get all groups for the logged-in user
// @Description Retrieve the groups the logged-in user belongs to
// @Tags groups
// @Accept json
// @Produce json
// @Success 200 {array} model.Group "List of groups the user belongs to"
// @Failure 401 {object} map[string]string "Unauthorized access"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/groups/me [get]
func (h *GroupHandler) GetMyGroups(c *gin.Context) {
	ctx := c.Request.Context()

	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		util.RespondUnauthorized(c, "Unauthorized access")
		return
	}

	groups, err := h.service.GetGroupsByUserID(ctx, userID)
	if err != nil {
		util.RespondInternalError(c, err.Error())
		return
	}

	util.RespondSuccess(c, groups)
}


// GET /api/groups/{id}/members
// @Summary Get members of a specific group
// @Description Retrieve the members of a group by its ID
// @Tags groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {object} model.GroupMembersResponse "List of members in the group"
// @Failure 400 {object} map[string]string "Invalid group ID format"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/groups/{id}/members [get]
func (h *GroupHandler) GetGroupMembers(c *gin.Context) {
	ctx := c.Request.Context()

	// extract 'GroupID'
	groupIDStr := c.Param("id")
	groupID, err := primitive.ObjectIDFromHex(groupIDStr)
	if err != nil {
		util.RespondBadRequest(c, "Invalid group ID format")
		return
	}

	// group에 포함된 member 가져오기
	groupMembers, err := h.service.GetGroupMembers(ctx, groupID)
	if err != nil {
		util.RespondInternalError(c, err.Error())
		return
	}

	util.RespondSuccess(c, groupMembers)
}