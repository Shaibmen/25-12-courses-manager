package pg

import (
	"context"
	"database/sql"
	"log/slog"
	"online-courses/internal/apperrors"
	"online-courses/internal/database"
	"online-courses/internal/domain/dto"
	"online-courses/internal/domain/entity"
	repoutils "online-courses/internal/repo/pg/repo_utils"

	"github.com/google/uuid"
)

type ListenerRepo struct {
	repo   database.DB
	logger *slog.Logger
}

func NewListenerRepo(db database.DB, logger *slog.Logger) *ListenerRepo {
	return &ListenerRepo{repo: db, logger: logger}
}

func (l *ListenerRepo) CreateInTx(ctx context.Context, tx database.Tx, m *entity.Listener, i dto.ListenerIDDTO) error {
	query := `INSERT INTO listener (id_listener, first_name, second_name, middle_name, date_of_birth, snils, contact_phone, email, id_passport, id_regaddress, id_educationlistener, id_placework)
	          Values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	_, err := tx.ExecContext(ctx, query, m.ID_Listener, m.FirstName, m.SecondName, m.MiddleName, m.DateOfBirth, m.SNILS, m.ContactPhone, m.Email, i.ID_Passport, i.ID_RegAddress, i.ID_EducationListener, i.ID_PlaceWork)
	if err != nil {
		l.logger.Error("database error",
			"operation", "insert_listener",
			"type", "exec",
			"err", err,
		)
		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (l *ListenerRepo) ReadListener(ctx context.Context, page int, filter string) ([]entity.Listener, error) {
	query := `
	SELECT id_listener, first_name, second_name, middle_name, date_of_birth, snils, contact_phone, email
	FROM listener
	where ($1::text is null or second_name ilike '%' || $1::text || '%' )
	LIMIT $2 OFFSET $3`

	limit, offset := repoutils.Pagination(page)

	rows, err := l.repo.QueryContext(ctx, query, filter, limit, offset)
	if err != nil {

		l.logger.Error("database error",
			"operation", "read_listener",
			"page", page,
			"type", "query",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}
	defer rows.Close()

	var listeners []entity.Listener

	for rows.Next() {
		var list entity.Listener

		if err = rows.Scan(
			&list.ID_Listener,
			&list.FirstName,
			&list.SecondName,
			&list.MiddleName,
			&list.DateOfBirth,
			&list.SNILS,
			&list.ContactPhone,
			&list.Email,
		); err != nil {

			l.logger.Error("database error",
				"operation", "read_mapping_listener",
				"page", page,
				"type", "query",
				"err", err,
			)

			return nil, repoutils.HandleRepoErr(err)
		}
		listeners = append(listeners, list)
	}

	return listeners, nil
}

func (l *ListenerRepo) ReadFullData(ctx context.Context, id uuid.UUID) (*entity.Listener, error) {

	exists, err := repoutils.Exists(ctx, l.repo, "listener", "id_listener", id)
	if err != nil {
		return nil, repoutils.HandleRepoErr(err)
	}

	if !exists {
		return nil, repoutils.HandleRepoErr(sql.ErrNoRows)
	}

	query := `
	select
	l.*, 
	r.mail_index, r.region, r.city, r.street, r.house, r.building, r.apartment
	from listener as l 
	inner join registrationaddress r on l.id_regaddress = r.id_regaddress

	where id_listener = $1;`

	rows, err := l.repo.QueryContext(ctx, query, id)
	if err != nil {

		l.logger.Error("database error",
			"operation", "read_full_listener",
			"id_listener", id,
			"type", "query",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}
	defer rows.Close()

	var data entity.Listener

	for rows.Next() {

		if err = rows.Scan(
			&data.ID_Listener,
			&data.FirstName,
			&data.SecondName,
			&data.MiddleName,
			&data.DateOfBirth,
			&data.SNILS,
			&data.ContactPhone,
			&data.Email,
			&data.ID_Passport,
			&data.ID_RegAddress,
			&data.ID_EducationListener,
			&data.ID_PlaceWork,
			&data.ID_Legalentity,
			&data.ID_Contractor,
			&data.RegistrationAddress.MailIndex,
			&data.RegistrationAddress.Region,
			&data.RegistrationAddress.City,
			&data.RegistrationAddress.Street,
			&data.RegistrationAddress.House,
			&data.RegistrationAddress.Building,
			&data.RegistrationAddress.Apartment,
		); err != nil {

			l.logger.Error("database error",
				"operation", "read_mapping_full_listener",
				"id_listener", id,
				"type", "query",
				"err", err,
			)

			return nil, repoutils.HandleRepoErr(err)
		}
	}

	return &data, nil
}

func (l *ListenerRepo) FindIdById(ctx context.Context, id uuid.UUID) (*entity.Listener, error) {

	exists, err := repoutils.Exists(ctx, l.repo, "listener", "id_listener", id)
	if err != nil {
		return nil, repoutils.HandleRepoErr(err)
	}

	if !exists {
		return nil, repoutils.HandleRepoErr(sql.ErrNoRows)
	}

	query := "SELECT id_listener, id_passport, id_regaddress, id_educationlistener, id_placework FROM listener WHERE id_listener = $1"

	rows, err := l.repo.QueryContext(ctx, query, id)
	if err != nil {

		l.logger.Error("database error",
			"operation", "find_by_id_listener",
			"id_listener", id,
			"type", "query",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}
	defer rows.Close()

	var listener entity.Listener

	for rows.Next() {

		if err = rows.Scan(
			&listener.ID_Listener,
			&listener.ID_Passport,
			&listener.ID_RegAddress,
			&listener.ID_EducationListener,
			&listener.ID_PlaceWork,
		); err != nil {

			l.logger.Error("database error",
				"operation", "find_by_id_mapping_listener",
				"id_listener", id,
				"type", "query",
				"err", err,
			)

			return nil, repoutils.HandleRepoErr(err)
		}
	}

	return &listener, nil
}

func (l *ListenerRepo) DeleteInTx(ctx context.Context, tx database.Tx, id uuid.UUID) error {

	exists, err := repoutils.Exists(ctx, l.repo, "listener", "id_listener", id)
	if err != nil {

		l.logger.Debug("database error",
			"operation", "check_unique",
			"table", "listener",
			"row", "id_listener",
			"type", "exist",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	if !exists {
		return apperrors.ErrNoExists
	}

	query := "DELETE FROM listener WHERE id_listener = $1"

	if _, err := tx.ExecContext(ctx, query, id); err != nil {

		l.logger.Error("database error",
			"operation", "delete_listener",
			"id_listener", id,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (l *ListenerRepo) UpdateInTx(ctx context.Context, tx database.Tx, listener entity.Listener) error {

	exists, err := repoutils.Exists(ctx, l.repo, "listener", "id_listener", listener.ID_Listener)
	if err != nil {

		l.logger.Debug("database error",
			"operation", "check_unique",
			"table", "listener",
			"row", "id_listener",
			"type", "exist",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	if !exists {
		return apperrors.ErrNoExists
	}

	query := `
	update listener
	set
	first_name = coalesce($1, first_name),
	second_name = coalesce($2, second_name),
	middle_name = coalesce($3, middle_name),
	date_of_birth = coalesce($4, date_of_birth),
	snils = coalesce($5, snils),
	contact_phone = coalesce($6, contact_phone),
	email = coalesce($7, email)
	where id_listener = $8;`

	if _, err := tx.ExecContext(ctx, query, listener.FirstName,
		listener.SecondName,
		listener.MiddleName,
		listener.DateOfBirth,
		listener.SNILS,
		listener.ContactPhone,
		listener.Email,
		listener.ID_Listener); err != nil {

		l.logger.Error("database error",
			"operation", "update_listener",
			"id_listener", listener.ID_Listener,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil

}
