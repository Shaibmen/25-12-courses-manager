package repository

import (
	"context"
	"online-courses/internal/domain/entity"

	"github.com/google/uuid"
)

type GroupRepository interface {
	ReadByID(ctx context.Context, id uuid.UUID) (*entity.Group, error)
	Create(ctx context.Context, m entity.Group) error
	Update(ctx context.Context, m entity.Group) error
	Delete(ctx context.Context, id uuid.UUID) error
	Read(ctx context.Context, page int, filter string) ([]entity.Group, error)
}
