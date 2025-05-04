// user/repository/user_repository.go

package repository

import (
	"context"
)

type UserRepository interface {
	// User methods
	GetUserInfo(ctx context.Context, userID int) (*UserInfoResponse, error)
	UpdateCountry(ctx context.Context, userID int, countryCode string) (*UserInfoResponse, error)

	// Favorite methods
	AddFavorite(ctx context.Context, userID int, situationIndex int) error
	GetFavorites(ctx context.Context, userID int) ([]int, error)
	DeleteFavorite(ctx context.Context, userID int, situationIndex int) error
}
