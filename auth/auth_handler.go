// auth/auth_handler.go

package auth

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
)

// generate a random state string
func generateState() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	state := make([]byte, 16)
	for i := range state {
		state[i] = charset[rand.Intn(len(charset))]
	}
	return string(state)
}


// GET /api/auth/login

// LoginHandler handles the Google OAuth2 login flow.
// @Summary Google OAuth2 Login
// @Description Redirects the user to Google login page for authentication.
// @Tags auth
// @Produce json
// @Success 307 {string} string "Redirect to Google OAuth2 login"
// @Router /api/auth/login [get]
func LoginHandler(c *gin.Context) {
	state := generateState()
	url := GoogleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Println("OAuth URL:", url)
	c.Redirect(http.StatusTemporaryRedirect, url)
}


// GET /api/auth/callback

// CallbackHandler handles the OAuth2 callback after Google authentication.
// @Summary Google OAuth2 Callback
// @Description Handles the callback from Google after user grants permission, generates JWT.
// @Tags auth
// @Produce json
// @Param code query string true "Authorization Code"
// @Success 200 {object} map[string]string "token"
// @Failure 400 {object} map[string]string "Missing code"
// @Failure 500 {object} map[string]string "Failed to fetch user info or generate JWT"
// @Router /api/auth/callback [get]
func CallbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code in URL"})
		return
	}

	userInfo, err := GetGoogleUserInfo(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}

	user, err := GetUserByEmail(userInfo.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}

	if user == nil {
		newUser := &model.User{
			Name:  userInfo.Name,
			Email: userInfo.Email,
		}

		ctx := c.Request.Context()
		user, err = CreateUser(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	token, err := GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}


// GET /api/auth/protected

// ProtectedHandler handles JWT-authenticated access.
// @Summary Protected route with JWT
// @Description Returns user info if the provided JWT token is valid.
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string "Authorized"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/auth/protected [get]
func ProtectedHandler(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "user not found",
        })
        return
    }

	switch v := user.(type) {
	case string:
		c.JSON(http.StatusOK, gin.H{
			"message": "Authorized",
			"user":    v,
		})
	case float64:
		c.JSON(http.StatusOK, gin.H{
			"message": "Authorized",
			"user":    fmt.Sprintf("%v", v),
		})
	default:
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user type",
		})
	}
}