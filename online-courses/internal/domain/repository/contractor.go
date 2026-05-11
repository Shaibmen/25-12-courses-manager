package repository

import (
	"context"
	"online-courses/internal/database"
	"online-courses/internal/domain/entity"

	"github.com/google/uuid"
)

type ContractorRepository interface {
	CreateInTx(ctx context.Context, tx database.Tx, m *entity.Contractor, idPassport, idRegAddress uuid.UUID) error
	UpdateInTx(ctx context.Context, tx database.Tx, contractor entity.Contractor) error
	DeleteInTx(ctx context.Context, tx database.Tx, id uuid.UUID) error
	FindById(ctx context.Context, id uuid.UUID) (*uuid.UUID, *uuid.UUID, error)
}
