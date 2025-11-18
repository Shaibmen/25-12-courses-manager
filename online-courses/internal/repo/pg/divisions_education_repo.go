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

type DivisionsEducationRepo struct {
	repo   database.DB
	logger *slog.Logger
}

func NewDivisionsEducationRepo(db database.DB, logger *slog.Logger) *DivisionsEducationRepo {
	return &DivisionsEducationRepo{repo: db, logger: logger}
}

func (d *DivisionsEducationRepo) Create(ctx context.Context, divisions entity.DivisionsEducation) error {
	query :=
		`
	insert into divisionseducation(id_divisionseducation, divisions)
	values ($1, $2);
	`

	if _, err := d.repo.ExecContext(ctx, query, divisions.ID_DivisionsEducation, divisions.Divisions); err != nil {

		d.logger.Error("database error",
			"operation", "insert_divisions_education",
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (d *DivisionsEducationRepo) Read(ctx context.Context, filter string) ([]entity.DivisionsEducation, error) {

	query :=
		`
	select id_divisionseducation, divisions
	from divisionseducation
	where ($1::text is null or divisions ilike '%' || $1::text || '%' )
	`

	rows, err := d.repo.QueryContext(ctx, query, filter)
	if err != nil {

		d.logger.Error("database error",
			"operation", "read_divisions_education",
			"type", "query",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}
	defer rows.Close()

	var divisions []entity.DivisionsEducation
	for rows.Next() {
		var list entity.DivisionsEducation
		if err = rows.Scan(
			&list.ID_DivisionsEducation,
			&list.Divisions,
		); err != nil {

			d.logger.Error("database error",
				"operation", "read_mapping_divisions_education",
				"query", query,
				"type", "query",
				"err", err,
			)

			return nil, repoutils.HandleRepoErr(err)
		}
		divisions = append(divisions, list)
	}

	return divisions, nil
}

// func (d *DivisionsEducationRepo) ReadByID(ctx context.Context, id uuid.UUID) (*entity.DivisionsEducation, error) {

// 	query :=
// 		`
// 	select id_divisionseducation, divisions
// 	from divisionseducation
// 	where id_divisionseducation = $1
// 	`

// 	rows, err := d.repo.QueryContext(ctx, query, id)
// 	if err != nil {

// 		d.logger.Error("database error",
// 			"operation", "read_divisions_education",
// 			"type", "query",
// 			"err", err,
// 		)

// 		return nil, repoutils.HandleRepoErr(err)
// 	}
// 	defer rows.Close()

// 	var divisions entity.DivisionsEducation
// 	for rows.Next() {
// 		if err = rows.Scan(
// 			&divisions.ID_DivisionsEducation,
// 			&divisions.Divisions,
// 		); err != nil {

// 			d.logger.Error("database error",
// 				"operation", "read_mapping_divisions_education",
// 				"query", query,
// 				"type", "query",
// 				"err", err,
// 			)

// 			return nil, repoutils.HandleRepoErr(err)
// 		}
// 	}

// 	return &divisions, nil
// }

func (d *DivisionsEducationRepo) Update(ctx context.Context, divisions entity.DivisionsEducation) error {

	exists, err := repoutils.Exists(ctx, d.repo, "divisionseducation", "id_divisionseducation", divisions.ID_DivisionsEducation)
	if err != nil {
		d.logger.Debug("database error",
			"operation", "check_unique",
			"table", "divisionseducation",
			"row", "id_divisionseducation",
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
	update divisionseducation
	set
	divisions = coalesce($1, divisions)
	where id_divisionseducation = $2;
	`

	if _, err := d.repo.ExecContext(ctx, query, divisions.Divisions, divisions.ID_DivisionsEducation); err != nil {

		d.logger.Error("database error",
			"operation", "update_divisions_education",
			"id_divisionseducation", divisions.ID_DivisionsEducation,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (d *DivisionsEducationRepo) Delete(ctx context.Context, id uuid.UUID) error {

	exists, err := repoutils.Exists(ctx, d.repo, "divisionseducation", "id_divisionseducation", id)
	if err != nil {

		d.logger.Debug("database error",
			"operation", "check_unique",
			"table", "divisionseducation",
			"row", "id_divisionseducation",
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
	delete from divisionseducation where id_divisionseducation = $1;
	`

	if _, err := d.repo.ExecContext(ctx, query, id); err != nil {

		d.logger.Error("database error",
			"operation", "delete_divisions_eucation",
			"id_divisionseducation", id,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}
