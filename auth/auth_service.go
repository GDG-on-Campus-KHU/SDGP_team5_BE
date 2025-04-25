// auth/auth_service.go

package auth

import (
	"context"
	"log"

	dbConfig "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/config"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := dbConfig.UserCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}


func CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	userID, err := util.GetNextID(ctx, dbConfig.Client.Database("resq"), "users")
	if err != nil {
		log.Printf("Error getting next user ID: %v", err)
		return nil, err
	}

	user.UserID = userID
	collection := dbConfig.UserCollection

	// User 생성
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	return user, nil
}