package request

import (
	"github.com/google/uuid"
)

type PassportRequest struct {
	PlaceBirth    string `json:"place_birth" validate:"required,max=50"`
	Citizenship   string `json:"citizenship" validate:"required,max=50"`
	Gender        string `json:"gender" validate:"required,max=7"`
	Seria         string `json:"seria" validate:"required,len=4,numeric"`
	Number        string `json:"number" validate:"required,len=6,numeric"`
	PassportGiven string `json:"passport_given" validate:"required,max=100"`
	DateGiven     string `json:"date_given" validate:"required"`
	Code          string `json:"code" validate:"required,len=7"`
}

type RegAddressRequest struct {
	MailIndex string `json:"mail_index" validate:"required,len=6"`
	Region    string `json:"region" validate:"required,max=100"`
	City      string `json:"city" validate:"required,max=100"`
	Street    string `json:"street" validate:"required,max=100"`
	House     string `json:"house" validate:"required,max=25"`
	Building  string `json:"building" validate:"required,max=50"`
	Apartment string `json:"apartment" validate:"required,max=50"`
}

type ListenerEducationRequest struct {
	Education string `json:"education" validate:"required,max=50"`
}

type EducationListenerRequest struct {
	DiplomSeria            string `json:"diplom_seria" validate:"required,len=6,numeric"`
	DiplomNumber           string `json:"diplom_number" validate:"required,len=7,numeric"`
	DateGiven              string `json:"date_given" validate:"required"`
	City                   string `json:"city" validate:"required,max=50"`
	Region                 string `json:"region" validate:"required,max=100"`
	EducationalInstitution string `json:"educational_institution" validate:"required,max=100"`
	Speciality             string `json:"speciality" validate:"required,max=100"`
	LevelEducation         string `json:"level_education" validate:"required,max=50"`
}

type PlaceWorkRequest struct {
	NameCompany        string `json:"name_company" validate:"required,max=100"`
	JobTitle           string `json:"job_title" validate:"required,max=100"`
	AllExperience      int    `json:"all_experience" validate:"required,max=20"`
	JobTitleExpirience int    `json:"job_title_expirience" validate:"required,max=20"`
}

type ProgramEducationRequest struct {
	NameProfEducation     string    `json:"name_prof_education" validate:"required,max=100"`
	TimeEducation         int       `json:"time_education" validate:"required"`
	IndividualPrice       float32   `json:"individual_price" validate:"required"`
	GroupPrice            float32   `json:"group_price" validate:"required"`
	CampusPrice           float32   `json:"campus_price" validate:"required"`
	ID_DivisionsEducation uuid.UUID `json:"id_divisionseducation" validate:"required,uuid"`
	ID_EducationType      uuid.UUID `json:"id_educationtype" validate:"required,uuid"`
}

type ListenerRequest struct {
	FirstName    string `json:"first_name" validate:"required,max=50"`
	SecondName   string `json:"second_name" validate:"required,max=50"`
	MiddleName   string `json:"middle_name" validate:"max=50"`
	DateOfBirth  string `json:"date_of_birth" validate:"required"`
	SNILS        string `json:"snils" validate:"required,len=14"`
	ContactPhone string `json:"contact_phone" validate:"required,len=12"`
	Email        string `json:"email" validate:"required,max=50,email"`
}

type FullListenerRequest struct {
	Listener            ListenerRequest          `json:"listener" validate:"required"`
	Passport            PassportRequest          `json:"passport" validate:"required"`
	RegistrationAddress RegAddressRequest        `json:"registration_address" validate:"required"`
	EducationListener   EducationListenerRequest `json:"education" validate:"omitempty"`
	PlaceWork           PlaceWorkRequest         `json:"placeWork" validate:"omitempty"`
}

type DivisionsEducationRequest struct {
	Divisions string `json:"divisions" validate:"required,max=100"`
}
type EducationTypeRequest struct {
	TypeName string `json:"type_name" validate:"required,max=50"`
}

type EnrollmentListenerRequest struct {
	ID_Listener         uuid.UUID `json:"id_listener" validate:"required,uuid"`
	ID_ProgramEducation uuid.UUID `json:"id_program" validate:"required,uuid"`
	StartDate           string    `json:"start_date" validate:"required"`
	EndDate             string    `json:"end_date" validate:"required"`
	CurrentPrice        float32   `json:"current_price" validate:"required"`
	Is_active           bool      `json:"is_active" validate:"required"`
}

type EnrollemenUpdateRequest struct {
	ID_ProgramEducation uuid.UUID `json:"id_program" validate:"required,uuid"`
	StartDate           string    `json:"start_date" validate:"required"`
	EndDate             string    `json:"end_date" validate:"required"`
	CurrentPrice        float32   `json:"current_price" validate:"required"`
}

type CreateCardRequest struct {
	ID_listener uuid.UUID `json:"id_listener" validate:"required,uuid"`
	ID_program  uuid.UUID `json:"id_program" validate:"required,uuid"`
}
