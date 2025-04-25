// db/model/group.go

package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Group struct {
	ID			primitive.ObjectID	`bson:"_id,omitempty" json:"id,omitempty"`                 // MongoDB ObjectId
	GroupID		int      			`bson:"group_id" json:"group_id"`
	GroupName	string				`bson:"group_name" json:"group_name"`
	Members		[]int				`bson:"members" json:"members"`								// ref. struct 'User'
	CreatedAt   primitive.DateTime 	`bson:"created_at,omitempty" json:"created_at,omitempty"`	// timestamp
}