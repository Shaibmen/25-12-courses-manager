package repository

import (
	"context"
	"online-courses/internal/database"
	"online-courses/internal/domain/entity"

	"github.com/google/uuid"
)

type RegistrationAddressRepository interface {
	CreateInTx(ctx context.Context, tx database.Tx, m entity.RegistrationAddress) error
	DeleteInTx(ctx context.Context, tx database.Tx, id uuid.UUID) error
	UpdateInTx(ctx context.Context, tx database.Tx, address entity.RegistrationAddress) error
}
