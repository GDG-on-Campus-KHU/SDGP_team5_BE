// situation/service/situation_service.go

package service

import (
	"context"
	"fmt"

	dbConfig "github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/config"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/situation/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GET by index and language
func GetSituationByIndex(index int, language string) (*model.Situation, error) {
	collection := dbConfig.Client.Database("resq").Collection("situation")

	var situation model.Situation
	filter := bson.M{"index": index}

	err := collection.FindOne(context.Background(), filter).Decode(&situation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("situation not found for index %d", index)
		}
		return nil, fmt.Errorf("error retrieving situation: %v", err)
	}

	situation.Description = filterLanguageContent(situation.Description, language)
	situation.Actions = filterActionSteps(situation.Actions, language)

	return &situation, nil
}

// GET by slug and language
func GetSituationBySlug(slug string, language string) (*model.Situation, error) {
	collection := dbConfig.Client.Database("resq").Collection("situation")

	var situation model.Situation
	filter := bson.M{"slug": slug}

	err := collection.FindOne(context.Background(), filter).Decode(&situation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("situation not found for slug %s", slug)
		}
		return nil, fmt.Errorf("error retrieving situation: %v", err)
	}

	situation.Description = filterLanguageContent(situation.Description, language)
	situation.Actions = filterActionSteps(situation.Actions, language)

	return &situation, nil
}

// filter based on the requested language
func filterLanguageContent(content model.MultilingualArray, language string) model.MultilingualArray {
	if content == nil {
		return nil
	}

	if langContent, exists := content[language]; exists {
		return model.MultilingualArray{language: langContent}
	}

	// default 'en'
	if langContent, exists := content["en"]; exists {
		return model.MultilingualArray{"en": langContent}
	}

	return nil
}

// filter based on the requested language
func filterActionSteps(content model.MultiLangActions, language string) model.MultiLangActions {
	if content == nil {
		return nil
	}

	if langContent, exists := content[language]; exists {
		return model.MultiLangActions{language: langContent}
	}

	// default 'en'
	if langContent, exists := content["en"]; exists {
		return model.MultiLangActions{"en": langContent}
	}

	return nil
}
