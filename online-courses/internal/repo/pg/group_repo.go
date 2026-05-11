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

type GroupDB struct {
	repo   database.DB
	logger *slog.Logger
}

func NewGroupDB(db database.DB, logger *slog.Logger) *GroupDB {
	return &GroupDB{repo: db, logger: logger}
}

func (g *GroupDB) Create(ctx context.Context, m entity.Group) error {

	query := `insert into groups (id_groups, name_group, raspisanie) values ($1, $2, $3)`

	_, err := g.repo.ExecContext(ctx, query, m.ID_Group, m.NameGroup, m.Raspisanie)
	if err != nil {
		g.logger.Error("database error",
			"operation", "insert_group",
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (g *GroupDB) Update(ctx context.Context, m entity.Group) error {

	exists, err := repoutils.Exists(ctx, g.repo, "groups", "id_groups", m.ID_Group)
	if err != nil {

		g.logger.Debug("database error",
			"operation", "check_unique",
			"table", "groups",
			"row", "id_groups",
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
	update groups
	set
	name_group = coalesce($1, name_group),
	raspisanie = coalesce($2, raspisanie)
	where id_groups = $3
	`

	if _, err := g.repo.ExecContext(
		ctx,
		query,
		m.NameGroup,
		m.Raspisanie,
		m.ID_Group); err != nil {

		g.logger.Error("database error",
			"operation", "update_groups",
			"id_groups", m.ID_Group,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (g *GroupDB) Delete(ctx context.Context, id uuid.UUID) error {

	exists, err := repoutils.Exists(ctx, g.repo, "groups", "id_groups", id)
	if err != nil {

		g.logger.Debug("database error",
			"operation", "check_unique",
			"table", "groups",
			"row", "id_groups",
			"type", "exist",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	if !exists {
		return apperrors.ErrNoExists
	}

	query := "DELETE FROM groups WHERE id_groups = $1"

	if _, err := g.repo.ExecContext(ctx, query, id); err != nil {

		g.logger.Error("database error",
			"operation", "delete_group",
			"id_groups", id,
			"type", "exec",
			"err", err,
		)

		return repoutils.HandleRepoErr(err)
	}

	return nil
}

func (g *GroupDB) ReadByID(ctx context.Context, id uuid.UUID) (*entity.Group, error) {
	query :=
		`
	select * from groups
	where id_groups = $1
	`

	group := entity.Group{}
	err := g.repo.QueryRowContext(ctx, query, id).Scan(
		&group.ID_Group,
		&group.NameGroup,
		&group.Raspisanie,
	)
	if err != nil {

		g.logger.Error("database error",
			"operation", "read_groups",
			"id_groups", id,
			"type", "exec",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}

	return &group, nil
}

func (g *GroupDB) Read(ctx context.Context, page int, filter string) ([]entity.Group, error) {

	limit, offset := repoutils.Pagination(page)

	query :=
		`
	select * from groups
	where ($1::text is null or name_group ilike '%' || $1::text || '%' )
	limit $2 offset $3
	`

	rows, err := g.repo.QueryContext(ctx, query, filter, limit, offset)
	if err != nil {

		g.logger.Error("database error",
			"operation", "read_groups",
			"page", page,
			"type", "query",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}
	defer rows.Close()

	var groups []entity.Group

	for rows.Next() {
		var group entity.Group
		if err := rows.Scan(
			&group.ID_Group,
			&group.NameGroup,
			&group.Raspisanie,
		); err != nil {

			g.logger.Error("database error",
				"operation", "read_mapping_program_education",
				"page", page,
				"type", "query",
				"err", err,
			)

			return nil, err
		}
		groups = append(groups, group)
	}

	return groups, nil
}
