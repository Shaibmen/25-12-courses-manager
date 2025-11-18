package pg

import (
	"context"
	"fmt"
	"log/slog"
	"online-courses/internal/database"
	"os"
	"os/exec"
	"time"
)

type BackupRepo struct {
	repo       database.DB
	dbUser     string
	dbPassword string
	dbName     string
	host       string
	logger     *slog.Logger
}

var (
	pathBackup = "backups/pg/"
)

func NewBackupRepo(db database.DB, dbUser, dbPassword, dbName, host string, logger *slog.Logger) *BackupRepo {
	return &BackupRepo{repo: db, dbUser: dbUser, dbPassword: dbPassword, dbName: dbName, host: host, logger: logger}
}

func (b *BackupRepo) BackupDB(ctx context.Context) error {

	time := time.Now()

	backupFile := "backup-postgres"

	NameBackup := fmt.Sprintf("%v-%s.dump", backupFile, time.Format("2006-01-02_15-04-05"))

	filePath := pathBackup + NameBackup

	outfile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outfile.Close()

	b.logger.Info(
		"создание бекапа базы данных",
		"time", time,
		"name_file", NameBackup,
	)

	cmd := exec.CommandContext(ctx, "pg_dump", "-U", b.dbUser, "-h", b.host, "-F", "c", b.dbName)

	cmd.Stdout = outfile

	info, err := outfile.Stat()
	if err != nil {
		return err
	}

	cmd.Env = append(os.Environ(), "PGPASSWORD="+b.dbPassword)

	err = cmd.Run()
	if err != nil {

		b.logger.Error("ошибка восстанолвения бекапа",
			"name_file", NameBackup,
			"error", err.Error(),
		)

		return err
	}

	b.logger.Info("данные бекапа",
		"name", info.Name(),
		"size", info.Size(),
	)

	return nil

}

func (b *BackupRepo) RestoreDB(ctx context.Context, backupFile string) error {

	backupPath := pathBackup + backupFile

	b.logger.Info("восстановление базы данных",
		"name_file", backupFile,
		"time", time.Now().String(),
	)

	cmd := exec.CommandContext(ctx, "pg_restore", "-U", b.dbUser, "-h", b.host, "-d", b.dbName, "-c", backupPath)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+b.dbPassword)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {

		b.logger.Error("ошибка восстанолвения бекапа",
			"name_file", backupFile,
			"error", err.Error(),
		)

		return err
	}

	return nil
}

func (b *BackupRepo) AllBackup() ([]string, error) {

	file, err := os.ReadDir(pathBackup)
	if err != nil {
		return nil, err
	}

	var backups []string
	for _, files := range file {
		backups = append(backups, files.Name())
	}

	return backups, nil
}
