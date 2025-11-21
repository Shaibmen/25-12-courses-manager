package service

import (
	"context"
	"online-courses/internal/domain/dto"

	"github.com/google/uuid"
)

type ContractorService interface {
	Create(ctx context.Context, dto *dto.ContractorCreateDTO) error
	UpdateInTx(ctx context.Context, dto *dto.ContractorCreateDTO, id uuid.UUID) error
	Delete(ctx context.Context, contractID uuid.UUID) error
}
