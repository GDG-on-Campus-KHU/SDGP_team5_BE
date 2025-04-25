// auth/auth_service.go

package auth

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	dbConfig "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/config"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/util"
)

func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := dbConfig.UserCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}


func CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	// user_id (auto-increment)
	userID, err := util.GetNextID(ctx, dbConfig.Client.Database("resq"), "users")
	if err != nil {
		log.Printf("Error getting next user ID: %v", err)
		return nil, err
	}

	user.UserID = userID

	// set default values
	if user.GroupIDs == nil {
		user.GroupIDs = []int{}
	}
	if user.Favorites == nil {
		user.Favorites = []int{}
	}
	if user.AppLang == "" {
		user.AppLang = "ko"
	}
	if user.CountryCode == "" {
		user.CountryCode = "KR"
	}
	user.InfoID = userID

	// save user
	collection := dbConfig.UserCollection
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	return user, nil
}


// generate a random state string
func GenerateState() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	state := make([]byte, 16)
	for i := range state {
		state[i] = charset[rand.Intn(len(charset))]
	}
	return string(state)
}



func ValidateRefreshToken(tokenString string) (*CustomClaims, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// HMAC: Hash-based Message Authentication Code
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(*CustomClaims); ok {
		return claims, nil
	}
	return nil, fmt.Errorf("could not parse claims")
}