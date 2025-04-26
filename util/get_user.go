// util/get_user.go

package util

import (
	"context"
	"errors"

	dbConfig "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/config"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"

	"go.mongodb.org/mongo-driver/bson"
)

func GetUserByUserID(ctx context.Context, userID int) (*model.User, error) {
	filter := bson.M{"user_id": userID}

	var user model.User
	err := dbConfig.UserCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}