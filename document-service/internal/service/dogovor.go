package service

import (
	"bytes"
	"context"
	"document-service/internal/domain/dto"
	"fmt"
	"time"

	"github.com/nguyenthenguyen/docx"
)

var dogovorPP_3_path = "./internal/documents/PP-FIZ-3.docx"
var dogovorPP_2_path = "./internal/documents/PP-FIZ-2.docx"
var dogovorPK_3_path = "./internal/documents/PK-FIZ-3.docx"
var dogovorPK_2_path = "./internal/documents/PK-FIZ-2.docx"
var dogovorDO_3_path = "./internal/documents/DO-FIZ-3.docx"

const (
	PP_3_FIZ = "pp_3_fiz"
	PP_2_FIZ = "pp_2_fiz"
	PK_3_FIZ = "pk_3_fiz"
	PK_2_FIZ = "pk_2_fiz"
	DO_3_FIZ = "do_3_fiz"
)

type DogovorService struct {
	s3client S3ClientInterface
}

func NewDogovorService(s3client S3ClientInterface) *DogovorService {
	return &DogovorService{s3client}
}

func (s *DogovorService) CreateDogovor(dogovor *dto.DogovorDTO, dogovorType string) error {

	var doc *docx.Docx
	var err error

	switch dogovorType {
	case PP_3_FIZ:
		r, err := docx.ReadDocxFile(dogovorPP_3_path)
		if err != nil {
			return err
		}
		defer r.Close()

		doc := r.Editable()

		replacePP3FIZ(doc, dogovor)
	case PP_2_FIZ:
		r, err := docx.ReadDocxFile(dogovorPP_2_path)
		if err != nil {
			return err
		}
		defer r.Close()

		doc := r.Editable()

		replacePP2FIZ(doc, dogovor)
	case PK_3_FIZ:
		r, err := docx.ReadDocxFile(dogovorPK_3_path)
		if err != nil {
			return err
		}
		defer r.Close()

		doc := r.Editable()

		replacePK3FIZ(doc, dogovor)
	case PK_2_FIZ:
		r, err := docx.ReadDocxFile(dogovorPK_2_path)
		if err != nil {
			return err
		}
		defer r.Close()

		doc := r.Editable()

		replacePK2FIZ(doc, dogovor)
	case DO_3_FIZ:
		r, err := docx.ReadDocxFile(dogovorDO_3_path)
		if err != nil {
			return err
		}
		defer r.Close()

		doc := r.Editable()

		replaceDO3FIZ(doc, dogovor)
	}

	var buffer bytes.Buffer
	err = doc.Write(&buffer)
	if err != nil {
		return err
	}

	uniqueParam := dogovor.ListenerData.SNILS
	nameFile := "Личное-дело-" + dogovor.ProgramEducation.NameProfEducation + "_" + uniqueParam + ".docx"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = s.s3client.SendDocumentToS3(ctx, buffer, nameFile)
	if err != nil {
		return err
	}

	return nil
}

func replacePP3FIZ(doc *docx.Docx, model *dto.DogovorDTO) {
	doc.Replace("STATUS", model.Executor.Status, -1)

	executorFio := model.Executor.ExecutorSurname + " " + model.Executor.ExecutorName + " " + model.Executor.ExecutorMiddlename
	doc.Replace("EXECUTORFIO", executorFio, -1)

	listenerFio := model.ListenerData.SecondName + " " + model.ListenerData.FirstName + " " + model.ListenerData.MiddleName
	doc.Replace("LISTENERFIO", listenerFio, -1)

	contractorFio := model.Contractor.SecondName + " " + model.Contractor.FirstName + " " + model.Contractor.MiddleName
	doc.Replace("CONTRACTORFIO", contractorFio, -1)

	doc.Replace("NAMEPROFEDUCATION", model.ProgramEducation.NameProfEducation, -1)
	doc.Replace("SINCE", model.Enrollment.StartDate.Format("02.01.2006"), -1)
	doc.Replace("FOR", model.Enrollment.StartDate.Format("02.01.2006"), -1)

	doc.Replace("OPTIOND", model.OptionDocument, -1)
	doc.Replace("OPTIONH", model.OptionNagruzka, -1)

	doc.Replace("PRICE", fmt.Sprintf("%.2f", model.Enrollment.CurrentPrice), -1)

	doc.Replace("OPTIONPRICE", model.OptionPrice, -1)

	dob, err := time.Parse(time.RFC3339, model.ListenerData.DateOfBirth)
	if err == nil {
		doc.Replace("DATEBIRTH", fmt.Sprintf("%02d.%02d.%02d", dob.Day(), dob.Month(), dob.Year()), -1)
	}

	doc.Replace("SNILS", model.ListenerData.SNILS, -1)
	doc.Replace("PHONEL", model.ListenerData.ContactPhone, -1)
	doc.Replace("EMAILL", model.ListenerData.Email, -1)

	doc.Replace("SERIAE", model.Contractor.Passport.Seria, -1)
	doc.Replace("NUMBERE", model.Contractor.Passport.Number, -1)
	doc.Replace("GIVENE", model.Contractor.Passport.PassportGiven, -1)

	doc.Replace("CITYE", model.Contractor.RegistrationAddress.City, -1)
	doc.Replace("STREETE", model.Contractor.RegistrationAddress.Street, -1)
	doc.Replace("HOUSE", model.Contractor.RegistrationAddress.House, -1)
	doc.Replace("BUILDINGE", model.Contractor.RegistrationAddress.Building, -1)
	doc.Replace("APARTMENTE", model.Contractor.RegistrationAddress.Apartment, -1)
	doc.Replace("PHONEE", model.Contractor.Contact_phone, -1)
	doc.Replace("EMAILE", model.Contractor.Email, -1)

}

func replacePP2FIZ(doc *docx.Docx, model *dto.DogovorDTO) {
	doc.Replace("STATUS", model.Executor.Status, -1)

	executorFio := model.Executor.ExecutorSurname + " " + model.Executor.ExecutorName + " " + model.Executor.ExecutorMiddlename
	doc.Replace("EXECUTORFIO", executorFio, -1)

	listenerFio := model.ListenerData.SecondName + " " + model.ListenerData.FirstName + " " + model.ListenerData.MiddleName
	doc.Replace("LISTENERFIO", listenerFio, -1)

	doc.Replace("NAMEPROFEDUCATION", model.ProgramEducation.NameProfEducation, -1)
	doc.Replace("SINCE", model.Enrollment.StartDate.Format("02.01.2006"), -1)
	doc.Replace("FOR", model.Enrollment.StartDate.Format("02.01.2006"), -1)

	doc.Replace("OPTIOND", model.OptionDocument, -1)
	doc.Replace("OPTIONH", model.OptionNagruzka, -1)

	doc.Replace("PRICE", fmt.Sprintf("%.2f", model.Enrollment.CurrentPrice), -1)

	doc.Replace("OPTIONPRICE", model.OptionPrice, -1)

	doc.Replace("SERIAE", model.Passport.Seria, -1)
	doc.Replace("NUMBERE", model.Passport.Number, -1)
	doc.Replace("GIVENE", model.Passport.PassportGiven, -1)

	doc.Replace("SNILS", model.ListenerData.SNILS, -1)

	doc.Replace("PHONEL", model.Contractor.Contact_phone, -1)
	doc.Replace("EMAILL", model.Contractor.Email, -1)
}

func replacePK3FIZ(doc *docx.Docx, model *dto.DogovorDTO) {
	doc.Replace("STATUS", model.Executor.Status, -1)

	executorFio := model.Executor.ExecutorSurname + " " + model.Executor.ExecutorName + " " + model.Executor.ExecutorMiddlename
	doc.Replace("EXECUTORFIO", executorFio, -1)

	listenerFio := model.ListenerData.SecondName + " " + model.ListenerData.FirstName + " " + model.ListenerData.MiddleName
	doc.Replace("LISTENERFIO", listenerFio, -1)

	contractorFio := model.Contractor.SecondName + " " + model.Contractor.FirstName + " " + model.Contractor.MiddleName
	doc.Replace("CONTRACTORFIO", contractorFio, -1)

	doc.Replace("NAMEPROFEDUCATION", model.ProgramEducation.NameProfEducation, -1)
	doc.Replace("SINCE", model.Enrollment.StartDate.Format("02.01.2006"), -1)
	doc.Replace("FOR", model.Enrollment.StartDate.Format("02.01.2006"), -1)

	doc.Replace("OPTIOND", model.OptionDocument, -1)
	doc.Replace("OPTIONH", model.OptionNagruzka, -1)

	doc.Replace("PRICE", fmt.Sprintf("%.2f", model.Enrollment.CurrentPrice), -1)

	doc.Replace("OPTIONPRICE", model.OptionPrice, -1)

	doc.Replace("SNILS", model.ListenerData.SNILS, -1)
	doc.Replace("PHONEL", model.ListenerData.ContactPhone, -1)
	doc.Replace("EMAILL", model.ListenerData.Email, -1)

	doc.Replace("SERIAE", model.Contractor.Passport.Seria, -1)
	doc.Replace("NUMBERE", model.Contractor.Passport.Number, -1)
	doc.Replace("GIVENE", model.Contractor.Passport.PassportGiven, -1)

	doc.Replace("CITYE", model.Contractor.RegistrationAddress.City, -1)
	doc.Replace("STREETE", model.Contractor.RegistrationAddress.Street, -1)
	doc.Replace("HOUSE", model.Contractor.RegistrationAddress.House, -1)
	doc.Replace("BUILDINGE", model.Contractor.RegistrationAddress.Building, -1)
	doc.Replace("APARTMENTE", model.Contractor.RegistrationAddress.Apartment, -1)
	doc.Replace("PHONEE", model.Contractor.Contact_phone, -1)
	doc.Replace("EMAILE", model.Contractor.Email, -1)

}

func replacePK2FIZ(doc *docx.Docx, model *dto.DogovorDTO) {
	doc.Replace("STATUS", model.Executor.Status, -1)

	executorFio := model.Executor.ExecutorSurname + " " + model.Executor.ExecutorName + " " + model.Executor.ExecutorMiddlename
	doc.Replace("EXECUTORFIO", executorFio, -1)

	listenerFio := model.ListenerData.SecondName + " " + model.ListenerData.FirstName + " " + model.ListenerData.MiddleName
	doc.Replace("LISTENERFIO", listenerFio, -1)

	doc.Replace("NAMEPROFEDUCATION", model.ProgramEducation.NameProfEducation, -1)
	doc.Replace("SINCE", model.Enrollment.StartDate.Format("02.01.2006"), -1)
	doc.Replace("FOR", model.Enrollment.StartDate.Format("02.01.2006"), -1)

	doc.Replace("OPTIOND", model.OptionDocument, -1)
	doc.Replace("OPTIONH", model.OptionNagruzka, -1)

	doc.Replace("PRICE", fmt.Sprintf("%.2f", model.Enrollment.CurrentPrice), -1)

	doc.Replace("OPTIONPRICE", model.OptionPrice, -1)

	doc.Replace("SERIAE", model.Passport.Seria, -1)
	doc.Replace("NUMBERE", model.Passport.Number, -1)
	doc.Replace("GIVENE", model.Passport.PassportGiven, -1)

	doc.Replace("SNILS", model.ListenerData.SNILS, -1)

	doc.Replace("PHONEL", model.Contractor.Contact_phone, -1)
	doc.Replace("EMAILL", model.Contractor.Email, -1)
}

func replaceDO3FIZ(doc *docx.Docx, model *dto.DogovorDTO) {
	doc.Replace("STATUS", model.Executor.Status, -1)

	executorFio := model.Executor.ExecutorSurname + " " + model.Executor.ExecutorName + " " + model.Executor.ExecutorMiddlename
	doc.Replace("EXECUTORFIO", executorFio, -1)

	listenerFio := model.ListenerData.SecondName + " " + model.ListenerData.FirstName + " " + model.ListenerData.MiddleName
	doc.Replace("LISTENERFIO", listenerFio, -1)

	contractorFio := model.Contractor.SecondName + " " + model.Contractor.FirstName + " " + model.Contractor.MiddleName
	doc.Replace("CONTRACTORFIO", contractorFio, -1)

	doc.Replace("NAMEPROFEDUCATION", model.ProgramEducation.NameProfEducation, -1)
	doc.Replace("SINCE", model.Enrollment.StartDate.Format("02.01.2006"), -1)
	doc.Replace("FOR", model.Enrollment.StartDate.Format("02.01.2006"), -1)

	doc.Replace("OPTIOND", model.OptionDocument, -1)
	doc.Replace("OPTIONH", model.OptionNagruzka, -1)

	doc.Replace("PRICE", fmt.Sprintf("%.2f", model.Enrollment.CurrentPrice), -1)

	doc.Replace("OPTIONPRICE", model.OptionPrice, -1)

	doc.Replace("SNILS", model.ListenerData.SNILS, -1)
	doc.Replace("PHONEL", model.ListenerData.ContactPhone, -1)
	doc.Replace("EMAILL", model.ListenerData.Email, -1)

	doc.Replace("SERIAE", model.Contractor.Passport.Seria, -1)
	doc.Replace("NUMBERE", model.Contractor.Passport.Number, -1)
	doc.Replace("GIVENE", model.Contractor.Passport.PassportGiven, -1)

	doc.Replace("CITYE", model.Contractor.RegistrationAddress.City, -1)
	doc.Replace("STREETE", model.Contractor.RegistrationAddress.Street, -1)
	doc.Replace("HOUSE", model.Contractor.RegistrationAddress.House, -1)
	doc.Replace("BUILDINGE", model.Contractor.RegistrationAddress.Building, -1)
	doc.Replace("APARTMENTE", model.Contractor.RegistrationAddress.Apartment, -1)
	doc.Replace("PHONEE", model.Contractor.Contact_phone, -1)
	doc.Replace("EMAILE", model.Contractor.Email, -1)
}
