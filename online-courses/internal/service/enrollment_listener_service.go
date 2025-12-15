package service

import (
	"context"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/entity"
	"online-courses/internal/domain/repository"
	"online-courses/internal/utils"

	"github.com/google/uuid"
)

type enrollmentListenerService struct {
	service repository.EnrollmentListenerRepository
}

func NewEnrollmentListenerService(repo repository.EnrollmentListenerRepository) *enrollmentListenerService {
	return &enrollmentListenerService{service: repo}
}

func (e *enrollmentListenerService) Create(ctx context.Context, model dto.EnrollmentListenerDTO) error {

	startDate, err := utils.TimeParse(model.StartDate)
	if err != nil {
		return err
	}

	endDate, err := utils.TimeParse(model.EndDate)
	if err != nil {
		return err
	}

	entityModel := entity.EnrollmentListener{
		ID_Listener:         model.ID_Listener,
		ID_ProgramEducation: model.ID_Program,
		StartDate:           *startDate,
		EndDate:             *endDate,
		CurrentPrice:        model.CurrentPrice,
		Is_active:           model.Is_active,
		Group:               model.Group,
		TypeOfRetraining:    model.TypeOfRetraining,
	}

	if err := e.service.Create(ctx, entityModel); err != nil {
		return err
	}

	program, err := e.service.GetProgram(ctx, model.ID_Program)
	if err != nil {
		return nil
	}
	AccurateProgram := entity.AccurateProgram{
		ID_Listener:       model.ID_Listener,
		NameProfEducation: program.NameProfEducation,
		TimeEducation:     program.TimeEducation,
		IndividualPrice:   program.IndividualPrice,
		GroupPrice:        program.GroupPrice,
		CampusPrice:       program.CampusPrice,
		EducationType:     program.EducationType,
		Division:          program.Division,
	}

	if err = e.service.SaveInfo(ctx, AccurateProgram); err != nil {
		return err
	}

	return nil
}

func (e *enrollmentListenerService) Read(ctx context.Context, page int, filter string) ([]dto.EnrollmentListenerDetailsDTO, error) {

	data, err := e.service.Read(ctx, page, filter)
	if err != nil {
		return nil, err
	}

	var enrollments []dto.EnrollmentListenerDetailsDTO

	for _, i := range data {
		enrollments = append(enrollments, dto.EnrollmentListenerDetailsDTO{
			ID_Listener:       i.ID_Listener,
			FirstName:         i.FirstName,
			SecondName:        i.SecondName,
			MiddleName:        i.MiddleName,
			NameProfEducation: i.NameProfEducation,
			StartDate:         i.StartDate.String(),
			EndDate:           i.EndDate.String(),
			CurrentPrice:      i.CurrentPrice,
			Group:             i.Group,
			TypeOfRetraining:  i.TypeOfRetraining,
		})
	}

	return enrollments, nil
}

func (e *enrollmentListenerService) Update(ctx context.Context, idListener, idProgram uuid.UUID, model dto.EnrollmentListenerDTO) error {

	startDate, err := utils.TimeParse(model.StartDate)
	if err != nil {
		return err
	}

	endDate, err := utils.TimeParse(model.EndDate)
	if err != nil {
		return err
	}

	entity := entity.EnrollmentListener{
		ID_Listener:         model.ID_Listener,
		ID_ProgramEducation: model.ID_Program,
		StartDate:           *startDate,
		EndDate:             *endDate,
		CurrentPrice:        model.CurrentPrice,
		Group:               model.Group,
		TypeOfRetraining:    model.TypeOfRetraining,
	}

	if err := e.service.Update(ctx, idListener, idProgram, entity); err != nil {
		return err
	}

	return nil
}

func (e *enrollmentListenerService) Delete(ctx context.Context, id_listener, id_program uuid.UUID) error {

	if err := e.service.Delete(ctx, id_listener, id_program); err != nil {
		return err
	}

	return nil
}

func (e *enrollmentListenerService) ReadDetailListener(ctx context.Context, id uuid.UUID) ([]dto.EnrollmentProgramDetailsDTO, error) {

	data, err := e.service.ReadDetailListener(ctx, id)
	if err != nil {
		return nil, err
	}

	enrollments := []dto.EnrollmentProgramDetailsDTO{}
	for _, i := range data {
		enrollments = append(enrollments, dto.EnrollmentProgramDetailsDTO{
			ID_Listener:         i.ID_Listener,
			ID_ProgramEducation: i.ID_ProgramEducation,
			NameProfEducation:   i.NameProfEducation,
			TimeEducation:       i.TimeEducation,
			IndividualPrice:     i.IndividualPrice,
			GroupPrice:          i.GroupPrice,
			CampusPrice:         i.CampusPrice,
			EducationType:       i.DivisionEducation,
			DivisionEducation:   i.DivisionEducation,
			StartDate:           i.StartDate.String(),
			EndDate:             i.EndDate.String(),
			CurrentPrice:        i.CurrentPrice,
			Group:               i.Group,
			TypeOfRetraining:    i.TypeOfRetraining,
		})

	}

	return enrollments, nil

}

func (e *enrollmentListenerService) ReadByProgram(ctx context.Context, id uuid.UUID, page int) ([]dto.EnrollmentListenerDetailsDTO, error) {
	data, err := e.service.ReadByProgram(ctx, id, page)
	if err != nil {
		return nil, err
	}

	var enrollments []dto.EnrollmentListenerDetailsDTO

	for _, i := range data {
		enrollments = append(enrollments, dto.EnrollmentListenerDetailsDTO{
			ID_Listener:       i.ID_Listener,
			FirstName:         i.FirstName,
			SecondName:        i.SecondName,
			MiddleName:        i.MiddleName,
			NameProfEducation: i.NameProfEducation,
			StartDate:         i.StartDate.String(),
			EndDate:           i.EndDate.String(),
			CurrentPrice:      i.CurrentPrice,
			Group:             i.Group,
			TypeOfRetraining:  i.TypeOfRetraining,
		})
	}

	return enrollments, nil
}

// func (e *enrollmentListenerService) InfoToPersonalCard(ctx context.Context, listenerID, programID uuid.UUID) (*dto.PersonalCardInfoDTO, error) {
// 	data, err := e.service.InfoToPersonalCard(ctx, listenerID, programID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	dto := &dto.PersonalCardInfoDTO{
// 		Listener: dto.ListenerDTO{
// 			FirstName:    data.FirstName,
// 			SecondName:   data.SecondName,
// 			MiddleName:   data.MiddleName,
// 			DateOfBirth:  data.DateOfBirth,
// 			SNILS:        data.SNILS,
// 			ContactPhone: data.ContactPhone,
// 			Email:        data.Email,
// 		},
// 		Passport: dto.PassportCardDTO{
// 			PlaceBirth:    data.PlaceBirth,
// 			Citizenship:   data.Citizenship,
// 			Gender:        data.Gender,
// 			Seria:         data.Seria,
// 			Number:        data.Number,
// 			PassportGiven: data.PassportGiven,
// 			DateGiven:     data.DateGiven,
// 			Code:          data.Code,
// 		},
// 		RegistrationAddress: dto.RegistrationAddressCardDTO{
// 			MailIndex: data.MailIndex,
// 			Region:    data.RegRegion,
// 			City:      data.RegCity,
// 			Street:    data.RegStreet,
// 			House:     data.House,
// 			Building:  data.Building,
// 			Apartment: data.Apartment,
// 		},
// 		EducationListener: dto.EducationListenerCardDTO{
// 			DiplomSeria:            data.DiplomSeria.String,
// 			DiplomNumber:           data.DiplomNumber.String,
// 			DateGiven:              data.DiplomDateGiven.String,
// 			City:                   data.DiplomCity.String,
// 			Region:                 data.DiplomRegion.String,
// 			EducationalInstitution: data.EducationalInstitution.String,
// 			Speciality:             data.Speciality.String,
// 			LevelEducation:         data.LevelEducation.String,
// 		},
// 		PlaceWork: dto.PlaceWorkDTO{
// 			NameCompany:        data.NameCompany.String,
// 			JobTitle:           data.JobTitle.String,
// 			AllExperience:      int(data.AllExperience.Int32),
// 			JobTitleExpirience: int(data.JobTitleExpirience.Int32),
// 		},
// 		ProgramEducation: dto.ProgramEducationToCardDTO{
// 			NameProfEducation: data.NameProfEducation,
// 			TimeEducation:     data.TimeEducation,
// 			DivisionEducation: data.DivisionEducation,
// 			EducationType:     data.EducationType,
// 		},
// 	}

// 	return dto, nil
// }

func (e *enrollmentListenerService) GetListenerFIO(ctx context.Context, listenerID uuid.UUID) (*dto.ListenerFIODTO, error) {
	data, err := e.service.GetListenerFIO(ctx, listenerID)
	if err != nil {
		return nil, err
	}

	dto := dto.ListenerFIODTO{
		FirstName:  data.FirstName,
		SecondName: data.SecondName,
		MiddleName: data.MiddleName,
	}

	return &dto, nil
}
