// user/repository/user_repository.go

package repository

import (
	"context"
)

type UserRepository interface {
	AddFavorite(ctx context.Context, userID int, situationIndex int) error
	GetFavorites(ctx context.Context, userID int) ([]int, error)
	DeleteFavorite(ctx context.Context, userID int, situationIndex int) error
}
