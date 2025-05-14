// country/country_service.go

package country

import (
	"context"
	"log"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/country/repository"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
)

type CountryService struct {
	repo repository.CountryRepository
}

// constructor function
func NewCountryService(repo repository.CountryRepository) *CountryService {
	if repo == nil {
		log.Fatal("CountryRepository cannot be nil")
	}

	return &CountryService{repo: repo}
}

func (s *CountryService) GetCountryInfo(ctx context.Context, countryCode string) (*model.Country, error) {
	return s.repo.GetCountryInfo(ctx, countryCode)
}
