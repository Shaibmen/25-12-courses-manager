package pg

import (
	"context"
	"log/slog"
	"online-courses/internal/database"
	"online-courses/internal/domain/entity"
	repoutils "online-courses/internal/repo/pg/repo_utils"
	"time"
)

type ReportRepo struct {
	repo   database.DB
	logger *slog.Logger
}

func NewReportRepo(db database.DB, logger *slog.Logger) *ReportRepo {
	return &ReportRepo{repo: db, logger: logger}
}

func (r *ReportRepo) ReportPeriod(ctx context.Context, startDate, endDate time.Time) ([]entity.ReportPeriod, error) {

	query :=
		`
	SELECT
    l.second_name || ' ' || l.first_name || ' ' || l.middle_name AS full_name,
    p.name_prof_education AS program_name,
    e.start_date,
    e.end_date,
    e.current_price AS payment,
    SUM(e.current_price) OVER () AS total_income_period
	FROM enrollmentlistener e
	JOIN listener l ON l.id_listener = e.id_listener
	JOIN programeducation p ON p.id_programeducation = e.id_programeducation
	WHERE e.start_date >= $1
	AND e.start_date < $2
	ORDER BY e.start_date;
	`

	rows, err := r.repo.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		r.logger.Error("database error",
			"operation", "read_reports",
			"type", "query",
			"err", err,
		)
		return nil, repoutils.HandleRepoErr(err)
	}
	defer rows.Close()

	var report []entity.ReportPeriod

	for rows.Next() {
		var reports entity.ReportPeriod
		if err = rows.Scan(
			&reports.FullName,
			&reports.ProgramName,
			&reports.StartDate,
			&reports.EndDate,
			&reports.Payment,
			&reports.TotalIncomePeriod,
		); err != nil {
			r.logger.Error("database error",
				"operation", "read_rows_reports",
				"type", "query",
				"err", err,
			)
			return nil, repoutils.HandleRepoErr(err)
		}

		report = append(report, reports)
	}

	return report, nil
}

func (r *ReportRepo) MostExpensiveProgram(ctx context.Context, startDate, endDate time.Time) ([]entity.MostExpensiveReport, error) {

	query :=
		`
	SELECT
    p.name_prof_education AS program_name,
    COUNT(e.id_listener) AS total_listeners,
    SUM(e.current_price) AS total_revenue
	FROM enrollmentlistener e
	JOIN programeducation p ON e.id_programeducation = p.id_programeducation
	WHERE e.start_date >= $1
	AND e.start_date < $2
	GROUP BY p.name_prof_education
	ORDER BY total_revenue DESC;
	`

	rows, err := r.repo.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		r.logger.Error("database error",
			"operation", "read_most_exp_reports",
			"type", "query",
			"err", err,
		)
		return nil, repoutils.HandleRepoErr(err)
	}

	defer rows.Close()

	var reports []entity.MostExpensiveReport
	for rows.Next() {
		var report entity.MostExpensiveReport
		if err = rows.Scan(
			&report.ProgramName,
			&report.TotalListeners,
			&report.TotalRevenue,
		); err != nil {
			r.logger.Error("database error",
				"operation", "read_rows_most_exp_reports",
				"type", "query",
				"err", err,
			)
			return nil, repoutils.HandleRepoErr(err)
		}
		reports = append(reports, report)
	}

	return reports, nil
}
