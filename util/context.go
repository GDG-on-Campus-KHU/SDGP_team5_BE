// util/context.go

package util

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"

func GetUserIDFromContext(c *gin.Context) (int, error) {
	value, exists := c.Get(UserIDKey)
	if !exists {
		return 0, fmt.Errorf("user_id not found in context")
	}

	userID, ok := value.(int)
	if !ok {
		return 0, fmt.Errorf("user_id is not of type int")
	}

	return userID, nil
}