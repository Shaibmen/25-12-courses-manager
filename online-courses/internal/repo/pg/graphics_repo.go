package pg

import (
	"context"
	"log/slog"
	"online-courses/internal/database"
	"online-courses/internal/domain/entity"
)

type GraphicsRepo struct {
	repo   database.DB
	logger *slog.Logger
}

func NewGraphicsRepo(repo database.DB, logger *slog.Logger) *GraphicsRepo {
	return &GraphicsRepo{repo: repo, logger: logger}
}

func (g *GraphicsRepo) CountListenersOnProgram(ctx context.Context) ([]entity.CountListenersOnProgramStruct, error) {

	var data []entity.CountListenersOnProgramStruct

	query :=
		`
		SELECT
		p.name_prof_education,
		COUNT(e.id_listener) AS listeners
		FROM enrollmentlistener e
		JOIN programeducation p ON e.id_programeducation = p.id_programeducation
		GROUP BY p.name_prof_education
		ORDER BY listeners DESC;
	`

	if err := g.repo.SelectContext(ctx, &data, query); err != nil {
		return nil, err
	}

	return data, nil
}

func (g GraphicsRepo) PopularProgramType(ctx context.Context) ([]entity.PopularProgramTypeStruct, error) {

	var data []entity.PopularProgramTypeStruct

	query :=
		`
		SELECT
		name_prof_education,
		educationtype,
		COUNT(*) AS listeners
		FROM accurateprogram
		GROUP BY name_prof_education, educationtype
		ORDER BY listeners DESC;
	`

	if err := g.repo.SelectContext(ctx, &data, query); err != nil {
		return nil, err
	}

	return data, nil
}

func (g GraphicsRepo) CountListenersOnProgramAccurate(ctx context.Context) ([]entity.CountListenersOnProgramStruct, error) {

	var data []entity.CountListenersOnProgramStruct

	query :=
		`
		SELECT
		name_prof_education,
		COUNT(*) AS listeners
		FROM accurateprogram
		GROUP BY name_prof_education
		ORDER BY listeners DESC;
	`

	if err := g.repo.SelectContext(ctx, &data, query); err != nil {
		return nil, err
	}

	return data, nil
}

func (g GraphicsRepo) WorthProgramAccurate(ctx context.Context) ([]entity.WorthProgramAccurateStruct, error) {

	var data []entity.WorthProgramAccurateStruct

	query :=
		`
		SELECT
    name_prof_education,
    educationtype,
    SUM(price) AS total_expected_revenue
	FROM accurateprogram
	GROUP BY name_prof_education, educationtype
	ORDER BY total_expected_revenue DESC;

	`

	if err := g.repo.SelectContext(ctx, &data, query); err != nil {
		return nil, err
	}

	return data, nil
}

func (g GraphicsRepo) AgeDiff(ctx context.Context) ([]entity.AgeDiffStruct, error) {

	var data []entity.AgeDiffStruct

	query :=
		`
		SELECT
		ap.name_prof_education,
		CONCAT(FLOOR(EXTRACT(YEAR FROM AGE(CURRENT_DATE, l.date_of_birth))/10)*10, '-',
		FLOOR(EXTRACT(YEAR FROM AGE(CURRENT_DATE, l.date_of_birth))/10)*10+9) AS age_range,
		COUNT(*) AS listeners
		FROM listener l
		JOIN accurateprogram ap ON l.id_listener = ap.id_listener
		GROUP BY ap.name_prof_education, age_range
		ORDER BY ap.name_prof_education, MIN(EXTRACT(YEAR FROM AGE(CURRENT_DATE, l.date_of_birth)));
	`

	if err := g.repo.SelectContext(ctx, &data, query); err != nil {
		return nil, err
	}

	return data, nil
}

func (g GraphicsRepo) WhoEnrolled(ctx context.Context) ([]entity.WhoEnrolledStruct, error) {

	var data []entity.WhoEnrolledStruct

	query :=
		`
		SELECT
		DATE_TRUNC('month', e.start_date) AS month,
		CASE
		WHEN l.id_legalentity IS NOT NULL THEN 'Юрлицо'
		WHEN l.id_contractor IS NOT NULL THEN 'Физ. лицо(3 стороны)'
		ELSE 'Физ. лицо'
		END AS source,
		COUNT(*) AS cnt
		FROM listener l
		JOIN enrollmentlistener e ON l.id_listener = e.id_listener
		GROUP BY month, source
		ORDER BY month, source;
	`

	if err := g.repo.SelectContext(ctx, &data, query); err != nil {
		return nil, err
	}

	return data, nil
}

func (g GraphicsRepo) GroupMembers(ctx context.Context) ([]entity.GroupMembersStruct, error) {

	var data []entity.GroupMembersStruct

	query :=
		`
		SELECT
		g.name_group,
		COUNT(e.id_listener) FILTER (WHERE e.is_active) AS active_enrolled
		FROM groups g
		LEFT JOIN enrollmentlistener e ON g.id_groups = e.id_group
		GROUP BY g.name_group
		ORDER BY active_enrolled DESC;

	`

	if err := g.repo.SelectContext(ctx, &data, query); err != nil {
		return nil, err
	}

	return data, nil
}

func (g GraphicsRepo) DivisionMember(ctx context.Context) ([]entity.DivisionMemberStruct, error) {

	var data []entity.DivisionMemberStruct

	query :=
		`
	SELECT
	divisionseducation,
	COUNT(*) AS listeners
	FROM accurateprogram
	GROUP BY divisionseducation;

	`

	if err := g.repo.SelectContext(ctx, &data, query); err != nil {
		return nil, err
	}

	return data, nil
}
