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
	// 혈액형
	BloodTypeAPlus   BloodType = "A+"
	BloodTypeAMinus  BloodType = "A-"
	BloodTypeBPlus   BloodType = "B+"
	BloodTypeBMinus  BloodType = "B-"
	BloodTypeOPlus   BloodType = "O+"
	BloodTypeOMinus  BloodType = "O-"
	BloodTypeABPlus  BloodType = "AB+"
	BloodTypeABMinus BloodType = "AB-"
	BloodTypeOthers  BloodType = "Others"

	// 키 단위 (cm, in, ft)
	HeightUnitCm   HeightUnit = "cm"
	HeightUnitInch HeightUnit = "in"
	HeightUnitFeet HeightUnit = "ft"
	// 몸무게 단위
	WeightUnitKg    WeightUnit = "kg"
	WeightUnitPound WeightUnit = "lb"
)

type MedicalInfo struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // MongoDB ObjectId
	InfoID     int                `bson:"info_id" json:"info_id"`
	BloodType  BloodType          `bson:"blood_type" json:"blood_type"`
	Allergy    string             `bson:"allergy" json:"allergy"`
	Medication string             `bson:"medication" json:"medication"`
	Height     float64            `bson:"height,omitempty" json:"height,omitempty"`
	HeightUnit HeightUnit         `bson:"height_unit" json:"height_unit"`
	Weight     float64            `bson:"weight,omitempty" json:"weight,omitempty"`
	WeightUnit WeightUnit         `bson:"weight_unit" json:"weight_unit"`
	BirthDate  string             `bson:"birth_date" json:"birth_date"`
	Notes      string             `bson:"notes" json:"notes"`
	CreatedAt  time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"` // timestamp
	UpdatedAt  time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"` // timestamp
}

// BloodType 유효성 검사
func IsValidBloodType(bt string) bool {
	switch BloodType(bt) {
	case BloodTypeAPlus, BloodTypeAMinus,
		BloodTypeBPlus, BloodTypeBMinus,
		BloodTypeOPlus, BloodTypeOMinus,
		BloodTypeABPlus, BloodTypeABMinus,
		BloodTypeOthers:
		return true
	default:
		return false
	}
}

// HeightUnit 유효성 검사
func IsValidHeightUnit(unit string) bool {
	switch HeightUnit(unit) {
	case HeightUnitCm, HeightUnitInch, HeightUnitFeet:
		return true
	default:
		return false
	}
}

// WeightUnit 유효성 검사
func IsValidWeightUnit(unit string) bool {
	switch WeightUnit(unit) {
	case WeightUnitKg, WeightUnitPound:
		return true
	default:
		return false
	}
}
