package repository

import (
	"context"
	"online-courses/internal/database"
	"online-courses/internal/domain/entity"

	"github.com/google/uuid"
)

type EducationListenerRepository interface {
	CreateInTx(ctx context.Context, tx database.Tx, m entity.EducationListener) error
	DeleteInTx(ctx context.Context, tx database.Tx, id uuid.UUID) error
	Read(ctx context.Context, id uuid.UUID) (*entity.EducationListener, error)
	UpdateInTx(ctx context.Context, tx database.Tx, education entity.EducationListener) error
}
