package pg

import (
	"context"
	"log/slog"
	"online-courses/internal/apperrors"
	"online-courses/internal/database"
	"online-courses/internal/domain/entity"
	repoutils "online-courses/internal/repo/pg/repo_utils"

	"github.com/google/uuid"
)

type ExecutorRepo struct {
	repo   database.DB
	logger *slog.Logger
}

func NewExecutorRepo(repo database.DB, logger *slog.Logger) *ExecutorRepo {
	return &ExecutorRepo{repo: repo, logger: logger}
}

func (e *ExecutorRepo) Create(ctx context.Context, m entity.Executor) error {
	query :=
		`
	insert into executor (id_executor, status, first_name, second_name, middle_name, doverenost)
	values ($1, $2, $3, $4, $5, $6)
	`

	_, err := e.repo.ExecContext(ctx, query, m.ID_Executor, m.Status, m.FirstName, m.SecondName, m.MiddleName, m.Doverenost)
	if err != nil {
		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (e *ExecutorRepo) Read(ctx context.Context, filter string) ([]entity.Executor, error) {
	query :=
		`
	select * from executor
	where ($1::text is null or second_name ilike '%' || $1::text || '%' )
	`

	rows, err := e.repo.QueryContext(ctx, query, filter)
	if err != nil {

		e.logger.Error("database error",
			"operation", "read_executor",
			"type", "query",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}
	defer rows.Close()

	var executors []entity.Executor
	for rows.Next() {
		var l entity.Executor
		if err := rows.Scan(
			&l.ID_Executor,
			&l.Status,
			&l.FirstName,
			&l.SecondName,
			&l.MiddleName,
			&l.Doverenost,
		); err != nil {

			e.logger.Error("database error",
				"operation", "read_mapping_executor",
				"type", "query",
				"err", err,
			)

			return nil, repoutils.HandleRepoErr(err)
		}
		executors = append(executors, l)
	}

	return executors, nil
}

func (e *ExecutorRepo) Delete(ctx context.Context, id uuid.UUID) error {
	exists, err := repoutils.Exists(ctx, e.repo, "executor", "id_executor", id)
	if err != nil {

		e.logger.Debug("database error",
			"operation", "check_unique",
			"table", "executor",
			"row", "id_executor",
			"type", "exist",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	if !exists {
		return apperrors.ErrNoExists
	}

	query :=
		`
	delete from executor where id_executor = $1
	`

	if _, err := e.repo.ExecContext(ctx, query, id); err != nil {

		e.logger.Error("database error",
			"operation", "delete_executor",
			"id_executor", id,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil

}
