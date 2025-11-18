package repository

import "context"

type BackupRepository interface {
	BackupDB(ctx context.Context) error
	RestoreDB(ctx context.Context, backupFile string) error
	AllBackup() ([]string, error)
}
