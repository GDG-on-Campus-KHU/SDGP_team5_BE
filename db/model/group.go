// db/model/group.go

package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupMember struct {
	UserID    int    `bson:"user_id" json:"user_id"`
	Email     string `bson:"email" json:"email"`
	Status    string `bson:"status" json:"status"`
	InvitedBy int    `bson:"invited_by" json:"invited_by"`
}


type Group struct {
	ID			primitive.ObjectID	`bson:"_id,omitempty" json:"id,omitempty"`                 // MongoDB ObjectId
	GroupName	string				`bson:"group_name" json:"group_name"`
	Members		[]GroupMember		`bson:"members" json:"members"`
	CreatedAt   time.Time	 		`bson:"created_at,omitempty" json:"created_at,omitempty"`	// timestamp
}


type GroupMembersResponse struct {
	GroupID primitive.ObjectID `json:"group_id"`
	Name    string             `json:"name"`
	Members []GroupMember      `json:"members"`
}