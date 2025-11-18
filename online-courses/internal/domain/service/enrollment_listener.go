package service

import (
	"context"
	"online-courses/internal/domain/dto"

	"github.com/google/uuid"
)

type EnrollmentListenerService interface {
	Create(ctx context.Context, model dto.EnrollmentListenerDTO) error
	Read(ctx context.Context, page int, filter string) ([]dto.EnrollmentListenerDetailsDTO, error)
	Update(ctx context.Context, idListener, idProgram uuid.UUID, model dto.EnrollmentListenerDTO) error
	Delete(ctx context.Context, id_listener, id_program uuid.UUID) error
	ReadDetailListener(ctx context.Context, id uuid.UUID) ([]dto.EnrollmentProgramDetailsDTO, error)
	ReadByProgram(ctx context.Context, id uuid.UUID, page int) ([]dto.EnrollmentListenerDetailsDTO, error)
	InfoToPersonalCard(ctx context.Context, listenerID, programID uuid.UUID) (*dto.PersonalCardInfoDTO, error)
	GetListenerFIO(ctx context.Context, listenerID uuid.UUID) (*dto.ListenerFIODTO, error)
}
