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

type User struct {
	Favorites []int `bson:"favorites"`
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

func (r *userRepositoryMongo) GetFavorites(ctx context.Context, userID int) ([]int, error) {
	var user User
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user.Favorites, nil
}
