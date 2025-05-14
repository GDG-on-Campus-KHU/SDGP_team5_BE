// country/repository/country_repository.go

package repository

import (
	"context"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
)

type CountryRepository interface {
	GetCountryInfo(ctx context.Context, countryCode string) (*model.Country, error)
}
