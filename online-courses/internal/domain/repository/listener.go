package repository

import (
	"context"
	"online-courses/internal/database"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/entity"

	"github.com/google/uuid"
)

type ListenerRepository interface {
	CreateInTx(ctx context.Context, tx database.Tx, m *entity.Listener, i dto.ListenerIDDTO) error
	ReadListener(ctx context.Context, page int, filter string) ([]entity.Listener, error)
	ReadFullData(ctx context.Context, id uuid.UUID) (*entity.Listener, error)
	FindIdById(ctx context.Context, id uuid.UUID) (*entity.Listener, error)
	DeleteInTx(ctx context.Context, x database.Tx, id uuid.UUID) error
	UpdateInTx(ctx context.Context, tx database.Tx, listener entity.Listener) error
	UpdateContractor(ctx context.Context, tx database.Tx, idListener uuid.UUID, idContractor *uuid.UUID) error
	FindByLegalEntity(ctx context.Context, id uuid.UUID) ([]entity.ListenerLegalEntity, error)
}
