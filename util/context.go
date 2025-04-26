// util/context.go

package util

import (
	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"

func GetUserIDFromContext(c *gin.Context) (int, bool) {
	value, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}

	userID, ok := value.(int)
	return userID, ok
}