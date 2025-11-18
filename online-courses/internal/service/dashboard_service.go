package service

import (
	"context"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/repository"
)

type DashboardService struct {
	service repository.DashboardRepository
}

func NewDashboardService(repo repository.DashboardRepository) *DashboardService {
	return &DashboardService{service: repo}
}

func (d *DashboardService) UserDashboard(ctx context.Context) (*dto.UserDashBoardDTO, error) {

	data, err := d.service.UserDashboard(ctx)
	if err != nil {
		return nil, err
	}

	var dtoProgram []dto.ProgramEndingSoonDTO

	for _, l := range data.ProgramEndingSoon {
		dtoProgram = append(dtoProgram, dto.ProgramEndingSoonDTO{
			NameProfEducation: l.NameProfEducation,
			EndDate:           l.EndDate,
			TotalListeners:    l.TotalListeners,
		})
	}
	if dtoProgram == nil {
		dtoProgram = []dto.ProgramEndingSoonDTO{}
	}

	dto := dto.UserDashBoardDTO{
		TotalListener:     data.TotalListener,
		TotalProgram:      data.TotalProgram,
		ActiveEnrollments: data.ActiveEnrollments,
		ProgramEndingSoon: dtoProgram,
	}

	return &dto, nil
}

func (d *DashboardService) AdminDashboard(ctx context.Context) (*dto.AdminDashboardDTO, error) {

	data, err := d.service.AdminDashboard(ctx)
	if err != nil {
		return nil, err
	}

	var dtoPgStat []dto.PgStatDTO
	for _, l := range data.PgStat {
		dtoPgStat = append(dtoPgStat, dto.PgStatDTO{
			Pid:             l.Pid,
			ApplicationName: l.ApplicationName,
			ClientAddr:      l.ClientAddr,
			ClientPort:      l.ClientPort,
			State:           l.State,
			QueryStart:      l.QueryStart,
		})
	}

	var role []dto.RoleDTO
	for _, l := range data.Role {
		role = append(role, dto.RoleDTO{
			ID:   l.ID,
			Role: l.Role,
		})
	}

	dto := dto.AdminDashboardDTO{
		PgStat: dtoPgStat,
		AdminStat: dto.AdminStatDTO{
			DatabaseName:     data.AdminStat.DatabaseName,
			Totalsize:        data.AdminStat.Totalsize,
			ActiveConnection: data.AdminStat.ActiveConnection,
			CommittedTx:      data.AdminStat.CommittedTx,
			RolledbackTx:     data.AdminStat.RolledbackTx,
			DiskBlockRead:    data.AdminStat.DiskBlockRead,
			BufferHits:       data.AdminStat.BufferHits,
		},
		Role: role,
	}

	return &dto, nil
}
