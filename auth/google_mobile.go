// auth/google_mobile.go

package auth

import (
	"context"
	"fmt"
	"os"
	"time"
	"net/http"
	"encoding/json"
	"strconv"
	
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gin-gonic/gin"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
	dbConfig "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/config"
)

// Google OAuth2 for Android mobile application
var GoogleMobileOAuthConfig = &oauth2.Config{
    ClientID:     os.Getenv("GOOGLE_WEB_AUTH_CLIENT_ID"),
    ClientSecret: os.Getenv("GOOGLE_WEB_AUTH_CLIENT_SECRET"),
    RedirectURL:  "",
    Scopes: []string{
        "openid", "email", "profile",
    },
    Endpoint: google.Endpoint,
}


type TokenClaims struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	jwt.StandardClaims
}


func generateGoogleToken(userID, name, email string, expiration time.Duration) (string, error) {
    secretKey := os.Getenv("JWT_SECRET_KEY")
    claims := &TokenClaims{
        UserID: userID,
        Name:   name,
        Email:  email,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(expiration).Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secretKey))
}


// POST /api/auth/login/google
func GoogleLoginHandler(c *gin.Context) {
    // `serverAuthCode` from Android mobile client
    var req struct {
        ServerAuthCode string `json:"serverAuthCode"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    // token exchange with 'serverAuthCode'
    token, err := GoogleMobileOAuthConfig.Exchange(context.Background(), req.ServerAuthCode)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to exchange token: %v", err)})
        return
    }

	idToken := ""
	if v := token.Extra("id_token"); v != nil {
		if str, ok := v.(string); ok {
			idToken = str
		}
	}

    client := GoogleMobileOAuthConfig.Client(context.Background(), token)
    resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to get user info: %v", err)})
        return
    }
    defer resp.Body.Close()

    var userInfo struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to decode user info: %v", err)})
        return
    }

    // 사용자가 이미 DB에 존재하는지 확인
    user, err := GetUserByEmail(userInfo.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to fetch user: %v", err)})
        return
    }

    if user == nil {
        // 사용자 정보가 없다면 새로운 사용자 생성
        newUser := &model.User{
            Name:  userInfo.Name,
            Email: userInfo.Email,
        }

        ctx := c.Request.Context()
        user, err = CreateUser(ctx, newUser)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create user: %v", err)})
            return
        }
    }

    // JWT 생성
    userID := strconv.Itoa(user.UserID)
    accessToken, err := generateGoogleToken(userID, user.Name, user.Email, time.Hour*72)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to generate access token: %v", err)})
        return
    }

    refreshToken := token.RefreshToken

    authToken := &model.AuthToken{
        GoogleUserID:	idToken,
        AccessToken:	accessToken,
        RefreshToken:	refreshToken,
        CreatedAt:    	time.Now(),
        UpdatedAt:		time.Now(),
        ExpireAt:		time.Now().Add(time.Hour * 72),
        Email:			userInfo.Email,
        UserID:       	userID,
    }

    ctx := c.Request.Context()
    if _, err := dbConfig.AuthTokenCollection.InsertOne(ctx, authToken); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to save auth token: %v", err)})
        return
    }

    // return access token
    c.JSON(http.StatusOK, gin.H{
        "access_token":  accessToken,
    })
}
