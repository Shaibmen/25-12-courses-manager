package models

type PassportRequest struct {
	PlaceBirth    string `json:"place_birth"`
	Citizenship   string `json:"citizenship"`
	Gender        string `json:"gender"`
	Seria         string `json:"seria"`
	Number        string `json:"number"`
	PassportGiven string `json:"passport_given"`
	Code          string `json:"code"`
	DateGiven     string `json:"date_given"`
}

type RegAddressRequest struct {
	MailIndex string `json:"mail_index"`
	Region    string `json:"region"`
	City      string `json:"city"`
	Street    string `json:"street"`
	House     string `json:"house"`
	Building  string `json:"building"`
	Apartment string `json:"apartment"`
}

type EducationListenerRequest struct {
	DateGiven              string `json:"date_given"`
	DiplomSeria            string `json:"diplom_seria"`
	DiplomNumber           string `json:"diplom_number"`
	City                   string `json:"city"`
	Region                 string `json:"region"`
	EducationalInstitution string `json:"educational_institution"`
	Speciality             string `json:"speciality"`
	LevelEducation         string `json:"level_education"`
}

type PlaceWorkRequest struct {
	NameCompany        string `json:"name_company"`
	JobTitle           string `json:"job_title"`
	AllExperience      int    `json:"all_experience"`
	JobTitleExpirience int    `json:"job_title_expirience"`
}

type ProgramEducationRequest struct {
	NameProfEducation string `json:"name_prof_education"`
	DivisionEducation string `json:"divisions_education"`
	EducationType     string `json:"education_type"`
	TimeEducation     int    `json:"time_education"`
}

type ListenerRequest struct {
	FirstName    string `json:"first_name"`
	SecondName   string `json:"second_name"`
	MiddleName   string `json:"middle_name"`
	SNILS        string `json:"snils"`
	ContactPhone string `json:"contact_phone"`
	Email        string `json:"email"`
	DateOfBirth  string `json:"date_of_birth"`
}

type EnrollmentListenerRequest struct {
	StartDate        string  `json:"start_date"`
	EndDate          string  `json:"end_date"`
	Group            string  `json:"group"`
	TypeOfRetraining string  `json:"type_of_retraining"`
	CurrentPrice     float32 `json:"current_price"`
	Is_active        bool    `json:"is_active"`
}

type FullListenerRequest struct {
	Listener            ListenerRequest           `json:"listener"`
	Passport            PassportRequest           `json:"passport"`
	RegistrationAddress RegAddressRequest         `json:"registration_address"`
	EducationListener   EducationListenerRequest  `json:"education"`
	PlaceWork           PlaceWorkRequest          `json:"placeWork"`
	ProgramEducation    ProgramEducationRequest   `json:"program_education"`
	EnrollmentListener  EnrollmentListenerRequest `json:"enrollment_listener"`
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
	Doverennost 	   string `json:"doverenost"`
}

type ZakazchikRequest struct {
	Listeners   []ListenerRequest `json:"listeners"` 
	Address     RegAddressRequest `json:"reg_address"`
	CompanyName string            `json:"company_name"`
	FIO         string            `json:"zakazchikfio"`
	Status      string            `json:"status"` 
	INN         string            `json:"inn"`
	KPP         string            `json:"kpp"`
	OGRN        string            `json:"ogrn"`
	Phone       string            `json:"phone"`
	Email       string            `json:"email"`
}

type ZayavlenieRequest struct {
	ProgramEducation   ProgramEducationRequest   `json:"program_education"`
	Listener           ListenerRequest           `json:"listener"`
	EnrollmentListener EnrollmentListenerRequest `json:"enrollment_listener"`
	Contractor         ContractorRequest         `json:"contractor"`
	Executor           ExecutorRequest           `json:"executor"`
	Passport           PassportRequest           `json:"passport"`
	Registration       RegAddressRequest         `json:"reg_address"`
	DogovorAgeType     string                    `json:"dogovor_type"`
	Variant            int                       `json:"variant"`
}

type DogovorRequest struct {
	ZakazchikData    ZakazchikRequest          `json:"zakazchik"`
	ProgramEducation ProgramEducationRequest   `json:"program_education"`
	ListenerData     ListenerRequest           `json:"listener" validate:"omitempty"`
	Contractor       ContractorRequest         `json:"contractor"`
	Executor         ExecutorRequest           `json:"executor"`
	Passport         PassportRequest           `json:"passport" validate:"omitempty"`
	Registration     RegAddressRequest         `json:"reg_address" validate:"omitempty"`
	Enrollment       EnrollmentListenerRequest `json:"enrollment_listener"`
	OptionNagruzka   int                       `json:"opion_nagruz"`
	OptionDocument   int                       `json:"opt_document"`
	OptionPrice      string                    `json:"opt_price"`
	DogovorType      string                    `json:"dogovor_type"`
}

/*
когда отправляешь юрика, кидаешь только "dogovor_card": {}
затем кидаешь на каждого пользователя "personal_card": {}, "zayavlenie_card": {}

когда отправляешь просто слушателя, либо слушателя через третье лицо, кидаешь
"personal_card": {}, "zayavlenie_card": {}, БЕЗ ZAKAZCHIK!!!!! "dogovor_card": {}
*/
type FullRequest struct {
	PersonalCardData FullListenerRequest `json:"personal_card"`
	ZayavlenieData   ZayavlenieRequest   `json:"zayavlenie_card"`
	DogovorData      DogovorRequest      `json:"dogovor_card"`
}
