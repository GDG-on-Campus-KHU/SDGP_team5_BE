// group/repository/group_repository.go

package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
)

type GroupRepository interface {
	Create(ctx context.Context, group *model.Group) (*model.Group, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.Group, error)
	Update(ctx context.Context, group *model.Group) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	List(ctx context.Context) ([]*model.Group, error)

	GetGroupMembers(ctx context.Context, groupID primitive.ObjectID) (*model.GroupMembersResponse, error)
	
	InviteUser(ctx context.Context, groupID primitive.ObjectID, inviterUserID int, inviteeEmail string) error
	AcceptInvite(ctx context.Context, groupID primitive.ObjectID, userID int) error
	RejectInvite(ctx context.Context, groupID primitive.ObjectID, userID int) error
	
	LeaveGroup(ctx context.Context, groupID primitive.ObjectID, userID int) error
}