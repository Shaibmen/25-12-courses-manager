package service

import (
	"context"
	"online-courses/internal/domain/dto"

	"github.com/google/uuid"
)

type LegalEntityService interface {
	Create(ctx context.Context, dto dto.LegalEntityCreateDTO) error
	Read(ctx context.Context, page int, filter string) ([]dto.LegalEntityDTO, error)
	Update(ctx context.Context, dto *dto.LegalEntityCreateDTO, id uuid.UUID) error
	Delete(ctx context.Context, legalEntityID uuid.UUID) error
}
