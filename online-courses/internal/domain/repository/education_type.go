package repository

import (
	"context"
	"online-courses/internal/domain/entity"

	"github.com/google/uuid"
)

type EducationTypeRepository interface {
	Create(ctx context.Context, model entity.EducationTypes) error
	Read(ctx context.Context, filter string) ([]entity.EducationTypes, error)
	// ReadByID(ctx context.Context, id uuid.UUID) (*entity.EducationTypes, error)
	Update(ctx context.Context, model entity.EducationTypes) error
	Delete(ctx context.Context, id uuid.UUID) error
}
