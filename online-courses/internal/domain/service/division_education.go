package service

import (
	"context"
	"online-courses/internal/domain/dto"

	"github.com/google/uuid"
)

type DivisionsEducationService interface {
	CreateDivision(ctx context.Context, models dto.DivisionsDTO) error
	ReadDivisions(ctx context.Context, filter string) ([]dto.DivisionsDTO, error)
	// ReadByID(ctx context.Context, id uuid.UUID) (*dto.DivisionsDTO, error)
	UpdateDivision(ctx context.Context, id uuid.UUID, models dto.DivisionsDTO) error
	DeleteDivision(ctx context.Context, id uuid.UUID) error
}
