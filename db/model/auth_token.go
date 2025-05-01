// db/model/auth_token.go

package model

import "time"

type AuthToken struct {
	GoogleUserID string			`bson:"google_user_id"`
	AccessToken  string  	  	`bson:"access_token"`
	RefreshToken string   	 	`bson:"refresh_token"`
	CreatedAt    time.Time 		`bson:"created_at"`
	UpdatedAt    time.Time 		`bson:"updated_at"`
	ExpireAt     time.Time 		`bson:"expire_at"`
	Email        string			`bson:"email"`
	UserID       string   		`bson:"user_id"`
}