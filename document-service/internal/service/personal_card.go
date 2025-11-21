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

var (
	templatePath = "./internal/documents/personal_card.docx"
)

type PersonalCardService struct {
	s3Client S3ClientInterface
}

func NewPersonalCardService(client S3ClientInterface) *PersonalCardService {
	return &PersonalCardService{s3Client: client}
}

func (p *PersonalCardService) CreatePersonalCard(ListenerData *dto.FullListenerDataDTO) error {

	r, err := docx.ReadDocxFile(templatePath)
	if err != nil {
		return err
	}
	defer r.Close()

	doc := r.Editable()
	replace(doc, ListenerData)

	var buffer bytes.Buffer
	err = doc.Write(&buffer)
	if err != nil {
		return err
	}

	uniqueParam := ListenerData.Listener.SNILS
	nameFile := "Личное-дело-" + ListenerData.ProgramEducation.NameProfEducation + "_" + uniqueParam + ".docx"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = p.s3Client.SendDocumentToS3(ctx, buffer, nameFile)
	if err != nil {
		return err
	}

	return nil
}

func replace(doc *docx.Docx, model *dto.FullListenerDataDTO) {
	doc.Replace("}}", "", -1)
	doc.Replace("{{", "", -1)
	doc.Replace("ProgramEducation", model.ProgramEducation.NameProfEducation, -1)
	doc.Replace("TypeOfEducation", model.ProgramEducation.EducationType, -1)
	doc.Replace("Hour", strconv.Itoa(model.ProgramEducation.TimeEducation), -1)

	fio := model.Listener.SecondName + " " + model.Listener.FirstName + " " + model.Listener.MiddleName

	doc.Replace("FIO", fio, -1)

	dob, err := time.Parse(time.RFC3339, model.Listener.DateOfBirth)
	if err == nil {
		doc.Replace("DateBirth", fmt.Sprintf("%02d", dob.Day()), -1)
		doc.Replace("MonthBirth", fmt.Sprintf("%02d", dob.Month()), -1)
		doc.Replace("YearBirth", fmt.Sprintf("%d", dob.Year()), -1)
	}

	doc.Replace("City", model.Passport.PlaceBirth, -1)
	doc.Replace("Citizenship", model.Passport.Citizenship, -1)

	switch model.Passport.Gender {
	case "Мужской":
		doc.Replace("мужской ☐", "мужской ☒", -1)
	case "Женский":
		doc.Replace("женский ☐", "женский ☒", -1)
	}

	doc.Replace("Seria", model.Passport.Seria, -1)
	doc.Replace("Number", model.Passport.Number, -1)
	doc.Replace("WhoGiven", model.Passport.PassportGiven, -1)

	given, err := time.Parse(time.RFC3339, model.Passport.DateGiven)
	if err == nil {
		doc.Replace("WhenGiven", given.Format("02.01.2006"), -1)
	}

	doc.Replace("MainIndex", model.RegistrationAddress.MailIndex, -1)
	doc.Replace("Region", model.RegistrationAddress.Region, -1)
	doc.Replace("City", model.RegistrationAddress.City, -1)
	doc.Replace("Street", model.RegistrationAddress.Street, -1)
	doc.Replace("House", model.RegistrationAddress.House, -1)
	doc.Replace("Building", model.RegistrationAddress.Building, -1)
	doc.Replace("Appartaments", model.RegistrationAddress.Apartment, -1)
	doc.Replace("SNILS", model.Listener.SNILS, -1)
	doc.Replace("Phone", model.Listener.ContactPhone, -1)
	doc.Replace("Email", model.Listener.Email, -1)

	switch model.EducationListener.LevelEducation {
	case "Среднее":
		doc.Replace("среднее ☐", "среднее ☒", -1)
	case "Среднее профессиональное":
		doc.Replace("среднее профессиональное ☐", "среднее профессиональное ☒", -1)
	case "Высшее":
		doc.Replace("высшее ☐", "высшее ☒", -1)
	}

	if model.EducationListener != (dto.EducationListenerDTO{}) {
		doc.Replace("диплом ☐", "диплом ☒", -1)
		doc.Replace("SDip", model.EducationListener.DiplomSeria, -1)
		doc.Replace("NDip", model.EducationListener.DiplomNumber, -1)
		dmy, err := time.Parse(time.RFC3339, model.Listener.DateOfBirth)
		if err == nil {
			doc.Replace("DDip", fmt.Sprintf("%02d", dmy.Day()), -1)
			doc.Replace("MDip", fmt.Sprintf("%02d", dmy.Month()), -1)
			doc.Replace("YDip", fmt.Sprintf("%d", dmy.Year()), -1)
		}
		doc.Replace("CDip", model.EducationListener.City, -1)
		doc.Replace("RDip", model.EducationListener.Region, -1)
		doc.Replace("WhoDip", model.EducationListener.EducationalInstitution, -1)
		doc.Replace("Speciality", model.EducationListener.Speciality, -1)
	} else {
		doc.Replace("SDip", "_______", -1)
		doc.Replace("NDip", "_______", -1)
		doc.Replace("DDip", "_______", -1)
		doc.Replace("MDip", "_______", -1)
		doc.Replace("YDip", "_______", -1)

		doc.Replace("CDip", "_______", -1)
		doc.Replace("RDip", "_______", -1)
		doc.Replace("WhoDip", "_______", -1)
		doc.Replace("Speciality", "_______", -1)
	}

	if model.PlaceWork != (dto.PlaceWorkDTO{}) {
		doc.Replace("PlaceWork", model.PlaceWork.NameCompany, -1)
		doc.Replace("Post", model.PlaceWork.JobTitle, -1)
		doc.Replace("AllP", strconv.Itoa(model.PlaceWork.AllExperience), -1)
		doc.Replace("OnP", strconv.Itoa(model.PlaceWork.JobTitleExpirience), -1)
	} else {
		doc.Replace("PlaceWork", "_______", -1)
		doc.Replace("Post", "_______", -1)
		doc.Replace("AllP", "_______", -1)
		doc.Replace("OnP", "_______", -1)
	}

	doc.Replace("Division", model.ProgramEducation.DivisionEducation, -1)

}

func (p *PersonalCardService) ExistsPersonalCard(fileName string) ([]string, error) {

	// dir := "./personal_card"

	// files, err := os.ReadDir(dir)
	// if err != nil {
	// 	return nil, err
	// }

	// for _, file := range files {
	// 	if !file.IsDir() && strings.Contains(file.Name(), fileName) {
	// 		matches = append(matches, file.Name())

	// 	}
	// }

	// return matches, nil
	// // _, err := os.Stat(fileName)
	// // if err != nil {
	// // 	return false, err
	// // }

	// // return true, nil
	var matches []string

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	doc, err := p.s3Client.GetDocumentFromS3(ctx, fileName)
	if err != nil {
		return []string{}, err
	}

	if len(doc) > 0 {
		matches = append(matches, fileName)
	}

	return matches, nil
}

func (p *PersonalCardService) DeletePersonalCard(fileName string) error {

	// err := os.Remove(fileName)
	// if err != nil {
	// 	return err
	// }

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := p.s3Client.DeleteDocumentFromS3(ctx, fileName)
	if err != nil {
		return err
	}

	return nil
}

func (p *PersonalCardService) DownloadPersonalCard(fileName string) ([]byte, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	
	doc, err := p.s3Client.GetDocumentFromS3(ctx, fileName)
	if err != nil {
		return []byte{}, err
	}

	return doc, nil

}
