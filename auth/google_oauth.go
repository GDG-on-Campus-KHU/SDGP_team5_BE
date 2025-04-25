// auth/google_oauth.go

package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOAuthConfig *oauth2.Config

func InitGoogleOAuthConfig() {
	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_AUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_AUTH_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_AUTH_REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	fmt.Println("OAuth Config Initialized: ", GoogleOAuthConfig.ClientID)
}

func GetGoogleAuthURL(state string) string {
	fmt.Println("Client ID:", GoogleOAuthConfig.ClientID)
	return GoogleOAuthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

type GoogleUserInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func GetGoogleUserInfo(code string) (*GoogleUserInfo, error) {
	token, err := GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	client := GoogleOAuthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}