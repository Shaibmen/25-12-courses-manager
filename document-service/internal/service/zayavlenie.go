package service

import (
	"bytes"
	"context"
	"document-service/internal/domain/dto"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/nguyenthenguyen/docx"
)

var zayavlenieFourteenPath = "./internal/documents/zayavlenie-fourteen.docx"
var zayavlenieEighteenPath = "./internal/documents/zayavlenie-eighteen.docx"
var zayavlenieBelowEighteenPath = "./internal/documents/zayavlenie-below-eighteen.docx"

const (
	BELOW_EIGHTEEN = "BELOW_EIGHTEEN"
	FOURTEEN       = "FOURTEEN"
	EIGHTEEN       = "EIGHTEEN"
)

type ZayavlenieService struct {
	s3Client S3ClientInterface
}

func NewZayavlenieService(client S3ClientInterface) *ZayavlenieService {
	return &ZayavlenieService{client}
}

func (s *ZayavlenieService) CreateZayavlenie(zayavlenieData *dto.ZayavlenieDTO, dogovorType string) error {

	var doc *docx.Docx
	var err error

	switch dogovorType {
	case BELOW_EIGHTEEN:

		r, err := docx.ReadDocxFile(zayavlenieBelowEighteenPath)
		if err != nil {
			return err
		}
		defer r.Close()

		doc = r.Editable()

		replaceZayavlenieBetweenEighteen(doc, zayavlenieData)

	case FOURTEEN:

		r, err := docx.ReadDocxFile(zayavlenieFourteenPath)
		if err != nil {
			return err
		}
		defer r.Close()

		doc = r.Editable()

		replaceZayavlenieFourteen(doc, zayavlenieData)

	case EIGHTEEN:

		r, err := docx.ReadDocxFile(zayavlenieEighteenPath)
		if err != nil {
			return err
		}
		defer r.Close()

		doc = r.Editable()

		replaceZayavlenieEighteen(doc, zayavlenieData)
	}

	err = doc.Replace("DIVISONEDUCATION", zayavlenieData.ProgramEducation.DivisionEducation, -1)
	err = doc.Replace("NAMEPROFEDUCATION", zayavlenieData.ProgramEducation.DivisionEducation, -1)

	variantsIntToString := make(map[int]string)

	variantsIntToString[1] = "ONE"
	variantsIntToString[2] = "TWO"
	variantsIntToString[3] = "THREE"
	variantsIntToString[4] = "FOUR"
	variantsIntToString[5] = "FIVE"
	variantsIntToString[6] = "SIX"

	for number, token := range variantsIntToString {
		if number != zayavlenieData.Variant {
			doc.Replace(token, "[ ]", -1)
		} else {
			doc.Replace(token, "[x]", -1)
		}
	}

	var buffer bytes.Buffer
	err = doc.Write(&buffer)
	if err != nil {
		return err
	}

	uniqueParam := zayavlenieData.ListenerData.SNILS
	nameFile := "Заявление-" + zayavlenieData.ProgramEducation.NameProfEducation + "_" + uniqueParam + ".docx"

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

	doc.Replace("FCI", model.Passport.PlaceBirth, -1)

	seriaNumber := fmt.Sprintf("%s %s ", model.Passport.Seria, model.Passport.Number)

	doc.Replace("SERIAL NUMBERL ", seriaNumber, -1)
	doc.Replace("GIVENL", model.Passport.PassportGiven, -1)

	doc.Replace("SNILS", model.ListenerData.SNILS, -1)
	doc.Replace("PHONEL", model.ListenerData.ContactPhone, -1)

	doc.Replace("CITY", model.Registration.City, -1)
	streetHouseBuildingApartment := fmt.Sprintf("%s, %s, %s, %s", model.Registration.Street, model.Registration.House, model.Registration.Building, model.Registration.Apartment)
	doc.Replace("STREET, HOUSE, BUILDING, APARTMENT", streetHouseBuildingApartment, -1)

	doc.Replace("EMAILL", model.ListenerData.Email, -1)

	contractorFio := model.Contractor.SecondName + " " + model.Contractor.FirstName + " " + model.Contractor.MiddleName
	doc.Replace("CONTRACTORFIO", contractorFio, -1)

	contractorSeriaNumber := fmt.Sprintf("%s %s ", model.Contractor.Passport.Seria, model.Contractor.Passport.Number)
	doc.Replace("SERIAE NUMBERE", contractorSeriaNumber, -1)

	dob, err = time.Parse(time.RFC3339, model.Contractor.Passport.DateGiven)
	if err == nil {
		doc.Replace("DATEGIVENE", fmt.Sprintf("%02d.%02d.%02d", dob.Day(), dob.Month(), dob.Year()), -1)
	}

	dob, err = time.Parse(time.RFC3339, model.Passport.DateGiven)
	if err == nil {
		doc.Replace("DATEGIVEN", fmt.Sprintf("%02d.%02d.%02d", dob.Day(), dob.Month(), dob.Year()), -1)
	}

	doc.Replace(": GIVENE", ":"+" "+model.Contractor.Passport.PassportGiven, -1)

	doc.Replace("PHONEE", model.Contractor.Contact_phone, -1)
	doc.Replace("EMAILE", model.Contractor.Email, -1)

	doc.Replace("HOUR", strconv.Itoa(model.ProgramEducation.TimeEducation), -1)

	doc.Replace("TYPEOFRETRAINING", model.EnrollmentListener.TypeOfRetraining, -1)

}

func replaceZayavlenieFourteen(doc *docx.Docx, model *dto.ZayavlenieDTO) {
	doc.Replace("STATUS", model.Executor.Status, -1)

	executorFio := model.Executor.ExecutorSurname + " " + model.Executor.ExecutorName + " " + model.Executor.ExecutorMiddlename
	doc.Replace("EXECUTORFIO", executorFio, -1)

	listenerFio := model.ListenerData.SecondName + " " + model.ListenerData.FirstName + " " + model.ListenerData.MiddleName
	doc.Replace("LISTENERFIO", listenerFio, -1)

	contractorFio := model.Contractor.SecondName + " " + model.Contractor.FirstName + " " + model.Contractor.MiddleName
	doc.Replace("CONTRACTORFIO", contractorFio, -1)

	seriaNumber := fmt.Sprintf(": %s %s ", model.Contractor.Passport.Seria, model.Contractor.Passport.Number)

	doc.Replace(": SERIAE NUMBERE ", seriaNumber, -1)
	doc.Replace("GIVENE", model.Contractor.Passport.PassportGiven, -1)

	dob, err := time.Parse(time.RFC3339, model.Contractor.Passport.DateGiven)
	if err == nil {
		doc.Replace("DATEGIVEN", fmt.Sprintf("%02d.%02d.%02d", dob.Day(), dob.Month(), dob.Year()), -1)
	}

	doc.Replace("PHONEE", model.Contractor.Contact_phone, -1)
	doc.Replace("EMAILE", model.Contractor.Email, -1)

	dob, err = time.Parse(time.RFC3339, model.ListenerData.DateOfBirth)
	if err == nil {
		doc.Replace("DATEBIRTH", fmt.Sprintf("%02d.%02d.%02d", dob.Day(), dob.Month(), dob.Year()), -1)
	}

	doc.Replace("TYPEOFRETRAINING", model.EnrollmentListener.TypeOfRetraining, -1)

	// variant := "VARIANT" + strconv.Itoa(model.Variant)

	// i := 1

	// for ; i < 7; i++ {
	// 	tempVar := "VARIANT" + strconv.Itoa(i)
	// 	if tempVar != variant {
	// 		doc.Replace(tempVar, "[ ]", -1)
	// 	} else {
	// 		doc.Replace(tempVar, "[x]", -1)
	// 	}
	// }
}

func replaceZayavlenieEighteen(doc *docx.Docx, model *dto.ZayavlenieDTO) {
	doc.Replace("STATUS", model.Executor.Status, -1)

	executorFio := model.Executor.ExecutorSurname + " " + model.Executor.ExecutorName + " " + model.Executor.ExecutorMiddlename
	doc.Replace("EXECUTORFIO", executorFio, -1)

	listenerFio := model.ListenerData.SecondName + " " + model.ListenerData.FirstName + " " + model.ListenerData.MiddleName
	doc.Replace("LISTENERFIO", listenerFio, -1)

	doc.Replace("FCI", model.Registration.City, -1)

	numberSeria := fmt.Sprintf(": %s %s ", model.Passport.Seria, model.Passport.Number)

	doc.Replace(": SERIAE NUMBERE ", numberSeria, -1)
	doc.Replace("GIVEN", model.Passport.PassportGiven, -1)
	doc.Replace("SNILS", model.ListenerData.SNILS, -1)

	dob, err := time.Parse(time.RFC3339, model.Passport.DateGiven)
	if err == nil {
		doc.Replace("DATEGIVEN", fmt.Sprintf("%02d.%02d.%02d", dob.Day(), dob.Month(), dob.Year()), -1)
	}

	dob, err = time.Parse(time.RFC3339, model.ListenerData.DateOfBirth)
	if err == nil {
		doc.Replace("ROZHD", fmt.Sprintf("%02d.%02d.%02d", dob.Day(), dob.Month(), dob.Year()), -1)
	}

	doc.Replace("PHONE", model.ListenerData.ContactPhone, -1)

	doc.Replace("CITY", model.Registration.City, -1)
	doc.Replace("STREET", model.Registration.Street, -1)
	doc.Replace("HOUSE", model.Registration.House, -1)
	doc.Replace("BUILDING", model.Registration.Building, -1)
	doc.Replace("APARTMENT", model.Registration.Apartment, -1)
	doc.Replace("EMAIL", model.ListenerData.Email, -1)

	// i := 1

	// for ; i < 7; i++ {
	// 	tempVar := "VAR" + strconv.Itoa(i)
	// 	if tempVar != variant {
	// 		doc.Replace(tempVar, "[ ]", -1)
	// 	} else {
	// 		doc.Replace(tempVar, "[x]", -1)
	// 	}
	// }

	doc.Replace("RETRAINTYPE", model.EnrollmentListener.TypeOfRetraining, -1)
	doc.Replace("EDUNAME", model.ProgramEducation.NameProfEducation, -1)
	doc.Replace("OBEM", strconv.Itoa(model.ProgramEducation.TimeEducation), -1)
}

func (s *ZayavlenieService) ExistsZayavlenie(fileName string) ([]string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	documents, err := s.s3Client.GetDocumentListFromS3(ctx, fileName)
	if err != nil {
		return []string{}, err
	}

	if len(documents) > 0 {
		return documents, nil
	} else {
		return []string{}, errors.New("Не найдено ни одного документа")
	}
}

func (s *ZayavlenieService) DeleteZayavlenie(fileName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := s.s3Client.DeleteDocumentFromS3(ctx, fileName)
	if err != nil {
		return err
	}

	return nil
}

func (s *ZayavlenieService) DownloadZayavlenie(param string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	document, err := s.s3Client.GetDocumentFromS3(ctx, param)
	if err != nil {
		return []byte{}, err
	}

	return document, nil
}
