package service

import (
	"document-service/internal/domain/dto"
	"fmt"
	"time"

	"github.com/nguyenthenguyen/docx"
)

func replacePP3FIZ(doc *docx.Docx, model *dto.DogovorDTO) {

	if model.Executor.Doverennost == "Устав" {
		doc.Replace("YST", "Устав", -1)
	} else {
		textDoverennost := "Доверенности " + model.Executor.Doverennost
		doc.Replace("YST", textDoverennost, -1)
	}

	doc.Replace("STATUS", model.Executor.Status, -1)

	executorFio := model.Executor.ExecutorSurname + " " + model.Executor.ExecutorName + " " + model.Executor.ExecutorMiddlename
	doc.Replace("EXECUTORFIO", executorFio, -1)

	listenerFio := model.ListenerData.SecondName + " " + model.ListenerData.FirstName + " " + model.ListenerData.MiddleName
	doc.Replace("LISTENERFIO", listenerFio, -1)

	contractorFio := model.Contractor.SecondName + " " + model.Contractor.FirstName + " " + model.Contractor.MiddleName
	doc.Replace("CONTRACTORFIO", contractorFio, -1)

	doc.Replace("NAMEPROFEDUCATION", model.ProgramEducation.NameProfEducation, -1)
	doc.Replace("SINCE", model.Enrollment.StartDate.Format("02.01.2006"), -1)
	doc.Replace("FOR", model.Enrollment.EndDate.Format("02.01.2006"), -1)

	dob, err := time.Parse(time.RFC3339, model.ListenerData.DateOfBirth)
	if err == nil {
		doc.Replace("DATEBIRTH", fmt.Sprintf("%02d.%02d.%02d", dob.Day(), dob.Month(), dob.Year()), -1)
	}

	doc.Replace("SNILS", model.ListenerData.SNILS, -1)
	doc.Replace("PHONEL", model.ListenerData.ContactPhone, -1)
	doc.Replace("EMAILL", model.ListenerData.Email, -1)

	seriaNumberGiven := fmt.Sprintf("%s %s, %s", model.Contractor.Passport.Seria, model.Contractor.Passport.Number, model.Contractor.Passport.PassportGiven)

	doc.Replace("SERIAE NUMBERE, GIVENE", seriaNumberGiven, -1)

	cityStreetHouseBuildingApartment := fmt.Sprintf("%s, %s, %s, %s, %s. ",
		model.Contractor.RegistrationAddress.City,
		model.Contractor.RegistrationAddress.Street,
		model.Contractor.RegistrationAddress.House,
		model.Contractor.RegistrationAddress.Building,
		model.Contractor.RegistrationAddress.Apartment)

	doc.Replace("CITYE, STREETE, HOUSEE, BUILDINGE, APARTMENTE. ", cityStreetHouseBuildingApartment, -1)

	doc.Replace("PHONEE", model.Contractor.Contact_phone, -1)
	doc.Replace("EMAILE ", model.Contractor.Email+" ", -1)
}

func replacePP2FIZ(doc *docx.Docx, model *dto.DogovorDTO) {

	if model.Executor.Doverennost == "Устав" {
		doc.Replace("YST", "Устав", -1)
	} else {
		textDoverennost := "Доверенности " + model.Executor.Doverennost
		doc.Replace("YST", textDoverennost, -1)
	}

	doc.Replace("STATUS", model.Executor.Status, -1)

	executorFio := model.Executor.ExecutorSurname + " " + model.Executor.ExecutorName + " " + model.Executor.ExecutorMiddlename
	doc.Replace("EXECUTORFIO", executorFio, -1)

	listenerFio := model.ListenerData.SecondName + " " + model.ListenerData.FirstName + " " + model.ListenerData.MiddleName
	doc.Replace("LISTENERFIO", listenerFio, -1)

	doc.Replace("NAMEPROFEDUCATION", model.ProgramEducation.NameProfEducation, -1)
	doc.Replace("SINCE", model.Enrollment.StartDate.Format("02.01.2006"), -1)
	doc.Replace("FOR", model.Enrollment.EndDate.Format("02.01.2006"), -1)

	seriaNumberGiven := fmt.Sprintf("%s %s, %s", model.Passport.Seria, model.Passport.Number, model.Passport.PassportGiven)
	doc.Replace("SERIAE NUMBERE, GIVENE", seriaNumberGiven, -1)

	cityStreetHouseBuildingApartment := fmt.Sprintf("%s, %s, %s, %s, %s.",
		model.Registration.City,
		model.Registration.Street,
		model.Registration.House,
		model.Registration.Building,
		model.Registration.Apartment)

	doc.Replace("CITYE, STREETE, HOUSEE, BUILDINGE, APARTMENTE.", cityStreetHouseBuildingApartment, -1)

	doc.Replace("SNILS", model.ListenerData.SNILS, -1)

	doc.Replace("PHONEL", model.ListenerData.ContactPhone, -1)
	doc.Replace("EMAILL", model.ListenerData.Email, -1)

	dob, err := time.Parse(time.RFC3339, model.ListenerData.DateOfBirth)
	if err == nil {
		doc.Replace("DATEBIRTH", fmt.Sprintf("%02d.%02d.%02d", dob.Day(), dob.Month(), dob.Year()), -1)
	}
}

func replacePK3FIZ(doc *docx.Docx, model *dto.DogovorDTO) {

	if model.Executor.Doverennost == "Устав" {
		doc.Replace("YST", "Устав", -1)
	} else {
		textDoverennost := "Доверенности " + model.Executor.Doverennost
		doc.Replace("YST", textDoverennost, -1)
	}

	doc.Replace("STATUS", model.Executor.Status, -1)

	executorFio := model.Executor.ExecutorSurname + " " + model.Executor.ExecutorName + " " + model.Executor.ExecutorMiddlename
	doc.Replace("EXECUTORFIO", executorFio, -1)

	listenerFio := model.ListenerData.SecondName + " " + model.ListenerData.FirstName + " " + model.ListenerData.MiddleName
	doc.Replace("LISTENERFIO", listenerFio, -1)

	contractorFio := model.Contractor.SecondName + " " + model.Contractor.FirstName + " " + model.Contractor.MiddleName
	doc.Replace("CONTRACTORFIO", contractorFio, -1)

	doc.Replace("NAMEPROFEDUCATION", model.ProgramEducation.NameProfEducation, -1)
	doc.Replace("SINCE", model.Enrollment.StartDate.Format("02.01.2006"), -1)
	doc.Replace("FOR", model.Enrollment.EndDate.Format("02.01.2006"), -1)

	doc.Replace("SNILS", model.ListenerData.SNILS, -1)
	doc.Replace("PHONEL", model.ListenerData.ContactPhone, -1)
	doc.Replace("EMAILL", model.ListenerData.Email, -1)

	doc.Replace("SERIAE", model.Contractor.Passport.Seria, -1)
	doc.Replace("NUMBERE", model.Contractor.Passport.Number, -1)
	doc.Replace("GIVENE", model.Contractor.Passport.PassportGiven, -1)

	cityTadada := fmt.Sprintf("%s, %s, %s, %s, %s",
		model.Contractor.RegistrationAddress.City,
		model.Contractor.RegistrationAddress.Street,
		model.Contractor.RegistrationAddress.House,
		model.Contractor.RegistrationAddress.Building,
		model.Contractor.RegistrationAddress.Apartment)

	doc.Replace("CITYE, STREETE, HOUSEE, BUILDINGE, APARTMENTE", cityTadada, -1)

	doc.Replace("PHONEE", model.Contractor.Contact_phone, -1)
	doc.Replace("EMAILE", model.Contractor.Email, -1)

	dob, err := time.Parse(time.RFC3339, model.ListenerData.DateOfBirth)
	if err == nil {
		doc.Replace("DATEBIRTH", fmt.Sprintf("%02d.%02d.%02d", dob.Day(), dob.Month(), dob.Year()), -1)
	}

}

func replacePK2FIZ(doc *docx.Docx, model *dto.DogovorDTO) {

	if model.Executor.Doverennost == "Устав" {
		doc.Replace("YST", "Устав", -1)
	} else {
		textDoverennost := "Доверенности " + model.Executor.Doverennost
		doc.Replace("YST", textDoverennost, -1)
	}

	doc.Replace("STATUS", model.Executor.Status, -1)

	executorFio := model.Executor.ExecutorSurname + " " + model.Executor.ExecutorName + " " + model.Executor.ExecutorMiddlename
	doc.Replace("EXECUTORFIO", executorFio, -1)

	listenerFio := model.ListenerData.SecondName + " " + model.ListenerData.FirstName + " " + model.ListenerData.MiddleName
	doc.Replace("LISTENERFIO", listenerFio, -1)

	doc.Replace("NAMEPROFEDUCATION", model.ProgramEducation.NameProfEducation, -1)
	doc.Replace("SINCE", model.Enrollment.StartDate.Format("02.01.2006"), -1)
	doc.Replace("FOR", model.Enrollment.EndDate.Format("02.01.2006"), -1)

	seriaNumberGiven := fmt.Sprintf("%s %s, %s", model.Passport.Seria, model.Passport.Number, model.Passport.PassportGiven)
	doc.Replace("SERIAE NUMBERE, GIVENE", seriaNumberGiven, -1)

	doc.Replace("SNILS", model.ListenerData.SNILS, -1)

	cityStreetHouseBuildingApartment := fmt.Sprintf("%s, %s, %s, %s, %s.",
		model.Registration.City,
		model.Registration.Street,
		model.Registration.House,
		model.Registration.Building,
		model.Registration.Apartment)

	doc.Replace("CITYE, STREETE, HOUSEE, BUILDINGE, APARTMENTE.", cityStreetHouseBuildingApartment, -1)

	doc.Replace("PHONEL", model.ListenerData.ContactPhone, -1)
	doc.Replace("EMAILL", model.ListenerData.Email, -1)

	dob, err := time.Parse(time.RFC3339, model.ListenerData.DateOfBirth)
	if err == nil {
		doc.Replace("DATEBIRTH", fmt.Sprintf("%02d.%02d.%02d", dob.Day(), dob.Month(), dob.Year()), -1)
	}
}

func replaceDO3FIZ(doc *docx.Docx, model *dto.DogovorDTO) {

	if model.Executor.Doverennost == "Устав" {
		doc.Replace("YST", "Устав", -1)
	} else {
		textDoverennost := "Доверенности " + model.Executor.Doverennost
		doc.Replace("YST", textDoverennost, -1)
	}

	doc.Replace("STATUS", model.Executor.Status, -1)

	executorFio := model.Executor.ExecutorSurname + " " + model.Executor.ExecutorName + " " + model.Executor.ExecutorMiddlename
	doc.Replace("EXECUTORFIO", executorFio, -1)

	listenerFio := model.ListenerData.SecondName + " " + model.ListenerData.FirstName + " " + model.ListenerData.MiddleName
	doc.Replace("LISTENER", listenerFio, -1)
	doc.Replace("FIO", "", -1)

	contractorFio := model.Contractor.SecondName + " " + model.Contractor.FirstName + " " + model.Contractor.MiddleName
	doc.Replace("CONTRACTORFIO", contractorFio, -1)

	doc.Replace("NAMEPROFEDUCATION", model.ProgramEducation.NameProfEducation, -1)
	doc.Replace("SIN", model.Enrollment.StartDate.Format("02.01.2006"), -1)
	doc.Replace("FOR", model.Enrollment.EndDate.Format("02.01.2006"), -1)

	dob, err := time.Parse(time.RFC3339, model.ListenerData.DateOfBirth)
	if err == nil {
		doc.Replace("DATEBIRTH", fmt.Sprintf("%02d.%02d.%02d", dob.Day(), dob.Month(), dob.Year()), -1)
	}

	doc.Replace("SNILS", model.ListenerData.SNILS, -1)
	doc.Replace("PHONEL", model.ListenerData.ContactPhone, -1)

	doc.Replace("SERIAE", model.Contractor.Passport.Seria, -1)
	doc.Replace("NUMBERE", model.Contractor.Passport.Number, -1)
	doc.Replace("GIVENE", model.Contractor.Passport.PassportGiven, -1)

	doc.Replace("CITYE", model.Contractor.RegistrationAddress.City, -1)
	doc.Replace("STREETE", model.Contractor.RegistrationAddress.Street, -1)
	doc.Replace("HOUSE", model.Contractor.RegistrationAddress.House, -1)
	doc.Replace("BUILDINGE, APARTMENTE", model.Contractor.RegistrationAddress.Building+", "+model.Contractor.RegistrationAddress.Apartment, -1)
	doc.Replace("PHONEE", model.Contractor.Contact_phone, -1)
	doc.Replace("EMAILE", model.Contractor.Email, -1)

	doc.Replace("EMA", model.ListenerData.Email, -1)
}
