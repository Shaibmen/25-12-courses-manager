package repository

import (
	"context"
	"online-courses/internal/domain/entity"
	"time"
)

type ReportRepository interface {
	ReportPeriod(ctx context.Context, startDate, endDate time.Time) ([]entity.ReportPeriod, error)
	MostExpensiveProgram(ctx context.Context, startDate, endDate time.Time) ([]entity.MostExpensiveReport, error)
}
