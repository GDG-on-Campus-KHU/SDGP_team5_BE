// medical_info/repository/medical_info_repository.go

package repository

import (
	"context"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
)

type MedicalInfoRepository interface {
	CreateMedicalInfo(ctx context.Context, info *model.MedicalInfo) error
	GetMedicalInfo(ctx context.Context, userID int) (*model.MedicalInfo, error)
	UpdateMedicalInfo(ctx context.Context, userID int, info *model.MedicalInfo) (*model.MedicalInfo, error)
}
