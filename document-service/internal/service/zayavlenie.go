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

func (s *ZayavlenieService) CreateZayavlenie(zayavlenieData *dto.ZayavlenieDTO, dogovorType int) error {

	r, err := docx.ReadDocxFile(personalCardPath)
	if err != nil {
		return err
	}
	defer r.Close()

	doc := r.Editable()
	replaceZayavlenieBetweenEighteen(doc, zayavlenieData)

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

	doc.Replace("CITY", model.Passport.PlaceBirth, -1)

	doc.Replace("SERIAL NUMBERL", model.Passport.Seria, -1)
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
	doc.Replace("DATEGIVENE", model.Contractor.Passport.DateGiven, -1)

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
