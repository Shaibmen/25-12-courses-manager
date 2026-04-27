package service

import (
	"bytes"
	"context"
	"encoding/json"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/entity"
	"online-courses/internal/domain/repository"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
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

func (g *GroupService) ExportExcel(ctx context.Context, id uuid.UUID) ([]byte, error) {

	data, err := g.service.ReadByID(ctx, id)
	if err != nil {
		return nil, err
	}

	file := excelize.NewFile()

	idx, err := file.NewSheet("Расписание")
	if err != nil {
		return nil, err
	}

	file.SetActiveSheet(idx)

	if err := file.DeleteSheet("Sheet1"); err != nil {
		return nil, err
	}

	err = file.SetCellValue("Расписание", "A1", "Группа")
	if err != nil {
		return nil, err
	}

	err = file.SetCellValue("Расписание", "B1", data.NameGroup)
	if err != nil {
		return nil, err
	}

	err = file.SetCellValue("Расписание", "B3", "Тема")
	if err != nil {
		return nil, err
	}
	err = file.SetCellValue("Расписание", "C3", "Дата")
	if err != nil {
		return nil, err
	}

	var raspisanie []map[string]interface{}
	if err := json.Unmarshal(data.Raspisanie, &raspisanie); err != nil {
		return nil, err
	}

	counPos := 4
	countNum := 1

	for _, item := range raspisanie {

		num := strconv.Itoa(countNum) + "."

		parseDate, err := time.Parse("2006-01-02", item["date"].(string))
		if err != nil {
			return nil, err
		}
		formattedDate := parseDate.Format("02.01.2006")

		if err := file.SetCellValue("Расписание", "B"+strconv.Itoa(counPos), item["theme"]); err != nil {
			return nil, err
		}

		if err := file.SetCellValue("Расписание", "C"+strconv.Itoa(counPos), formattedDate); err != nil {
			return nil, err
		}
		if err := file.SetCellValue("Расписание", "A"+strconv.Itoa(counPos), num); err != nil {
			return nil, err
		}
		counPos++
		countNum++
	}

	buf := new(bytes.Buffer)
	if err := file.Write(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
