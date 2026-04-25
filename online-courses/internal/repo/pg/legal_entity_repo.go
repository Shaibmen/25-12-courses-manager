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

type LegalEntityRepo struct {
	repo   database.DB
	logger *slog.Logger
}

func NewLegalEntity(db database.DB, logger *slog.Logger) *LegalEntityRepo {
	return &LegalEntityRepo{repo: db, logger: logger}
}

func (l *LegalEntityRepo) CreateInTx(ctx context.Context, tx database.Tx, m *entity.LegalEntity) error {

	query :=
		`
	insert into legal_entity (id_legalentity, name_company, inn, kpp, ogrn, phone, email, first_name, second_name, middle_name, id_regaddress, status)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := tx.ExecContext(ctx, query, m.ID_Legalentity, m.NameCompany, m.Inn, m.Kpp, m.Ogrn, m.Phone, m.Email, m.FirstName, m.SecondName, m.MiddleName, m.ID_RegAddress, m.Status)
	if err != nil {
		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (l *LegalEntityRepo) Read(ctx context.Context, page int, filter string) ([]entity.LegalEntity, error) {

	query :=
		`
	select 
	l.*
	from legal_entity as l
	where ($1::text is null or name_company ilike '%' || $1::text || '%' )
	LIMIT $2 OFFSET $3
	`

	limit, offset := repoutils.Pagination(page)

	rows, err := l.repo.QueryContext(ctx, query, filter, limit, offset)
	if err != nil {
		return nil, repoutils.HandleRepoErr(err)
	}

	defer rows.Close()

	var entities []entity.LegalEntity

	for rows.Next() {
		var data entity.LegalEntity
		if err = rows.Scan(
			&data.ID_Legalentity,
			&data.NameCompany,
			&data.Inn,
			&data.Kpp,
			&data.Ogrn,
			&data.Phone,
			&data.Email,
			&data.FirstName,
			&data.SecondName,
			&data.MiddleName,
			&data.Status,
			&data.ID_RegAddress,
		); err != nil {
			return nil, repoutils.HandleRepoErr(err)
		}
		entities = append(entities, data)
	}

	return entities, nil
}

func (l *LegalEntityRepo) UpdateInTx(ctx context.Context, tx database.Tx, m entity.LegalEntity) error {

	exists, err := repoutils.Exists(ctx, l.repo, "legal_entity", "id_legalentity", m.ID_Legalentity)
	if err != nil {

		l.logger.Debug("database error",
			"operation", "check_unique",
			"table", "legal_entity",
			"row", "id_legalentity",
			"type", "exist",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	if !exists {

		return apperrors.ErrNoExists
	}

	query := `
	update legal_entity 
	set
	name_company = coalesce($1, name_company),
	inn = coalesce($2, inn),
	kpp = coalesce($3, kpp),
	ogrn = coalesce($4, ogrn),
	phone = coalesce($5, phone),
	email = coalesce($6, email),
	first_name = coalesce($7, first_name),
	second_name = coalesce($8, second_name),
	middle_name = coalesce($9, middle_name),
	status = coalesce($10, status)
	where id_legalentity = $11;`

	if _, err := tx.ExecContext(ctx, query, m.NameCompany, m.Inn, m.Kpp, m.Ogrn, m.Phone, m.Email, m.FirstName, m.SecondName, m.MiddleName, m.Status, m.ID_Legalentity); err != nil {

		l.logger.Error("database error",
			"operation", "update_legal_entity",
			"id_legal_entity", m.ID_Legalentity,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil

}

func (l *LegalEntityRepo) DeleteInTx(ctx context.Context, tx database.Tx, id uuid.UUID) error {

	exists, err := repoutils.Exists(ctx, l.repo, "legal_entity", "id_legalentity", id)
	if err != nil {

		l.logger.Debug("database error",
			"operation", "check_unique",
			"table", "legal_entity",
			"row", "id_legalentity",
			"type", "exist",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	if !exists {
		return apperrors.ErrNoExists
	}

	query := "DELETE FROM legal_entity WHERE id_legalentity = $1"

	if _, err := tx.ExecContext(ctx, query, id); err != nil {

		l.logger.Error("database error",
			"operation", "delete_legalentity",
			"id_legalentity", id,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (c *LegalEntityRepo) ReadFullData(ctx context.Context, id uuid.UUID) (*entity.LegalEntity, error) {

	exists, err := repoutils.Exists(ctx, c.repo, "legal_entity", "id_legalentity", id)
	if err != nil {
		return nil, repoutils.HandleRepoErr(err)
	}

	if !exists {
		return nil, repoutils.HandleRepoErr(sql.ErrNoRows)
	}

	query :=
		`
	select 
	l.*,
	r.mail_index, r.region, r.city, r.street, r.house, r.building, r.apartment
	from legal_entity as l
	inner join registrationaddress r on l.id_regaddress = r.id_regaddress
	
	where id_legalentity = $1
	`

	rows, err := c.repo.QueryContext(ctx, query, id)
	if err != nil {

		c.logger.Error("database error",
			"operation", "read_full_legal_entity",
			"id_legal_entity", id,
			"type", "query",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}
	defer rows.Close()

	var data entity.LegalEntity

	for rows.Next() {
		if err = rows.Scan(
			&data.ID_Legalentity,
			&data.NameCompany,
			&data.Inn,
			&data.Kpp,
			&data.Ogrn,
			&data.Phone,
			&data.Email,
			&data.FirstName,
			&data.SecondName,
			&data.MiddleName,
			&data.Status,
			&data.ID_RegAddress,
			&data.RegistrationAddress.MailIndex,
			&data.RegistrationAddress.Region,
			&data.RegistrationAddress.City,
			&data.RegistrationAddress.Street,
			&data.RegistrationAddress.House,
			&data.RegistrationAddress.Building,
			&data.RegistrationAddress.Apartment,
		); err != nil {

			c.logger.Error("database error",
				"operation", "read_mapping_full_legal_entity",
				"id_legalentity", id,
				"type", "query",
				"err", err,
			)

			return nil, repoutils.HandleRepoErr(err)
		}
	}
	return &data, nil
}

func (c *LegalEntityRepo) FindById(ctx context.Context, id uuid.UUID) (*uuid.UUID, error) {
	exists, err := repoutils.Exists(ctx, c.repo, "legal_entity", "id_legalentity", id)
	if err != nil {
		return nil, repoutils.HandleRepoErr(err)
	}

	if !exists {

		return nil, repoutils.HandleRepoErr(sql.ErrNoRows)
	}

	query := "SELECT id_regaddress FROM legal_entity WHERE id_legalentity = $1"

	rows, err := c.repo.QueryContext(ctx, query, id)
	if err != nil {
		c.logger.Error("database error",
			"operation", "find_by_id_legalentity",
			"legal_entity", id,
			"type", "query",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}

	defer rows.Close()

	var regAdderess uuid.UUID

	for rows.Next() {

		if err = rows.Scan(&regAdderess); err != nil {

			c.logger.Error("database error",
				"operation", "find_by_id_mapping_legal_entity",
				"id_legalentity", id,
				"type", "query",
				"err", err,
			)

			return nil, repoutils.HandleRepoErr(err)
		}
	}

	return &regAdderess, nil
}
