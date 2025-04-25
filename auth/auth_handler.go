// auth/auth_handler.go

package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
)

// GET /api/auth/login

// LoginHandler handles the Google OAuth2 login flow.
// @Summary Google OAuth2 Login
// @Description Redirects the user to Google login page for authentication.
// @Tags auth
// @Produce json
// @Success 307 {string} string "Redirect to Google OAuth2 login"
// @Router /api/auth/login [get]
func LoginHandler(c *gin.Context) {
	state := GenerateState()
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

	userID := strconv.Itoa(user.UserID)
	token, err := generateToken(userID, user.Name, user.Email, time.Hour*72)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}


// GET /api/auth/protected

// ProtectedHandler handles JWT-authenticated access.
// @Summary Protected route that requires a valid JWT token
// @Description Returns user info if the provided JWT token is valid. If the token is invalid or expired, the request is unauthorized.
// @Tags auth
// @Security BearerAuth  // This indicates the need for a Bearer token in the Authorization header
// @Produce json
// @Success 200 {object} map[string]string "Authorized"  // Successful response with user info
// @Failure 401 {object} map[string]string "Unauthorized"  // Unauthorized if token is invalid or missing
// @Router /api/auth/protected [get]
func ProtectedHandler(c *gin.Context) {
    user, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": "user not found",
        })
        return
    }

	if userStr, ok := user.(string); ok {
		c.JSON(http.StatusOK, gin.H{
			"message": "Authorized",
			"user":    userStr,
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user type",
		})
	}
}


// POST /api/auth/refresh-token

// RefreshTokenHandler handles the refresh token logic.
// @Summary Refreshes the access token using a valid refresh token
// @Description Accepts a valid refresh token and issues a new access token if the refresh token is valid. The refresh token should be passed in the request body.
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body string true "The refresh token used to generate a new access token"
// @Success 200 {object} map[string]string "access_token"  // Returns the newly generated access token
// @Failure 400 {object} map[string]string "Invalid request"  // Invalid or malformed request body
// @Failure 401 {object} map[string]string "Invalid refresh token"  // If the refresh token is invalid
// @Failure 500 {object} map[string]string "Could not generate access token"  // If there is an issue generating the access token
// @Router /api/auth/refresh-token [post]
func RefreshTokenHandler(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// refresh token validation
	claims, err := ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// issue new token
	accessToken, err := GenerateAccessToken(claims.Subject, claims.Name, claims.Email, time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}