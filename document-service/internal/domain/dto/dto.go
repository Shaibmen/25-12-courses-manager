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
	AllExperience      int
	JobTitleExpirience int
	NameCompany        string
	JobTitle           string
}

type ProgramEducationDTO struct {
	TimeEducation     int
	NameProfEducation string
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

type ContractorDTO struct {
	Passport            PassportDTO
	RegistrationAddress RegistrationAddressDTO
	FirstName           string
	SecondName          string
	MiddleName          string
	Contact_phone       string
	Email               string
}

type ZayavlenieDTO struct {
	ProgramEducation ProgramEducationDTO
	Listener         ListenerDTO
	Contractor       ContractorDTO
	Variant          int
}
