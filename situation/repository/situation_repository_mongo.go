// situation/repository/situation_repository_mongo.go

package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	dbConfig "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/config"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
)

type situationRepositoryMongo struct {
	collection *mongo.Collection
}

type FavoriteResponse struct {
	SituationIndex int                    `bson:"situation_index"`
	Situation      model.MultilingualText `bson:"situation"`
	Slug           string                 `bson:"slug"`
}

// constructor function
func NewSituationRepository() SituationRepository {
	return &situationRepositoryMongo{
		collection: dbConfig.SituationCollection,
	}
}

func (r *situationRepositoryMongo) GetSituationByIndex(ctx context.Context, situationIndex int) (*model.Situation, error) {
	var situation model.Situation
	err := r.collection.FindOne(ctx, bson.M{"index": situationIndex}).Decode(&situation)
	if err != nil {
		log.Println("Trying to find situation with index:", situationIndex)
		return nil, err
	}
	return &situation, nil
}
