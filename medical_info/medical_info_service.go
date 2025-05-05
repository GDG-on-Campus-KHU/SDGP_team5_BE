// medical_info/medical_info_service.go

package medical_info

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/medical_info/repository"
)

type MedicalInfoService struct {
	repo repository.MedicalInfoRepository
}

// constructor function
func NewMedicalInfoService(repo repository.MedicalInfoRepository) *MedicalInfoService {
	if repo == nil {
		log.Fatal("UserRepository cannot be nil")
		return nil
	}

	return &MedicalInfoService{repo: repo}
}

// CreateMedicalInfo creates a new medical info
func (s *MedicalInfoService) CreateMedicalInfo(ctx context.Context, userID int, req MedicalInfoRequest) (*model.MedicalInfo, error) {
	if !model.IsValidBloodType(req.BloodType) {
		return nil, fmt.Errorf("invalid blood type: %s", req.BloodType)
	}
	if !model.IsValidHeightUnit(req.HeightUnit) {
		return nil, fmt.Errorf("invalid height unit: %s", req.HeightUnit)
	}
	if !model.IsValidWeightUnit(req.WeightUnit) {
		return nil, fmt.Errorf("invalid weight unit: %s", req.WeightUnit)
	}

	now := time.Now()

	info := &model.MedicalInfo{
		InfoID:     userID,
		BloodType:  model.BloodType(req.BloodType),
		Allergy:    req.Allergy,
		Medication: req.Medication,
		Height:     req.Height,
		HeightUnit: model.HeightUnit(req.HeightUnit),
		Weight:     req.Weight,
		WeightUnit: model.WeightUnit(req.WeightUnit),
		BirthDate:  req.BirthDate,
		Notes:      req.Notes,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	err := s.repo.CreateMedicalInfo(ctx, info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

// GetMedicalInfo retrieves medical info by user ID
func (s *MedicalInfoService) GetMedicalInfo(ctx context.Context, userID int) (*model.MedicalInfo, error) {
	info, err := s.repo.GetMedicalInfo(ctx, userID)
	if err != nil {
		return nil, err
	}

	return info, nil
}

// UpdateMedicalInfo updates medical info by user ID
func (s *MedicalInfoService) UpdateMedicalInfo(ctx context.Context, userID int, req MedicalInfoRequest) (*model.MedicalInfo, error) {
	if !model.IsValidBloodType(req.BloodType) {
		return nil, fmt.Errorf("invalid blood type: %s", req.BloodType)
	}
	if !model.IsValidHeightUnit(req.HeightUnit) {
		return nil, fmt.Errorf("invalid height unit: %s", req.HeightUnit)
	}
	if !model.IsValidWeightUnit(req.WeightUnit) {
		return nil, fmt.Errorf("invalid weight unit: %s", req.WeightUnit)
	}
	now := time.Now()

	info := &model.MedicalInfo{
		InfoID:     userID,
		BloodType:  model.BloodType(req.BloodType),
		Allergy:    req.Allergy,
		Medication: req.Medication,
		Height:     req.Height,
		HeightUnit: model.HeightUnit(req.HeightUnit),
		Weight:     req.Weight,
		WeightUnit: model.WeightUnit(req.WeightUnit),
		BirthDate:  req.BirthDate,
		Notes:      req.Notes,
		UpdatedAt:  now,
	}

	result, err := s.repo.UpdateMedicalInfo(ctx, userID, info)
	if err != nil {
		return nil, err
	}

	return result, nil
}
