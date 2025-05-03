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

type Favorites struct {
	Favorites []int `bson:"favorites"`
}

type UserInfoResponse struct {
	UserID      int    `bson:"user_id"`
	CountryCode string `bson:"country_code"`
	Name        string `bson:"name"`
	Email       string `bson:"email"`
	AppLang     string `bson:"app_lang"`
}

// constructor function
func NewUserRepository() UserRepository {
	return &userRepositoryMongo{
		collection: dbConfig.UserCollection,
	}
}

// User methods
func (r *userRepositoryMongo) GetUserInfo(ctx context.Context, userID int) (*UserInfoResponse, error) {
	var userInfoResponse UserInfoResponse
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&userInfoResponse)
	if err != nil {
		return nil, err
	}
	return &userInfoResponse, nil
}

func (r *userRepositoryMongo) UpdateCountry(ctx context.Context, userID int, countryCode string) (*UserInfoResponse, error) {
	update := bson.M{
		"$set": bson.M{
			"country_code": countryCode,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"user_id": userID}, update)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("user with user_id %d not found", userID)
	}
	if result.ModifiedCount == 0 {
		return nil, fmt.Errorf("user with user_id %d already has country_code %s", userID, countryCode)
	}
	var userInfoResponse UserInfoResponse
	err = r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&userInfoResponse)
	if err != nil {
		return nil, err
	}
	return &userInfoResponse, nil
}

// Favorite methods
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
	var favorites Favorites
	err := r.collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&favorites)
	if err != nil {
		return nil, err
	}

	return favorites.Favorites, nil
}

func (r *userRepositoryMongo) DeleteFavorite(ctx context.Context, userID int, situationIndex int) error {
	filter := bson.M{
		"user_id": userID,
	}

	update := bson.M{
		"$pull": bson.M{
			"favorites": situationIndex,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)

	if result.ModifiedCount == 0 {
		return fmt.Errorf("not in favorites")
	}
	return err
}
