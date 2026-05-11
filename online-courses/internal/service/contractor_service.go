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

type ContractorService struct {
	db             database.DB
	service        repository.ContractorRepository
	passportRepo   repository.PassportRepository
	regAddressRepo repository.RegistrationAddressRepository
	listenerRepo   repository.ListenerRepository
}

func NewContractorService(db database.DB, service repository.ContractorRepository, listenerRepo repository.ListenerRepository, passportRepo repository.PassportRepository, regAddress repository.RegistrationAddressRepository) *ContractorService {
	return &ContractorService{db: db, service: service, listenerRepo: listenerRepo, passportRepo: passportRepo, regAddressRepo: regAddress}
}

func (c *ContractorService) Create(ctx context.Context, dto *dto.ContractorCreateDTO, idListener uuid.UUID) error {

	tx, err := c.db.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	passportID := uuid.New()
	passport, err := mapping.MapPassportToEntity(dto.Passport, passportID)
	if err != nil {
		return err
	}

	if err = c.passportRepo.CreateInTx(ctx, tx, passport); err != nil {
		return err
	}

	regAddressID := uuid.New()
	regAddress := mapping.MapRegistrationAddressToEntity(dto.RegAddress, regAddressID)

	if err = c.regAddressRepo.CreateInTx(ctx, tx, regAddress); err != nil {
		return err
	}

	idContractor := uuid.New()
	entityContractor := &entity.Contractor{
		ID_Contractor: idContractor,
		FirstName:     dto.Contractor.FirstName,
		SecondName:    dto.Contractor.SecondName,
		MiddleName:    dto.Contractor.MiddleName,
		Contact_phone: dto.Contractor.Contact_phone,
		Email:         dto.Contractor.Email,
	}

	if err = c.service.CreateInTx(ctx, tx, entityContractor, passport.ID_Passport, regAddress.ID_RegAddress); err != nil {
		return err
	}

	if err = c.listenerRepo.UpdateContractor(ctx, tx, idListener, &idContractor); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (c *ContractorService) UpdateInTx(ctx context.Context, dto *dto.ContractorCreateDTO, id uuid.UUID) error {

	tx, err := c.db.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	passportID, regAddressID, err := c.service.FindById(ctx, id)
	if err != nil {
		return err
	}

	passport, err := mapping.MapPassportToEntity(dto.Passport, *passportID)
	if err != nil {
		return err
	}
	if err = c.passportRepo.UpdateInTx(ctx, tx, passport); err != nil {
		return err
	}

	regAddress := mapping.MapRegistrationAddressToEntity(dto.RegAddress, *regAddressID)
	if err = c.regAddressRepo.UpdateInTx(ctx, tx, regAddress); err != nil {
		return err
	}

	contractorEntity := entity.Contractor{
		ID_Contractor: id,
		FirstName:     dto.Contractor.FirstName,
		SecondName:    dto.Contractor.SecondName,
		MiddleName:    dto.Contractor.MiddleName,
		Contact_phone: dto.Contractor.Contact_phone,
		Email:         dto.Contractor.Email,
	}

	if err = c.service.UpdateInTx(ctx, tx, contractorEntity); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (c *ContractorService) Delete(ctx context.Context, contractID uuid.UUID) error {

	tx, err := c.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	passportID, regAddressID, err := c.service.FindById(ctx, contractID)
	if err != nil {
		return err
	}

	if err = c.service.DeleteInTx(ctx, tx, contractID); err != nil {
		return err
	}

	if err := c.passportRepo.DeleteInTx(ctx, tx, *passportID); err != nil {
		return err
	}

	if err = c.regAddressRepo.DeleteInTx(ctx, tx, *regAddressID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
