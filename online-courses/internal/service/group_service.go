package service

import (
	"context"
	"encoding/json"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/entity"
	"online-courses/internal/domain/repository"

	"github.com/google/uuid"
)

type GroupService struct {
	service repository.GroupRepository
}

func NewGroupService(service repository.GroupRepository) *GroupService {
	return &GroupService{service: service}
}

func (g *GroupService) Create(ctx context.Context, m dto.GroupDTO) error {

	raspisanie, err := json.Marshal(m.Raspisanie)
	if err != nil {
		return err
	}
	id := uuid.New()
	entity := entity.Group{
		ID_Group:   id,
		NameGroup:  m.NameGroup,
		Raspisanie: raspisanie,
	}

	if err := g.service.Create(ctx, entity); err != nil {
		return err
	}

	return nil
}

func (g *GroupService) Update(ctx context.Context, m dto.GroupDTO) error {

	raspisanie, err := json.Marshal(m.Raspisanie)
	if err != nil {
		return err
	}

	entity := entity.Group{
		ID_Group:   m.ID_Group,
		NameGroup:  m.NameGroup,
		Raspisanie: raspisanie,
	}

	if err := g.service.Update(ctx, entity); err != nil {
		return err
	}

	return nil
}

func (g *GroupService) Delete(ctx context.Context, id uuid.UUID) error {

	if err := g.service.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (g *GroupService) Read(ctx context.Context, page int, filter string) ([]dto.GroupDTO, error) {

	data, err := g.service.Read(ctx, page, filter)
	if err != nil {
		return nil, err
	}

	var groups []dto.GroupDTO
	for _, group := range data {

		var raspisanie []map[string]interface{}
		if err := json.Unmarshal(group.Raspisanie, &raspisanie); err != nil {
			return nil, err
		}
		groups = append(groups, dto.GroupDTO{
			ID_Group:   group.ID_Group,
			NameGroup:  group.NameGroup,
			Raspisanie: raspisanie,
		})
	}

	return groups, nil
}

func (g *GroupService) ReadByID(ctx context.Context, id uuid.UUID) (*dto.GroupDTO, error) {

	data, err := g.service.ReadByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var raspisanie []map[string]interface{}
	if err := json.Unmarshal(data.Raspisanie, &raspisanie); err != nil {
		return nil, err
	}

	dto := dto.GroupDTO{
		ID_Group:   data.ID_Group,
		NameGroup:  data.NameGroup,
		Raspisanie: raspisanie,
	}

	return &dto, nil
}
