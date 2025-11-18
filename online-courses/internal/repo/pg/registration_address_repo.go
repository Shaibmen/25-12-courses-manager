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

type RegistrationAddressRepo struct {
	repo   database.DB
	logger *slog.Logger
}

func NewRegistrationAddressRepo(db database.DB, logger *slog.Logger) *RegistrationAddressRepo {
	return &RegistrationAddressRepo{repo: db, logger: logger}
}

func (r *RegistrationAddressRepo) CreateInTx(ctx context.Context, tx database.Tx, m entity.RegistrationAddress) error {
	query := `INSERT INTO registrationaddress (id_regaddress, mail_index, region, city, street, house, building, apartment)
			  Values ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := tx.ExecContext(ctx, query, m.ID_RegAddress, m.MailIndex, m.Region, m.City, m.Street, m.House, m.Building, m.Apartment)
	if err != nil {

		r.logger.Error("database error",
			"operation", "insert_registration_address",
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (r *RegistrationAddressRepo) DeleteInTx(ctx context.Context, tx database.Tx, id uuid.UUID) error {

	exists, err := repoutils.Exists(ctx, r.repo, "registrationaddress", "id_regaddress", id)
	if err != nil {

		r.logger.Debug("database error",
			"operation", "check_unique",
			"table", "registrationaddress",
			"row", "id_regaddress",
			"type", "exist",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	if !exists {
		return apperrors.ErrNoExists
	}

	query := "DELETE FROM registrationaddress WHERE id_regaddress = $1"

	if _, err := tx.ExecContext(ctx, query, id); err != nil {

		r.logger.Error("database error",
			"operation", "delete_registration_address",
			"id_registration_address", id,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (r *RegistrationAddressRepo) UpdateInTx(ctx context.Context, tx database.Tx, address entity.RegistrationAddress) error {

	exists, err := repoutils.Exists(ctx, r.repo, "registrationaddress", "id_regaddress", address.ID_RegAddress)
	if err != nil {

		r.logger.Debug("database error",
			"operation", "check_unique",
			"table", "registrationaddress",
			"row", "id_regaddress",
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
	update registrationaddress
	set 
	mail_index = coalesce($1, mail_index),
	region = coalesce($2, region),
	city = coalesce($3, city),
	street = coalesce($4, street),
	house = coalesce($5, house),
	building = coalesce($6, building),
	apartment = coalesce($7, apartment)
	where id_regaddress = $8;
	`

	if _, err := tx.ExecContext(
		ctx,
		query,
		address.MailIndex,
		address.Region,
		address.City,
		address.Street,
		address.House,
		address.Building,
		address.Building,
		address.ID_RegAddress); err != nil {

		r.logger.Error("database error",
			"operation", "read_mapping_registration_address",
			"id_regaddress", address.ID_RegAddress,
			"type", "query",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}
	return nil
}
