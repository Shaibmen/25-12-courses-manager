package service

import (
	"context"
	"online-courses/internal/domain/repository"
)

type BackupService struct {
	service repository.BackupRepository
}

func NewBackupService(service repository.BackupRepository) *BackupService {
	return &BackupService{service: service}
}

func (b *BackupService) BackupDB(ctx context.Context) error {

	err := b.service.BackupDB(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (b *BackupService) RestoreDB(ctx context.Context, backupFile string) error {

	err := b.service.RestoreDB(ctx, backupFile)
	if err != nil {
		return err
	}

	return nil
}

func (b *BackupService) AllBackup() ([]string, error) {
	data, err := b.service.AllBackup()
	if err != nil {
		return nil, err
	}

	return data, nil
}
