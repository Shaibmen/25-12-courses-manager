package service

import (
	"context"
	"online-courses/internal/domain/dto"

	"github.com/google/uuid"
)

type LegalEntityService interface {
	Create(ctx context.Context, dto dto.LegalEntityFullDTO) error
	Read(ctx context.Context, page int, filter string) ([]dto.LegalEntityDTO, error)
	Update(ctx context.Context, dto *dto.LegalEntityFullDTO, id uuid.UUID) error
	Delete(ctx context.Context, legalEntityID uuid.UUID) error
	ReadFullData(ctx context.Context, id uuid.UUID) (*dto.LegalEntityWithListenersDTO, error)
}
