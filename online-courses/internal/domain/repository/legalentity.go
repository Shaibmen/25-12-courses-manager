package repository

import (
	"context"
	"online-courses/internal/database"
	"online-courses/internal/domain/entity"

	"github.com/google/uuid"
)

type LegalEntityRepository interface {
	CreateInTx(ctx context.Context, tx database.Tx, m *entity.LegalEntity) error
	Read(ctx context.Context, page int, filter string) ([]entity.LegalEntity, error)
	UpdateInTx(ctx context.Context, tx database.Tx, m entity.LegalEntity) error
	DeleteInTx(ctx context.Context, tx database.Tx, id uuid.UUID) error
	ReadFullData(ctx context.Context, id uuid.UUID) (*entity.LegalEntity, error)
	FindById(ctx context.Context, id uuid.UUID) (*uuid.UUID, error)
}
