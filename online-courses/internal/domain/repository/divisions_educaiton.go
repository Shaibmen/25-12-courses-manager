package repository

import (
	"context"
	"online-courses/internal/domain/entity"

	"github.com/google/uuid"
)

type DivisionsEducationRepository interface {
	Create(ctx context.Context, divisions entity.DivisionsEducation) error
	Read(ctx context.Context, filter string) ([]entity.DivisionsEducation, error)
	// ReadByID(ctx context.Context, id uuid.UUID) (*entity.DivisionsEducation, error)
	Update(ctx context.Context, divisions entity.DivisionsEducation) error
	Delete(ctx context.Context, id uuid.UUID) error
}
