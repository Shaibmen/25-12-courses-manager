package service

import (
	"context"
	"online-courses/internal/domain/dto"

	"github.com/google/uuid"
)

type ListenerService interface {
	CreateFullListener(ctx context.Context, models *dto.CreateListenerDTO) error
	ReadListener(ctx context.Context, page int, filter string) ([]dto.ListenerDTOWithID, error)
	ReadFullListener(ctx context.Context, id uuid.UUID) (*dto.FullListenerDataDTO, error)
	UpdateListener(ctx context.Context, models *dto.CreateListenerDTO, id uuid.UUID) error
	DeleteListener(ctx context.Context, id uuid.UUID) error
}
