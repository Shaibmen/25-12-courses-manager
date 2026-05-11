package service

import (
	"context"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/repository"
)

type GraphicsService struct {
	service repository.GraphicsRepository
}

func NewGraphicsService(service repository.GraphicsRepository) *GraphicsService {
	return &GraphicsService{service: service}
}

func (g *GraphicsService) CountListenersOnProgram(ctx context.Context) ([]dto.CountListenersOnProgramDTO, error) {

	entity, err := g.service.CountListenersOnProgram(ctx)
	if err != nil {
		return nil, err
	}

	var counts []dto.CountListenersOnProgramDTO

	for _, i := range entity {
		counts = append(counts, dto.CountListenersOnProgramDTO{
			NameProfEducation: i.NameProfEducation,
			Listeners:         i.Listeners,
		})

	}

	return counts, nil
}

func (g *GraphicsService) PopularProgramType(ctx context.Context) ([]dto.PopularProgramTypeDTO, error) {

	entity, err := g.service.PopularProgramType(ctx)
	if err != nil {
		return nil, err
	}

	var counts []dto.PopularProgramTypeDTO

	for _, i := range entity {
		counts = append(counts, dto.PopularProgramTypeDTO{
			NameProfEducation: i.NameProfEducation,
			EducationType:     i.EducationType,
			Listeners:         i.Listeners,
		})

	}

	return counts, nil
}

func (g *GraphicsService) CountListenersOnProgramAccurate(ctx context.Context) ([]dto.CountListenersOnProgramDTO, error) {

	entity, err := g.service.CountListenersOnProgramAccurate(ctx)
	if err != nil {
		return nil, err
	}

	var counts []dto.CountListenersOnProgramDTO

	for _, i := range entity {
		counts = append(counts, dto.CountListenersOnProgramDTO{
			NameProfEducation: i.NameProfEducation,
			Listeners:         i.Listeners,
		})

	}

	return counts, nil
}

func (g *GraphicsService) WorthProgramAccurate(ctx context.Context) ([]dto.WorthProgramDTO, error) {

	entity, err := g.service.WorthProgramAccurate(ctx)
	if err != nil {
		return nil, err
	}

	var counts []dto.WorthProgramDTO

	for _, i := range entity {
		counts = append(counts, dto.WorthProgramDTO{
			NameProfEducation:    i.NameProfEducation,
			EducationType:        i.EducationType,
			TotalExpectedRevenue: i.Totalrevenue,
		})

	}

	return counts, nil
}

func (g *GraphicsService) AgeDiff(ctx context.Context) ([]dto.AgeDiffDTO, error) {

	entity, err := g.service.AgeDiff(ctx)
	if err != nil {
		return nil, err
	}

	var counts []dto.AgeDiffDTO

	for _, i := range entity {
		counts = append(counts, dto.AgeDiffDTO{
			NameProfEducation: i.NameProfEducation,
			AgeRange:          i.AgeRange,
			Listeners:         i.Listeners,
		})

	}

	return counts, nil
}

func (g *GraphicsService) WhoEnrolled(ctx context.Context) ([]dto.WhoEnrolledDTO, error) {

	entity, err := g.service.WhoEnrolled(ctx)
	if err != nil {
		return nil, err
	}

	var counts []dto.WhoEnrolledDTO

	for _, i := range entity {
		counts = append(counts, dto.WhoEnrolledDTO{
			Month:  i.Month,
			Source: i.Source,
			Cnt:    i.Cnt,
		})

	}

	return counts, nil
}

func (g *GraphicsService) GroupMembers(ctx context.Context) ([]dto.GroupMembersDTO, error) {

	entity, err := g.service.GroupMembers(ctx)
	if err != nil {
		return nil, err
	}

	var counts []dto.GroupMembersDTO

	for _, i := range entity {
		counts = append(counts, dto.GroupMembersDTO{
			NameGroup:      i.NameGroup,
			ActiveEnrolled: i.ActiveEnrolled,
		})

	}

	return counts, nil
}

func (g *GraphicsService) DivisionMember(ctx context.Context) ([]dto.DivisionMemberDTO, error) {

	entity, err := g.service.DivisionMember(ctx)
	if err != nil {
		return nil, err
	}

	var counts []dto.DivisionMemberDTO

	for _, i := range entity {
		counts = append(counts, dto.DivisionMemberDTO{
			Divisioneducation: i.Divisioneducation,
			Listeners:         i.Listeners,
		})

	}

	return counts, nil
}
