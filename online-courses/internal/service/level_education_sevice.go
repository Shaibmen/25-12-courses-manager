package service

import (
	"context"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/repository"
)

type levelEducationService struct {
	service repository.LevelEducationRepository
}

func NewLevelEducationService(repo repository.LevelEducationRepository) *levelEducationService {
	return &levelEducationService{service: repo}
}

func (lvl *levelEducationService) ReadAll(ctx context.Context, filter string) ([]dto.LevelEducationDTO, error) {

	data, err := lvl.service.Read(ctx, filter)
	if err != nil {
		return nil, err
	}

	var levels []dto.LevelEducationDTO

	for _, l := range data {
		levels = append(levels, dto.LevelEducationDTO{
			ID_LevelEducation: l.ID_LevelEducation,
			Education:         l.Education,
		})
	}

	return levels, nil
}
