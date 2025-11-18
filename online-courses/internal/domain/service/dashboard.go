package service

import (
	"context"
	"online-courses/internal/domain/dto"
)

type DashboardService interface {
	UserDashboard(ctx context.Context) (*dto.UserDashBoardDTO, error)
	AdminDashboard(ctx context.Context) (*dto.AdminDashboardDTO, error)
}
