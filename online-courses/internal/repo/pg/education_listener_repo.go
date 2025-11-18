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

type EducationListenerRepo struct {
	repo   database.DB
	logger *slog.Logger
}

func NewEducationListenerRepo(db database.DB, logger *slog.Logger) *EducationListenerRepo {
	return &EducationListenerRepo{repo: db, logger: logger}
}

func (e *EducationListenerRepo) CreateInTx(ctx context.Context, tx database.Tx, m entity.EducationListener) error {
	query := `INSERT INTO educationlistener (id_educationlistener, diplom_seria, diplom_number, date_given, city, region, educational_institution, speciality, level_education)
	          VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := tx.ExecContext(ctx, query, m.ID_EducationListener, m.DiplomSeria, m.DiplomNumber, m.DateGiven, m.City, m.Region, m.EducationalInstitution, m.Speciality, m.LevelEducation)
	if err != nil {

		e.logger.Error("database error",
			"operation", "insert_education_listener",
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (e *EducationListenerRepo) DeleteInTx(ctx context.Context, tx database.Tx, id uuid.UUID) error {

	exists, err := repoutils.Exists(ctx, e.repo, "educationlistener", "id_educationlistener", id)
	if err != nil {

		e.logger.Debug("database error",
			"operation", "check_unique",
			"table", "educationlistener",
			"row", "id_educationlistener",
			"type", "exist",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	if !exists {
		return apperrors.ErrNoExists
	}

	query := "DELETE FROM educationlistener WHERE id_educationlistener = $1"

	if _, err := tx.ExecContext(ctx, query, id); err != nil {

		e.logger.Error("database error",
			"operation", "delete_education_listener",
			"id_educaitonlisener", id,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (e *EducationListenerRepo) Read(ctx context.Context, id uuid.UUID) (*entity.EducationListener, error) {
	query := `
	select diplom_seria, diplom_number, date_given, city, region, educational_institution, speciality, lvl.education
	from educationlistener as e
	left join leveleducation lvl on e.level_education  = lvl.id_leveleducation
	where e.id_educationlistener = $1;`

	rows, err := e.repo.QueryContext(ctx, query, id)
	if err != nil {

		e.logger.Error("database error",
			"operation", "read_education_listener",
			"id_educationlistener", id,
			"type", "query",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}
	defer rows.Close()

	var education entity.EducationListener
	for rows.Next() {
		if err = rows.Scan(
			&education.DiplomSeria,
			&education.DiplomNumber,
			&education.DateGiven,
			&education.City,
			&education.Region,
			&education.EducationalInstitution,
			&education.Speciality,
			&education.LevelEducation,
		); err != nil {

			e.logger.Error("database error",
				"operation", "read_mapping_education_listener",
				"id_educationlistener", id,
				"type", "query",
				"err", err,
			)

			return nil, repoutils.HandleRepoErr(err)
		}
	}

	return &education, nil
}

func (e *EducationListenerRepo) UpdateInTx(ctx context.Context, tx database.Tx, education entity.EducationListener) error {

	exists, err := repoutils.Exists(ctx, e.repo, "educationlistener", "id_educationlistener", education.ID_EducationListener)
	if err != nil {

		e.logger.Debug("database error",
			"operation", "check_unique",
			"table", "educationlistener",
			"row", "id_educationlistener",
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
	update educationlistener
	set 
	diplom_seria = coalesce($1, diplom_seria),
	diplom_number = coalesce($2, diplom_number),
	date_Given = coalesce($3, date_Given),
	city = coalesce($4, city),
	region = coalesce($5, region),
	educational_institution = coalesce($6, educational_institution),
	speciality = coalesce($7, speciality),
	level_education = coalesce($8, level_education)
	where id_educationlistener = $9;
	`

	if _, err := tx.ExecContext(ctx, query,
		education.DiplomSeria,
		education.DiplomNumber,
		education.DateGiven,
		education.City,
		education.Region,
		education.EducationalInstitution,
		education.Speciality,
		education.LevelEducation,
		education.ID_EducationListener); err != nil {

		e.logger.Error("database error",
			"operation", "update_education_listener",
			"id_educationlistener", education.ID_EducationListener,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}
