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

type PassportRepo struct {
	repo   database.DB
	logger *slog.Logger
}

func NewPassportRepo(db database.DB, logger *slog.Logger) *PassportRepo {
	return &PassportRepo{repo: db, logger: logger}
}

func (p *PassportRepo) CreateInTx(ctx context.Context, tx database.Tx, m entity.Passport) error {
	query := `INSERT INTO passport (id_passport, place_birth, citizenship, gender, seria, number, passport_given, date_given, code)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			  `

	_, err := tx.ExecContext(ctx, query, m.ID_Passport, m.PlaceBirth, m.Citizenship, m.Gender, m.Seria, m.Number, m.PassportGiven, m.DateGiven, m.Code)
	if err != nil {

		p.logger.Error("database error",
			"operation", "insert_passport",
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (p *PassportRepo) DeleteInTx(ctx context.Context, tx database.Tx, id uuid.UUID) error {

	exists, err := repoutils.Exists(ctx, p.repo, "passport", "id_passport", id)
	if err != nil {

		p.logger.Debug("database error",
			"operation", "check_unique",
			"table", "passport",
			"row", "id_passport",
			"type", "exist",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	if !exists {
		return apperrors.ErrNoExists
	}

	query := `DELETE FROM passport WHERE id_passport = $1`

	if _, err := tx.ExecContext(ctx, query, id); err != nil {

		p.logger.Error("database error",
			"operation", "delete_passport",
			"id_passport", id,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (p *PassportRepo) UpdateInTx(ctx context.Context, tx database.Tx, passport entity.Passport) error {

	exists, err := repoutils.Exists(ctx, p.repo, "passport", "id_passport", passport.ID_Passport)
	if err != nil {

		p.logger.Debug("database error",
			"operation", "check_unique",
			"table", "passport",
			"row", "id_passport",
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
	update passport
	set 
	place_birth = coalesce($1, place_birth),
	citizenship = coalesce($2, citizenship),
	gender = coalesce($3, gender),
	seria = coalesce($4, seria),
	number = coalesce($5, number),
	passport_given = coalesce($6, passport_given),
	date_given = coalesce($7, date_given),
	code = coalesce($8, code)
	where id_passport = $9;
	`

	if _, err := tx.ExecContext(
		ctx,
		query,
		passport.PlaceBirth,
		passport.Citizenship,
		passport.Gender,
		passport.Seria,
		passport.Number,
		passport.PassportGiven,
		passport.DateGiven,
		passport.Code,
		passport.ID_Passport); err != nil {

		p.logger.Error("database error",
			"operation", "update_education_listener",
			"id_passport", passport.ID_Passport,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}
	return nil
}
