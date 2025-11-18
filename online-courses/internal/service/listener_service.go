package service

import (
	"context"
	"online-courses/internal/database"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/repository"
	"online-courses/internal/mapping"

	"github.com/google/uuid"
)

type listenerService struct {
	db               database.DB
	passportRepo     repository.PassportRepository
	educationRepo    repository.EducationListenerRepository
	placeWorkRepo    repository.PlaceWorkRepository
	registrationRepo repository.RegistrationAddressRepository
	listenerRepo     repository.ListenerRepository
}

func NewListenerService(db database.DB,
	passport repository.PassportRepository,
	education repository.EducationListenerRepository,
	placework repository.PlaceWorkRepository,
	registration repository.RegistrationAddressRepository,
	listener repository.ListenerRepository) *listenerService {
	return &listenerService{
		db:               db,
		passportRepo:     passport,
		educationRepo:    education,
		placeWorkRepo:    placework,
		registrationRepo: registration,
		listenerRepo:     listener,
	}
}

func (l *listenerService) CreateFullListener(ctx context.Context, models *dto.CreateListenerDTO) error {

	tx, err := l.db.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	passportID := uuid.New()
	passport, err := mapping.MapPassportToEntity(models.Passport, passportID)
	if err != nil {
		return err
	}

	if err = l.passportRepo.CreateInTx(ctx, tx, passport); err != nil {
		return err
	}

	var placeWorkID *uuid.UUID = nil

	if models.PlaceWork != (dto.PlaceWorkDTO{}) {
		newID := uuid.New()
		placeWorkID = &newID

		placeWork := mapping.MapPlaceWorkToEntity(models.PlaceWork, *placeWorkID)

		if err = l.placeWorkRepo.CreateInTx(ctx, tx, placeWork); err != nil {
			return err
		}
	}

	var educationID *uuid.UUID = nil
	if models.EducationListener != (dto.EducationListenerDTO{}) {
		newID := uuid.New()
		educationID = &newID
		education, err := mapping.MapEducationListenerToEntity(models.EducationListener, *educationID)
		if err != nil {
			return err
		}

		if err = l.educationRepo.CreateInTx(ctx, tx, education); err != nil {
			return err
		}
	}

	registrationID := uuid.New()
	registration := mapping.MapRegistrationAddressToEntity(models.RegistrationAddress, registrationID)

	if err = l.registrationRepo.CreateInTx(ctx, tx, registration); err != nil {
		return err
	}

	listenerID := uuid.New()
	listener, err := mapping.MapListenerToEntity(models.Listener, listenerID)
	if err != nil {
		return err
	}

	optionalID := dto.ListenerIDDTO{
		ID_Passport:          passportID,
		ID_RegAddress:        registrationID,
		ID_EducationListener: educationID,
		ID_PlaceWork:         placeWorkID,
	}

	if err = l.listenerRepo.CreateInTx(ctx, tx, listener, optionalID); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (l *listenerService) ReadListener(ctx context.Context, page int, filter string) ([]dto.ListenerDTOWithID, error) {
	data, err := l.listenerRepo.ReadListener(ctx, page, filter)
	if err != nil {
		return nil, err
	}

	var dtoList []dto.ListenerDTOWithID

	for _, l := range data {
		dtoList = append(dtoList, dto.ListenerDTOWithID{
			ID_Listener:  l.ID_Listener,
			FirstName:    l.FirstName,
			SecondName:   l.SecondName,
			MiddleName:   l.MiddleName,
			DateOfBirth:  l.DateOfBirth.String(),
			SNILS:        l.SNILS,
			ContactPhone: l.ContactPhone,
			Email:        l.Email,
		})
	}

	return dtoList, nil
}

func (l *listenerService) ReadFullListener(ctx context.Context, id uuid.UUID) (*dto.FullListenerDataDTO, error) {
	result, err := l.listenerRepo.ReadFullData(ctx, id)
	if err != nil {
		return nil, err
	}

	if result.ID_EducationListener != nil {
		educaiton, err := l.educationRepo.Read(ctx, *result.ID_EducationListener)
		if err != nil {
			return nil, err
		}
		result.EducationListener = *educaiton
	}

	if result.ID_PlaceWork != nil {
		placework, err := l.placeWorkRepo.Read(ctx, *result.ID_PlaceWork)
		if err != nil {
			return nil, err
		}

		result.PlaceWork = *placework
	}

	return mapping.MapListenerEntityToDTO(result), nil
}

func (l *listenerService) UpdateListener(ctx context.Context, models *dto.CreateListenerDTO, id uuid.UUID) error {
	tx, err := l.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ids, err := l.listenerRepo.FindIdById(ctx, id)
	if err != nil {
		return err
	}

	listener, err := mapping.MapListenerToEntity(models.Listener, ids.ID_Listener)
	if err != nil {
		return err
	}

	passport, err := mapping.MapPassportToEntity(models.Passport, ids.ID_Passport)
	if err != nil {
		return err
	}

	regaddress := mapping.MapRegistrationAddressToEntity(models.RegistrationAddress, ids.ID_RegAddress)

	if err := l.listenerRepo.UpdateInTx(ctx, tx, *listener); err != nil {
		return err
	}

	if err := l.passportRepo.UpdateInTx(ctx, tx, passport); err != nil {
		return err
	}

	if err := l.registrationRepo.UpdateInTx(ctx, tx, regaddress); err != nil {
		return err
	}

	if ids.ID_EducationListener != nil {
		education, err := mapping.MapEducationListenerToEntity(models.EducationListener, *ids.ID_EducationListener)
		if err != nil {
			return err
		}

		if err := l.educationRepo.UpdateInTx(ctx, tx, education); err != nil {
			return err
		}
	}

	if ids.ID_PlaceWork != nil {
		placework := mapping.MapPlaceWorkToEntity(models.PlaceWork, *ids.ID_PlaceWork)

		if err := l.placeWorkRepo.UpdateInTx(ctx, tx, placework); err != nil {
			return err
		}

	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (l *listenerService) DeleteListener(ctx context.Context, id uuid.UUID) error {

	tx, err := l.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ids, err := l.listenerRepo.FindIdById(ctx, id)
	if err != nil {
		return err
	}

	if err = l.listenerRepo.DeleteInTx(ctx, tx, ids.ID_Listener); err != nil {
		return err
	}

	if err = l.passportRepo.DeleteInTx(ctx, tx, ids.ID_Passport); err != nil {
		return err
	}

	if err = l.registrationRepo.DeleteInTx(ctx, tx, ids.ID_RegAddress); err != nil {
		return err
	}

	if ids.ID_EducationListener != nil {
		if err = l.educationRepo.DeleteInTx(ctx, tx, *ids.ID_EducationListener); err != nil {
			return err
		}
	}

	if ids.ID_PlaceWork != nil {
		if err = l.placeWorkRepo.DeleteInTx(ctx, tx, *ids.ID_PlaceWork); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil

}
