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
	IndividualPrice       float32
	GroupPrice            float32
	CampusPrice           float32
	ID_EducationType      uuid.UUID
	EducationType         EducationTypes
	ID_DivisionsEducation uuid.UUID
	Division              DivisionsEducation
}

func (ProgramEducation) TableName() string {
	return "programeducation"
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
	ID_Passport          uuid.UUID
	Passport             Passport
	ID_RegAddress        uuid.UUID
	RegistrationAddress  RegistrationAddress
	ID_EducationListener *uuid.UUID
	EducationListener    EducationListener
	ID_PlaceWork         *uuid.UUID
	PlaceWork            PlaceWork
}

func (Listener) TableName() string {
	return "listener"
}

type EnrollmentListener struct {
	ID_Listener         uuid.UUID
	ID_ProgramEducation uuid.UUID
	StartDate           time.Time
	EndDate             time.Time
	CurrentPrice        float32
	Is_active           bool
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
	CurrentPrice      float32
}

type EnrollmentProgramDetails struct {
	ID_Listener         uuid.UUID
	ID_ProgramEducation uuid.UUID
	NameProfEducation   string
	TimeEducation       int
	IndividualPrice     float32
	GroupPrice          float32
	CampusPrice         float32
	EducationType       string
	DivisionEducation   string
	StartDate           time.Time
	EndDate             time.Time
	CurrentPrice        float32
}

type PersonalCardInfo struct {
	FirstName    string
	SecondName   string
	MiddleName   string
	DateOfBirth  string
	SNILS        string
	ContactPhone string
	Email        string

	PlaceBirth    string
	Citizenship   string
	Gender        string
	Seria         string
	Number        string
	PassportGiven string
	DateGiven     string
	Code          string

	MailIndex string
	RegRegion string
	RegCity   string
	RegStreet string
	House     string
	Building  string
	Apartment string

	DiplomSeria            sql.NullString
	DiplomNumber           sql.NullString
	DiplomDateGiven        sql.NullString
	DiplomCity             sql.NullString
	DiplomRegion           sql.NullString
	EducationalInstitution sql.NullString
	Speciality             sql.NullString
	LevelEducation         sql.NullString

	NameCompany        sql.NullString
	JobTitle           sql.NullString
	AllExperience      sql.NullInt32
	JobTitleExpirience sql.NullInt32

	NameProfEducation string
	TimeEducation     int
	DivisionEducation string
	EducationType     string
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
