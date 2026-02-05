package service

import (
	"bytes"
	"context"
	"document-service/internal/domain/dto"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"

	"github.com/nguyenthenguyen/docx"
)

const (
	dogovorPP_3_path     = "./internal/documents/PP-FIZ-3.docx"
	dogovorPP_2_path     = "./internal/documents/PP-FIZ-2.docx"
	dogovorPK_3_path     = "./internal/documents/PK-FIZ-3.docx"
	dogovorPK_2_path     = "./internal/documents/PK-FIZ-2.docx"
	dogovorDO_3_path     = "./internal/documents/DO-FIZ-3.docx"
	dogovorPP_3_YUR_path = "./internal/documents/PP-YUR-3.docx"
	dogovorPK_3_YUR_path = "./internal/documents/PK-YUR-3.docx"
)

const (
	PP_3_FIZ = "PP_3_FIZ"
	PP_2_FIZ = "PP_2_FIZ"
	PK_3_FIZ = "PK_3_FIZ"
	PK_2_FIZ = "PK_2_FIZ"
	DO_3_FIZ = "DO_3_FIZ"
	PK_3_YUR = "PK_3_YUR"
	PP_3_YUR = "PP_3_YUR"
)

type DogovorService struct {
	s3client    S3ClientInterface
	nagruzkaMap map[int]string
	diplomMap   map[int]string
	doMap       map[int]string
}

func NewDogovorService(s3client S3ClientInterface) *DogovorService {

	nagruzkaMap := make(map[int]string)
	nagruzkaMap[1] = "Недельная учебная нагрузка по настоящему договору составляет 3 академических часа в неделю, включая 2 академических часа взаимодействия с преподавателем и 1 академический час самостоятельной работы; общая продолжительность освоения — 81 неделя"
	nagruzkaMap[2] = "Недельная учебная нагрузка по настоящему договору составляет 6 академических часов в неделю, включая 4 академических часа взаимодействия с преподавателем и 2 академических часа самостоятельной работы; общая продолжительность освоения — 41 неделя."
	nagruzkaMap[3] = "Недельная учебная нагрузка по настоящему договору составляет 12 академических часов в неделю, включая 8 академических часов взаимодействия с преподавателем и 4 академических часа самостоятельной работы; общая продолжительность освоения — 21 неделя."
	nagruzkaMap[4] = "Недельная учебная нагрузка по настоящему договору составляет 15 академических часов в неделю, включая 10 академических часов взаимодействия с преподавателем и 5 академических часов самостоятельной работы; общая продолжительность освоения — 17 недель."
	nagruzkaMap[5] = "Недельная учебная нагрузка по настоящему договору составляет 30 академических часов в неделю, включая 20 академических часов взаимодействия с преподавателем и 10 академических часов самостоятельной работы; общая продолжительность освоения — 9 недель."
	nagruzkaMap[6] = "Недельная учебная нагрузка по настоящему договору составляет 32 академических часа в неделю, включая 20 академических часов взаимодействия с преподавателем и 12 академических часов самостоятельной работы; общая продолжительность освоения — 8 недель."

	diplomMap := make(map[int]string)
	diplomMap[1] = "диплом о профпереподготовке вручается по окончании;"
	diplomMap[2] = "диплом выдаётся одновременно с дипломом СПО/ВО (ч. 16 ст. 76 ФЗ-273). До этого момента диплом хранится у Исполнителя."

	doMap := make(map[int]string)
	doMap[1] = "Недельная учебная нагрузка по настоящему договору составляет 1 академический час в неделю; общая продолжительность освоения — 20 недель."
	doMap[2] = "Недельная учебная нагрузка по настоящему договору составляет 2 академических часа в неделю; общая продолжительность освоения — 10 недель."
	doMap[3] = "Недельная учебная нагрузка по настоящему договору составляет 4 академических часа в неделю; общая продолжительность освоения — 5 недель."
	doMap[4] = "Недельная учебная нагрузка по настоящему договору составляет 8 академических часов в неделю; общая продолжительность освоения — 2,5 недели."
	doMap[5] = "Недельная учебная нагрузка по настоящему договору составляет 10 академических часов в неделю; общая продолжительность освоения — 2 недели."

	return &DogovorService{s3client, nagruzkaMap, diplomMap, doMap}
}

func (s *DogovorService) CreateDogovor(dogovor *dto.DogovorDTO, dogovorType string) error {

	var doc *docx.Docx
	var r *docx.ReplaceDocx
	var err error
	var nameFile string

	type listenersTableData struct {
		FIO       string `json:"fio"`
		DateBirth string `json:"date_birth"`
		Document  string `json:"document"`
		SNILS     string `json:"SNILS"`
		Email     string `json:"email"`
		Period    string `json:"period"`
		Obem      string `json:"obem"`
		Cost      string `json:"cost"`
	}

	allListenersFormatData := make([]listenersTableData, 0)

	for _, listener := range dogovor.Zakazchik.Listeners {

		allListenersFormatData = append(allListenersFormatData, listenersTableData{
			FIO:       listener.SecondName + " " + listener.FirstName + " " + listener.MiddleName,
			DateBirth: listener.DateOfBirth,
			Document:  "-",
			SNILS:     listener.SNILS,
			Email:     listener.Email,
			Period:    dogovor.Enrollment.StartDate.Format("02.01.2006") + " - " + dogovor.Enrollment.EndDate.Format("02.01.2006"),
			Obem:      strconv.Itoa(dogovor.ProgramEducation.TimeEducation),
			Cost:      fmt.Sprintf("%.2f", dogovor.Enrollment.CurrentPrice),
		})
	}

	jsonListenersData, err := json.Marshal(allListenersFormatData)
	if err != nil {
		log.Println("произошла ошибка при маршаллинге:", err)
		return err
	}

	switch dogovorType {
	case PP_3_FIZ:
		r, err = docx.ReadDocxFile(dogovorPP_3_path)
		if err != nil {
			return err
		}
		defer r.Close()

		nameFile = "Договор-" + dogovor.ProgramEducation.NameProfEducation + "_" + dogovor.ListenerData.SNILS + ".docx"

		doc = r.Editable()

		replacePP3FIZ(doc, dogovor)
	case PP_2_FIZ:
		r, err = docx.ReadDocxFile(dogovorPP_2_path)
		if err != nil {
			return err
		}
		defer r.Close()

		nameFile = "Договор-" + dogovor.ProgramEducation.NameProfEducation + "_" + dogovor.ListenerData.SNILS + ".docx"

		doc = r.Editable()

		replacePP2FIZ(doc, dogovor)
	case PK_3_FIZ:
		r, err = docx.ReadDocxFile(dogovorPK_3_path)
		if err != nil {
			return err
		}
		defer r.Close()

		nameFile = "Договор-" + dogovor.ProgramEducation.NameProfEducation + "_" + dogovor.ListenerData.SNILS + ".docx"

		doc = r.Editable()

		replacePK3FIZ(doc, dogovor)
	case PK_2_FIZ:
		r, err = docx.ReadDocxFile(dogovorPK_2_path)
		if err != nil {
			return err
		}
		defer r.Close()

		nameFile = "Договор-" + dogovor.ProgramEducation.NameProfEducation + "_" + dogovor.ListenerData.SNILS + ".docx"

		doc = r.Editable()

		replacePK2FIZ(doc, dogovor)
	case DO_3_FIZ:
		r, err = docx.ReadDocxFile(dogovorDO_3_path)
		if err != nil {
			return err
		}
		defer r.Close()

		nameFile = "Договор-" + dogovor.ProgramEducation.NameProfEducation + "_" + dogovor.ListenerData.SNILS + ".docx"

		doc = r.Editable()

		replaceDO3FIZ(doc, dogovor)
	case PP_3_YUR:

		cmd := exec.Command("python", "scripts/make_listeners_tables.py", "/app/internal/documents/PP-YUR-3.docx", "/app/internal/documents/PP-YUR-3-RAW.docx", string(jsonListenersData))

		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Println("пиздец!", err, string(output))
		}

		r, err = docx.ReadDocxFile(dogovorPP_3_YUR_path)
		if err != nil {
			return err
		}
		defer r.Close()

		nameFile = "Договор-" + dogovor.ProgramEducation.NameProfEducation + "_" + dogovor.Zakazchik.CompanyName + "_" + time.Now().Format("2001.02.01") + ".docx"

		doc = r.Editable()

		replacePP3YUR(doc, dogovor)
	case PK_3_YUR:

		cmd := exec.Command("python", "scripts/make_listeners_tables.py", "/app/internal/documents/PK-YUR-3.docx", "/app/internal/documents/PK-YUR-3-RAW.docx", string(jsonListenersData))

		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Println("пиздец!", err, string(output))
		}

		r, err = docx.ReadDocxFile(dogovorPK_3_YUR_path)
		if err != nil {
			return err
		}
		defer r.Close()

		nameFile = "Договор-" + dogovor.ProgramEducation.NameProfEducation + "_" + dogovor.Zakazchik.CompanyName + "_" + time.Now().Format("2001.02.01") + ".docx"

		doc = r.Editable()

		replacePK3YUR(doc, dogovor)
	default:
		return errors.New("бро ты натворил ъуйни, ожидай последствия")
	}

	doc.Replace("OPTIOND", s.diplomMap[dogovor.OptionDocument], -1)
	doc.Replace("OPTIONH", s.nagruzkaMap[dogovor.OptionNagruzka], -1)
	doc.Replace("OPTIONPRICE", dogovor.OptionPrice, -1)
	doc.Replace("PRICE", fmt.Sprintf("%.2f руб.", dogovor.Enrollment.CurrentPrice), -1)
	doc.Replace("OPTION", s.doMap[dogovor.OptionNagruzka], -1)

	var buffer bytes.Buffer
	err = doc.Write(&buffer)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = s.s3client.SendDocumentToS3(ctx, buffer, nameFile)
	if err != nil {
		return err
	}

	return nil
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
		return []string{}, errors.New("Не найдено ни одного документа")
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
