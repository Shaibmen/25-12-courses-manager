package service

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3ClientInterface interface {
	SendDocumentToS3(ctx context.Context, buffer bytes.Buffer, fileName string) error
	GetDocumentFromS3(ctx context.Context, fileName string) ([]byte, error)
}

type S3Client struct {
	client *minio.Client
}

func InitS3Client() S3Client {

	client, err := minio.New(os.Getenv("S3_ENDPOINT"), &minio.Options{
		Creds:  credentials.NewStaticV4(os.Getenv("S3_LOGIN"), os.Getenv("S3_PASSWORD"), ""),
		Secure: false, // поменять
	})

	if err != nil {
		panic(fmt.Sprintf("Произошла ошибка при запуске S3 хранилища: %s", err.Error()))
	}

	return S3Client{client}
}

func (c S3Client) SendDocumentToS3(ctx context.Context, buffer bytes.Buffer, fileName string) error {

	_, err := c.client.PutObject(ctx,
		os.Getenv("S3_BUCKET_NAME"),
		fileName,
		bytes.NewReader(buffer.Bytes()),
		int64(buffer.Len()),
		minio.PutObjectOptions{})

	if err != nil {
		return err
	}

	return nil
}

func (c S3Client) GetDocumentFromS3(ctx context.Context, fileName string) ([]byte, error) {

	object, err := c.client.GetObject(ctx,
		os.Getenv("S3_BUCKET_NAME"),
		fileName,
		minio.GetObjectOptions{})

	if err != nil {
		return []byte{}, err
	}

	readInfo, err := object.Stat()
	if err != nil {
		return []byte{}, err
	}

	buffer := make([]byte, readInfo.Size)

	_, err = object.Read(buffer)
	if err != nil {
		return []byte{}, err
	}

	return buffer, nil
}
