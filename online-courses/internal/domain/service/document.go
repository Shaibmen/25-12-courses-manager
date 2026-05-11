package service

import (
	"context"
	"online-courses/internal/domain/dto"
	"online-courses/internal/server/http/request"

	"github.com/google/uuid"
)

type DocumentService interface {
	FormingDataDocument(ctx context.Context, idListener, idProgram, idExecutor uuid.UUID, frontData request.FrontDataDeliverRequest) (*dto.FullDocumentInfoDTO, error)
}
