package postgres

import (
	"context"
	"database/sql"
	"log/slog"
	"online-courses/internal/database"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresSQL struct {
	DBX *sqlx.DB
}

func MustNewConnectionPostgresSQL(cfg string, logger *slog.Logger) *PostgresSQL {

	connx, err := sqlx.Open("postgres", cfg)
	if err != nil {
		logger.Error("ошибка подключения к базе данных x",
			"error", err,
		)
		panic(1)
	}

	if err = connx.Ping(); err != nil {
		logger.Error("ошибка при пинге бд",
			"error", err,
		)
		panic(1)
	}

	logger.Info("база данных подключена")
	return &PostgresSQL{DBX: connx}
}

func (p *PostgresSQL) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {

	rows, err := p.DBX.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil

}

func (p *PostgresSQL) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	return p.DBX.GetContext(ctx, dest, query, args...)
}

func (p *PostgresSQL) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	return p.DBX.SelectContext(ctx, dest, query, args...)
}

func (p *PostgresSQL) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	rows, err := p.DBX.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (p *PostgresSQL) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	row := p.DBX.QueryRowContext(ctx, query, args...)
	return row
}

func (p *PostgresSQL) Ping() error {
	return p.DBX.Ping()
}

func (p *PostgresSQL) Close() error {
	return p.DBX.Close()
}

func (p *PostgresSQL) BeginTx(ctx context.Context) (database.Tx, error) {
	tx, err := p.DBX.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &postgresTx{Tx: tx}, nil
}

type postgresTx struct {
	Tx *sql.Tx
}

func (t *postgresTx) Commit() error {
	return t.Tx.Commit()
}

func (t *postgresTx) Rollback() error {
	return t.Tx.Rollback()
}

func (t *postgresTx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return t.Tx.ExecContext(ctx, query, args...)
}

func (t *postgresTx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return t.Tx.QueryContext(ctx, query, args...)
}
