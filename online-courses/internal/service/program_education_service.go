package service

import (
	"context"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/entity"
	"online-courses/internal/domain/repository"

	"github.com/google/uuid"
)

type programEducationService struct {
	service repository.ProgramEducationRepository
}

func NewProgramEducationService(repo repository.ProgramEducationRepository) *programEducationService {
	return &programEducationService{service: repo}
}

func (p *programEducationService) CreateProgram(ctx context.Context, dto dto.ProgramEducationDTO) error {

	id := uuid.New()

	entity := entity.ProgramEducation{
		ID_ProgramEducation:   id,
		NameProfEducation:     dto.NameProfEducation,
		TimeEducation:         dto.TimeEducation,
		IndividualPrice:       dto.IndividualPrice,
		GroupPrice:            dto.GroupPrice,
		CampusPrice:           dto.CampusPrice,
		ID_EducationType:      dto.ID_EducationType,
		ID_DivisionsEducation: dto.ID_DivisionsEducation,
	}

	if err := p.service.Create(ctx, entity); err != nil {
		return err
	}

	return nil
}

func (p *programEducationService) ReadProgram(ctx context.Context, page int, filter string) ([]dto.ProgramEducationDTO, error) {

	data, err := p.service.Read(ctx, page, filter)
	if err != nil {
		return nil, err
	}

	var programs []dto.ProgramEducationDTO
	for _, prog := range data {
		programs = append(programs, dto.ProgramEducationDTO{
			ID_ProgramEducation:   prog.ID_ProgramEducation,
			NameProfEducation:     prog.NameProfEducation,
			TimeEducation:         prog.TimeEducation,
			IndividualPrice:       prog.IndividualPrice,
			GroupPrice:            prog.GroupPrice,
			CampusPrice:           prog.CampusPrice,
			ID_EducationType:      prog.ID_EducationType,
			ID_DivisionsEducation: prog.ID_DivisionsEducation,
		})
	}

	return programs, err
}

func (p *programEducationService) ReadByID(ctx context.Context, id uuid.UUID) (*dto.ProgramEducationDTO, error) {

	data, err := p.service.ReadByID(ctx, id)
	if err != nil {
		return nil, err
	}

	dto := dto.ProgramEducationDTO{
		ID_ProgramEducation:   data.ID_ProgramEducation,
		NameProfEducation:     data.NameProfEducation,
		TimeEducation:         data.TimeEducation,
		IndividualPrice:       data.IndividualPrice,
		GroupPrice:            data.GroupPrice,
		CampusPrice:           data.CampusPrice,
		ID_EducationType:      data.ID_EducationType,
		ID_DivisionsEducation: data.ID_DivisionsEducation,
	}

	return &dto, nil
}

func (p *programEducationService) UpdateProgram(ctx context.Context, dto dto.ProgramEducationDTO) error {

	entity := entity.ProgramEducation{
		ID_ProgramEducation:   dto.ID_ProgramEducation,
		NameProfEducation:     dto.NameProfEducation,
		TimeEducation:         dto.TimeEducation,
		IndividualPrice:       dto.IndividualPrice,
		GroupPrice:            dto.GroupPrice,
		CampusPrice:           dto.CampusPrice,
		ID_EducationType:      dto.ID_EducationType,
		ID_DivisionsEducation: dto.ID_DivisionsEducation,
	}

	if err := p.service.Update(ctx, entity); err != nil {
		return err
	}

	return nil
}

func (p *programEducationService) DeleteProgram(ctx context.Context, id uuid.UUID) error {

	if err := p.service.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
