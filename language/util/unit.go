// language/util/unit.go

package util

import (
	"math"
	
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
)


const (
	HeightUnitCm   = model.HeightUnitCm
	HeightUnitFeet = model.HeightUnitFeet
	HeightUnitInch = model.HeightUnitInch

	WeightUnitKg    = model.WeightUnitKg
	WeightUnitPound = model.WeightUnitPound
)


func DefaultHeightUnit(countryCode string) model.HeightUnit {
	switch countryCode {
	case "US", "GB":
		return HeightUnitFeet
	default:
		return HeightUnitCm
	}
}

func DefaultWeightUnit(countryCode string) model.WeightUnit {
	switch countryCode {
	case "US", "GB":
		return WeightUnitPound
	default:
		return WeightUnitKg
	}
}


func ConvertFeetToCm(feet float64) float64 {
	return feet * 30.48
}

func ConvertLbToKg(lb float64) float64 {
	return lb * 0.453592
}


func ConvertCmToFeet(cm float64) float64 {
	return cm / 30.48
}

func ConvertKgToLb(kg float64) float64 {
	return kg / 0.453592
}



func roundToTwoDecimal(val float64) float64 {
	return math.Round(val*100) / 100
}


func ConvertHeight(medicalHeight float64, medicalHeightUnit model.HeightUnit, targetHeightUnit model.HeightUnit) (float64, model.HeightUnit) {
	if medicalHeightUnit == HeightUnitCm && targetHeightUnit != HeightUnitCm {
		if targetHeightUnit == HeightUnitFeet {
			medicalHeight = ConvertCmToFeet(medicalHeight)
		}
	} else if medicalHeightUnit == HeightUnitFeet && targetHeightUnit == HeightUnitCm {
		medicalHeight = ConvertFeetToCm(medicalHeight)
	}

	return roundToTwoDecimal(medicalHeight), targetHeightUnit
}

func ConvertWeight(medicalWeight float64, medicalWeightUnit model.WeightUnit, targetWeightUnit model.WeightUnit) (float64, model.WeightUnit) {
	if medicalWeightUnit == WeightUnitKg && targetWeightUnit != WeightUnitKg {
		if targetWeightUnit == WeightUnitPound {
			medicalWeight = ConvertKgToLb(medicalWeight)
		}
	} else if medicalWeightUnit == WeightUnitPound && targetWeightUnit == WeightUnitKg {
		medicalWeight = ConvertLbToKg(medicalWeight)
	}

	return roundToTwoDecimal(medicalWeight), targetWeightUnit
}