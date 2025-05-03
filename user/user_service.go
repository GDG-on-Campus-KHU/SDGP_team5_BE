// user/user_service.go

package user

import (
	"context"
	"log"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/user/repository"
)

type UserService struct {
	repo repository.UserRepository
}

// constructor function
func NewUserService(repo repository.UserRepository) *UserService {
	if repo == nil {
		log.Fatal("UserRepository cannot be nil")
	}

	return &UserService{repo: repo}
}

func (s *UserService) GetUserInfo(ctx context.Context, userID int) (*repository.UserInfoResponse, error) {
	return s.repo.GetUserInfo(ctx, userID)
}

func (s *UserService) UpdateCountry(ctx context.Context, userID int, countryCode string) (*repository.UserInfoResponse, error) {
	return s.repo.UpdateCountry(ctx, userID, countryCode)
}
