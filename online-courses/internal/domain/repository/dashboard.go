package repository

import (
	"context"
	"online-courses/internal/domain/entity"
)

type DashboardRepository interface {
	UserDashboard(ctx context.Context) (*entity.UserDashBoard, error)
	AdminDashboard(ctx context.Context) (*entity.AdminDashboard, error)
}
