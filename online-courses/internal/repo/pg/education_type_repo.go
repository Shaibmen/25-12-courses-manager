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

type EducationTypeRepo struct {
	repo   database.DB
	logger *slog.Logger
}

func NewEducationTypeRepo(db database.DB, logger *slog.Logger) *EducationTypeRepo {
	return &EducationTypeRepo{repo: db, logger: logger}
}

func (e *EducationTypeRepo) Create(ctx context.Context, model entity.EducationTypes) error {
	query :=
		`
	insert into educationtypes (id_educationtype, type_name)
	values
	($1, $2)
	`

	if _, err := e.repo.ExecContext(ctx, query, model.ID_EducationType, model.TypeName); err != nil {

		e.logger.Error("database error",
			"operation", "insert_education_types",
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (e *EducationTypeRepo) Read(ctx context.Context, filter string) ([]entity.EducationTypes, error) {

	query :=
		`
	select id_educationtype, type_name
	from educationtypes
	where ($1::text is null or type_name ilike '%' || $1::text || '%' )
	`
	rows, err := e.repo.QueryContext(ctx, query, filter)
	if err != nil {

		e.logger.Error("database error",
			"operation", "read_education_types",
			"type", "query",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}
	defer rows.Close()

	var types []entity.EducationTypes
	for rows.Next() {
		var list entity.EducationTypes
		if err = rows.Scan(
			&list.ID_EducationType,
			&list.TypeName,
		); err != nil {

			e.logger.Error("database error",
				"operation", "read_mapping_education_types",
				"type", "query",
				"err", err,
			)

			return nil, repoutils.HandleRepoErr(err)
		}
		types = append(types, list)
	}

	return types, nil
}

// func (e *EducationTypeRepo) ReadByID(ctx context.Context, id uuid.UUID) (*entity.EducationTypes, error) {

// 	query :=
// 		`
// 	select id_educationtype, type_name
// 	from educationtypes
// 	where id_educationtype = $1
// 	`
// 	rows, err := e.repo.QueryContext(ctx, query, id)
// 	if err != nil {

// 		e.logger.Error("database error",
// 			"operation", "read_education_types",
// 			"type", "query",
// 			"err", err,
// 		)

// 		return nil, repoutils.HandleRepoErr(err)
// 	}
// 	defer rows.Close()

// 	var types entity.EducationTypes
// 	for rows.Next() {
// 		if err = rows.Scan(
// 			&types.ID_EducationType,
// 			&types.TypeName,
// 		); err != nil {

// 			e.logger.Error("database error",
// 				"operation", "read_mapping_education_types",
// 				"type", "query",
// 				"err", err,
// 			)

// 			return nil, repoutils.HandleRepoErr(err)
// 		}
// 	}

// 	return &types, nil
// }

func (e *EducationTypeRepo) Update(ctx context.Context, model entity.EducationTypes) error {

	exists, err := repoutils.Exists(ctx, e.repo, "educationtypes", "id_educationtype", model.ID_EducationType)
	if err != nil {

		e.logger.Debug("database error",
			"operation", "check_unique",
			"table", "educationtypes",
			"row", "id_educationtype",
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
	update educationtypes
	set
	type_name = coalesce($1, type_name)
	where id_educationtype = $2
	`

	if _, err := e.repo.ExecContext(ctx, query, model.TypeName, model.ID_EducationType); err != nil {

		e.logger.Error("database error",
			"operation", "update_education_types",
			"id_educationlistener", model.ID_EducationType,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (e *EducationTypeRepo) Delete(ctx context.Context, id uuid.UUID) error {

	exists, err := repoutils.Exists(ctx, e.repo, "educationtypes", "id_educationtype", id)
	if err != nil {

		e.logger.Debug("database error",
			"operation", "check_unique",
			"table", "educationtypes",
			"row", "id_educationtype",
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
	delete from educationtypes where id_educationtype = $1
	`

	if _, err := e.repo.ExecContext(ctx, query, id); err != nil {

		e.logger.Error("database error",
			"operation", "delete_education_type",
			"id_educationtype", id,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}
