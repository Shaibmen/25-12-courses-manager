package service

import (
	"context"
	"online-courses/internal/domain/dto"

	"github.com/google/uuid"
)

type GroupService interface {
	Create(ctx context.Context, m dto.GroupDTO) error
	Update(ctx context.Context, m dto.GroupDTO) error
	Delete(ctx context.Context, id uuid.UUID) error
	Read(ctx context.Context, page int, filter string) ([]dto.GroupDTO, error)
	ReadByID(ctx context.Context, id uuid.UUID) (*dto.GroupDTO, error)
	ExportExcel(ctx context.Context, id uuid.UUID) ([]byte, error)
}
