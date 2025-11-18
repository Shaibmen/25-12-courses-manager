package mapper

import (
	"document-service/internal/domain/dto"
	"document-service/internal/server/models"
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
		ProgramEducation: dto.ProgramEducationDTO{
			NameProfEducation: request.ProgramEducation.NameProfEducation,
			TimeEducation:     request.ProgramEducation.TimeEducation,
			DivisionEducation: request.ProgramEducation.DivisionEducation,
			EducationType:     request.ProgramEducation.EducationType,
		},
	}

	return dto, nil
}
