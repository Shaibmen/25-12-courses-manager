package repository

import (
	"context"
	"online-courses/internal/domain/entity"

	"github.com/google/uuid"
)

type ProgramEducationRepository interface {
	Create(ctx context.Context, model entity.ProgramEducation) error
	Read(ctx context.Context, page int, filter string) ([]entity.ProgramEducation, error)
	ReadByID(ctx context.Context, id uuid.UUID) (*entity.ProgramEducation, error)
	Update(ctx context.Context, model entity.ProgramEducation) error
	Delete(ctx context.Context, id uuid.UUID) error
}
