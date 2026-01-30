package service

import (
	"bytes"
	"context"
	"document-service/internal/domain/dto"
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
