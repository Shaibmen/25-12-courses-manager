package service

import "context"

type BackupService interface {
	BackupDB(ctx context.Context) error
	RestoreDB(ctx context.Context, backupFile string) error
	AllBackup() ([]string, error)
}
