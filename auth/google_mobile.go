// auth/google_mobile.go

package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	dbConfig "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/config"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
	dbUtil "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/util"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/util"
)

type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
}

type GoogleIDTokenPayload struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type GoogleCustomClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func generateAccessToken(userID, name, email string, expiration time.Duration) (string, error) {
	claims := &GoogleCustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		},
		Name:  name,
		Email: email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_SECRET_KEY")

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// renew the access token
func generateNewToken(userID, name, email string) (string, error) {
	return generateAccessToken(userID, name, email, 72*time.Hour)
}

// POST /api/auth/login/google
func GoogleLoginHandler(c *gin.Context) {
	var req struct {
		ServerAuthCode string `json:"serverAuthCode"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	clientID := os.Getenv("GOOGLE_AUTH_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_AUTH_CLIENT_SECRET")

	// fmt.Println("🔄 Exchanging authorization code for token...")

	tokenResp, err := exchangeAuthCode(req.ServerAuthCode, clientID, clientSecret)
	if err != nil {
		fmt.Printf("❌ Token exchange failed: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token exchange failed"})
		return
	}

	// fmt.Printf("✅ Token Response: %+v\n", tokenResp)

	userInfo, err := parseIDToken(tokenResp.IdToken)
	if err != nil {
		fmt.Printf("❌ Failed to parse id_token: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse id_token"})
		return
	}

	// 데이터베이스에 저장된 사용자인지 확인
	user, err := util.GetUserByEmail(c, userInfo.Email)
	if err != nil {

		// 새로운 사용자 추가
		user, err = AddNewUser(c, userInfo)
		if err != nil {
			fmt.Printf("❌ Failed to create new user: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create new user"})
			return
		}
	}

	accessToken, err := generateAccessToken(userInfo.Sub, userInfo.Name, userInfo.Email, 24*time.Hour)
	if err != nil {
		fmt.Printf("❌ Failed to generate access token: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	refreshToken, err := generateNewToken(userInfo.Sub, userInfo.Name, userInfo.Email)
	if err != nil {
		fmt.Printf("❌ Failed to generate refresh token: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	fmt.Printf("✅ Access Token: %s\n", accessToken)
	fmt.Printf("✅ Refresh Token: %s\n", refreshToken)

	// token 저장
	authToken := model.AuthToken{
		GoogleUserID: userInfo.Sub,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ExpireAt:     time.Now().Add(24 * time.Hour),
		Email:        userInfo.Email,
		UserID:       user.UserID,
	}

	if err := saveOrUpdateAuthToken(authToken); err != nil {
		log.Printf("Failed to save auth token: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store auth token"})
		return
	}

	// API response
	c.JSON(http.StatusOK, gin.H{
		"token": accessToken,
	})
}

func exchangeAuthCode(authCode, clientID, clientSecret string) (*GoogleTokenResponse, error) {
	values := map[string]string{
		"code":          authCode,
		"client_id":     clientID,
		"client_secret": clientSecret,
		"redirect_uri":  os.Getenv("GOOGLE_AUTH_REDIRECT_URL"),
		// "redirect_uri": "http://localhost:5100/oauth2callback",  // local test
		"grant_type": "authorization_code",
	}

	log.Println("🔄 Payload to Google:", values)
	data, _ := json.Marshal(values)

	fmt.Println("🔄 Sending request to Google API to exchange auth code for token...")
	resp, err := http.Post("https://oauth2.googleapis.com/token", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Google API: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("👉 Google API Response Status Code: %d\n", resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Println("Google API Response Body:", string(body))
		return nil, fmt.Errorf("token exchange failed: %s", string(body))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed: %s", string(body))
	}

	var tokenRes GoogleTokenResponse
	if err := json.Unmarshal(body, &tokenRes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token response: %w", err)
	}

	return &tokenRes, nil
}

func parseIDToken(idToken string) (*GoogleIDTokenPayload, error) {
	token, _, err := jwt.NewParser().ParseUnverified(idToken, &GoogleIDTokenPayload{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse id_token: %w", err)
	}

	payload, ok := token.Claims.(*GoogleIDTokenPayload)
	if !ok {
		return nil, fmt.Errorf("invalid id_token payload")
	}
	return payload, nil
}

func storeRefreshToken(userID, token string) error {
	fmt.Printf("storing refresh token for user %s: %s\n", userID, token)
	return nil
}

// skipping validation check
func (p *GoogleIDTokenPayload) Valid() error {
	return nil
}

// POST /api/auth/refresh-token
func GoogleRefreshTokenHandler(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"` // previous access token
	}
	if err := c.BindJSON(&req); err != nil || req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, _, err := jwt.NewParser().ParseUnverified(req.RefreshToken, &GoogleCustomClaims{})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
		return
	}

	claims, ok := token.Claims.(*GoogleCustomClaims)
	if !ok || claims.Subject == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	userID := claims.Subject

	// storing tokens

	newAccessToken, err := generateAccessToken(userID, claims.Name, claims.Email, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate new access token"})
		return
	}

	// 2. Update token in DB
	err = UpdateAccessTokenByGoogleUserID(userID, newAccessToken, time.Now().Add(24*time.Hour))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update token in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}

func AddNewUser(ctx context.Context, userInfo *GoogleIDTokenPayload) (*model.User, error) {

	// user_id (auto-increment)
	userID, err := dbUtil.GetNextID(ctx, dbConfig.Client.Database("resq"), "users")
	if err != nil {
		return nil, fmt.Errorf("failed to get next user ID: %v", err)
	}

	// 신규 사용자 저장을 위한 객체 초기화
	user := &model.User{
		UserID:      userID,
		Name:        userInfo.Name,
		Email:       userInfo.Email,
		AppLang:     "ko", // default
		CountryCode: "KR", // default
		GroupIDs:    []string{},
		Favorites:   []int32{},
	}

	// 데이터베이스에 신규 사용자 저장
	_, err = dbConfig.UserCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to insert new user: %v", err)
	}

	return user, nil
}

func saveOrUpdateAuthToken(token model.AuthToken) error {
	filter := bson.M{"email": token.Email}

	update := bson.M{
		"$set": bson.M{
			"access_token":   token.AccessToken,
			"refresh_token":  token.RefreshToken,
			"updated_at":     time.Now(),
			"expire_at":      token.ExpireAt,
			"google_user_id": token.GoogleUserID,
			"user_id":        token.UserID,
		},
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
	}

	opts := options.Update().SetUpsert(true)

	_, err := dbConfig.AuthTokenCollection.UpdateOne(context.TODO(), filter, update, opts)
	return err
}

func UpdateAccessTokenByGoogleUserID(googleUserID, newToken string, expireAt time.Time) error {
	ctx := context.TODO()
	filter := bson.M{"google_user_id": googleUserID}
	update := bson.M{
		"$set": bson.M{
			"access_token": newToken,
			"updated_at":   time.Now(),
			"expire_at":    expireAt,
		},
	}

	_, err := dbConfig.AuthTokenCollection.UpdateOne(ctx, filter, update)
	return err
}
