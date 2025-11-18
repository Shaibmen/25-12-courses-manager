package dto

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
	TimeEducation     int
	DivisionEducation string
	EducationType     string
}

type FullListenerDataDTO struct {
	Listener            ListenerDTO
	Passport            PassportDTO
	RegistrationAddress RegistrationAddressDTO
	EducationListener   EducationListenerDTO
	PlaceWork           PlaceWorkDTO
	ProgramEducation    ProgramEducationDTO
}
