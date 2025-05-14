// country/repository/country_repository.go

package repository

import (
	"context"
	"fmt"

	dbConfig "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/config"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type countryRepositoryMongo struct {
	collection *mongo.Collection
}

// constructor function
func NewCountryRepository() CountryRepository {
	return &countryRepositoryMongo{
		collection: dbConfig.CountryCollection,
	}
}

// User methods
func (r *countryRepositoryMongo) GetCountryInfo(ctx context.Context, countryCode string) (*model.Country, error) {
	var country model.Country
	filter := bson.M{"country_code": countryCode}
	err := r.collection.FindOne(ctx, filter).Decode(&country)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("country with country_code %s not found", countryCode)
		}
		return nil, err
	}
	return &country, nil
}
