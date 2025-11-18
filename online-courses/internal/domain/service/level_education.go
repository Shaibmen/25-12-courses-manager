package service

import (
	"context"
	"online-courses/internal/domain/dto"
)

type LevelEducationService interface {
	ReadAll(ctx context.Context, filter string) ([]dto.LevelEducationDTO, error)
}
