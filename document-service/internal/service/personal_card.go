package service

import (
	"bytes"
	"context"
	"document-service/internal/domain/dto"
	"errors"
	"fmt"

	// "log"
	"strconv"
	"time"

	"github.com/nguyenthenguyen/docx"
)

var (
	personalCardPath = "./internal/documents/personal-card.docx"
)

type PersonalCardService struct {
	s3Client S3ClientInterface
}

func NewPersonalCardService(client S3ClientInterface) *PersonalCardService {
	return &PersonalCardService{s3Client: client}
}

func (p *PersonalCardService) CreatePersonalCard(ListenerData *dto.FullListenerDataDTO) error {

	r, err := docx.ReadDocxFile(personalCardPath)
	if err != nil {
		return err
	}
	defer r.Close()

	doc := r.Editable()
	replaceFIZ(doc, ListenerData) // сделать так же проверку на юрика

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

func replaceFIZ(doc *docx.Docx, model *dto.FullListenerDataDTO) {
	doc.Replace("NAMEPROFEDUCATION", model.ProgramEducation.NameProfEducation, -1)
	doc.Replace("TYPEOFRETRAINING", model.EnrollmentListener.TypeOfRetraining, -1)
	doc.Replace("HOUR", strconv.Itoa(model.ProgramEducation.TimeEducation), -1)

	fio := model.Listener.SecondName + " " + model.Listener.FirstName + " " + model.Listener.MiddleName

	doc.Replace("LISTENERFIO", fio, -1)

	dob, err := time.Parse(time.RFC3339, model.Listener.DateOfBirth)
	if err == nil {
		doc.Replace("BIRTHD", fmt.Sprintf("%02d.%02d.%02d", dob.Day(), dob.Month(), dob.Year()), -1)
		// doc.Replace("", fmt.Sprintf("%02d", dob.Month()), -1)
		// doc.Replace("YearBirth", fmt.Sprintf("%d", dob.Year()), -1)
	}

	doc.Replace("CITY", model.Passport.PlaceBirth, -1)
	doc.Replace("CITIZENSHIP", model.Passport.Citizenship, -1)

	switch model.Passport.Gender {
	case "Мужской":
		doc.Replace("MUZH" /*"мужской ☐"*/, "мужской [x]", -1)
		doc.Replace("ZHEN", "женский [ ]", -1)
		// log.Println("Мужской пол")
	case "Женский":
		doc.Replace("MUZH", "мужской [ ]", -1)
		doc.Replace("ZHEN", "женский [x]", -1)
		// log.Println("Женский пол")
	default:
		doc.Replace("MUZH", "мужской [ ]", -1)
		doc.Replace("ZHEN", "женский [ ]", -1)
		// log.Println("другой пол:", model.Passport.Gender)
	}

	doc.Replace("SERIA", model.Passport.Seria, -1)
	doc.Replace("NUMBER", model.Passport.Number, -1)
	doc.Replace("GIVEN", model.Passport.PassportGiven, -1)

	given, err := time.Parse(time.RFC3339, model.Passport.DateGiven)
	if err == nil {
		doc.Replace("PASSPORTD", given.Format("02.01.2006"), -1)
	}

	doc.Replace("INDEX", model.RegistrationAddress.MailIndex, -1)
	doc.Replace("REGION", model.RegistrationAddress.Region, -1)
	doc.Replace("CITY", model.RegistrationAddress.City, -1)
	doc.Replace("STREET", model.RegistrationAddress.Street, -1)
	doc.Replace("HOUSE", model.RegistrationAddress.House, -1)
	doc.Replace("BUILDING", model.RegistrationAddress.Building, -1)
	doc.Replace("APARTMENT", model.RegistrationAddress.Apartment, -1)
	doc.Replace("SNILS", model.Listener.SNILS, -1)
	doc.Replace("PHONE", model.Listener.ContactPhone, -1)
	doc.Replace("EMAIL", model.Listener.Email, -1)

	switch model.EducationListener.LevelEducation {
	case "Среднее":
		// doc.Replace("SREDN", "среднее ☒", -1)
		// doc.Replace("SPROF", "среднее профессиональное ☐", -1)
		// doc.Replace("VISH", "высшее ☐", -1)
		// log.Println("Среднее образование")
		doc.Replace("PACANI", "среднее [x] среднее профессиональное [ ] высшее [ ]", -1)
	case "Среднее специальное":
		// doc.Replace("SREDN", "среднее ☐", -1)
		// doc.Replace("SPROF", "среднее профессиональное ☒", -1)
		// doc.Replace("VISH", "высшее ☐", -1)
		// log.Println("Среднее профессиональное образование")
		doc.Replace("PACANI", "среднее [ ] среднее профессиональное [x] высшее [ ]", -1)
	case "Бакалавр", "Магистр", "Кандидат наук":
		// doc.Replace("SREDN", "среднее ☐", -1)
		// doc.Replace("SPROF", "среднее профессиональное ☐", -1)
		// doc.Replace("VISH", "высшее ☒", -1)
		// log.Println("Высшее образование")
		doc.Replace("PACANI", "среднее [ ] среднее профессиональное [ ] высшее [x]", -1)
	default:
		// log.Println("Другой тип образования:", model.EducationListener.LevelEducation)
		doc.Replace("PACANI", "среднее [ ] среднее профессиональное [ ] высшее [ ]", -1)
	}

	if model.EducationListener != (dto.EducationListenerDTO{}) {
		doc.Replace("DIPLOM", "диплом [x]", -1)
		doc.Replace("DIPS", model.EducationListener.DiplomSeria, -1)
		doc.Replace("DIPN", model.EducationListener.DiplomNumber, -1)
		dmy, _ := time.Parse(time.RFC3339, model.EducationListener.DateGiven)
		// if err == nil {
		doc.Replace("IPD", fmt.Sprintf("%02d.%02d.%02d", dmy.Day(), dmy.Month(), dmy.Year()), -1)
		// doc.Replace("MDip", fmt.Sprintf("%02d", dmy.Month()), -1)
		// doc.Replace("YDip", fmt.Sprintf("%d", dmy.Year()), -1)
		// }
		doc.Replace("DIPC", model.EducationListener.City, -1)
		doc.Replace("DIPR", model.EducationListener.Region, -1)
		doc.Replace("INSTITUTION", model.EducationListener.EducationalInstitution, -1)
		// doc.Replace("Speciality", model.EducationListener.Speciality, -1)
	} else {
		doc.Replace("DIPLOM", "диплом [ ]", -1) // диплом ☐
		doc.Replace("DIPS", "_______", -1)
		doc.Replace("DIPN", "_______", -1)
		doc.Replace("IPD", "_______", -1)
		doc.Replace("MDip", "_______", -1)
		doc.Replace("YDip", "_______", -1)

		doc.Replace("DIPC", "_______", -1)
		doc.Replace("DIPR", "_______", -1)
		doc.Replace("INSTITUTION", "_______", -1)
		// doc.Replace("Speciality", "_______", -1)
	}

	if model.PlaceWork != (dto.PlaceWorkDTO{}) {
		doc.Replace("NAMECOM", model.PlaceWork.NameCompany, -1)
		doc.Replace("JOBTITLE", model.PlaceWork.JobTitle, -1)
		doc.Replace("AEXP", strconv.Itoa(model.PlaceWork.AllExperience), -1)
		doc.Replace("JEXP", strconv.Itoa(model.PlaceWork.JobTitleExpirience), -1)
	} else {
		doc.Replace("NAMECOM", "_______", -1)
		doc.Replace("JOBTITLE", "_______", -1)
		doc.Replace("AEXP", "_______", -1)
		doc.Replace("JEXP", "_______", -1)
	}

	// doc.Replace("Division", model.ProgramEducation.DivisionEducation, -1)

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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	documents, err := p.s3Client.GetDocumentListFromS3(ctx, fileName)
	if err != nil {
		return []string{}, err
	}

	if len(documents) > 0 {
		return documents, nil
	} else {
		return []string{}, errors.New("Не найдено ни одного документа")
	}
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

func (p *PersonalCardService) DownloadPersonalCard(param string) ([]byte, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	document, err := p.s3Client.GetDocumentFromS3(ctx, param)
	if err != nil {
		return []byte{}, err
	}

	return document, nil
}
