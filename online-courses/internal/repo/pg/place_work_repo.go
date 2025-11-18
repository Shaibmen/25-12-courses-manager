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

type PlaceWorkRepo struct {
	repo   database.DB
	logger *slog.Logger
}

func NewPlaceWorkRepo(db database.DB, logger *slog.Logger) *PlaceWorkRepo {
	return &PlaceWorkRepo{repo: db, logger: logger}
}

func (p *PlaceWorkRepo) CreateInTx(ctx context.Context, tx database.Tx, m entity.PlaceWork) error {
	query := `INSERT INTO placework(id_placework, name_company, job_title, all_experience, job_title_experience)
	          Values ($1, $2, $3, $4, $5)
			  `

	_, err := tx.ExecContext(ctx, query, m.ID_PlaceWork, m.NameCompany, m.JobTitle, m.AllExperience, m.JobTitleExperience)
	if err != nil {

		p.logger.Error("database error",
			"operation", "insert_place_work",
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (p *PlaceWorkRepo) DeleteInTx(ctx context.Context, tx database.Tx, id uuid.UUID) error {

	exists, err := repoutils.Exists(ctx, p.repo, "placework", "id_placework", id)
	if err != nil {

		p.logger.Debug("database error",
			"operation", "check_unique",
			"table", "placework",
			"row", "id_placework",
			"type", "exist",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	if !exists {
		return apperrors.ErrNoExists
	}

	query := "DELETE FROM placework WHERE id_placework = $1"

	if _, err := tx.ExecContext(ctx, query, id); err != nil {

		p.logger.Error("database error",
			"operation", "delete_place_work",
			"id_placework", id,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (p *PlaceWorkRepo) Read(ctx context.Context, id uuid.UUID) (*entity.PlaceWork, error) {

	query := `
	select name_company, job_title, all_experience, job_title_experience 
	from placework 
	where id_placework = $1`

	rows, err := p.repo.QueryContext(ctx, query, id)
	if err != nil {

		p.logger.Error("database error",
			"operation", "read_place_work",
			"id_place_work", id,
			"type", "query",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}
	defer rows.Close()

	var data entity.PlaceWork
	for rows.Next() {
		if err = rows.Scan(
			&data.NameCompany,
			&data.JobTitle,
			&data.AllExperience,
			&data.JobTitleExperience,
		); err != nil {

			p.logger.Error("database error",
				"operation", "read_mapping_place_work",
				"id_place_work", id,
				"type", "query",
				"err", err,
			)

			return nil, repoutils.HandleRepoErr(err)
		}

	}

	return &data, nil
}

func (p *PlaceWorkRepo) UpdateInTx(ctx context.Context, tx database.Tx, placework entity.PlaceWork) error {

	exists, err := repoutils.Exists(ctx, p.repo, "placework", "id_placework", placework.ID_PlaceWork)
	if err != nil {

		p.logger.Debug("database error",
			"operation", "check_unique",
			"table", "placework",
			"row", "id_placework",
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
	update placework
	set 
	name_company = coalesce($1, name_company),
	job_title = coalesce($2, job_title),
	all_experience = coalesce($3, all_experience),
	job_title_experience = coalesce($4, job_title_experience)
	where id_placework = $5;
	`

	if _, err := tx.ExecContext(
		ctx,
		query,
		placework.NameCompany,
		placework.JobTitle,
		placework.AllExperience,
		placework.JobTitleExperience,
		placework.ID_PlaceWork); err != nil {

		p.logger.Error("database error",
			"operation", "update_place_work",
			"id_placework", placework.ID_PlaceWork,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}
