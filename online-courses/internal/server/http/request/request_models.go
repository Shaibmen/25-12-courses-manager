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
	FirstName        string    `json:"first_name" validate:"required,max=50"`
	SecondName       string    `json:"second_name" validate:"required,max=50"`
	MiddleName       string    `json:"middle_name" validate:"max=50"`
	DateOfBirth      string    `json:"date_of_birth" validate:"required"`
	SNILS            string    `json:"snils" validate:"required,len=14"`
	ContactPhone     string    `json:"contact_phone" validate:"required,len=12"`
	Email            string    `json:"email" validate:"required,max=50,email"`
	ID_LegalEntity   uuid.UUID `json:"id_legalentity"`
	ID_Contractor    uuid.UUID `json:"id_contractor"`
	LootingEducation bool      `json:"looting_education"`
}

type FullListenerRequest struct {
	Listener            ListenerRequest          `json:"listener" validate:"required"`
	Passport            PassportRequest          `json:"passport" validate:"omitempty"`
	RegistrationAddress RegAddressRequest        `json:"registration_address" validate:"required"`
	EducationListener   EducationListenerRequest `json:"education" validate:"omitempty"`
	PlaceWork           PlaceWorkRequest         `json:"placeWork" validate:"omitempty"`
}

type ContractorRequest struct {
	FirstName     string `json:"first_name" validate:"required,max=255"`
	SecondName    string `json:"second_name" validate:"required,max=255"`
	MiddleName    string `json:"middle_name" validate:"required,max=255"`
	Contact_phone string `json:"contact_phone" validate:"required,max=15"`
	Email         string `json:"email" validate:"required,max=255"`
}

type FullContractorRequest struct {
	Contractor ContractorRequest `json:"contractor" validate:"required"`
	Passport   PassportRequest   `json:"passport" validate:"required"`
	RegAddress RegAddressRequest `json:"reg_address" validate:"required"`
}

type LegalEntityRequest struct {
	NameCompany string `json:"name_company" validate:"required,max=255"`
	Inn         string `json:"inn" validate:"required,max=20"`
	Kpp         string `json:"kpp" validate:"required,max=20"`
	Ogrn        string `json:"ogrn" validate:"required,max=20"`
	Phone       string `json:"phone" validate:"required,max=20"`
	Email       string `json:"email" validate:"required,max=255"`
	FirstName   string `json:"first_name" validate:"required,max=100"`
	SecondName  string `json:"second_name" validate:"required,max=100"`
	MiddleName  string `json:"middle_name" validate:"required,max=100"`
	Status      string `json:"status" validate:"required,max=100"`
}

type FullLegalEntityRequest struct {
	LegalEntity LegalEntityRequest `json:"legal_entity" validate:"required"`
	RegAddress  RegAddressRequest  `json:"reg_address" validate:"required"`
}

type ExecutorRequest struct {
	Status     string `json:"status" validate:"required"`
	FirstName  string `json:"first_name" validate:"required"`
	SecondName string `json:"second_name" validate:"required"`
	MiddleName string `json:"middle_name" validate:"required"`
	Doverenost string `json:"doverenost"`
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
	Group               string    `json:"group"`
	TypeOfRetraining    string    `json:"type_of_retraining"`
}

type EnrollemenUpdateRequest struct {
	ID_ProgramEducation uuid.UUID `json:"id_program" validate:"required,uuid"`
	StartDate           string    `json:"start_date" validate:"required"`
	EndDate             string    `json:"end_date" validate:"required"`
	CurrentPrice        float32   `json:"current_price" validate:"required"`
	Group               string    `json:"group" validate:"required"`
	TypeOfRetraining    string    `json:"type_of_retraining" validate:"required"`
}

type CreateCardRequest struct {
	ID_listener uuid.UUID `json:"id_listener" validate:"required,uuid"`
	ID_program  uuid.UUID `json:"id_program" validate:"required,uuid"`
}

type ListenerInLegalEntity struct {
	FirstName   string `json:"first_name"`
	SecondName  string `json:"second_name"`
	MiddleName  string `json:"middle_name"`
	SNILS       string `json:"snils"`
	DateOfBirth string `json:"date_of_birth"`
	Email       string `json:"email"`
}

type LegalEntity struct {
	Listeners   []ListenerInLegalEntity `json:"listeners"`
	Address     RegAddressRequest       `json:"reg_address"`
	CompanyName string                  `json:"company_name"`
	FIO         string                  `json:"zakazchikfio"`
	Status      string                  `json:"status"`
	INN         string                  `json:"inn"`
	KPP         string                  `json:"kpp"`
	OGRN        string                  `json:"ogrn"`
	Phone       string                  `json:"phone"`
	Email       string                  `json:"email"`
}

type DocumentsDataRequest struct {
	ID_Listener uuid.UUID               `json:"id_listener"`
	ID_Program  uuid.UUID               `json:"id_program"`
	ID_Executor uuid.UUID               `json:"id_executor"`
	FrontData   FrontDataDeliverRequest `json:"front_data"`
}

type FrontDataDeliverRequest struct {
	LegalEntity       LegalEntity `json:"legal_entity"`
	Variant           int         `json:"variant"`
	DogovorType       string      `json:"dogovor_type"`
	OptionNagruzka    int         `json:"opt_nagruz"`
	OptionDocument    int         `json:"opt_document"`
	DogovorAgeType    string      `json:"dogovor_age"`
	OptionPrice       string      `json:"opt_price"`
	StartDate         string      `json:"start_date"`
	EndDate           string      `json:"end_date"`
	NameProfEducation string      `json:"program_name"`
	DivisionEducation string 	  `json:"division_education"`
	CurrentPrice      float32     `json:"price_enrollment"`
	TimeEducation     int         `json:"time_education"`
}
