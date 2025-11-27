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

type ContractorRepo struct {
	repo   database.DB
	logger *slog.Logger
}

func NewContractorRepo(db database.DB, logger *slog.Logger) *ContractorRepo {
	return &ContractorRepo{repo: db, logger: logger}
}

func (l *ContractorRepo) CreateInTx(ctx context.Context, tx database.Tx, m *entity.Contractor, idPassport, idRegAddress uuid.UUID) error {

	query :=
		`
	insert into contractor (id_contractor, first_name, second_name, middle_name, contact_phone, email, id_passport, id_regaddress)
	values ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := tx.ExecContext(ctx, query, m.ID_Contractor, m.FirstName, m.SecondName, m.MiddleName, m.Contact_phone, m.Email, idPassport, idRegAddress)
	if err != nil {

		l.logger.Error("database error",
			"operation", "insert_contractor",
			"id_contractor", m.ID_Contractor,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (c *ContractorRepo) UpdateInTx(ctx context.Context, tx database.Tx, contractor entity.Contractor) error {

	exists, err := repoutils.Exists(ctx, c.repo, "contractor", "id_contractor", contractor.ID_Contractor)
	if err != nil {

		c.logger.Debug("database error",
			"operation", "check_unique",
			"table", "contractor",
			"row", "id_contractor",
			"type", "exist",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	if !exists {
		return apperrors.ErrNoExists
	}

	query := `
	update contractor 
	set
	first_name = coalesce($1, first_name),
	second_name = coalesce($2, second_name),
	middle_name = coalesce($3, middle_name),
	contact_phone = coalesce($4, contact_phone),
	email = coalesce($5, email)
	where id_contractor = $6;`

	if _, err := tx.ExecContext(ctx, query, contractor.FirstName, contractor.SecondName, contractor.MiddleName, contractor.Contact_phone, contractor.Email, contractor.ID_Contractor); err != nil {

		c.logger.Error("database error",
			"operation", "update_contractor",
			"id_listener", contractor.ID_Contractor,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil

}

func (c *ContractorRepo) DeleteInTx(ctx context.Context, tx database.Tx, id uuid.UUID) error {

	exists, err := repoutils.Exists(ctx, c.repo, "contractor", "id_contractor", id)
	if err != nil {

		c.logger.Debug("database error",
			"operation", "check_unique",
			"table", "contractor",
			"row", "id_contractor",
			"type", "exist",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	if !exists {
		return apperrors.ErrNoExists
	}

	query := "DELETE FROM contractor WHERE id_contractor = $1"

	if _, err := tx.ExecContext(ctx, query, id); err != nil {

		c.logger.Error("database error",
			"operation", "delete_contractor",
			"id_contractor", id,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (c *ContractorRepo) FindById(ctx context.Context, id uuid.UUID) (*uuid.UUID, *uuid.UUID, error) {
	exists, err := repoutils.Exists(ctx, c.repo, "contractor", "id_contractor", id)
	if err != nil {
		return nil, nil, repoutils.HandleRepoErr(err)
	}

	if !exists {
		return nil, nil, repoutils.HandleRepoErr(sql.ErrNoRows)
	}

	query := "SELECT id_passport, id_regaddress FROM contractor WHERE id_contractor = $1"

	rows, err := c.repo.QueryContext(ctx, query, id)
	if err != nil {

		c.logger.Error("database error",
			"operation", "find_by_id_contractor",
			"id_contractor", id,
			"type", "query",
			"err", err,
		)

		return nil, nil, repoutils.HandleRepoErr(err)
	}

	defer rows.Close()

	var passportID, regAdderess uuid.UUID

	for rows.Next() {

		if err = rows.Scan(&passportID, &regAdderess); err != nil {

			c.logger.Error("database error",
				"operation", "find_by_id_mapping_contractor",
				"id_contractor", id,
				"type", "query",
				"err", err,
			)

			return nil, nil, repoutils.HandleRepoErr(err)
		}
	}

	return &passportID, &regAdderess, nil
}
