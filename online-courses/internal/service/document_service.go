package service

import (
	"context"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/repository"
	"online-courses/internal/server/http/request"

	"github.com/google/uuid"
)

type DocumentService struct {
	repo repository.DocumentRepository
}

func NewDocumentService(repo repository.DocumentRepository) *DocumentService {
	return &DocumentService{repo: repo}
}

func (d *DocumentService) FormingDataDocument(ctx context.Context, idListener, idProgram, idExecutor uuid.UUID, frontData request.FrontDataDeliverRequest) (*dto.FullDocumentInfoDTO, error) {

	prepareInfo, err := d.repo.PrepareDataDocument(ctx, idListener, idProgram, idExecutor)
	if err != nil {
		return nil, err
	}

	// mainIndexLegalEntity, err := strconv.Atoi(frontData.LegalEntity.Address.MailIndex)
	// if err != nil {
	// 	return nil, err
	// }

	var dtoListenerLegalEntity []dto.ListenerInLegalEntity

	for _, i := range frontData.LegalEntity.Listeners {
		dtoListenerLegalEntity = append(dtoListenerLegalEntity, dto.ListenerInLegalEntity{
			FirstName:   i.FirstName,
			SecondName:  i.SecondName,
			MiddleName:  i.MiddleName,
			SNILS:       i.SNILS,
			DateOfBirth: i.DateOfBirth,
			Email:       i.Email,
		})
	}

	dto := dto.FullDocumentInfoDTO{
		PersonalCardInfo: dto.PersonalCardInfoDTO{
			Listener: dto.ListenerDTO{
				FirstName:    prepareInfo.PersonalInfo.FirstName,
				SecondName:   prepareInfo.PersonalInfo.SecondName,
				MiddleName:   prepareInfo.PersonalInfo.MiddleName,
				DateOfBirth:  prepareInfo.PersonalInfo.DateOfBirth,
				SNILS:        prepareInfo.PersonalInfo.SNILS,
				ContactPhone: prepareInfo.PersonalInfo.ContactPhone,
				Email:        prepareInfo.PersonalInfo.Email,
			},
			Passport: dto.PassportCardDTO{
				PlaceBirth:    prepareInfo.PersonalInfo.PlaceBirth.String,
				Citizenship:   prepareInfo.PersonalInfo.Citizenship.String,
				Gender:        prepareInfo.PersonalInfo.Gender.String,
				Seria:         prepareInfo.PersonalInfo.Seria.String,
				Number:        prepareInfo.PersonalInfo.Number.String,
				PassportGiven: prepareInfo.PersonalInfo.PassportGiven.String,
				DateGiven:     prepareInfo.PersonalInfo.DateGiven.String,
				Code:          prepareInfo.PersonalInfo.Code.String,
			},
			RegistrationAddress: dto.RegistrationAddressCardDTO{
				MailIndex: prepareInfo.PersonalInfo.MailIndex,
				Region:    prepareInfo.PersonalInfo.RegRegion,
				City:      prepareInfo.PersonalInfo.RegCity,
				Street:    prepareInfo.PersonalInfo.RegStreet,
				House:     prepareInfo.PersonalInfo.House,
				Building:  prepareInfo.PersonalInfo.Building,
				Apartment: prepareInfo.PersonalInfo.Apartment,
			},
			EducationListener: dto.EducationListenerCardDTO{
				DiplomSeria:            prepareInfo.PersonalInfo.DiplomSeria.String,
				DiplomNumber:           prepareInfo.PersonalInfo.DiplomNumber.String,
				City:                   prepareInfo.PersonalInfo.City.String,
				Region:                 prepareInfo.PersonalInfo.Region.String,
				EducationalInstitution: prepareInfo.PersonalInfo.EducationalInstitution.String,
				Speciality:             prepareInfo.PersonalInfo.Speciality.String,
				LevelEducation:         prepareInfo.PersonalInfo.LevelEducation.String,
			},
			PlaceWork: dto.PlaceWorkDTO{
				NameCompany:        prepareInfo.PersonalInfo.NameCompany.String,
				JobTitle:           prepareInfo.PersonalInfo.JobTitle.String,
				AllExperience:      int(prepareInfo.PersonalInfo.AllExperience.Int32),
				JobTitleExpirience: int(prepareInfo.PersonalInfo.JobTitleExpirience.Int32),
			},
			ProgramEducation: dto.ProgramEducationToCardDTO{
				NameProfEducation: prepareInfo.PersonalInfo.NameProfEducation,
				TimeEducation:     prepareInfo.PersonalInfo.TimeEducation,
				DivisionEducation: prepareInfo.PersonalInfo.DivisionEducation,
				EducationType:     prepareInfo.PersonalInfo.EducationType,
			},
			EnrollmentListener: dto.EnrollmentListenerToCard{
				StartDate:        prepareInfo.PersonalInfo.StartDate,
				EndDate:          prepareInfo.PersonalInfo.EndDate,
				CurrentPrice:     prepareInfo.PersonalInfo.CurrentPrice,
				Is_active:        prepareInfo.PersonalInfo.Is_active,
				Group:            prepareInfo.PersonalInfo.Group,
				TypeOfRetraining: prepareInfo.PersonalInfo.TypeOfRetraining,
			},
		},
		ZayavlenieCardInfo: dto.ZayavlenieCardInfo{
			ProgramEducation: dto.ProgramEducationToCardDTO{
				NameProfEducation: prepareInfo.PersonalInfo.NameProfEducation,
				TimeEducation:     prepareInfo.PersonalInfo.TimeEducation,
				DivisionEducation: prepareInfo.PersonalInfo.DivisionEducation,
				EducationType:     prepareInfo.PersonalInfo.EducationType,
			},
			Listener: dto.ListenerDTO{
				FirstName:    prepareInfo.PersonalInfo.FirstName,
				SecondName:   prepareInfo.PersonalInfo.SecondName,
				MiddleName:   prepareInfo.PersonalInfo.MiddleName,
				DateOfBirth:  prepareInfo.PersonalInfo.DateOfBirth,
				SNILS:        prepareInfo.PersonalInfo.SNILS,
				ContactPhone: prepareInfo.PersonalInfo.ContactPhone,
				Email:        prepareInfo.PersonalInfo.Email,
			},
			EnrollmentListener: dto.EnrollmentListenerToCard{
				StartDate:        prepareInfo.PersonalInfo.StartDate,
				EndDate:          prepareInfo.PersonalInfo.EndDate,
				CurrentPrice:     prepareInfo.PersonalInfo.CurrentPrice,
				Is_active:        prepareInfo.PersonalInfo.Is_active,
				Group:            prepareInfo.PersonalInfo.Group,
				TypeOfRetraining: prepareInfo.PersonalInfo.TypeOfRetraining,
			},
			Contractor: dto.ContractorCardInfo{
				Passport: dto.PassportCardDTO{
					PlaceBirth:    prepareInfo.Contractor.PlaceBirth.String,
					Citizenship:   prepareInfo.Contractor.Citizenship.String,
					Gender:        prepareInfo.Contractor.Gender.String,
					Seria:         prepareInfo.Contractor.Seria.String,
					Number:        prepareInfo.Contractor.Number.String,
					PassportGiven: prepareInfo.Contractor.PassportGiven.String,
					DateGiven:     prepareInfo.Contractor.DateGiven.String,
					Code:          prepareInfo.Contractor.Code.String,
				},
				RegistrationAddress: dto.RegistrationAddressCardDTO{
					MailIndex: prepareInfo.Contractor.MailIndex.String,
					Region:    prepareInfo.Contractor.RegRegion.String,
					City:      prepareInfo.Contractor.RegCity.String,
					Street:    prepareInfo.Contractor.RegStreet.String,
					House:     prepareInfo.Contractor.House.String,
					Building:  prepareInfo.Contractor.Building.String,
					Apartment: prepareInfo.Contractor.Apartment.String,
				},
				FirstName:     prepareInfo.Contractor.FirstName.String,
				SecondName:    prepareInfo.Contractor.SecondName.String,
				MiddleName:    prepareInfo.Contractor.MiddleName.String,
				Contact_phone: prepareInfo.Contractor.ContactPhone.String,
				Email:         prepareInfo.Contractor.Email.String,
			},
			Executor: dto.ExecutorCardInfo{
				Status:             prepareInfo.Executor.Status,
				ExecutorName:       prepareInfo.Executor.FirstName,
				ExecutorSurname:    prepareInfo.Executor.SecondName,
				ExecutorMiddlename: prepareInfo.Executor.MiddleName,
				Doverenost:         prepareInfo.Executor.Doverenost,
			},
			Passport: dto.PassportCardDTO{
				PlaceBirth:    prepareInfo.PersonalInfo.PlaceBirth.String,
				Citizenship:   prepareInfo.PersonalInfo.Citizenship.String,
				Gender:        prepareInfo.PersonalInfo.Gender.String,
				Seria:         prepareInfo.PersonalInfo.Seria.String,
				Number:        prepareInfo.PersonalInfo.Number.String,
				PassportGiven: prepareInfo.PersonalInfo.PassportGiven.String,
				DateGiven:     prepareInfo.PersonalInfo.DateGiven.String,
				Code:          prepareInfo.PersonalInfo.Code.String,
			},
			Registration: dto.RegistrationAddressCardDTO{
				MailIndex: prepareInfo.PersonalInfo.MailIndex,
				Region:    prepareInfo.PersonalInfo.RegRegion,
				City:      prepareInfo.PersonalInfo.RegCity,
				Street:    prepareInfo.PersonalInfo.RegStreet,
				House:     prepareInfo.PersonalInfo.House,
				Building:  prepareInfo.PersonalInfo.Building,
				Apartment: prepareInfo.PersonalInfo.Apartment,
			},
			Variant:     frontData.Variant,
			DogovorType: frontData.DogovorAgeType,
		},
		DogovorRequest: dto.DogovorCardInfo{
			LegalEntity: dto.LegalEntity{
				Listeners: dtoListenerLegalEntity,
				Address: dto.RegistrationAddressLegalEntity{
					MailIndex: frontData.LegalEntity.Address.MailIndex,
					Region:    frontData.LegalEntity.Address.Region,
					City:      frontData.LegalEntity.Address.City,
					Street:    frontData.LegalEntity.Address.Street,
					House:     frontData.LegalEntity.Address.House,
					Building:  frontData.LegalEntity.Address.Building,
					Apartment: frontData.LegalEntity.Address.Apartment,
				},
				CompanyName: frontData.LegalEntity.CompanyName,
				FIO:         frontData.LegalEntity.FIO,
				Status:      frontData.LegalEntity.Status,
				INN:         frontData.LegalEntity.INN,
				KPP:         frontData.LegalEntity.KPP,
				OGRN:        frontData.LegalEntity.OGRN,
				Phone:       frontData.LegalEntity.Phone,
				Email:       frontData.LegalEntity.Email,
			},
			ProgramEducation: dto.ProgramEducationToCardDTO{
				NameProfEducation: prepareInfo.PersonalInfo.NameProfEducation,
				TimeEducation:     prepareInfo.PersonalInfo.TimeEducation,
				DivisionEducation: prepareInfo.PersonalInfo.DivisionEducation,
				EducationType:     prepareInfo.PersonalInfo.EducationType,
			},
			ListenerData: dto.ListenerDTO{
				FirstName:    prepareInfo.PersonalInfo.FirstName,
				SecondName:   prepareInfo.PersonalInfo.SecondName,
				MiddleName:   prepareInfo.PersonalInfo.MiddleName,
				DateOfBirth:  prepareInfo.PersonalInfo.DateOfBirth,
				SNILS:        prepareInfo.PersonalInfo.SNILS,
				ContactPhone: prepareInfo.PersonalInfo.ContactPhone,
				Email:        prepareInfo.PersonalInfo.Email,
			},
			Contractor: dto.ContractorCardInfo{
				Passport: dto.PassportCardDTO{
					PlaceBirth:    prepareInfo.Contractor.PlaceBirth.String,
					Citizenship:   prepareInfo.Contractor.Citizenship.String,
					Gender:        prepareInfo.Contractor.Gender.String,
					Seria:         prepareInfo.Contractor.Seria.String,
					Number:        prepareInfo.Contractor.Number.String,
					PassportGiven: prepareInfo.Contractor.PassportGiven.String,
					DateGiven:     prepareInfo.Contractor.DateGiven.String,
					Code:          prepareInfo.Contractor.Code.String,
				},
				RegistrationAddress: dto.RegistrationAddressCardDTO{
					MailIndex: prepareInfo.Contractor.MailIndex.String,
					Region:    prepareInfo.Contractor.RegRegion.String,
					City:      prepareInfo.Contractor.RegCity.String,
					Street:    prepareInfo.Contractor.RegStreet.String,
					House:     prepareInfo.Contractor.House.String,
					Building:  prepareInfo.Contractor.Building.String,
					Apartment: prepareInfo.Contractor.Apartment.String,
				},
				FirstName:     prepareInfo.Contractor.FirstName.String,
				SecondName:    prepareInfo.Contractor.SecondName.String,
				MiddleName:    prepareInfo.Contractor.MiddleName.String,
				Contact_phone: prepareInfo.Contractor.ContactPhone.String,
				Email:         prepareInfo.Contractor.Email.String,
			},
			Executor: dto.ExecutorCardInfo{
				Status:             prepareInfo.Executor.Status,
				ExecutorName:       prepareInfo.Executor.FirstName,
				ExecutorSurname:    prepareInfo.Executor.SecondName,
				ExecutorMiddlename: prepareInfo.PersonalInfo.MiddleName,
				Doverenost:         prepareInfo.Executor.Doverenost,
			},
			Passport: dto.PassportCardDTO{
				PlaceBirth:    prepareInfo.PersonalInfo.PlaceBirth.String,
				Citizenship:   prepareInfo.PersonalInfo.Citizenship.String,
				Gender:        prepareInfo.PersonalInfo.Gender.String,
				Seria:         prepareInfo.PersonalInfo.Seria.String,
				Number:        prepareInfo.PersonalInfo.Number.String,
				PassportGiven: prepareInfo.PersonalInfo.PassportGiven.String,
				DateGiven:     prepareInfo.PersonalInfo.DateGiven.String,
				Code:          prepareInfo.PersonalInfo.Code.String,
			},
			Registration: dto.RegistrationAddressCardDTO{
				MailIndex: prepareInfo.PersonalInfo.MailIndex,
				Region:    prepareInfo.PersonalInfo.RegRegion,
				City:      prepareInfo.PersonalInfo.RegCity,
				Street:    prepareInfo.PersonalInfo.RegStreet,
				House:     prepareInfo.PersonalInfo.House,
				Building:  prepareInfo.PersonalInfo.Building,
				Apartment: prepareInfo.PersonalInfo.Apartment,
			},
			Enrollment: dto.EnrollmentListenerToCard{
				NameProfEducation: prepareInfo.PersonalInfo.NameProfEducation,
				StartDate:         prepareInfo.PersonalInfo.StartDate,
				EndDate:           prepareInfo.PersonalInfo.EndDate,
				CurrentPrice:      prepareInfo.PersonalInfo.CurrentPrice,
				Is_active:         prepareInfo.PersonalInfo.Is_active,
				Group:             prepareInfo.PersonalInfo.Group,
				TypeOfRetraining:  prepareInfo.PersonalInfo.TypeOfRetraining,
			},
			OptionNagruzka: frontData.OptionNagruzka,
			OptionDocument: frontData.OptionDocument,
			OptionPrice:    frontData.OptionPrice,
			DogovorType:    frontData.DogovorType,
		},
	}

	return &dto, nil
}
