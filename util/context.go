// util/context.go

package util

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	dbConfig "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/config"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
)

// const UserIDKey = "user_id"

func GetUserIDFromContext(c *gin.Context) (string, error) {
	name, exists := c.Get("user")
	if !exists {
		return "", fmt.Errorf("user not found in context")
	}

	nameStr, ok := name.(string)
	if !ok {
		return "", fmt.Errorf("user name is not a string")
	}

	var user model.User
	err := dbConfig.UserCollection.FindOne(context.TODO(), bson.M{"name": nameStr}).Decode(&user)
	if err != nil {
		return "", fmt.Errorf("failed to find user by name: %v", err)
	}

	// return user.UserID, nil
	return strconv.Itoa(user.UserID), nil
}