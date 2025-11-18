package service

import (
	"context"
	"online-courses/internal/domain/dto"

	"github.com/google/uuid"
)

type EducationTypeService interface {
	CreateType(ctx context.Context, models dto.EducationTypeDTO) error
	ReadEducationType(ctx context.Context, filter string) ([]dto.EducationTypeDTO, error)
	// ReadEducationTypeByID(ctx context.Context, id uuid.UUID) (*dto.EducationTypeDTO, error)
	UpdateEducationType(ctx context.Context, id uuid.UUID, models dto.EducationTypeDTO) error
	DeleteEducationType(ctx context.Context, id uuid.UUID) error
}
