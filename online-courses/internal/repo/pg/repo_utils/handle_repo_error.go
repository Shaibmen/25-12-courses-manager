package repoutils

import (
	"database/sql"
	"errors"
	"log/slog"
	"online-courses/internal/apperrors"

	"github.com/lib/pq"
)

var logger *slog.Logger

func InitRepoLogger(l *slog.Logger) {
	logger = l
}

func HandleRepoErr(err error) error {

	if errors.Is(err, sql.ErrNoRows) {
		return apperrors.ErrEmptyData
	}

	var pgError *pq.Error
	if errors.As(err, &pgError) {
		switch pgError.Code {
		case "23505":
			return apperrors.ErrUnique
		case "23503":
			return apperrors.ErrForeignKey
		case "23502":
			return apperrors.ErrNotNull
		case "23514":
			return apperrors.ErrCheckViolation
		case "42P01":
			return apperrors.ErrNoTable
		case "42703":
			return apperrors.ErrNoRow

		}
		// добавлять по необходимости
		logger.Debug(
			"DEBUG DB ERROR",
			"code", pgError.Code,
			"detail", pgError.Detail,
			"message", pgError.Message,
			"table", pgError.Table,
		)
	}

	return err
}
