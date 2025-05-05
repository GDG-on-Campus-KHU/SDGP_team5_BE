// favorite/favorite_service.go

package favorite

import (
	"context"
	"log"

	situationRepository "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/situation/repository"
	userRepository "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/user/repository"
)

type FavoriteService struct {
	userRepo      userRepository.UserRepository
	situationRepo situationRepository.SituationRepository
}

// constructor function
func NewFavoriteService(userRepo userRepository.UserRepository, situationRepo situationRepository.SituationRepository) *FavoriteService {
	if userRepo == nil {
		log.Fatal("UserRepository cannot be nil")
		return nil
	}

	if situationRepo == nil {
		log.Fatal("SituationRepository cannot be nil")
		return nil
	}

	return &FavoriteService{userRepo: userRepo, situationRepo: situationRepo}
}

// AddFavorite add a favorite to a favorite list
func (s *FavoriteService) AddFavorite(ctx context.Context, userID int, situationIndex int) error {
	return s.userRepo.AddFavorite(ctx, userID, situationIndex)
}

// GetFavorites get a favorite list
func (s *FavoriteService) GetFavorites(ctx context.Context, userID int) (*[]situationRepository.FavoriteResponse, error) {
	favorites, err := s.userRepo.GetFavorites(ctx, userID)
	if err != nil {
		return nil, err
	}

	situations := make([]situationRepository.FavoriteResponse, len(favorites))
	for i, situationIndex := range favorites {
		situation, err := s.situationRepo.GetSituationByIndex(ctx, situationIndex)
		if err != nil {
			return nil, err
		}
		situations[i] = situationRepository.FavoriteResponse{
			SituationIndex: situation.Index,
			Situation:      situation.EmerTitle,
			Slug:           situation.Slug,
		}
	}

	return &situations, nil
}

// DeleteFavorite get a favorite list
func (s *FavoriteService) DeleteFavorite(ctx context.Context, userID int, situationIndex int) error {
	return s.userRepo.DeleteFavorite(ctx, userID, situationIndex)
}
