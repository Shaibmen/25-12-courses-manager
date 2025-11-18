package repository

import (
	"context"
	"online-courses/internal/domain/entity"
)

type LevelEducationRepository interface {
	Read(ctx context.Context, filter string) ([]entity.LevelEducation, error)
}
