package service

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3ClientInterface interface {
	SendDocumentToS3(ctx context.Context, buffer bytes.Buffer, fileName string) error
	GetDocumentFromS3(ctx context.Context, fileName string) ([]byte, error)
	GetDocumentListFromS3(ctx context.Context, searchString string) ([]string, error)
	DeleteDocumentFromS3(ctx context.Context, fileName string) error
}

type S3Client struct {
	client *minio.Client
}

func MustInitS3Client() S3Client {

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

	defer object.Close()

	
	readInfo, err := object.Stat()
	if err != nil {
		return []byte{}, err
	}
	// log.Printf("Object stat: %v\nObject size: %v\n", readInfo, readInfo.Size)

	buffer := make([]byte, readInfo.Size)

	// log.Println("Buffer size:", len(buffer))

	_, err = object.Read(buffer)
	if err != nil && err != io.EOF {
		return []byte{}, err
	}

	return buffer, nil
}

func (c S3Client) DeleteDocumentFromS3(ctx context.Context, fileName string) error {

	err := c.client.RemoveObject(ctx,
		os.Getenv("S3_BUCKET_NAME"),
		fileName,
		minio.RemoveObjectOptions{})

	if err != nil {
		return err
	}

	return nil
}

func (c S3Client) GetDocumentListFromS3(ctx context.Context, searchString string) ([]string, error) {
	
	result := make([]string, 0)

	objectList := c.client.ListObjects(ctx, os.Getenv("S3_BUCKET_NAME"), minio.ListObjectsOptions{})

	for object := range objectList {
		if err := object.Err; err != nil {
			return result, err
		}

		if strings.Contains(object.Key, searchString) {
			result = append(result, object.Key)
		}
	}

	return result, nil
}
