// db/util/counter.go

package util

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
)

func GetNextID(ctx context.Context, db *mongo.Database, collectionName string) (int, error) {
	collection := db.Collection(collectionName)

	var lastUser model.User
	cursor, err := collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "user_id", Value: -1}}).SetLimit(1))
	if err != nil {
		return 0, fmt.Errorf("error finding the last user ID: %v", err)
	}
	defer cursor.Close(ctx)

	// DB에 document가 존재하지 않으면 1부터 시작
	if !cursor.Next(ctx) {
		fmt.Println("No documents found, starting user_id from 1")
		return 1, nil
	}

	if err := cursor.Decode(&lastUser); err != nil {
		return 0, fmt.Errorf("error decoding last user document: %v", err)
	}

	fmt.Printf("Last user ID: %v\n", lastUser.UserID)

	return int(lastUser.UserID) + 1, nil
}