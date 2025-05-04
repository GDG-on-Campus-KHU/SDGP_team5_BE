// util/get_user.go

package util

import (
	"context"
	"errors"
	"fmt"

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


func GetMedicalInfoByID(ctx context.Context, infoID int32) (*model.MedicalInfo, error) {
	filter := bson.M{"info_id": infoID}

	var medicalInfo model.MedicalInfo
	err := dbConfig.MedicalInfoCollection.FindOne(ctx, filter).Decode(&medicalInfo)
	if err != nil {
		return nil, errors.New("medical info not found")
	}

	return &medicalInfo, nil
}


func GetUserByIDFromRequest(ctx context.Context, userID int) (*model.User, error) {
	var user model.User

	err := dbConfig.UserCollection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	return &user, nil
}