package mapping

import (
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/entity"
	"online-courses/internal/server/http/request"
	"online-courses/internal/utils"
	"strconv"

	"github.com/google/uuid"
)

func MapPassportToEntity(dto dto.PassportDTO, id uuid.UUID) (entity.Passport, error) {

	dateGiven, err := utils.TimeParse(dto.DateGiven)
	if err != nil {
		return entity.Passport{}, err
	}

	return entity.Passport{
		ID_Passport:   id,
		PlaceBirth:    dto.PlaceBirth,
		Citizenship:   dto.Citizenship,
		Gender:        dto.Gender,
		Seria:         dto.Seria,
		Number:        dto.Number,
		PassportGiven: dto.PassportGiven,
		DateGiven:     *dateGiven,
		Code:          dto.Code,
	}, nil
}

func MapEducationListenerToEntity(dto dto.EducationListenerDTO, id uuid.UUID) (entity.EducationListener, error) {

	dateGiven, err := utils.TimeParse(dto.DateGiven)
	if err != nil {
		return entity.EducationListener{}, err
	}

	return entity.EducationListener{
		ID_EducationListener:   id,
		DiplomSeria:            dto.DiplomSeria,
		DiplomNumber:           dto.DiplomNumber,
		DateGiven:              *dateGiven,
		City:                   dto.City,
		Region:                 dto.Region,
		EducationalInstitution: dto.EducationalInstitution,
		Speciality:             dto.Speciality,
		LevelEducation:         dto.LevelEducation,
	}, nil
}

func MapPlaceWorkToEntity(dto dto.PlaceWorkDTO, id uuid.UUID) entity.PlaceWork {
	return entity.PlaceWork{
		ID_PlaceWork:       id,
		NameCompany:        dto.NameCompany,
		JobTitle:           dto.JobTitle,
		AllExperience:      dto.AllExperience,
		JobTitleExperience: dto.JobTitleExpirience,
	}
}

func MapRegistrationAddressToEntity(dto dto.RegistrationAddressDTO, id uuid.UUID) entity.RegistrationAddress {
	return entity.RegistrationAddress{
		ID_RegAddress: id,
		MailIndex:     dto.MailIndex,
		Region:        dto.Region,
		City:          dto.City,
		Street:        dto.Street,
		House:         dto.House,
		Building:      dto.Building,
		Apartment:     dto.Apartment,
	}
}

func MapListenerToEntity(dto dto.ListenerDTO, id uuid.UUID) (*entity.Listener, error) {

	DateOfBirth, err := utils.TimeParse(dto.DateOfBirth)
	if err != nil {
		return nil, err
	}

	return &entity.Listener{
		ID_Listener:  id,
		FirstName:    dto.FirstName,
		SecondName:   dto.SecondName,
		MiddleName:   dto.MiddleName,
		DateOfBirth:  *DateOfBirth,
		SNILS:        dto.SNILS,
		ContactPhone: dto.ContactPhone,
		Email:        dto.Email,
	}, nil
}

func MapListenerEntityToDTO(entity *entity.Listener) *dto.FullListenerDataDTO {

	listener := dto.RawOnlyListener{
		ID_Listener:          entity.ID_Listener,
		FirstName:            entity.FirstName,
		SecondName:           entity.SecondName,
		MiddleName:           entity.MiddleName,
		DateOfBirth:          entity.DateOfBirth.String(),
		SNILS:                entity.SNILS,
		ContactPhone:         entity.ContactPhone,
		Email:                entity.Email,
		ID_Passport:          entity.ID_Passport,
		ID_RegAddress:        entity.ID_RegAddress,
		ID_EducationListener: entity.ID_EducationListener,
		ID_PlaceWork:         entity.ID_PlaceWork,
	}

	var passport *dto.PassportDTO
	if entity.ID_Passport != nil {
		passport = &dto.PassportDTO{
			PlaceBirth:    entity.Passport.PlaceBirth,
			Citizenship:   entity.Passport.Citizenship,
			Gender:        entity.Passport.Gender,
			Seria:         entity.Passport.Seria,
			Number:        entity.Passport.Number,
			PassportGiven: entity.Passport.PassportGiven,
			DateGiven:     entity.Passport.DateGiven.String(),
			Code:          entity.Passport.Code,
		}
	}

	registratoinAddress := dto.RegistrationAddressDTO{
		MailIndex: entity.RegistrationAddress.MailIndex,
		Region:    entity.RegistrationAddress.Region,
		City:      entity.RegistrationAddress.City,
		Street:    entity.RegistrationAddress.Street,
		House:     entity.RegistrationAddress.House,
		Building:  entity.RegistrationAddress.Building,
		Apartment: entity.RegistrationAddress.Apartment,
	}

	var education *dto.EducationListenerDTO
	if entity.ID_EducationListener != nil {
		education = &dto.EducationListenerDTO{
			DiplomSeria:            entity.EducationListener.DiplomSeria,
			DiplomNumber:           entity.EducationListener.DiplomNumber,
			DateGiven:              entity.EducationListener.DateGiven.String(),
			City:                   entity.EducationListener.City,
			Region:                 entity.EducationListener.Region,
			EducationalInstitution: entity.EducationListener.EducationalInstitution,
			Speciality:             entity.EducationListener.Speciality,
			LevelEducation:         entity.EducationListener.LevelEducation,
		}
	}

	var placework *dto.PlaceWorkDTO
	if entity.ID_PlaceWork != nil {
		placework = &dto.PlaceWorkDTO{
			NameCompany:        entity.PlaceWork.NameCompany,
			JobTitle:           entity.PlaceWork.JobTitle,
			AllExperience:      entity.PlaceWork.AllExperience,
			JobTitleExpirience: entity.PlaceWork.JobTitleExperience,
		}
	}

	return &dto.FullListenerDataDTO{
		Listener:            listener,
		Passport:            passport,
		RegistrationAddress: registratoinAddress,
		EducationListener:   education,
		PlaceWork:           placework,
	}
}

func MapListenerReqToDto(request request.FullListenerRequest) (*dto.CreateListenerDTO, error) {

	var err error
	var seriaPassport, numberPassport int
	if request.Passport.Seria != "" {
		seriaPassport, err = strconv.Atoi(request.Passport.Seria)
		if err != nil {
			return &dto.CreateListenerDTO{}, err
		}

		numberPassport, err = strconv.Atoi(request.Passport.Number)
		if err != nil {
			return &dto.CreateListenerDTO{}, err
		}
	}

	passport := dto.PassportDTO{
		PlaceBirth:    request.Passport.PlaceBirth,
		Citizenship:   request.Passport.Citizenship,
		Gender:        request.Passport.Gender,
		Seria:         seriaPassport,
		Number:        numberPassport,
		PassportGiven: request.Passport.PassportGiven,
		DateGiven:     request.Passport.DateGiven,
		Code:          request.Passport.Code,
	}

	mailIndex, err := strconv.Atoi(request.RegistrationAddress.MailIndex)
	if err != nil {
		return &dto.CreateListenerDTO{}, err
	}
	regAddress := dto.RegistrationAddressDTO{
		MailIndex: mailIndex,
		Region:    request.RegistrationAddress.Region,
		City:      request.RegistrationAddress.City,
		Street:    request.RegistrationAddress.Street,
		House:     request.RegistrationAddress.House,
		Building:  request.RegistrationAddress.Building,
		Apartment: request.RegistrationAddress.Apartment,
	}

	var diplomSeria, diplomNumber int
	if request.EducationListener.DiplomSeria != "" || request.EducationListener.DiplomNumber != "" {
		diplomSeria, err = strconv.Atoi(request.EducationListener.DiplomSeria)
		if err != nil {
			return &dto.CreateListenerDTO{}, err
		}

		diplomNumber, err = strconv.Atoi(request.EducationListener.DiplomNumber)
		if err != nil {
			return &dto.CreateListenerDTO{}, err
		}
	}

	educaitonListener := dto.EducationListenerDTO{
		DiplomSeria:            diplomSeria,
		DiplomNumber:           diplomNumber,
		DateGiven:              request.EducationListener.DateGiven,
		City:                   request.EducationListener.City,
		Region:                 request.EducationListener.Region,
		EducationalInstitution: request.EducationListener.EducationalInstitution,
		Speciality:             request.EducationListener.Speciality,
		LevelEducation:         request.EducationListener.LevelEducation,
	}

	return &dto.CreateListenerDTO{
		Listener:            dto.ListenerDTO(request.Listener),
		Passport:            passport,
		RegistrationAddress: regAddress,
		EducationListener:   educaitonListener,
		PlaceWork:           dto.PlaceWorkDTO(request.PlaceWork),
	}, err
}

func MapContcratorReqToDTO(request request.FullContractorRequest) (*dto.ContractorCreateDTO, error) {
	passportSeria, err := strconv.Atoi(request.Passport.Seria)
	if err != nil {
		return nil, err
	}

	passportNumber, err := strconv.Atoi(request.Passport.Number)
	if err != nil {
		return nil, err
	}

	mainIndex, err := strconv.Atoi(request.RegAddress.MailIndex)
	if err != nil {
		return nil, err
	}

	dto := &dto.ContractorCreateDTO{
		Contractor: dto.ContractorDTO{
			FirstName:     request.Contractor.FirstName,
			SecondName:    request.Contractor.SecondName,
			MiddleName:    request.Contractor.MiddleName,
			Contact_phone: request.Contractor.Contact_phone,
			Email:         request.Contractor.Email,
		},
		Passport: dto.PassportDTO{
			PlaceBirth:    request.Passport.PlaceBirth,
			Citizenship:   request.Passport.Citizenship,
			Gender:        request.Passport.Gender,
			Seria:         passportSeria,
			Number:        passportNumber,
			PassportGiven: request.Passport.PassportGiven,
			DateGiven:     request.Passport.DateGiven,
			Code:          request.Passport.Code,
		},
		RegAddress: dto.RegistrationAddressDTO{
			MailIndex: mainIndex,
			Region:    request.RegAddress.Region,
			City:      request.RegAddress.City,
			Street:    request.RegAddress.Street,
			House:     request.RegAddress.House,
			Building:  request.RegAddress.Building,
			Apartment: request.RegAddress.Apartment,
		},
	}

	return dto, nil
}

func LegalEntityFullMappping(request request.FullLegalEntityRequest) (*dto.LegalEntityCreateDTO, error) {
	mailIndex, err := strconv.Atoi(request.RegAddress.MailIndex)
	if err != nil {
		return nil, err
	}

	data := &dto.LegalEntityCreateDTO{
		LegalEntity: dto.LegalEntityDTO{
			NameCompany: request.LegalEntity.NameCompany,
			Inn:         request.LegalEntity.Inn,
			Kpp:         request.LegalEntity.Kpp,
			Ogrn:        request.LegalEntity.Ogrn,
			Phone:       request.LegalEntity.Phone,
			Email:       request.LegalEntity.Email,
			FirstName:   request.LegalEntity.FirstName,
			SecondName:  request.LegalEntity.SecondName,
			MiddleName:  request.LegalEntity.MiddleName,
		},
		RegAddress: dto.RegistrationAddressDTO{
			MailIndex: mailIndex,
			Region:    request.RegAddress.Region,
			City:      request.RegAddress.City,
			Street:    request.RegAddress.Street,
			House:     request.RegAddress.House,
			Building:  request.RegAddress.Building,
			Apartment: request.RegAddress.Apartment,
		},
	}
	return data, nil
}

// func MapUpdateListenerReqToDto(request request.FullListenerRequest) (*dto.CreateListenerDTO, error) {

// 	seria, err := strconv.Atoi(request.Passport.Seria)
// 	if err != nil {
// 		return &dto.CreateListenerDTO{}, err
// 	}

// 	number, err := strconv.Atoi(request.Passport.Number)
// 	if err != nil {
// 		return &dto.CreateListenerDTO{}, err
// 	}

// 	passport := dto.PassportDTO{
// 		PlaceBirth:    request.Passport.PlaceBirth,
// 		Citizenship:   request.Passport.Citizenship,
// 		Gender:        request.Passport.Gender,
// 		Seria:         seria,
// 		Number:        number,
// 		PassportGiven: request.Passport.PassportGiven,
// 		DateGiven:     request.Passport.DateGiven,
// 		Code:          request.Passport.Code,
// 	}

// 	return &dto.CreateListenerDTO{
// 		Listener:            dto.ListenerDTO(request.Listener),
// 		Passport:            passport,
// 		RegistrationAddress: dto.RegistrationAddressDTO(request.RegistrationAddress),
// 		EducationListener:   dto.EducationListenerDTO(request.EducationListener),
// 		PlaceWork:           dto.PlaceWorkDTO(request.PlaceWork),
// 	}, nil
// }
