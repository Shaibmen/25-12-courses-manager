package service

import (
	"context"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/entity"
	"online-courses/internal/domain/repository"

	"github.com/google/uuid"
)

type divisionsEducationService struct {
	divisionsRepo repository.DivisionsEducationRepository
}

func NewDivisionsEducationService(divisions repository.DivisionsEducationRepository) *divisionsEducationService {
	return &divisionsEducationService{divisionsRepo: divisions}
}

func (d *divisionsEducationService) CreateDivision(ctx context.Context, models dto.DivisionsDTO) error {

	id := uuid.New()

	entity := entity.DivisionsEducation{
		ID_DivisionsEducation: id,
		Divisions:             models.Divisions,
	}

	if err := d.divisionsRepo.Create(ctx, entity); err != nil {
		return err
	}

	return nil
}

func (d *divisionsEducationService) ReadDivisions(ctx context.Context, filter string) ([]dto.DivisionsDTO, error) {

	data, err := d.divisionsRepo.Read(ctx, filter)
	if err != nil {
		return nil, err
	}

	var dtolist []dto.DivisionsDTO

	for _, l := range data {
		dtolist = append(dtolist, dto.DivisionsDTO{
			ID_DivisionsEducation: l.ID_DivisionsEducation,
			Divisions:             l.Divisions,
		})
	}

	return dtolist, nil
}

// func (d *divisionsEducationService) ReadByID(ctx context.Context, id uuid.UUID) (*dto.DivisionsDTO, error) {

// 	data, err := d.divisionsRepo.ReadByID(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	dto := dto.DivisionsDTO{
// 		ID_DivisionsEducation: data.ID_DivisionsEducation,
// 		Divisions:             data.Divisions,
// 	}

// 	return &dto, nil
// }

func (d *divisionsEducationService) UpdateDivision(ctx context.Context, id uuid.UUID, models dto.DivisionsDTO) error {

	entity := entity.DivisionsEducation{
		ID_DivisionsEducation: id,
		Divisions:             models.Divisions,
	}

	if err := d.divisionsRepo.Update(ctx, entity); err != nil {
		return err
	}

	return nil
}

func (d *divisionsEducationService) DeleteDivision(ctx context.Context, id uuid.UUID) error {

	if err := d.divisionsRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
