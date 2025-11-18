package service

import (
	"context"
	"online-courses/internal/domain/dto"

	"github.com/google/uuid"
)

type ProgramEducationService interface {
	CreateProgram(ctx context.Context, dto dto.ProgramEducationDTO) error
	ReadProgram(ctx context.Context, page int, filter string) ([]dto.ProgramEducationDTO, error)
	ReadByID(ctx context.Context, id uuid.UUID) (*dto.ProgramEducationDTO, error)
	UpdateProgram(ctx context.Context, dto dto.ProgramEducationDTO) error
	DeleteProgram(ctx context.Context, id uuid.UUID) error
}
