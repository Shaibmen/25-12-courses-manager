package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Username string
	Password string
}

type Passport struct {
	ID_Passport   uuid.UUID
	PlaceBirth    string
	Citizenship   string
	Gender        string
	Seria         int
	Number        int
	PassportGiven string
	DateGiven     time.Time
	Code          string
}

func (Passport) TableName() string {
	return "passport"
}

type RegistrationAddress struct {
	ID_RegAddress uuid.UUID
	MailIndex     int
	Region        string
	City          string
	Street        string
	House         string
	Building      string
	Apartment     string
}

func (RegistrationAddress) TableName() string {
	return "registrationaddress"
}

type LevelEducation struct {
	ID_LevelEducation uuid.UUID
	Education         string
}

func (LevelEducation) TableName() string {
	return "leveleducation"
}

type EducationListener struct {
	ID_EducationListener   uuid.UUID
	DiplomSeria            int
	DiplomNumber           int
	DateGiven              time.Time
	City                   string
	Region                 string
	EducationalInstitution string
	Speciality             string
	LevelEducation         string
}

type PlaceWork struct {
	ID_PlaceWork       uuid.UUID
	NameCompany        string
	JobTitle           string
	AllExperience      int
	JobTitleExperience int
}

func (PlaceWork) TableName() string {
	return "placework"
}

type DivisionsEducation struct {
	ID_DivisionsEducation uuid.UUID
	Divisions             string
}

func (DivisionsEducation) TableName() string {
	return "divisionseducation"
}

type EducationTypes struct {
	ID_EducationType uuid.UUID
	TypeName         string
}

func (EducationTypes) TableName() string {
	return "educationtypes"
}

type ProgramEducation struct {
	ID_ProgramEducation   uuid.UUID
	NameProfEducation     string
	TimeEducation         int
	Price                 float64
	ID_EducationType      uuid.UUID
	EducationType         EducationTypes
	ID_DivisionsEducation uuid.UUID
	Division              DivisionsEducation
}

func (ProgramEducation) TableName() string {
	return "programeducation"
}

type AccurateProgram struct {
	ID_Listener       uuid.UUID `db:"id_listener"`
	NameProfEducation string    `db:"name_prof_education"`
	TimeEducation     int       `db:"time_education"`
	Price             float64   `db:"price"`
	EducationType     string    `db:"educationtype"`
	Division          string    `db:"divisionseducation"`
}

type ProgramToAccurate struct {
	NameProfEducation string  `db:"name_prof_education"`
	TimeEducation     int     `db:"time_education"`
	Price             float64 `db:"price"`
	EducationType     string  `db:"type_name"`
	Division          string  `db:"divisions"`
}

func (AccurateProgram) TableName() string {
	return "accurateprogram"
}

type LegalEntity struct {
	ID_Legalentity      uuid.UUID
	NameCompany         string
	Inn                 string
	Kpp                 string
	Ogrn                string
	Phone               string
	Email               string
	FirstName           string
	SecondName          string
	MiddleName          string
	ID_RegAddress       uuid.UUID
	RegistrationAddress RegistrationAddress
	Status              string
}

func (LegalEntity) TableName() string {
	return "legal_entity"
}

type Contractor struct {
	ID_Contractor       uuid.UUID
	FirstName           string
	SecondName          string
	MiddleName          string
	Contact_phone       string
	Email               string
	ID_Passport         *uuid.UUID
	Passport            Passport
	ID_RegAddress       uuid.UUID
	RegistrationAddress RegistrationAddress
}

func (Contractor) TableName() string {
	return "contractor"
}

type Listener struct {
	ID_Listener          uuid.UUID
	FirstName            string
	SecondName           string
	MiddleName           string
	DateOfBirth          time.Time
	SNILS                string
	ContactPhone         string
	Email                string
	LootingEducation     bool
	ID_Passport          *uuid.UUID
	Passport             Passport
	ID_RegAddress        uuid.UUID
	RegistrationAddress  RegistrationAddress
	ID_EducationListener *uuid.UUID
	EducationListener    EducationListener
	ID_PlaceWork         *uuid.UUID
	PlaceWork            PlaceWork
	ID_Legalentity       *uuid.UUID
	LegalEntity          LegalEntity
	ID_Contractor        *uuid.UUID
	Contractor           Contractor
}

func (Listener) TableName() string {
	return "listener"
}

type ListenerLegalEntity struct {
	ID_Listener uuid.UUID
	FirstName   string
	SecondName  string
	MiddleName  string
	SNILS       string
	DateOfBirth string
	Email       string
}

type Executor struct {
	ID_Executor uuid.UUID
	Status      string `db:"status"`
	FirstName   string `db:"first_name"`
	SecondName  string `db:"second_name"`
	MiddleName  string `db:"middle_name"`
	Doverenost  string `db:"doverenost"`
}

func (Executor) TableName() string {
	return "executor"
}

type EnrollmentListener struct {
	ID_Listener         uuid.UUID
	ID_ProgramEducation uuid.UUID
	StartDate           time.Time
	EndDate             time.Time
	Is_active           bool
	Group               string
	TypeOfRetraining    string
}

func (EnrollmentListener) TableName() string {
	return "enrollmentlistener"
}

type EnrollmentListenerDetails struct {
	ID_Listener       uuid.UUID
	FirstName         string
	SecondName        string
	MiddleName        string
	NameProfEducation string
	StartDate         time.Time
	EndDate           time.Time
	Group             string
	TypeOfRetraining  string
}

type EnrollmentProgramDetails struct {
	ID_Listener         uuid.UUID
	ID_ProgramEducation uuid.UUID
	NameProfEducation   string
	TimeEducation       int
	Price               float64
	EducationType       string
	DivisionEducation   string
	StartDate           time.Time
	EndDate             time.Time
	CurrentPrice        float32
	Group               string
	TypeOfRetraining    string
}

type ListenerFIO struct {
	FirstName  string
	SecondName string
	MiddleName string
}

type ProgramEndingSoon struct {
	NameProfEducation string
	EndDate           time.Time
	TotalListeners    int
}

type UserDashBoard struct {
	TotalListener     int
	TotalProgram      int
	ActiveEnrollments int
	ProgramEndingSoon []ProgramEndingSoon
}

type ReportPeriod struct {
	FullName          string
	ProgramName       string
	StartDate         time.Time
	EndDate           time.Time
	Payment           float32
	TotalIncomePeriod float32
}

type MostExpensiveReport struct {
	ProgramName    string
	TotalListeners int
	TotalRevenue   float32
}

type PgStat struct {
	Pid             int
	ApplicationName string
	ClientAddr      string
	ClientPort      int
	State           string
	QueryStart      time.Time
}

type AdminStat struct {
	DatabaseName     string
	Totalsize        string
	ActiveConnection int
	CommittedTx      int
	RolledbackTx     int
	DiskBlockRead    int
	BufferHits       int
}

type AdminDashboard struct {
	PgStat    []PgStat
	AdminStat AdminStat
	Role      []Role
}

type Role struct {
	ID   uuid.UUID
	Role string
}
type PersonalCardInfo struct {
	FirstName    string `db:"first_name"`
	SecondName   string `db:"second_name"`
	MiddleName   string `db:"middle_name"`
	DateOfBirth  string `db:"date_of_birth"`
	SNILS        string `db:"snils"`
	ContactPhone string `db:"contact_phone"`
	Email        string `db:"email"`

	PlaceBirth    sql.NullString `db:"place_birth"`
	Citizenship   sql.NullString `db:"citizenship"`
	Gender        sql.NullString `db:"gender"`
	Seria         sql.NullString `db:"seria"`
	Number        sql.NullString `db:"number"`
	PassportGiven sql.NullString `db:"passport_given"`
	DateGiven     sql.NullString `db:"date_given"`
	Code          sql.NullString `db:"code"`

	MailIndex string `db:"mail_index"`
	RegRegion string `db:"region"`
	RegCity   string `db:"city"`
	RegStreet string `db:"street"`
	House     string `db:"house"`
	Building  string `db:"building"`
	Apartment string `db:"apartment"`

	DiplomSeria            sql.NullString `db:"diplom_seria"`
	DiplomNumber           sql.NullString `db:"diplom_number"`
	City                   sql.NullString `db:"educ_city"`
	Region                 sql.NullString `db:"educ_region"`
	EducationalInstitution sql.NullString `db:"educational_institution"`
	Speciality             sql.NullString `db:"speciality"`
	LevelEducation         sql.NullString `db:"education"`

	NameCompany        sql.NullString `db:"name_company"`
	JobTitle           sql.NullString `db:"job_title"`
	AllExperience      sql.NullInt32  `db:"all_experience"`
	JobTitleExpirience sql.NullInt32  `db:"job_title_experience"`

	NameProfEducation string `db:"name_prof_education"`
	TimeEducation     int    `db:"time_education"`
	DivisionEducation string `db:"divisions"`
	EducationType     string `db:"type_name"`

	StartDate        string  `db:"start_date"`
	EndDate          string  `db:"end_date"`
	CurrentPrice     float32 `db:"current_price"`
	Is_active        bool    `db:"is_active"`
	Group            string  `db:"group_number"`
	TypeOfRetraining string  `db:"type_of_retraining"`
}

type ContractorDoc struct {
	FirstName    sql.NullString `db:"first_name"`
	SecondName   sql.NullString `db:"second_name"`
	MiddleName   sql.NullString `db:"middle_name"`
	ContactPhone sql.NullString `db:"contact_phone"`
	Email        sql.NullString `db:"email"`

	PlaceBirth    sql.NullString `db:"place_birth"`
	Citizenship   sql.NullString `db:"citizenship"`
	Gender        sql.NullString `db:"gender"`
	Seria         sql.NullString `db:"seria"`
	Number        sql.NullString `db:"number"`
	PassportGiven sql.NullString `db:"passport_given"`
	DateGiven     sql.NullString `db:"date_given"`
	Code          sql.NullString `db:"code"`

	MailIndex sql.NullString `db:"mail_index"`
	RegRegion sql.NullString `db:"city"`
	RegCity   sql.NullString `db:"region"`
	RegStreet sql.NullString `db:"street"`
	House     sql.NullString `db:"house"`
	Building  sql.NullString `db:"building"`
	Apartment sql.NullString `db:"apartment"`
}

type FullDocument struct {
	PersonalInfo PersonalCardInfo
	Executor     Executor
	Contractor   ContractorDoc
}
