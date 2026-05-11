package repository

import (
	"context"
	"online-courses/internal/domain/entity"

	"github.com/google/uuid"
)

type ExecutorRepository interface {
	Create(ctx context.Context, m entity.Executor) error
	Read(ctx context.Context, filter string) ([]entity.Executor, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
