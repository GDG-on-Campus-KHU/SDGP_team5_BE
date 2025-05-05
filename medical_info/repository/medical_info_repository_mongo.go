// medical_info/repository/medical_info_repository_mongo.go

package repository

import (
	"context"
	"fmt"

	dbConfig "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/config"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type medicalInfoRepositoryMongo struct {
	collection *mongo.Collection
}

// constructor function
func NewMedicalInfoRepository() MedicalInfoRepository {
	return &medicalInfoRepositoryMongo{
		collection: dbConfig.MedicalInfoCollection,
	}
}

func (r *medicalInfoRepositoryMongo) CreateMedicalInfo(ctx context.Context, info *model.MedicalInfo) error {
	filter := bson.M{"info_id": info.InfoID}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("medical info with info_id %d already exists", info.InfoID)
	}

	result, err := r.collection.InsertOne(ctx, info)
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("failed to insert medical info")
	}

	return nil
}

func (r *medicalInfoRepositoryMongo) GetMedicalInfo(ctx context.Context, userID int) (*model.MedicalInfo, error) {
	var info model.MedicalInfo
	filter := bson.M{"info_id": userID}
	err := r.collection.FindOne(ctx, filter).Decode(&info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (r *medicalInfoRepositoryMongo) UpdateMedicalInfo(ctx context.Context, userID int, info *model.MedicalInfo) (*model.MedicalInfo, error) {
	update := bson.M{
		"$set": bson.M{
			"blood_type":  info.BloodType,
			"allergy":     info.Allergy,
			"medication":  info.Medication,
			"height":      info.Height,
			"height_unit": info.HeightUnit,
			"weight":      info.Weight,
			"weight_unit": info.WeightUnit,
			"birth_date":  info.BirthDate,
			"notes":       info.Notes,
			"updated_at":  info.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"info_id": userID}, update)
	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("no medical info found for user ID %d", userID)
	}

	var updatedInfo model.MedicalInfo
	err = r.collection.FindOne(ctx, bson.M{"info_id": userID}).Decode(&updatedInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated medical info: %w", err)
	}

	return &updatedInfo, nil
}
