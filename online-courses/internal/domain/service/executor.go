package service

import (
	"context"
	"online-courses/internal/domain/dto"

	"github.com/google/uuid"
)

type ExecutorService interface {
	Create(ctx context.Context, dto dto.ExecutorDTO) error
	Read(ctx context.Context, filter string) ([]dto.ExecutorDTO, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
