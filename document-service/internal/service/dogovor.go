package service

import (
	"bytes"
	"context"
	"document-service/internal/domain/dto"
	"errors"
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
	nameFile := "Договор-" + dogovor.ProgramEducation.NameProfEducation + "_" + uniqueParam + ".docx"

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
	doc.Replace("FOR", model.Enrollment.EndDate.Format("02.01.2006"), -1)

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

	seriaNumberGiven := fmt.Sprintf("%s %s, %s", model.Passport.Seria, model.Passport.Number, model.Passport.PassportGiven)

	doc.Replace("SERIAE, NUMBERE, GIVENE", seriaNumberGiven, -1)

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
	doc.Replace("STATUS", model.Executor.Status, -1)

	executorFio := model.Executor.ExecutorSurname + " " + model.Executor.ExecutorName + " " + model.Executor.ExecutorMiddlename
	doc.Replace("EXECUTORFIO", executorFio, -1)

	listenerFio := model.ListenerData.SecondName + " " + model.ListenerData.FirstName + " " + model.ListenerData.MiddleName
	doc.Replace("LISTENERFIO", listenerFio, -1)

	doc.Replace("NAMEPROFEDUCATION", model.ProgramEducation.NameProfEducation, -1)
	doc.Replace("SINCE", model.Enrollment.StartDate.Format("02.01.2006"), -1)
	doc.Replace("FOR", model.Enrollment.EndDate.Format("02.01.2006"), -1)

	doc.Replace("OPTIOND", model.OptionDocument, -1)
	doc.Replace("OPTIONH", model.OptionNagruzka, -1)

	doc.Replace("PRICE", fmt.Sprintf("%.2f", model.Enrollment.CurrentPrice), -1)

	doc.Replace("OPTIONPRICE", model.OptionPrice, -1)

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
	doc.Replace("STATUS", model.Executor.Status, -1)

	executorFio := model.Executor.ExecutorSurname + " " + model.Executor.ExecutorName + " " + model.Executor.ExecutorMiddlename
	doc.Replace("EXECUTORFIO", executorFio, -1)

	listenerFio := model.ListenerData.SecondName + " " + model.ListenerData.FirstName + " " + model.ListenerData.MiddleName
	doc.Replace("LISTENERFIO", listenerFio, -1)

	doc.Replace("NAMEPROFEDUCATION", model.ProgramEducation.NameProfEducation, -1)
	doc.Replace("SINCE", model.Enrollment.StartDate.Format("02.01.2006"), -1)
	doc.Replace("FOR", model.Enrollment.EndDate.Format("02.01.2006"), -1)

	doc.Replace("OPTIOND", model.OptionDocument, -1)
	doc.Replace("OPTIONH", model.OptionNagruzka, -1)

	doc.Replace("PRICE", fmt.Sprintf("%.2f", model.Enrollment.CurrentPrice), -1)

	doc.Replace("OPTIONPRICE", model.OptionPrice, -1)

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
	doc.Replace("STATUS", model.Executor.Status, -1)

	executorFio := model.Executor.ExecutorSurname + " " + model.Executor.ExecutorName + " " + model.Executor.ExecutorMiddlename
	doc.Replace("EXECUTORFIO", executorFio, -1)

	listenerFio := model.ListenerData.SecondName + " " + model.ListenerData.FirstName + " " + model.ListenerData.MiddleName
	doc.Replace("LISTENERFIO", listenerFio, -1)

	contractorFio := model.Contractor.SecondName + " " + model.Contractor.FirstName + " " + model.Contractor.MiddleName
	doc.Replace("CONTRACTORFIO", contractorFio, -1)

	doc.Replace("NAMEPROFEDUCATION", model.ProgramEducation.NameProfEducation, -1)
	doc.Replace("SIN", model.Enrollment.StartDate.Format("02.01.2006"), -1)
	doc.Replace("CE", "", -1)
	doc.Replace("FOR", model.Enrollment.EndDate.Format("02.01.2006"), -1)

	doc.Replace("OPTION", model.OptionNagruzka, -1)

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
	doc.Replace("BUILDINGE, APARTMENTE", model.Contractor.RegistrationAddress.Building+", "+model.Contractor.RegistrationAddress.Apartment, -1)
	doc.Replace("PHONEE", model.Contractor.Contact_phone, -1)
	doc.Replace("EMAILE", model.Contractor.Email, -1)
}

func (s *DogovorService) ExistsDogovor(fileName string) ([]string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	documents, err := s.s3client.GetDocumentListFromS3(ctx, fileName)
	if err != nil {
		return []string{}, err
	}

	if len(documents) > 0 {
		return documents, nil
	} else {
		return []string{}, errors.New("не найдено ни одного документа")
	}
}

func (s *DogovorService) DeleteDogovor(fileName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := s.s3client.DeleteDocumentFromS3(ctx, fileName)
	if err != nil {
		return err
	}

	return nil
}

func (s *DogovorService) DownloadDogovor(param string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	document, err := s.s3client.GetDocumentFromS3(ctx, param)
	if err != nil {
		return []byte{}, err
	}

	return document, nil
}
