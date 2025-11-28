package dto

import (
	"time"

	"github.com/google/uuid"
)

type ListenerIDDTO struct {
	ID_Listener          uuid.UUID  `json:"id_listener"`
	ID_Passport          *uuid.UUID `json:"id_passport"`
	ID_RegAddress        uuid.UUID  `json:"id_reg_address"`
	ID_EducationListener *uuid.UUID `json:"id_education_listener"`
	ID_PlaceWork         *uuid.UUID `json:"id_placework"`
	ID_LegalEntity       *uuid.UUID `json:"id_legalentity"`
	ID_Contractor        *uuid.UUID `json:"id_contractor"`
}

type ListenerDTO struct {
	FirstName    string `json:"first_name"`
	SecondName   string `json:"second_name"`
	MiddleName   string `json:"middle_name"`
	DateOfBirth  string `json:"date_of_birth"`
	SNILS        string `json:"snils"`
	ContactPhone string `json:"contact_phone"`
	Email        string `json:"email"`
}

type ListenerDTOWithID struct {
	ID_Listener  uuid.UUID `json:"id_listener"`
	FirstName    string    `json:"first_name"`
	SecondName   string    `json:"second_name"`
	MiddleName   string    `json:"middle_name"`
	DateOfBirth  string    `json:"date_of_birth"`
	SNILS        string    `json:"snils"`
	ContactPhone string    `json:"contact_phone"`
	Email        string    `json:"email"`
}

type RawOnlyListener struct {
	ID_Listener          uuid.UUID  `json:"id_listener"`
	FirstName            string     `json:"first_name"`
	SecondName           string     `json:"second_name"`
	MiddleName           string     `json:"middle_name"`
	DateOfBirth          string     `json:"date_of_birth"`
	SNILS                string     `json:"snils"`
	ContactPhone         string     `json:"contact_phone"`
	Email                string     `json:"email"`
	ID_Passport          *uuid.UUID `json:"id_passport"`
	ID_RegAddress        uuid.UUID  `json:"id_reg_address"`
	ID_EducationListener *uuid.UUID `json:"id_education_listener"`
	ID_PlaceWork         *uuid.UUID `json:"id_placework"`
}

type PassportDTO struct {
	PlaceBirth    string `json:"place_birth"`
	Citizenship   string `json:"citizenship"`
	Gender        string `json:"gender"`
	Seria         int    `json:"seria"`
	Number        int    `json:"number"`
	PassportGiven string `json:"passport_given"`
	DateGiven     string `json:"date_given"`
	Code          string `json:"code"`
}

type RegistrationAddressDTO struct {
	MailIndex int    `json:"mail_index"`
	Region    string `json:"region"`
	City      string `json:"city"`
	Street    string `json:"street"`
	House     string `json:"house"`
	Building  string `json:"building"`
	Apartment string `json:"apartment"`
}

type LevelEducationDTO struct {
	ID_LevelEducation uuid.UUID `json:"id_level_education"`
	Education         string    `json:"education"`
}

type EducationListenerDTO struct {
	DiplomSeria            int    `json:"diplom_seria"`
	DiplomNumber           int    `json:"diplom_number"`
	DateGiven              string `json:"date_given"`
	City                   string `json:"city"`
	Region                 string `json:"region"`
	EducationalInstitution string `json:"educational_institution"`
	Speciality             string `json:"speciality"`
	LevelEducation         string `json:"level_education"`
}

type PlaceWorkDTO struct {
	NameCompany        string `json:"name_company"`
	JobTitle           string `json:"job_title"`
	AllExperience      int    `json:"all_experience"`
	JobTitleExpirience int    `json:"job_title_expirience"`
}

type FullListenerDataDTO struct {
	Listener            RawOnlyListener        `json:"listener"`
	Passport            *PassportDTO           `json:"passport,omitempty"`
	RegistrationAddress RegistrationAddressDTO `json:"regaddress"`
	EducationListener   *EducationListenerDTO  `json:"education_listener,omitempty"`
	PlaceWork           *PlaceWorkDTO          `json:"placework,omitempty"`
	Contractor          *ContractorCreateDTO   `json:"contractor,omitempty"`
}

type CreateListenerDTO struct {
	Listener            ListenerDTO            `json:"listener"`
	Passport            PassportDTO            `json:"passport"`
	RegistrationAddress RegistrationAddressDTO `json:"reg_address"`
	EducationListener   EducationListenerDTO   `json:"education_listener"`
	PlaceWork           PlaceWorkDTO           `json:"placework"`
}
type ListenerLegalEntity struct {
	ID_Listener uuid.UUID `json:"id_listener"`
	FirstName   string    `json:"first_name"`
	SecondName  string    `json:"second_name"`
	MiddleName  string    `json:"middle_name"`
	SNILS       string    `json:"snils"`
}

type DivisionsDTO struct {
	ID_DivisionsEducation uuid.UUID `json:"id_divisionsEducation"`
	Divisions             string    `json:"divisions"`
}

type EducationTypeDTO struct {
	ID_EducationType uuid.UUID `json:"id_educationType"`
	TypeName         string    `json:"typeName"`
}

type ProgramEducationDTO struct {
	ID_ProgramEducation   uuid.UUID `json:"id_program_education"`
	NameProfEducation     string    `json:"name_prof_education"`
	TimeEducation         int       `json:"time_education"`
	IndividualPrice       float32   `json:"individual_price"`
	GroupPrice            float32   `json:"group_price"`
	CampusPrice           float32   `json:"campus_price"`
	ID_EducationType      uuid.UUID `json:"id_education_type"`
	ID_DivisionsEducation uuid.UUID `json:"id_divisions_education"`
}

type EnrollmentListenerDTO struct {
	ID_Listener      uuid.UUID `json:"id_distener"`
	ID_Program       uuid.UUID `json:"id_drogram"`
	StartDate        string    `json:"start_date"`
	EndDate          string    `json:"end_date"`
	CurrentPrice     float32   `json:"current_price"`
	Is_active        bool      `json:"is_active"`
	Group            string    `json:"group"`
	TypeOfRetraining string    `json:"type_of_retraining"`
}

type EnrollmentListenerDetailsDTO struct {
	ID_Listener       uuid.UUID `json:"id_listener"`
	FirstName         string    `json:"first_name"`
	SecondName        string    `json:"second_name"`
	MiddleName        string    `json:"middle_name"`
	NameProfEducation string    `json:"name_prof_education"`
	StartDate         string    `json:"start_date"`
	EndDate           string    `json:"end_date"`
	CurrentPrice      float32   `json:"current_price"`
	Group             string    `json:"group"`
	TypeOfRetraining  string    `json:"type_of_retraining"`
}

type EnrollmentProgramDetailsDTO struct {
	ID_Listener         uuid.UUID `json:"id_listener"`
	ID_ProgramEducation uuid.UUID `json:"id_program_education"`
	NameProfEducation   string    `json:"name_prof_education"`
	TimeEducation       int       `json:"time_education"`
	IndividualPrice     float32   `json:"individual_price"`
	GroupPrice          float32   `json:"group_price"`
	CampusPrice         float32   `json:"campus_price"`
	EducationType       string    `json:"education_type"`
	DivisionEducation   string    `json:"division_education"`
	StartDate           string    `json:"start_date"`
	EndDate             string    `json:"end_date"`
	CurrentPrice        float32   `json:"current_price"`
	Group               string    `json:"group"`
	TypeOfRetraining    string    `json:"type_of_retraining"`
}

type ProgramEducationToCardDTO struct {
	NameProfEducation string `json:"name_prof_education"`
	TimeEducation     int    `json:"time_education"`
	DivisionEducation string `json:"divisions_education"`
	EducationType     string `json:"education_type"`
}

type PassportCardDTO struct {
	PlaceBirth    string `json:"place_birth"`
	Citizenship   string `json:"citizenship"`
	Gender        string `json:"gender"`
	Seria         string `json:"seria"`
	Number        string `json:"number"`
	PassportGiven string `json:"passport_given"`
	DateGiven     string `json:"date_given"`
	Code          string `json:"code"`
}

type RegistrationAddressCardDTO struct {
	MailIndex string `json:"mail_index"`
	Region    string `json:"region"`
	City      string `json:"city"`
	Street    string `json:"street"`
	House     string `json:"house"`
	Building  string `json:"building"`
	Apartment string `json:"apartment"`
}

type EducationListenerCardDTO struct {
	DiplomSeria            string `json:"diplom_seria"`
	DiplomNumber           string `json:"diplom_number"`
	DateGiven              string `json:"date_given"`
	City                   string `json:"city"`
	Region                 string `json:"region"`
	EducationalInstitution string `json:"educational_institution"`
	Speciality             string `json:"speciality"`
	LevelEducation         string `json:"level_education"`
}

type PersonalCardInfoDTO struct {
	Listener            ListenerDTO                `json:"listener"`
	Passport            PassportCardDTO            `json:"passport"`
	RegistrationAddress RegistrationAddressCardDTO `json:"registration_address" `
	EducationListener   EducationListenerCardDTO   `json:"education"`
	PlaceWork           PlaceWorkDTO               `json:"placeWork"`
	ProgramEducation    ProgramEducationToCardDTO  `json:"program_education"`
}

type ListenerFIODTO struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	MiddleName string `json:"middle_name"`
}

type ProgramEndingSoonDTO struct {
	NameProfEducation string    `json:"name_prof_education"`
	EndDate           time.Time `json:"end_date"`
	TotalListeners    int       `json:"total_listeners"`
}

type ContractorDTO struct {
	ID_Contractor uuid.UUID `json:"id_contractor"`
	FirstName     string    `json:"first_name"`
	SecondName    string    `json:"second_name"`
	MiddleName    string    `json:"middle_name"`
	Contact_phone string    `json:"contact_phone"`
	Email         string    `json:"email"`
}

type ContractorCreateDTO struct {
	Contractor ContractorDTO          `json:"contractor"`
	Passport   PassportDTO            `json:"passport"`
	RegAddress RegistrationAddressDTO `json:"reg_address"`
}

type LegalEntityDTO struct {
	ID_Legalentity uuid.UUID `json:"id_legalentity"`
	NameCompany    string    `json:"name_company"`
	Inn            string    `json:"inn"`
	Kpp            string    `json:"kpp"`
	Ogrn           string    `json:"ogrn"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
	FirstName      string    `json:"first_name"`
	SecondName     string    `json:"second_name"`
	MiddleName     string    `json:"middle_name"`
	ID_RegAddress  uuid.UUID `json:"id_regaddress"`
}

type LegalEntityFullDTO struct {
	LegalEntity LegalEntityDTO         `json:"legal_entity"`
	RegAddress  RegistrationAddressDTO `json:"reg_address"`
}

type ExecutorDTO struct {
	ID_Executor uuid.UUID `json:"id_executor"`
	Status      string    `json:"status"`
	FirstName   string    `json:"first_name"`
	SecondName  string    `json:"second_name"`
	MiddleName  string    `json:"middle_name"`
}

type UserDashBoardDTO struct {
	TotalListener     int                    `json:"total_listener"`
	TotalProgram      int                    `json:"total_program"`
	ActiveEnrollments int                    `json:"active_enrollments"`
	ProgramEndingSoon []ProgramEndingSoonDTO `json:"program_ending_soon"`
}

type ReportPeriodDTO struct {
	FullName          string    `json:"full_name"`
	ProgramName       string    `json:"program_name"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	Payment           float32   `json:"payment"`
	TotalIncomePeriod float32   `json:"total_income_period"`
}

type MostExpensiveReportDTO struct {
	ProgramName    string  `json:"program_name"`
	TotalListeners int     `json:"total_listeners"`
	TotalRevenue   float32 `json:"totl_revenue"`
}

type PgStatDTO struct {
	Pid             int       `json:"pid"`
	ApplicationName string    `json:"application_name"`
	ClientAddr      string    `json:"client_addr"`
	ClientPort      int       `json:"client_port"`
	State           string    `json:"state"`
	QueryStart      time.Time `json:"query_start"`
}

type AdminStatDTO struct {
	DatabaseName     string `json:"database_name"`
	Totalsize        string `json:"total_size"`
	ActiveConnection int    `json:"active_connection"`
	CommittedTx      int    `json:"committed_tx"`
	RolledbackTx     int    `json:"rolledback_tx"`
	DiskBlockRead    int    `json:"disk_block_read"`
	BufferHits       int    `json:"buffer_hits"`
}

type AdminDashboardDTO struct {
	PgStat    []PgStatDTO  `json:"pg_stat"`
	AdminStat AdminStatDTO `json:"admin_stat"`
	Role      []RoleDTO    `json:"role"`
}

type RoleDTO struct {
	ID   uuid.UUID `json:"id"`
	Role string    `json:"role"`
}
