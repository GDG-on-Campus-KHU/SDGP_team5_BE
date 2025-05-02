// user/repository/user_repository.go

package repository

import (
	"context"
	"fmt"

	dbConfig "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepositoryMongo struct {
	collection *mongo.Collection
}

// constructor function
func NewUserRepository() UserRepository {
	return &userRepositoryMongo{
		collection: dbConfig.UserCollection,
	}
}

func (r *userRepositoryMongo) AddFavorite(ctx context.Context, userID int, situationIndex int) error {
	update := bson.M{
		"$addToSet": bson.M{
			"favorites": situationIndex,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"user_id": userID}, update)

	if result.ModifiedCount == 0 {
		return fmt.Errorf("already in favorites")
	}
	return err
}
