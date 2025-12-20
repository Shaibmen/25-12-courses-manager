package repository

import (
	"context"
	"online-courses/internal/domain/entity"

	"github.com/google/uuid"
)

type DocumentRepository interface {
	PrepareDataDocument(ctx context.Context, idListener, idProgram, idExecutor uuid.UUID) (*entity.FullDocument, error)
}
