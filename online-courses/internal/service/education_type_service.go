package service

import (
	"context"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/entity"
	"online-courses/internal/domain/repository"

	"github.com/google/uuid"
)

type educationTypeService struct {
	educationType repository.EducationTypeRepository
}

func NewEducationTypeService(educationTypes repository.EducationTypeRepository) *educationTypeService {
	return &educationTypeService{educationType: educationTypes}
}

func (e *educationTypeService) CreateType(ctx context.Context, models dto.EducationTypeDTO) error {

	id := uuid.New()

	entity := entity.EducationTypes{
		ID_EducationType: id,
		TypeName:         models.TypeName,
	}

	if err := e.educationType.Create(ctx, entity); err != nil {
		return err
	}

	return nil
}

func (e *educationTypeService) ReadEducationType(ctx context.Context, filter string) ([]dto.EducationTypeDTO, error) {

	data, err := e.educationType.Read(ctx, filter)
	if err != nil {
		return nil, err
	}

	var dtolist []dto.EducationTypeDTO

	for _, l := range data {
		dtolist = append(dtolist, dto.EducationTypeDTO{
			ID_EducationType: l.ID_EducationType,
			TypeName:         l.TypeName,
		})
	}

	return dtolist, nil
}

// func (e *educationTypeService) ReadEducationTypeByID(ctx context.Context, id uuid.UUID) (*dto.EducationTypeDTO, error) {

// 	data, err := e.educationType.ReadByID(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	dto := dto.EducationTypeDTO{
// 		ID_EducationType: data.ID_EducationType,
// 		TypeName:         data.TypeName,
// 	}

// 	return &dto, nil
// }

func (e *educationTypeService) UpdateEducationType(ctx context.Context, id uuid.UUID, models dto.EducationTypeDTO) error {

	entity := entity.EducationTypes{
		ID_EducationType: id,
		TypeName:         models.TypeName,
	}

	if err := e.educationType.Update(ctx, entity); err != nil {
		return err
	}

	return nil
}

func (e *educationTypeService) DeleteEducationType(ctx context.Context, id uuid.UUID) error {

	if err := e.educationType.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
