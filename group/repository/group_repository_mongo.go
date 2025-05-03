// group/repository/group_repository_mongo.go

package repository

import (
	"context"
	"time"
	"fmt"

	dbConfig "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/config"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type groupRepositoryMongo struct {
	collection *mongo.Collection
}

// constructor function
func NewGroupRepository() GroupRepository {
	return &groupRepositoryMongo{
		collection: dbConfig.GroupCollection,
	}
}

func (r *groupRepositoryMongo) Create(ctx context.Context, group *model.Group) (*model.Group, error) {
	group.ID = primitive.NewObjectID()
	group.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, group)
	return group, err
}

func (r *groupRepositoryMongo) GetByID(ctx context.Context, id primitive.ObjectID) (*model.Group, error) {
	var group model.Group
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&group)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *groupRepositoryMongo) Update(ctx context.Context, group *model.Group) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": group.ID},
		bson.M{"$set": group},
	)
	return err
}

func (r *groupRepositoryMongo) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *groupRepositoryMongo) List(ctx context.Context) ([]*model.Group, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
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
	return groups, nil
}

func (r *groupRepositoryMongo) GetGroupMembers(ctx context.Context, groupID primitive.ObjectID) (*model.GroupMembersResponse, error) {
	var group model.Group
	err := r.collection.FindOne(ctx, bson.M{"_id": groupID}).Decode(&group)
	if err != nil {
		return nil, err
	}

	return &model.GroupMembersResponse{
		GroupID: group.ID,
		Name:    group.GroupName,
		Members: group.Members,
	}, nil
}


func (r *groupRepositoryMongo) InviteUser(ctx context.Context, groupID primitive.ObjectID, inviterUserID int, inviteeEmail string) error {
	update := bson.M{
		"$push": bson.M{
			"members": bson.M{
				"user_id":   0,
				"email":     inviteeEmail,
				"status":    "pending",
				"invited_by": inviterUserID,
			},
		},
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": groupID}, update)
	return err
}


func (r *groupRepositoryMongo) AcceptInvite(ctx context.Context, groupID primitive.ObjectID, userID int) error {
	filter := bson.M{
		"_id": groupID,
		"members.email": bson.M{"$exists": true},
	}

	update := bson.M{
		"$set": bson.M{
			"members.$[elem].status": "accepted",
			"members.$[elem].user_id": userID,
		},
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.status": "pending", "elem.user_id": 0},
		},
	})

	_, err := r.collection.UpdateOne(ctx, filter, update, arrayFilters)
	return err
}


func (r *groupRepositoryMongo) RejectInvite(ctx context.Context, groupID primitive.ObjectID, userID int) error {
	filter := bson.M{
		"_id": groupID,
	}

	update := bson.M{
		"$pull": bson.M{
			"members": bson.M{
				"user_id": userID,
				"status":  "pending",
			},
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *groupRepositoryMongo) LeaveGroup(ctx context.Context, groupID primitive.ObjectID, userID int) error {
	filter := bson.M{
		"_id": groupID,
	}

	update := bson.M{
		"$pull": bson.M{
			"members": bson.M{
				"user_id": userID,
			},
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}


func (r *groupRepositoryMongo) GetUserEmailByID(ctx context.Context, groupID primitive.ObjectID, userID int) (string, error) {
	var group model.Group
	err := r.collection.FindOne(ctx, bson.M{"_id": groupID}).Decode(&group)
	if err != nil {
		return "", fmt.Errorf("failed to find group: %v", err)
	}

	for _, member := range group.Members {
		if member.UserID == userID {
			return member.Email, nil
		}
	}

	return "", fmt.Errorf("user with user_id %d not found in group", userID)
}