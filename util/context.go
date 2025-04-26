// util/context.go

package util

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	dbConfig "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/config"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
)


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
	err := dbConfig.UserCollection.FindOne(c.Request.Context(), bson.M{"name": nameStr}).Decode(&user)
	if err != nil {
		return "", fmt.Errorf("failed to find user by name: %v", err)
	}

	return strconv.Itoa(user.UserID), nil
}