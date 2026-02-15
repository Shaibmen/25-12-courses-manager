package service

import (
	"context"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/entity"
	"online-courses/internal/domain/repository"

	"github.com/google/uuid"
)

type ExecutorService struct {
	service repository.ExecutorRepository
}

func NewExecutorService(service repository.ExecutorRepository) *ExecutorService {
	return &ExecutorService{service: service}
}

func (e *ExecutorService) Create(ctx context.Context, dto dto.ExecutorDTO) error {

	id := uuid.New()
	entity := &entity.Executor{
		ID_Executor: id,
		Status:      dto.Status,
		FirstName:   dto.FirstName,
		SecondName:  dto.SecondName,
		MiddleName:  dto.MiddleName,
		Doverenost:  dto.Doverenost,
	}

	if err := e.service.Create(ctx, *entity); err != nil {
		return err
	}

	return nil
}

func (e *ExecutorService) Read(ctx context.Context, filter string) ([]dto.ExecutorDTO, error) {

	data, err := e.service.Read(ctx, filter)
	if err != nil {
		return nil, err
	}

	var executors []dto.ExecutorDTO
	for _, l := range data {
		executors = append(executors, dto.ExecutorDTO{
			ID_Executor: l.ID_Executor,
			Status:      l.Status,
			FirstName:   l.FirstName,
			SecondName:  l.SecondName,
			MiddleName:  l.MiddleName,
			Doverenost:  l.Doverenost,
		})
	}

	return executors, nil
}

func (e *ExecutorService) Delete(ctx context.Context, id uuid.UUID) error {

	if err := e.service.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
