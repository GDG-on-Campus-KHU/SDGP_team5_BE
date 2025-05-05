// db/model/medical_info.go

package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BloodType string

type WeightUnit string
type HeightUnit string

const (
	BloodTypeAPlus   BloodType = "A+"
	BloodTypeAMinus  BloodType = "A-"
	BloodTypeBPlus   BloodType = "B+"
	BloodTypeBMinus  BloodType = "B-"
	BloodTypeOPlus   BloodType = "O+"
	BloodTypeOMinus  BloodType = "O-"
	BloodTypeABPlus  BloodType = "AB+"
	BloodTypeABMinus BloodType = "AB-"
	BloodTypeOthers  BloodType = "Others"
)

const (
	WeightUnitKg    WeightUnit = "kg"
	WeightUnitPound WeightUnit = "pound"
	HeightUnitCm    HeightUnit = "cm"
	HeightUnitInch  HeightUnit = "inch"
	HeightUnitFeet  HeightUnit = "feet"
)

type MedicalInfo struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // MongoDB ObjectId
	InfoID     int                `bson:"info_id" json:"info_id"`
	BloodType  BloodType          `bson:"blood_type" json:"blood_type"`
	Allergy    string             `bson:"allergy" json:"allergy"`
	Medication string             `bson:"medication" json:"medication"`
	Height     float64            `bson:"height" json:"height"`
	HeightUnit HeightUnit         `bson:"height_unit" json:"height_unit"`
	Weight     float64            `bson:"weight" json:"weight"`
	WeightUnit WeightUnit         `bson:"weight_unit" json:"weight_unit"`
	BirthDate  string             `bson:"birth_date" json:"birth_date"`
	Notes      string             `bson:"notes" json:"notes"`
	CreatedAt  time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"` // timestamp
	UpdatedAt  time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"` // timestamp
}
