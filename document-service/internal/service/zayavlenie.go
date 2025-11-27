package service

import (
	"bytes"
	"context"
	"document-service/internal/domain/dto"
	"fmt"
	"strconv"
	"time"

	"github.com/nguyenthenguyen/docx"
)

var zayavleniePath = "./internal/documents/zayavlenie.docx"

type ZayavlenieService struct {
	s3Client S3ClientInterface
}

func NewZayavlenieService(client S3ClientInterface) *ZayavlenieService {
	return &ZayavlenieService{client}
}

func (s *ZayavlenieService) CreateZayavlenie(zayavlenieData *dto.ZayavlenieDTO, dogovorType string) error {

	r, err := docx.ReadDocxFile(personalCardPath)
	if err != nil {
		return err
	}
	defer r.Close()

	doc := r.Editable()
	switch dogovorType {
	case "betweenEighteen":
		replaceZayavlenieBetweenEighteen(doc, zayavlenieData)
	case "fourteen":
		replaceZayavlenieFourteen(doc, zayavlenieData)
	case "eighteen":
		replaceZayavlenieEighteen(doc, zayavlenieData)
	}

	var buffer bytes.Buffer
	err = doc.Write(&buffer)
	if err != nil {
		return err
	}

	uniqueParam := zayavlenieData.ListenerData.SNILS
	nameFile := "Личное-дело-" + zayavlenieData.ProgramEducation.NameProfEducation + "_" + uniqueParam + ".docx"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = s.s3Client.SendDocumentToS3(ctx, buffer, nameFile)
	if err != nil {
		return err
	}

	return nil
}

func replaceZayavlenieBetweenEighteen(doc *docx.Docx, model *dto.ZayavlenieDTO) {
	doc.Replace("STATUS", model.Executor.Status, -1)

	executorFio := model.Executor.ExecutorSurname + " " + model.Executor.ExecutorName + " " + model.Executor.ExecutorMiddlename
	doc.Replace("EXECUTORFIO", executorFio, -1)

	listenerFio := model.ListenerData.SecondName + " " + model.ListenerData.FirstName + " " + model.ListenerData.MiddleName
	doc.Replace("LISTENERFIO", listenerFio, -1)

	dob, err := time.Parse(time.RFC3339, model.ListenerData.DateOfBirth)
	if err == nil {
		doc.Replace("DATEBIRTH", fmt.Sprintf("%02d.%02d.%02d", dob.Day(), dob.Month(), dob.Year()), -1)
	}

	doc.Replace("CITYB", model.Passport.PlaceBirth, -1)

	doc.Replace("SERIAL", model.Passport.Seria, -1)
	doc.Replace("NUMBERL", model.Passport.Number, -1)
	doc.Replace("GIVENL", model.Passport.PassportGiven, -1)
	doc.Replace("DATEGIVENS", model.Passport.DateGiven, -1)

	doc.Replace("SNILS", model.ListenerData.SNILS, -1)
	doc.Replace("PHONEL", model.ListenerData.ContactPhone, -1)
	doc.Replace("CITY", model.Registration.City, -1)
	doc.Replace("STREET", model.Registration.Street, -1)
	doc.Replace("HOUSE", model.Registration.House, -1)
	doc.Replace("BUILDING", model.Registration.Building, -1)
	doc.Replace("APARTMENT", model.Registration.Apartment, -1)
	doc.Replace("EMAILL", model.ListenerData.Email, -1)

	contractorFio := model.Contractor.SecondName + " " + model.Contractor.FirstName + " " + model.Contractor.MiddleName
	doc.Replace("CONTRACTORFIO", contractorFio, -1)

	doc.Replace("SERIAE NUMBERE", model.Contractor.Passport.Seria, -1)
	doc.Replace("GIVENE", model.Contractor.Passport.PassportGiven, -1)
	doc.Replace("DATEGIVENS", model.Contractor.Passport.DateGiven, -1)

	doc.Replace("PHONEE", model.Contractor.Contact_phone, -1)
	doc.Replace("EMAILE", model.Contractor.Email, -1)

	variant := "VARIANT" + strconv.Itoa(model.Variant)

	i := 1

	for ; i < 7; i++ {
		tempVar := "VARIANT" + strconv.Itoa(i)
		if tempVar != variant {
			doc.Replace(tempVar, "[ ]", -1)
		} else {
			doc.Replace(tempVar, "[x]", -1)
		}
	}
}

func replaceZayavlenieFourteen(doc *docx.Docx, model *dto.ZayavlenieDTO) {
	doc.Replace("STATUS", model.Executor.Status, -1)

	executorFio := model.Executor.ExecutorSurname + " " + model.Executor.ExecutorName + " " + model.Executor.ExecutorMiddlename
	doc.Replace("EXECUTORFIO", executorFio, -1)

	listenerFio := model.ListenerData.SecondName + " " + model.ListenerData.FirstName + " " + model.ListenerData.MiddleName
	doc.Replace("LISTENERFIO", listenerFio, -1)

	contractorFio := model.Contractor.SecondName + " " + model.Contractor.FirstName + " " + model.Contractor.MiddleName
	doc.Replace("CONTRACTORFIO", contractorFio, -1)

	doc.Replace("SERIAE", model.Contractor.Passport.Seria, -1)
	doc.Replace("NUMBERE", model.Contractor.Passport.Number, -1)
	doc.Replace("GIVENE", model.Contractor.Passport.PassportGiven, -1)
	doc.Replace("DATEGIVENS", model.Contractor.Passport.DateGiven, -1)

	doc.Replace("PHONEE", model.Contractor.Contact_phone, -1)
	doc.Replace("EMAILE", model.Contractor.Email, -1)

	dob, err := time.Parse(time.RFC3339, model.ListenerData.DateOfBirth)
	if err == nil {
		doc.Replace("DATEBIRTH", fmt.Sprintf("%02d.%02d.%02d", dob.Day(), dob.Month(), dob.Year()), -1)
	}

	variant := "VARIANT" + strconv.Itoa(model.Variant)

	i := 1

	for ; i < 7; i++ {
		tempVar := "VARIANT" + strconv.Itoa(i)
		if tempVar != variant {
			doc.Replace(tempVar, "[ ]", -1)
		} else {
			doc.Replace(tempVar, "[x]", -1)
		}
	}
}

func replaceZayavlenieEighteen(doc *docx.Docx, model *dto.ZayavlenieDTO) {
	doc.Replace("STATUS", model.Executor.Status, -1)

	executorFio := model.Executor.ExecutorSurname + " " + model.Executor.ExecutorName + " " + model.Executor.ExecutorMiddlename
	doc.Replace("EXECUTORFIO", executorFio, -1)

	listenerFio := model.ListenerData.SecondName + " " + model.ListenerData.FirstName + " " + model.ListenerData.MiddleName
	doc.Replace("LISTENERFIO", listenerFio, -1)

	dob, err := time.Parse(time.RFC3339, model.ListenerData.DateOfBirth)
	if err == nil {
		doc.Replace("DATEBIRTH", fmt.Sprintf("%02d.%02d.%02d", dob.Day(), dob.Month(), dob.Year()), -1)
	}

	doc.Replace("FCITY", model.Registration.City, -1)

	doc.Replace("SERIAE", model.Passport.Seria, -1)
	doc.Replace("NUMBERE", model.Passport.Number, -1)
	doc.Replace("FGIVEN", model.Passport.PassportGiven, -1)
	doc.Replace("DATEGIVEN", model.Passport.DateGiven, -1)
	doc.Replace("SNILS", model.ListenerData.SNILS, -1)

	doc.Replace("PHONE", model.ListenerData.ContactPhone, -1)

	doc.Replace("CITY", model.Registration.City, -1)
	doc.Replace("STREET", model.Registration.Street, -1)
	doc.Replace("HOUSE", model.Registration.House, -1)
	doc.Replace("BUILDING", model.Registration.Building, -1)
	doc.Replace("APARTMENT", model.Registration.Apartment, -1)
	doc.Replace("EMAIL", model.ListenerData.Email, -1)

	variant := "VARIANT" + strconv.Itoa(model.Variant)

	i := 1

	for ; i < 7; i++ {
		tempVar := "VARIANT" + strconv.Itoa(i)
		if tempVar != variant {
			doc.Replace(tempVar, "[ ]", -1)
		} else {
			doc.Replace(tempVar, "[x]", -1)
		}
	}
}
