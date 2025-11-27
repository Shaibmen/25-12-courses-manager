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
	query := `INSERT INTO listener (id_listener, first_name, second_name, middle_name, date_of_birth, snils, contact_phone, email, id_passport, id_regaddress, id_educationlistener, id_placework, id_legalentity, id_contractor)
	          Values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err := tx.ExecContext(ctx, query, m.ID_Listener, m.FirstName, m.SecondName, m.MiddleName, m.DateOfBirth, m.SNILS, m.ContactPhone, m.Email, i.ID_Passport, i.ID_RegAddress, i.ID_EducationListener, i.ID_PlaceWork, i.ID_LegalEntity, i.ID_Contractor)
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
	r.mail_index, r.region, r.city, r.street, r.house, r.building, r.apartment,
	c.id_contractor, c.first_name, c.second_name, c.middle_name, c.contact_phone, c.email,
	p.place_birth, p.citizenship, p.gender, p.seria, p.number, p.passport_given, p.date_given, p.code,
	reg.mail_index, reg.region, reg.city, reg.street, reg.house, reg.building, reg.apartment
	from listener as l 
	inner join registrationaddress r on l.id_regaddress = r.id_regaddress
	left join contractor c on l.id_contractor = c.id_contractor
	left join passport p on c.id_passport = p.id_passport
	left join registrationaddress reg on c.id_regaddress = reg.id_regaddress

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

		var (
			contractorID                           sql.NullString
			contractorFirstName                    sql.NullString
			contractorSecondName                   sql.NullString
			contractorMiddleName                   sql.NullString
			contractorContact_phone                sql.NullString
			contractorEmail                        sql.NullString
			contractorPassportPlaceBirth           sql.NullString
			contractorPassportCitizenship          sql.NullString
			contractorPassportGender               sql.NullString
			contractorPassportSeria                sql.NullInt64
			contractorPassportNumber               sql.NullInt64
			contractorPassportPassportGiven        sql.NullString
			contractorPassportDateGiven            sql.NullTime
			contractorPassportCode                 sql.NullString
			contractorRegistrationAddressMailIndex sql.NullInt64
			contractorRegistrationAddressRegion    sql.NullString
			contractorRegistrationAddressCity      sql.NullString
			contractorRegistrationAddressStreet    sql.NullString
			contractorRegistrationAddressHouse     sql.NullString
			contractorRegistrationAddressBuilding  sql.NullString
			contractorRegistrationAddressApartment sql.NullString
		)
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
			&contractorID,
			&contractorFirstName,
			&contractorSecondName,
			&contractorMiddleName,
			&contractorContact_phone,
			&contractorEmail,
			&contractorPassportPlaceBirth,
			&contractorPassportCitizenship,
			&contractorPassportGender,
			&contractorPassportSeria,
			&contractorPassportNumber,
			&contractorPassportPassportGiven,
			&contractorPassportDateGiven,
			&contractorPassportCode,
			&contractorRegistrationAddressMailIndex,
			&contractorRegistrationAddressRegion,
			&contractorRegistrationAddressCity,
			&contractorRegistrationAddressStreet,
			&contractorRegistrationAddressHouse,
			&contractorRegistrationAddressBuilding,
			&contractorRegistrationAddressApartment,
		); err != nil {

			l.logger.Error("database error",
				"operation", "read_mapping_full_listener",
				"id_listener", id,
				"type", "query",
				"err", err,
			)

			return nil, repoutils.HandleRepoErr(err)
		}
		if contractorFirstName.Valid {
			uuidContractor, err := uuid.Parse(contractorID.String)
			if err != nil {
				return nil, repoutils.HandleRepoErr(err)
			}
			data.Contractor = entity.Contractor{
				ID_Contractor: uuidContractor,
				FirstName:     contractorFirstName.String,
				SecondName:    contractorSecondName.String,
				MiddleName:    contractorMiddleName.String,
				Contact_phone: contractorContact_phone.String,
				Email:         contractorEmail.String,
				Passport: entity.Passport{
					PlaceBirth:    contractorPassportPlaceBirth.String,
					Citizenship:   contractorPassportCitizenship.String,
					Gender:        contractorPassportGender.String,
					Seria:         int(contractorPassportSeria.Int64),
					Number:        int(contractorPassportNumber.Int64),
					PassportGiven: contractorPassportPassportGiven.String,
					DateGiven:     contractorPassportDateGiven.Time,
					Code:          contractorPassportCode.String,
				},
				RegistrationAddress: entity.RegistrationAddress{
					MailIndex: int(contractorRegistrationAddressMailIndex.Int64),
					Region:    contractorRegistrationAddressRegion.String,
					City:      contractorRegistrationAddressCity.String,
					Street:    contractorRegistrationAddressStreet.String,
					House:     contractorRegistrationAddressHouse.String,
					Building:  contractorRegistrationAddressBuilding.String,
					Apartment: contractorRegistrationAddressApartment.String,
				},
			}
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

func (l *ListenerRepo) UpdateContractor(ctx context.Context, tx database.Tx, idListener uuid.UUID, idContractor *uuid.UUID) error {
	exists, err := repoutils.Exists(ctx, l.repo, "listener", "id_listener", idListener)
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
	id_contractor = coalesce($1, id_contractor)
	where id_listener = $2;`

	if _, err := tx.ExecContext(ctx, query, idContractor, idListener); err != nil {

		l.logger.Error("database error",
			"operation", "update_listener",
			"id_listener", idListener,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}
