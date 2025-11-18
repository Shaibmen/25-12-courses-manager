package pg

import (
	"context"
	"log/slog"
	"online-courses/internal/database"
	"online-courses/internal/domain/entity"
	repoutils "online-courses/internal/repo/pg/repo_utils"
)

type LevelEducationRepo struct {
	repo   database.DB
	logger *slog.Logger
}

func NewLevelEducationRepo(db database.DB, logger *slog.Logger) *LevelEducationRepo {
	return &LevelEducationRepo{repo: db, logger: logger}
}

func (lvl *LevelEducationRepo) Read(ctx context.Context, filter string) ([]entity.LevelEducation, error) {
	query :=
		`
	select * from leveleducation
	where ($1::text is null or education ilike '%' || $1::text || '%' )
	`

	rows, err := lvl.repo.QueryContext(ctx, query, filter)
	if err != nil {

		lvl.logger.Error("database error",
			"operation", "read_level_education",
			"type", "query",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}
	defer rows.Close()

	var levels []entity.LevelEducation
	for rows.Next() {
		var l entity.LevelEducation
		if err := rows.Scan(
			&l.ID_LevelEducation,
			&l.Education,
		); err != nil {

			lvl.logger.Error("database error",
				"operation", "read_mapping_level_education",
				"type", "query",
				"err", err,
			)

			return nil, repoutils.HandleRepoErr(err)
		}
		levels = append(levels, l)
	}

	return levels, nil
}
