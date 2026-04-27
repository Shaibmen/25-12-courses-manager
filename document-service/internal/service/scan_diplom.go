package service

import (
	"bytes"
	"context"
	"errors"
	"time"
)

type ScanDiplomService struct {
	s3client S3ClientInterface
}

func NewScanDiplomService(s3client S3ClientInterface) *ScanDiplomService {
	return &ScanDiplomService{
		s3client: s3client,
	}
}

func (s *ScanDiplomService) CreateScanDiplom(ctx context.Context, file []byte, filename string) error {

	var buffer bytes.Buffer
	_, err := buffer.Write(file)
	if err != nil {
		return err
	}

	err = s.s3client.SendDocumentToS3(ctx, buffer, filename)
	if err != nil {
		return err
	}

	return nil
}

func (s *ScanDiplomService) DownloadScanDiplom(ctx context.Context, param string) ([]byte, error) {

	scan, err := s.s3client.GetDocumentFromS3(ctx, param)
	if err != nil {
		return []byte{}, err
	}

	return scan, nil
}

func (s *ScanDiplomService) ExistsScanDiplom(ctx context.Context, filename string) ([]string, error) {

	documents, err := s.s3client.GetDocumentListFromS3(ctx, filename)
	if err != nil {
		return []string{}, err
	}

	if len(documents) > 0 {
		return documents, nil
	} else {
		return []string{}, errors.New("Не найдено ни одного документа")
	}
}

func (s *ScanDiplomService) DeleteScanDiplom(fileName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := s.s3client.DeleteDocumentFromS3(ctx, fileName)
	if err != nil {
		return err
	}

	return nil
}
