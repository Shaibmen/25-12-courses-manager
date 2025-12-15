package service

import (
	"context"
	"online-courses/internal/domain/dto"

	"github.com/google/uuid"
)

type DocumentService interface {
	FormingDataDocument(ctx context.Context, idListener, idProgram, idExecutor uuid.UUID, frontData dto.FrontDataDeliver) (*dto.FullDocumentInfoDTO, error)
}
