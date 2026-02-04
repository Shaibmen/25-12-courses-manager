package service

import (
	"context"
	"online-courses/internal/database"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/entity"
	"online-courses/internal/domain/repository"
	"online-courses/internal/mapping"

	"github.com/google/uuid"
)

type LegalEntityService struct {
	db             database.DB
	service        repository.LegalEntityRepository
	regAddressRepo repository.RegistrationAddressRepository
	listenerRepo   repository.ListenerRepository
}

func NewLegalEntityService(db database.DB, service repository.LegalEntityRepository, regAddress repository.RegistrationAddressRepository, listenerRepo repository.ListenerRepository) *LegalEntityService {
	return &LegalEntityService{db: db, service: service, regAddressRepo: regAddress, listenerRepo: listenerRepo}
}

func (l *LegalEntityService) Create(ctx context.Context, dto dto.LegalEntityFullDTO) error {

	tx, err := l.db.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	regAddressID := uuid.New()
	regAddress := mapping.MapRegistrationAddressToEntity(dto.RegAddress, regAddressID)

	if err = l.regAddressRepo.CreateInTx(ctx, tx, regAddress); err != nil {
		return err
	}

	legalEntityID := uuid.New()
	legalEntity := &entity.LegalEntity{
		ID_Legalentity: legalEntityID,
		NameCompany:    dto.LegalEntity.NameCompany,
		Inn:            dto.LegalEntity.Inn,
		Kpp:            dto.LegalEntity.Kpp,
		Ogrn:           dto.LegalEntity.Ogrn,
		Phone:          dto.LegalEntity.Phone,
		Email:          dto.LegalEntity.Email,
		FirstName:      dto.LegalEntity.FirstName,
		SecondName:     dto.LegalEntity.SecondName,
		MiddleName:     dto.LegalEntity.MiddleName,
		ID_RegAddress:  regAddressID,
		Status:         dto.LegalEntity.Status,
	}

	if err = l.service.CreateInTx(ctx, tx, legalEntity); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (l *LegalEntityService) Read(ctx context.Context, page int, filter string) ([]dto.LegalEntityDTO, error) {

	data, err := l.service.Read(ctx, page, filter)
	if err != nil {
		return nil, err
	}

	var dtoList []dto.LegalEntityDTO

	for _, l := range data {
		dtoList = append(dtoList, dto.LegalEntityDTO{
			ID_Legalentity: l.ID_Legalentity,
			NameCompany:    l.NameCompany,
			Inn:            l.Inn,
			Kpp:            l.Kpp,
			Ogrn:           l.Ogrn,
			Phone:          l.Phone,
			Email:          l.Email,
			FirstName:      l.FirstName,
			SecondName:     l.SecondName,
			MiddleName:     l.MiddleName,
			ID_RegAddress:  l.ID_RegAddress,
			Status:         l.Status,
		})
	}

	return dtoList, nil
}

func (l *LegalEntityService) Update(ctx context.Context, dto *dto.LegalEntityFullDTO, id uuid.UUID) error {
	tx, err := l.db.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	regAddressID, err := l.service.FindById(ctx, id)
	if err != nil {
		return err
	}

	regAddress := mapping.MapRegistrationAddressToEntity(dto.RegAddress, *regAddressID)
	if err = l.regAddressRepo.UpdateInTx(ctx, tx, regAddress); err != nil {

		return err
	}
	legalEntity := entity.LegalEntity{
		ID_Legalentity: id,
		NameCompany:    dto.LegalEntity.NameCompany,
		Inn:            dto.LegalEntity.Inn,
		Kpp:            dto.LegalEntity.Kpp,
		Ogrn:           dto.LegalEntity.Ogrn,
		Phone:          dto.LegalEntity.Phone,
		Email:          dto.LegalEntity.Email,
		FirstName:      dto.LegalEntity.FirstName,
		SecondName:     dto.LegalEntity.SecondName,
		MiddleName:     dto.LegalEntity.MiddleName,
		Status:         dto.LegalEntity.Status,
	}

	if err = l.service.UpdateInTx(ctx, tx, legalEntity); err != nil {

		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (l *LegalEntityService) Delete(ctx context.Context, legalEntityID uuid.UUID) error {
	tx, err := l.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	regAddressID, err := l.service.FindById(ctx, legalEntityID)
	if err != nil {
		return err
	}

	if err = l.service.DeleteInTx(ctx, tx, legalEntityID); err != nil {
		return err
	}

	if err = l.regAddressRepo.DeleteInTx(ctx, tx, *regAddressID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (l *LegalEntityService) ReadFullData(ctx context.Context, id uuid.UUID) (*dto.LegalEntityWithListenersDTO, error) {

	result, err := l.service.ReadFullData(ctx, id)
	if err != nil {
		return nil, err
	}

	listener, err := l.listenerRepo.FindByLegalEntity(ctx, id)
	if err != nil {
		return nil, err
	}

	var dtoListeners []dto.ListenerForLegalEntity

	for _, i := range listener {
		dtoListeners = append(dtoListeners, dto.ListenerForLegalEntity{
			ID_Listener: i.ID_Listener,
			FirstName:   i.FirstName,
			SecondName:  i.SecondName,
			MiddleName:  i.MiddleName,
			SNILS:       i.SNILS,
			DateOfBirth: i.DateOfBirth,
			Email:       i.Email,
		})
	}

	dto := dto.LegalEntityWithListenersDTO{
		LegalEntity: dto.LegalEntityDTO{
			Listeners:      dtoListeners,
			ID_Legalentity: result.ID_Legalentity,
			NameCompany:    result.NameCompany,
			Inn:            result.Inn,
			Kpp:            result.Kpp,
			Ogrn:           result.Ogrn,
			Phone:          result.Phone,
			Email:          result.Email,
			FirstName:      result.FirstName,
			SecondName:     result.SecondName,
			MiddleName:     result.MiddleName,
			ID_RegAddress:  result.ID_RegAddress,
			Status:         result.Status,
		},
		RegAddress: dto.RegistrationAddressDTO{
			MailIndex: result.RegistrationAddress.MailIndex,
			Region:    result.RegistrationAddress.Region,
			City:      result.RegistrationAddress.City,
			Street:    result.RegistrationAddress.Street,
			House:     result.RegistrationAddress.House,
			Building:  result.RegistrationAddress.Building,
			Apartment: result.RegistrationAddress.Apartment,
		},
	}

	return &dto, nil
}
