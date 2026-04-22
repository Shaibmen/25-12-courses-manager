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

type ProgramEducationRepo struct {
	repo   database.DB
	logger *slog.Logger
}

func NewProgramEducationRepo(db database.DB, logger *slog.Logger) *ProgramEducationRepo {
	return &ProgramEducationRepo{repo: db, logger: logger}
}

func (p *ProgramEducationRepo) Create(ctx context.Context, model entity.ProgramEducation) error {

	query :=
		`
	insert into programeducation (id_programeducation, name_prof_education, time_education, price, id_educationtype, id_divisionseducation)
	values ($1, $2, $3, $4, $5, $6)
	`
	if _, err := p.repo.ExecContext(
		ctx,
		query,
		model.ID_ProgramEducation,
		model.NameProfEducation,
		model.TimeEducation,
		model.Price,
		model.ID_EducationType,
		model.ID_DivisionsEducation); err != nil {

		p.logger.Error("database error",
			"operation", "insert_program_education",
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (p *ProgramEducationRepo) Read(ctx context.Context, page int, filter string) ([]entity.ProgramEducation, error) {

	limit, offset := repoutils.Pagination(page)

	query :=
		`
	select * from programeducation
	where ($1::text is null or name_prof_education ilike '%' || $1::text || '%' )
	limit $2 offset $3
	`

	rows, err := p.repo.QueryContext(ctx, query, filter, limit, offset)
	if err != nil {

		p.logger.Error("database error",
			"operation", "read_program_education",
			"page", page,
			"type", "query",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}
	defer rows.Close()

	var programs []entity.ProgramEducation

	for rows.Next() {
		var program entity.ProgramEducation
		if err := rows.Scan(
			&program.ID_ProgramEducation,
			&program.NameProfEducation,
			&program.TimeEducation,
			&program.Price,
			&program.ID_EducationType,
			&program.ID_DivisionsEducation,
		); err != nil {

			p.logger.Error("database error",
				"operation", "read_mapping_program_education",
				"page", page,
				"type", "query",
				"err", err,
			)

			return nil, err
		}
		programs = append(programs, program)
	}

	return programs, nil
}

func (p *ProgramEducationRepo) ReadByID(ctx context.Context, id uuid.UUID) (*entity.ProgramEducation, error) {

	query :=
		`
	select * from programeducation
	where id_programeducation  = $1
	`

	rows, err := p.repo.QueryContext(ctx, query, id)
	if err != nil {

		p.logger.Error("database error",
			"operation", "read_program_education",
			"type", "query",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}
	defer rows.Close()

	var programs entity.ProgramEducation

	for rows.Next() {
		if err := rows.Scan(
			&programs.ID_ProgramEducation,
			&programs.NameProfEducation,
			&programs.TimeEducation,
			&programs.Price,
			&programs.ID_EducationType,
			&programs.ID_DivisionsEducation,
		); err != nil {

			p.logger.Error("database error",
				"operation", "read_by_name_prof_mapping_program_education",
				"type", "query",
				"err", err,
			)

			return nil, err
		}
	}

	return &programs, nil
}

func (p *ProgramEducationRepo) Update(ctx context.Context, model entity.ProgramEducation) error {

	exists, err := repoutils.Exists(ctx, p.repo, "programeducation", "id_programeducation", model.ID_ProgramEducation)
	if err != nil {

		p.logger.Debug("database error",
			"operation", "check_unique",
			"table", "programeducation",
			"row", "id_programeducation",
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
	update programeducation
	set 
	name_prof_education = coalesce($1, name_prof_education),
	time_education = coalesce($2, time_education),
	price = coalesce($3, price),
	id_educationtype = coalesce($4, id_educationtype),
	id_divisionseducation = coalesce($5, id_divisionseducation)
	where id_programeducation = $6;
	`

	if _, err := p.repo.ExecContext(
		ctx,
		query,
		model.NameProfEducation,
		model.TimeEducation,
		model.Price,
		model.ID_EducationType,
		model.ID_DivisionsEducation,
		model.ID_ProgramEducation,
	); err != nil {

		p.logger.Error("database error",
			"operation", "update_program_education",
			"id_programeducation", model.ID_ProgramEducation,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (p *ProgramEducationRepo) Delete(ctx context.Context, id uuid.UUID) error {

	exists, err := repoutils.Exists(ctx, p.repo, "programeducation", "id_programeducation", id)
	if err != nil {

		p.logger.Debug("database error",
			"operation", "check_unique",
			"table", "programeducation",
			"row", "id_programeducation",
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
	delete from programeducation where id_programeducation = $1
	`

	if _, err := p.repo.ExecContext(ctx, query, id); err != nil {

		p.logger.Error("database error",
			"operation", "delete_program_education",
			"id_programeducation", id,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}
	return nil
}
