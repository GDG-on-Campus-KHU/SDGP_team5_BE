// group/group_routes.go

package group

import (
	"github.com/gin-gonic/gin"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/auth"
)

func RegisterGroupRoutes(r *gin.Engine, groupHandler *GroupHandler) {

	// group routes
	r.POST("/api/groups", auth.JWTAuthMiddleware(), groupHandler.CreateGroup)
	r.GET("/api/groups/:id", auth.JWTAuthMiddleware(), groupHandler.GetGroupByGroupID)
	r.PATCH("/api/groups/:id", auth.JWTAuthMiddleware(), groupHandler.UpdateGroup)
	r.DELETE("/api/groups/:id", auth.JWTAuthMiddleware(), groupHandler.DeleteGroup)
	r.GET("/api/groups", auth.JWTAuthMiddleware(), groupHandler.ListGroups)

	// get all groups for logged-in user
	r.GET("/api/groups/me", auth.JWTAuthMiddleware(), groupHandler.GetMyGroups)

	// group members routes
	r.GET("/api/groups/:id/members", auth.JWTAuthMiddleware(), groupHandler.GetGroupMembers)

	// group invitation routes
	r.POST("/api/groups/:id/invite", auth.JWTAuthMiddleware(), groupHandler.InviteUser)
	r.POST("/api/groups/:id/accept", auth.JWTAuthMiddleware(), groupHandler.AcceptInvite)
	r.POST("/api/groups/:id/reject", auth.JWTAuthMiddleware(), groupHandler.RejectInvite)

	// leave group route
	r.DELETE("/api/groups/:id/members/me", auth.JWTAuthMiddleware(), groupHandler.LeaveGroup)

	// pending group route
	r.GET("/api/groups/pending/me", auth.JWTAuthMiddleware(), groupHandler.GetPendingGroups)
}