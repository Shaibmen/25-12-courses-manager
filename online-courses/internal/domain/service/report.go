package service

import (
	"context"
	"time"
)

type ReportService interface {
	ReportPeriod(ctx context.Context, startDate, endDate time.Time) (string, error)
	MostExpensiveProgram(ctx context.Context, startDate, endDate time.Time) (string, error)
}
