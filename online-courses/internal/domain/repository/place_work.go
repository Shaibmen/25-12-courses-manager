package repository

import (
	"context"
	"online-courses/internal/database"
	"online-courses/internal/domain/entity"

	"github.com/google/uuid"
)

type PlaceWorkRepository interface {
	CreateInTx(ctx context.Context, tx database.Tx, m entity.PlaceWork) error
	DeleteInTx(ctx context.Context, tx database.Tx, id uuid.UUID) error
	Read(ctx context.Context, id uuid.UUID) (*entity.PlaceWork, error)
	UpdateInTx(ctx context.Context, tx database.Tx, placework entity.PlaceWork) error
}
