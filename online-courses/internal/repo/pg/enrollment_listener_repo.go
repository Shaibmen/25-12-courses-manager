package pg

import (
	"context"
	"database/sql"
	"log/slog"
	"online-courses/internal/apperrors"
	"online-courses/internal/database"
	"online-courses/internal/domain/entity"
	repoutils "online-courses/internal/repo/pg/repo_utils"

	"github.com/google/uuid"
)

type enrollmentListenerRepo struct {
	repo   database.DB
	logger *slog.Logger
}

func NewEnrollmentListenerRepo(db database.DB, logger *slog.Logger) *enrollmentListenerRepo {
	return &enrollmentListenerRepo{repo: db, logger: logger}
}

func (e *enrollmentListenerRepo) Create(ctx context.Context, model entity.EnrollmentListener) error {

	query :=
		`
	insert into enrollmentlistener (id_listener, id_programeducation, start_date, end_date, is_active, type_of_retraining, id_group)
	values ($1, $2, $3, $4, $5, $6, $7)
	`
	if _, err := e.repo.ExecContext(ctx,
		query,
		model.ID_Listener,
		model.ID_ProgramEducation,
		model.StartDate,
		model.EndDate,
		model.Is_active,
		model.TypeOfRetraining,
		model.ID_Group,
	); err != nil {

		e.logger.Error("database error",
			"operation", "insert_enrollment_listener",
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (e *enrollmentListenerRepo) Read(ctx context.Context, page int, filter string) ([]entity.EnrollmentListenerDetails, error) {

	limit, offset := repoutils.Pagination(page)

	query :=
		`
		select
	e.id_listener,	 
	l.first_name, l.second_name, l.middle_name,
	p.name_prof_education,
	e.start_date, e.end_date, e.type_of_retraining, e.id_group
	from enrollmentlistener as e
	inner join listener l on e.id_listener = l.id_listener
	inner join programeducation p on e.id_programeducation = p.id_programeducation
	where ($1::text is null or l.second_name ilike '%' || $1::text || '%' ) and e.is_active = TRUE
	limit $2
	offset $3
	`

	rows, err := e.repo.QueryContext(ctx, query, filter, limit, offset)
	if err != nil {

		e.logger.Error("database error",
			"operation", "read_enrollment_listener",
			"page", page,
			"type", "query",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}
	defer rows.Close()

	var enrollments []entity.EnrollmentListenerDetails

	for rows.Next() {
		var enrollment entity.EnrollmentListenerDetails
		if err = rows.Scan(
			&enrollment.ID_Listener,
			&enrollment.FirstName,
			&enrollment.SecondName,
			&enrollment.MiddleName,
			&enrollment.NameProfEducation,
			&enrollment.StartDate,
			&enrollment.EndDate,
			&enrollment.TypeOfRetraining,
			&enrollment.ID_Group,
		); err != nil {

			e.logger.Error("database error",
				"operation", "read_mapping_enrollment_listener",
				"page", page,
				"type", "query",
				"err", err,
			)

			return nil, repoutils.HandleRepoErr(err)
		}
		enrollments = append(enrollments, enrollment)
	}

	return enrollments, nil
}

func (e *enrollmentListenerRepo) Update(ctx context.Context, idListener, idProgram uuid.UUID, model entity.EnrollmentListener) error {

	exists, err := repoutils.Exists(ctx, e.repo, "enrollmentlistener", "id_listener", idListener)
	if err != nil {

		e.logger.Debug("database error",
			"operation", "check_unique",
			"table", "enrollmentlistener",
			"row", "id_listener",
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
	update enrollmentlistener 
	set 
	id_programeducation = coalesce($1, id_programeducation),
	start_date = coalesce($2, start_date),
	end_date = coalesce($3, end_date),
	id_group = coalesce($4, group_number),
	type_of_retraining = coalesce($5, type_of_retraining)
	where id_listener = $6 and id_programeducation = $7
	`
	if _, err := e.repo.ExecContext(
		ctx,
		query,
		model.ID_ProgramEducation,
		model.StartDate,
		model.EndDate,
		model.ID_Group,
		model.TypeOfRetraining,
		idListener,
		idProgram,
	); err != nil {

		e.logger.Error("database error",
			"operation", "update_enrollment_listener",
			"id_listener", idListener,
			"id_program", idProgram,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil

}

func (e *enrollmentListenerRepo) Delete(ctx context.Context, ListenerID, ProgramID uuid.UUID) error {

	queryExists := `select exists (select 1 from enrollmentlistener where id_listener = $1 and id_programeducation = $2)`

	var exists bool

	err := e.repo.QueryRowContext(ctx, queryExists, ListenerID, ProgramID).Scan(&exists)
	if err != nil {

		e.logger.Debug("database error",
			"operation", "chech_exists",
			"table", "enrollmentlistener",
			"row", "id_listener",
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
	delete from enrollmentlistener
	WHERE id_listener = $1 AND id_programeducation = $2;
	`

	if _, err := e.repo.ExecContext(ctx, query, ListenerID, ProgramID); err != nil {

		e.logger.Error("database error",
			"operation", "delete_enrollment_listener",
			"id_listener", ListenerID,
			"id_program", ProgramID,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (e *enrollmentListenerRepo) ReadDetailListener(ctx context.Context, id uuid.UUID) ([]entity.EnrollmentProgramDetails, error) {

	query :=
		`
	select 
	e.id_listener, e.id_programeducation, e.type_of_retraining, e.id_group,
	p.name_prof_education, p.time_education, p.price, educ.type_name , d.divisions,
	e.start_date, e.end_date
	from enrollmentlistener as e
	inner join programeducation p on e.id_programeducation = p.id_programeducation
	inner join educationtypes educ on p.id_educationtype = educ.id_educationtype
	inner join divisionseducation d  on p.id_divisionseducation = d.id_divisionseducation
	where id_listener = $1 and e.is_active = TRUE
	`

	rows, err := e.repo.QueryContext(ctx, query, id)
	if err != nil {

		e.logger.Error("database error",
			"operation", "read_enrollment_listener",
			"id_listener", id,
			"type", "query",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}
	defer rows.Close()

	var enrollments []entity.EnrollmentProgramDetails

	for rows.Next() {
		enrollment := entity.EnrollmentProgramDetails{}
		if err = rows.Scan(
			&enrollment.ID_Listener,
			&enrollment.ID_ProgramEducation,
			&enrollment.TypeOfRetraining,
			&enrollment.ID_Group,
			&enrollment.NameProfEducation,
			&enrollment.TimeEducation,
			&enrollment.Price,
			&enrollment.EducationType,
			&enrollment.DivisionEducation,
			&enrollment.StartDate,
			&enrollment.EndDate,
		); err != nil {

			e.logger.Error("database error",
				"operation", "read_mapping_enrollment_listener",
				"id_listener", id,
				"type", "query",
				"err", err,
			)

			return nil, repoutils.HandleRepoErr(err)
		}
		enrollments = append(enrollments, enrollment)
	}

	if len(enrollments) == 0 {
		return nil, repoutils.HandleRepoErr(sql.ErrNoRows)
	}

	return enrollments, nil
}

func (e *enrollmentListenerRepo) ReadByProgram(ctx context.Context, id uuid.UUID, page int) ([]entity.EnrollmentListenerDetails, error) {

	exists, err := repoutils.Exists(ctx, e.repo, "enrollmentlistener", "id_programeducation", id)
	if err != nil {

		e.logger.Debug("database error",
			"operation", "check_unique",
			"table", "enrollmentlistener",
			"row", "id_listener",
			"type", "exist",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}

	if !exists {
		return nil, apperrors.ErrNoExists
	}

	limit, offset := repoutils.Pagination(page)

	query :=
		`
	select
	e.id_listener,	 
	l.first_name, l.second_name, l.middle_name,
	p.name_prof_education,
	e.start_date, e.end_date, e.type_of_retraining, e.id_group
	from enrollmentlistener as e
	inner join listener l on e.id_listener = l.id_listener
	inner join programeducation p on e.id_programeducation = p.id_programeducation
	where p.id_programeducation = $1 and e.is_active = TRUE
	limit $2
	offset $3

	`

	rows, err := e.repo.QueryContext(ctx, query, id, limit, offset)
	if err != nil {
		return nil, repoutils.HandleRepoErr(err)
	}
	defer rows.Close()

	var enrollments []entity.EnrollmentListenerDetails

	for rows.Next() {
		var enrollment entity.EnrollmentListenerDetails
		if err = rows.Scan(
			&enrollment.ID_Listener,
			&enrollment.FirstName,
			&enrollment.SecondName,
			&enrollment.MiddleName,
			&enrollment.NameProfEducation,
			&enrollment.StartDate,
			&enrollment.EndDate,
			&enrollment.TypeOfRetraining,
			&enrollment.ID_Group,
		); err != nil {

			e.logger.Error("database error",
				"operation", "read_mapping_enrollment_listener",
				"id_program", id,
				"type", "query",
				"err", err,
			)

			return nil, repoutils.HandleRepoErr(err)
		}
		enrollments = append(enrollments, enrollment)
	}

	return enrollments, nil
}

func (e *enrollmentListenerRepo) SaveInfo(ctx context.Context, m entity.AccurateProgram) error {

	query :=

		`
	insert into accurateprogram (id_listener, name_prof_education, time_education, price, educationtype, divisionseducation)
	values ($1, $2, $3, $4, $5, $6)
	`

	_, err := e.repo.ExecContext(ctx, query, m.ID_Listener, m.NameProfEducation, m.TimeEducation, m.Price, m.EducationType, m.Division)
	if err != nil {
		return repoutils.HandleRepoErr(err)
	}

	return nil

}

func (e *enrollmentListenerRepo) GetAccurateEnrollment(ctx context.Context) ([]entity.AccurateProgram, error) {
	program := []entity.AccurateProgram{}
	query :=
		`
	select * from accurateprogram
	`
	err := e.repo.SelectContext(ctx, &program, query)
	if err != nil {
		return nil, repoutils.HandleRepoErr(err)
	}

	return program, nil

}

func (e *enrollmentListenerRepo) GetProgram(ctx context.Context, id uuid.UUID) (*entity.ProgramToAccurate, error) {

	query :=
		`
	
	select 
	p.name_prof_education , p.time_education, p.price, e.type_name, d.divisions
	from programeducation as p
	inner join educationtypes e on p.id_educationtype = e.id_educationtype
	inner join divisionseducation d on p.id_divisionseducation = d.id_divisionseducation 
	where p.id_programeducation = $1
	`
	program := entity.ProgramToAccurate{}
	err := e.repo.GetContext(ctx, &program, query, id)
	if err != nil {
		return nil, repoutils.HandleRepoErr(err)
	}

	return &program, err
}

func (e *enrollmentListenerRepo) GetListenerFIO(ctx context.Context, listenerID uuid.UUID) (*entity.ListenerFIO, error) {
	query := `
	select
	l.first_name, l.second_name, l.middle_name
	from enrollmentlistener as e
	left join listener l on e.id_listener = l.id_listener
	where e.id_listener = $1;
	`

	rows, err := e.repo.QueryContext(ctx, query, listenerID)
	if err != nil {
		e.logger.Debug("database error",
			"operation", "query",
			"table", "enrollmentlistener",
			"id_listener", listenerID,
			"type", "exist",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}
	defer rows.Close()

	var fio entity.ListenerFIO

	for rows.Next() {
		if err = rows.Scan(
			&fio.FirstName,
			&fio.SecondName,
			&fio.MiddleName,
		); err != nil {

			e.logger.Debug("database error",
				"operation", "query",
				"table", "enrollmentlistener",
				"id_listener", listenerID,
				"type", "exist",
				"err", err,
			)

			return nil, repoutils.HandleRepoErr(err)
		}
	}

	return &fio, nil
}
