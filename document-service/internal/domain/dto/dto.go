package dto

import "time"

type ListenerDTO struct {
	FirstName    string
	SecondName   string
	MiddleName   string
	DateOfBirth  string
	SNILS        string
	ContactPhone string
	Email        string
}

type PassportDTO struct {
	PlaceBirth    string
	Citizenship   string
	Gender        string
	Seria         string
	Number        string
	PassportGiven string
	DateGiven     string
	Code          string
}

type RegistrationAddressDTO struct {
	MailIndex string
	Region    string
	City      string
	Street    string
	House     string
	Building  string
	Apartment string
}

type EducationListenerDTO struct {
	DiplomSeria            string
	DiplomNumber           string
	DateGiven              string
	City                   string
	Region                 string
	EducationalInstitution string
	Speciality             string
	LevelEducation         string
}

type PlaceWorkDTO struct {
	NameCompany        string
	JobTitle           string
	AllExperience      int
	JobTitleExpirience int
}

type ProgramEducationDTO struct {
	NameProfEducation string
	DivisionEducation string
	EducationType     string
	TimeEducation     int
}

type EnrollmentListenerDTO struct {
	StartDate        time.Time
	EndDate          time.Time
	Group            string
	TypeOfRetraining string
	CurrentPrice     float32
	Is_active        bool
}

type FullListenerDataDTO struct {
	Listener            ListenerDTO
	Passport            PassportDTO
	RegistrationAddress RegistrationAddressDTO
	EducationListener   EducationListenerDTO
	PlaceWork           PlaceWorkDTO
	ProgramEducation    ProgramEducationDTO
	EnrollmentListener  EnrollmentListenerDTO
}

type ContractorDTO struct {
	Passport            PassportDTO
	RegistrationAddress RegistrationAddressDTO
	FirstName           string
	SecondName          string
	MiddleName          string
	Contact_phone       string
	Email               string
}

type ExecutorDTO struct {
	Status             string
	ExecutorName       string
	ExecutorSurname    string
	ExecutorMiddlename string
}

type ZayavlenieDTO struct {
	ProgramEducation   ProgramEducationDTO
	ListenerData       ListenerDTO
	Contractor         ContractorDTO
	Executor           ExecutorDTO
	Passport           PassportDTO
	Registration       RegistrationAddressDTO
	EnrollmentListener EnrollmentListenerDTO
	Variant            int
}

type DogovorDTO struct {
	ProgramEducation ProgramEducationDTO
	ListenerData     ListenerDTO
	Contractor       ContractorDTO
	Executor         ExecutorDTO
	Passport         PassportDTO
	Registration     RegistrationAddressDTO
	Enrollment       EnrollmentListenerDTO
	OptionNagruzka   int
	OptionDocument   int
	OptionPrice      int
}
