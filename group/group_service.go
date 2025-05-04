// group/group_service.go

package group

import (
	"context"
	"log"
	"strconv"
	"time"
	"fmt"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	dbConfig "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/config"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/group/repository"
)

type GroupService struct {
	repo repository.GroupRepository
}

// constructor function
func NewGroupService(repo repository.GroupRepository) *GroupService {
	if repo == nil {
		log.Fatal("GroupRepository cannot be nil")
		return nil
	}
	return &GroupService{repo: repo}
}


// CreateGroup creates a new group with the given name and returns the created group
func (s *GroupService) CreateGroup(ctx context.Context, name string, creatorUserID string) (*model.Group, error) {
	// 'creatorUserID'를 int로 convert
	userIDInt, err := strconv.Atoi(creatorUserID)
	if err != nil {
		log.Printf("invalid creatorUserID: %v", err)
		return nil, err
	}

	// 'creatorUserID(int)'로 사용자의 email 조회
	var user model.User
	err = dbConfig.UserCollection.FindOne(ctx, bson.M{"user_id": userIDInt}).Decode(&user)
	if err != nil {
		log.Printf("failed to find user by user_id: %v", err)
		return nil, err
	}

	// 'Group' 객체 생성
	newGroup := &model.Group{
		ID:        primitive.NewObjectID(),
		GroupName: name,
		Members: []model.GroupMember{
			{
				UserID:    userIDInt,
				Email:     user.Email,
				Status:    "accepted",
				InvitedBy: userIDInt,
			},
		},
		CreatedAt: time.Now(),
	}

	// 생성된 Group 저장
	group, err := s.repo.Create(ctx, newGroup)
	if err != nil {
		return nil, err
	}

	// 'users' collection의 'group_ids'에 추가
	filter := bson.M{"user_id": userIDInt}
	update := bson.M{
		"$push": bson.M{
			"group_ids": group.ID.Hex(),
		},
	}
	_, err = dbConfig.UserCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("failed to update user's group_ids: %v", err)
	}

	return group, nil
}


// GetGroupByID retrieves a group by its ID
func (s *GroupService) GetGroupByGroupID(ctx context.Context, id primitive.ObjectID) (*model.Group, error) {
	return s.repo.GetByID(ctx, id)
}

// UpdateGroup updates an existing group
func (s *GroupService) UpdateGroup(ctx context.Context, group *model.Group) error {
	return s.repo.Update(ctx, group)
}

// DeleteGroup deletes a group by its ID
func (s *GroupService) DeleteGroup(ctx context.Context, id primitive.ObjectID) error {
	return s.repo.Delete(ctx, id)
}

// ListGroups retrieves all groups
func (s *GroupService) ListGroups(ctx context.Context) ([]*model.Group, error) {
	return s.repo.List(ctx)
}

// GetGroupsByUserID retrieves groups by user ID
func (s *GroupService) GetGroupsByUserID(ctx context.Context, userID string) ([]*model.Group, error) {
	
	// extract user info
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		log.Printf("Invalid userID: %v", err)
		return nil, err
	}

	// 사용자의 'group_ids' 가져오기
	var user model.User
	err = dbConfig.UserCollection.FindOne(ctx, bson.M{"user_id": userIDInt}).Decode(&user)
	if err != nil {
		log.Printf("failed to find user by user_id: %v", err)
		return nil, err
	}

	// GroupIDs에 대한 처리 (ObjectID로 변환)
	var objectIDs []primitive.ObjectID
    for _, idStr := range user.GroupIDs {
        oid, err := primitive.ObjectIDFromHex(idStr)
        if err != nil {
            log.Printf("Invalid GroupID: %s, error: %v", idStr, err)
            continue
        }
        objectIDs = append(objectIDs, oid)
    }

	cursor, err := dbConfig.GroupCollection.Find(ctx, bson.M{"_id": bson.M{"$in": objectIDs}})

	if err != nil {
		log.Printf("Error querying groups: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var groups []*model.Group

	for cursor.Next(ctx) {
		var group model.Group
		if err := cursor.Decode(&group); err != nil {
			return nil, err
		}
		groups = append(groups, &group)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

// GetGroupMembers retrieves members of a group by its ID
func (s *GroupService) GetGroupMembers(ctx context.Context, groupID primitive.ObjectID) (*model.GroupMembersResponse, error) {
	return s.repo.GetGroupMembers(ctx, groupID)
}


// InviteUserToGroup invites a user to a group
func (s *GroupService) InviteUserToGroup(ctx context.Context, groupID primitive.ObjectID, inviterUserID int, inviteeEmail string) error {
	group, err := s.repo.GetByID(ctx, groupID)
	if err != nil {
		return fmt.Errorf("failed to find group: %v", err)
	}

	// 초대하려는 사용자가 이미 그룹에 있는지 확인
	for _, member := range group.Members {
		if member.Email == inviteeEmail {
			return fmt.Errorf("user is already a member of the group")
		}
	}

	// 사용자가 이미 초대받았는지 확인
	for _, member := range group.Members {
		if member.Email == inviteeEmail && member.Status == "pending" {
			return fmt.Errorf("user is already invited to the group")
		}
	}

	err = s.repo.InviteUser(ctx, groupID, inviterUserID, inviteeEmail)
	if err != nil {
		return fmt.Errorf("failed to invite user: %v", err)
	}

	return nil
}

// AcceptGroupInvite accepts an invitation to join a group
func (s *GroupService) AcceptGroupInvite(ctx context.Context, groupID primitive.ObjectID, userID int) error {
	group, err := s.repo.GetByID(ctx, groupID)
	if err != nil {
		return err
	}

	found := false
	for _, member := range group.Members {
		if member.UserID == userID && member.Status == "pending" {
			found = true
			break
		}
	}

	if !found {
		return errors.New("no pending invitation found for this user")
	}

	err = s.repo.AcceptInvite(ctx, groupID, userID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$push": bson.M{
			"group_ids": groupID.Hex(),
		},
	}

	_, err = dbConfig.UserCollection.UpdateOne(ctx, bson.M{"user_id": userID}, update)
	return err
}


// RejectGroupInvite rejects an invitation to join a group
func (s *GroupService) RejectGroupInvite(ctx context.Context, groupID primitive.ObjectID, userID int) error {
	group, err := s.repo.GetByID(ctx, groupID)
	if err != nil {
		return err
	}

	found := false
	for _, member := range group.Members {
		if member.UserID == userID && member.Status == "pending" {
			found = true
			break
		}
	}

	if !found {
		return errors.New("no pending invitation found for this user")
	}

	return s.repo.RejectInvite(ctx, groupID, userID)
}


// LeaveGroup allows a user to leave a group
func (s *GroupService) LeaveGroup(ctx context.Context, groupID primitive.ObjectID, userID int) error {
	err := s.repo.LeaveGroup(ctx, groupID, userID)
	if err != nil {
		return err
	}

	return s.removeGroupFromUser(ctx, userID, groupID)
}

func (s *GroupService) removeGroupFromUser(ctx context.Context, userID int, groupID primitive.ObjectID) error {
	update := bson.M{
		"$pull": bson.M{
			"group_ids": groupID.Hex(),
		},
	}

	_, err := dbConfig.UserCollection.UpdateOne(ctx, bson.M{"user_id": userID}, update)
	return err
}


// GetPendingGroups retrieves groups with pending invites for a specific user
func (s *GroupService) GetPendingGroups(ctx context.Context, userID int) ([]*model.Group, error) {
    return s.repo.GetPendingGroups(ctx, userID)
}