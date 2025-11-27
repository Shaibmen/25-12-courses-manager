package models

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
	NameProfEducation string `json:"name_prof_education" validate:"required,max=100"`
	TimeEducation     int    `json:"time_education" validate:"required"`
	DivisionEducation string `json:"divisions_education" validate:"required"`
	EducationType     string `json:"education_type" validate:"required"`
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
	ProgramEducation    ProgramEducationRequest  `json:"program_education" validate:"required"`
}

// type ContractorRequest struct {
// 	Passport            PassportRequest
// 	RegistrationAddress RegAddressRequest
// 	FirstName           string
// 	SecondName          string
// 	MiddleName          string
// 	Contact_phone       string
// 	Email               string
// }

//	type Zayavlenie struct {
//		ProgramEducation ProgramEducationRequest
//		Listener         ListenerRequest
//		Contractor       ContractorRequest
//		Variant          int
//	}
type EnrollmentListenerRequest struct {
	StartDate    string  `json:"start_date"`
	EndDate      string  `json:"end_date"`
	CurrentPrice float32 `json:"current_price"`
}

type ContractorRequest struct {
	Passport            PassportRequest   `json:"passport"`
	RegistrationAddress RegAddressRequest `json:"registration_address"`
	FirstName           string            `json:"first_name"`
	SecondName          string            `json:"second_name"`
	MiddleName          string            `json:"middle_name"`
	Contact_phone       string            `json:"contact_phone"`
	Email               string            `json:"email"`
}

type ExecutorRequest struct {
	Status             string `json:"status"`
	ExecutorName       string `json:"executor_name"`
	ExecutorSurname    string `json:"executor_surname"`
	ExecutorMiddlename string `json:"executor_middlename"`
}

type ZayavlenieRequest struct {
	ProgramEducation ProgramEducationRequest `json:"program_education"`
	Listener         ListenerRequest         `json:"listener"`
	Contractor       ContractorRequest       `json:"contractor"`
	Executor         ExecutorRequest         `json:"executor"`
	Passport         PassportRequest         `json:"passport"`
	Registration     RegAddressRequest       `json:"reg_address"`
	Variant          int                     `json:"variant"`
}

type DogovorRequest struct {
	ProgramEducation ProgramEducationRequest   `json:"program_education"`
	ListenerData     ListenerRequest           `json:"listener"`
	Contractor       ContractorRequest         `json:"contractor"`
	Executor         ExecutorRequest           `json:"executor"`
	Passport         PassportRequest           `json:"passport"`
	Registration     RegAddressRequest         `json:"reg_address"`
	Enrollment       EnrollmentListenerRequest `json:"enrollment_listener"`
	OptionNagruzka   string                    `json:"opion_nagruz"`
	OptionDocument   string                    `json:"opt_document"`
	OptionPrice      string                    `json:"opt_price"`
	Price            float32                   `json:"price"`
}
