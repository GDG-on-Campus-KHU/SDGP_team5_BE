// situation/repository/situation_repository.go

package repository

import (
	"context"

	"github.com/GDG-on-Campus-KHU/SDGP_team5_BE/db/model"
)

type SituationRepository interface {
	GetSituationByIndex(ctx context.Context, situationIndex int) (*model.Situation, error)
}
