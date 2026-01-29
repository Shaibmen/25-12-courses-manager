package mapper

import (
	"document-service/internal/domain/dto"
	"document-service/internal/server/models"
	"time"
	"unsafe"
)

func FullListenerMapping(request models.FullListenerRequest) (*dto.FullListenerDataDTO, error) {

	// если понядобится int

	// passportSeria, err := strconv.Atoi(request.Passport.Seria)
	// if err != nil {
	// 	return nil, err
	// }

	// passportNumber, err := strconv.Atoi(request.Passport.Number)
	// if err != nil {
	// 	return nil, err
	// }

	startDate, err := time.Parse(time.RFC3339, request.EnrollmentListener.StartDate)
	if err != nil {
		return &dto.FullListenerDataDTO{}, err
	}
	endDate, err := time.Parse(time.RFC3339, request.EnrollmentListener.EndDate)
	if err != nil {
		return &dto.FullListenerDataDTO{}, err
	}

	dto := &dto.FullListenerDataDTO{
		Listener: dto.ListenerDTO{
			FirstName:    request.Listener.FirstName,
			SecondName:   request.Listener.SecondName,
			MiddleName:   request.Listener.MiddleName,
			DateOfBirth:  request.Listener.DateOfBirth,
			SNILS:        request.Listener.SNILS,
			ContactPhone: request.Listener.ContactPhone,
			Email:        request.Listener.ContactPhone,
		},
		Passport: dto.PassportDTO{
			PlaceBirth:    request.Passport.PlaceBirth,
			Citizenship:   request.Passport.Citizenship,
			Gender:        request.Passport.Gender,
			Seria:         request.Passport.Seria,
			Number:        request.Passport.Number,
			PassportGiven: request.Passport.PassportGiven,
			DateGiven:     request.Passport.DateGiven,
			Code:          request.Passport.Code,
		},
		RegistrationAddress: dto.RegistrationAddressDTO{
			MailIndex: request.RegistrationAddress.MailIndex,
			Region:    request.RegistrationAddress.Region,
			City:      request.RegistrationAddress.City,
			Street:    request.RegistrationAddress.Street,
			House:     request.RegistrationAddress.House,
			Building:  request.RegistrationAddress.Building,
			Apartment: request.RegistrationAddress.Apartment,
		},
		EducationListener: dto.EducationListenerDTO{
			DiplomSeria:            request.EducationListener.DiplomSeria,
			DiplomNumber:           request.EducationListener.DiplomNumber,
			DateGiven:              request.EducationListener.DateGiven,
			City:                   request.EducationListener.City,
			Region:                 request.EducationListener.Region,
			EducationalInstitution: request.EducationListener.EducationalInstitution,
			Speciality:             request.EducationListener.Speciality,
			LevelEducation:         request.EducationListener.LevelEducation,
		},
		PlaceWork: dto.PlaceWorkDTO{
			NameCompany:        request.PlaceWork.NameCompany,
			JobTitle:           request.PlaceWork.JobTitle,
			AllExperience:      request.PlaceWork.AllExperience,
			JobTitleExpirience: request.PlaceWork.JobTitleExpirience,
		},
		EnrollmentListener: dto.EnrollmentListenerDTO{
			StartDate:        startDate,
			EndDate:          endDate,
			TypeOfRetraining: request.EnrollmentListener.TypeOfRetraining,
			CurrentPrice:     request.EnrollmentListener.CurrentPrice,
		},
		ProgramEducation: dto.ProgramEducationDTO{
			NameProfEducation: request.ProgramEducation.NameProfEducation,
			TimeEducation:     request.ProgramEducation.TimeEducation,
			DivisionEducation: request.ProgramEducation.DivisionEducation,
			EducationType:     request.ProgramEducation.EducationType,
		},
	}

	return dto, nil
}

func ZayavlenieMapping(request models.ZayavlenieRequest) *dto.ZayavlenieDTO {
	dto := dto.ZayavlenieDTO{
		ProgramEducation: dto.ProgramEducationDTO{
			TimeEducation:     request.ProgramEducation.TimeEducation,
			NameProfEducation: request.ProgramEducation.NameProfEducation,
			DivisionEducation: request.ProgramEducation.DivisionEducation,
			EducationType:     request.ProgramEducation.EducationType,
		},
		ListenerData: dto.ListenerDTO{
			FirstName:    request.Listener.FirstName,
			SecondName:   request.Listener.SecondName,
			MiddleName:   request.Listener.MiddleName,
			DateOfBirth:  request.Listener.DateOfBirth,
			SNILS:        request.Listener.SNILS,
			ContactPhone: request.Listener.ContactPhone,
			Email:        request.Listener.Email,
		},
		Contractor: dto.ContractorDTO{
			Passport: dto.PassportDTO{
				PlaceBirth:    request.Contractor.Passport.PlaceBirth,
				Citizenship:   request.Contractor.Passport.Citizenship,
				Gender:        request.Contractor.Passport.Gender,
				Seria:         request.Contractor.Passport.Seria,
				Number:        request.Contractor.Passport.Number,
				PassportGiven: request.Contractor.Passport.PassportGiven,
				DateGiven:     request.Contractor.Passport.DateGiven,
				Code:          request.Contractor.Passport.Code,
			},
			RegistrationAddress: dto.RegistrationAddressDTO{
				MailIndex: request.Contractor.RegistrationAddress.MailIndex,
				Region:    request.Contractor.RegistrationAddress.Region,
				City:      request.Contractor.RegistrationAddress.City,
				Street:    request.Contractor.RegistrationAddress.Street,
				House:     request.Contractor.RegistrationAddress.House,
				Building:  request.Contractor.RegistrationAddress.Building,
				Apartment: request.Contractor.RegistrationAddress.Apartment,
			},
			FirstName:     request.Contractor.FirstName,
			SecondName:    request.Contractor.SecondName,
			MiddleName:    request.Contractor.MiddleName,
			Contact_phone: request.Contractor.Contact_phone,
			Email:         request.Contractor.Email,
		},
		Executor: dto.ExecutorDTO{
			Status:             request.Executor.Status,
			ExecutorName:       request.Executor.ExecutorName,
			ExecutorSurname:    request.Executor.ExecutorSurname,
			ExecutorMiddlename: request.Executor.ExecutorMiddlename,
		},
		Passport: dto.PassportDTO{
			PlaceBirth:    request.Passport.PlaceBirth,
			Citizenship:   request.Passport.Citizenship,
			Gender:        request.Passport.Gender,
			Seria:         request.Passport.Seria,
			Number:        request.Passport.Number,
			PassportGiven: request.Passport.PassportGiven,
			DateGiven:     request.Passport.DateGiven,
			Code:          request.Passport.Code,
		},
		Registration: dto.RegistrationAddressDTO{
			MailIndex: request.Registration.MailIndex,
			Region:    request.Registration.Region,
			City:      request.Registration.City,
			Street:    request.Registration.Street,
			House:     request.Registration.House,
			Building:  request.Registration.Building,
			Apartment: request.Registration.Apartment,
		},
		Variant: request.Variant,
	}
	return &dto
}

func DogovorMapping(request models.DogovorRequest) (*dto.DogovorDTO, error) {
	startDate, err := time.Parse(time.RFC3339, request.Enrollment.StartDate)
	if err != nil {
		return &dto.DogovorDTO{}, err
	}
	endDate, err := time.Parse(time.RFC3339, request.Enrollment.EndDate)
	if err != nil {
		return &dto.DogovorDTO{}, err
	}

	dto := dto.DogovorDTO{
		Zakazchik: dto.ZakazchikDTO{
			Listeners:   *(*[]dto.ListenerDTO)(unsafe.Pointer(&request.ZakazchikData.Listeners)),
			Address:     (dto.RegistrationAddressDTO)(request.ZakazchikData.Address),
			CompanyName: request.ZakazchikData.CompanyName,
			FIO:         request.ZakazchikData.FIO,
			Status:      request.ZakazchikData.Status,
			Osnovanie:   request.ZakazchikData.Osnovanie,
			INN:         request.ZakazchikData.INN,
			KPP:         request.ZakazchikData.KPP,
			OGRN:        request.ZakazchikData.OGRN,
			Phone:       request.ZakazchikData.Phone,
			Email:       request.Contractor.Email,
		},
		ProgramEducation: dto.ProgramEducationDTO{
			TimeEducation:     request.ProgramEducation.TimeEducation,
			NameProfEducation: request.ProgramEducation.NameProfEducation,
			DivisionEducation: request.ProgramEducation.DivisionEducation,
			EducationType:     request.ProgramEducation.EducationType,
		},
		ListenerData: dto.ListenerDTO{
			FirstName:    request.ListenerData.FirstName,
			SecondName:   request.ListenerData.SecondName,
			MiddleName:   request.ListenerData.MiddleName,
			DateOfBirth:  request.ListenerData.DateOfBirth,
			SNILS:        request.ListenerData.SNILS,
			ContactPhone: request.ListenerData.ContactPhone,
			Email:        request.ListenerData.Email,
		},
		Contractor: dto.ContractorDTO{
			Passport: dto.PassportDTO{
				PlaceBirth:    request.Contractor.Passport.PlaceBirth,
				Citizenship:   request.Contractor.Passport.Citizenship,
				Gender:        request.Contractor.Passport.Gender,
				Seria:         request.Contractor.Passport.Seria,
				Number:        request.Contractor.Passport.Number,
				PassportGiven: request.Contractor.Passport.PassportGiven,
				DateGiven:     request.Contractor.Passport.DateGiven,
				Code:          request.Contractor.Passport.Code,
			},
			RegistrationAddress: dto.RegistrationAddressDTO{
				MailIndex: request.Contractor.RegistrationAddress.MailIndex,
				Region:    request.Contractor.RegistrationAddress.Region,
				City:      request.Contractor.RegistrationAddress.City,
				Street:    request.Contractor.RegistrationAddress.Street,
				House:     request.Contractor.RegistrationAddress.House,
				Building:  request.Contractor.RegistrationAddress.Building,
				Apartment: request.Contractor.RegistrationAddress.Apartment,
			},
			FirstName:     request.Contractor.FirstName,
			SecondName:    request.Contractor.SecondName,
			MiddleName:    request.Contractor.MiddleName,
			Contact_phone: request.Contractor.Contact_phone,
			Email:         request.Contractor.Email,
		},
		Executor: dto.ExecutorDTO{
			Status:             request.Executor.Status,
			ExecutorName:       request.Executor.ExecutorName,
			ExecutorSurname:    request.Executor.ExecutorSurname,
			ExecutorMiddlename: request.Executor.ExecutorMiddlename,
		},
		Passport: dto.PassportDTO{
			PlaceBirth:    request.Passport.PlaceBirth,
			Citizenship:   request.Passport.Citizenship,
			Gender:        request.Passport.Gender,
			Seria:         request.Passport.Seria,
			Number:        request.Passport.Number,
			PassportGiven: request.Passport.PassportGiven,
			DateGiven:     request.Passport.DateGiven,
			Code:          request.Passport.Code,
		},
		Registration: dto.RegistrationAddressDTO{
			MailIndex: request.Registration.MailIndex,
			Region:    request.Registration.Region,
			City:      request.Registration.City,
			Street:    request.Registration.Street,
			House:     request.Registration.House,
			Building:  request.Registration.Building,
			Apartment: request.Registration.Apartment,
		},
		Enrollment: dto.EnrollmentListenerDTO{
			StartDate:        startDate,
			EndDate:          endDate,
			TypeOfRetraining: request.Enrollment.TypeOfRetraining,
			CurrentPrice:     request.Enrollment.CurrentPrice,
		},
		OptionNagruzka: request.OptionNagruzka,
		OptionDocument: request.OptionDocument,
		OptionPrice:    request.OptionPrice,
	}
	return &dto, nil
}
