package repository

import (
	"context"
	"online-courses/internal/domain/entity"

	"github.com/google/uuid"
)

type EnrollmentListenerRepository interface {
	Create(ctx context.Context, model entity.EnrollmentListener) error
	Read(ctx context.Context, page int, filter string) ([]entity.EnrollmentListenerDetails, error)
	Update(ctx context.Context, idListener, idProgram uuid.UUID, model entity.EnrollmentListener) error
	Delete(ctx context.Context, id_listener, id_program uuid.UUID) error
	ReadDetailListener(ctx context.Context, id uuid.UUID) ([]entity.EnrollmentProgramDetails, error)
	ReadByProgram(ctx context.Context, id uuid.UUID, page int) ([]entity.EnrollmentListenerDetails, error)
	InfoToPersonalCard(ctx context.Context, listenerID, programID uuid.UUID) (*entity.PersonalCardInfo, error)
	GetListenerFIO(ctx context.Context, listenerID uuid.UUID) (*entity.ListenerFIO, error)
	SaveInfo(ctx context.Context, m entity.AccurateProgram) error
	GetProgram(ctx context.Context, id uuid.UUID) (*entity.ProgramToAccurate, error)
}
