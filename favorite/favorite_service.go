// favorite/favorite_service.go

package favorite

import (
	"context"
	"log"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/user/repository"
)

type FavoriteService struct {
	repo repository.UserRepository
}

// constructor function
func NewFavoriteService(repo repository.UserRepository) *FavoriteService {
	if repo == nil {
		log.Fatal("UserRepository cannot be nil")
		return nil
	}

	return &FavoriteService{repo: repo}
}

// AddFavorite add a favorite to a favorite list
func (s *FavoriteService) AddFavorite(ctx context.Context, userID int, situationIndex int) error {
	return s.repo.AddFavorite(ctx, userID, situationIndex)
}
