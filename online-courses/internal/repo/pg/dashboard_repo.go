package pg

import (
	"context"
	"log/slog"
	"online-courses/internal/database"
	"online-courses/internal/domain/entity"
	repoutils "online-courses/internal/repo/pg/repo_utils"
)

type DashboardRepo struct {
	db     database.DB
	logger *slog.Logger
}

func NewDashboardRepo(db database.DB, logger *slog.Logger) *DashboardRepo {
	return &DashboardRepo{db: db, logger: logger}
}

func (d *DashboardRepo) UserDashboard(ctx context.Context) (*entity.UserDashBoard, error) {
	query :=
		`
	select * from v_total_listeners;
	`
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {

		d.logger.Error("database error",
			"operation", "query v_total_listeners",
			"type", "exec",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}
	defer rows.Close()

	var totalListener int
	for rows.Next() {
		if err = rows.Scan(&totalListener); err != nil {

			d.logger.Error("database error",
				"operation", "row.Next v_total_listeners",
				"type", "exec",
				"err", err,
			)

			return nil, repoutils.HandleRepoErr(err)
		}
	}

	query =
		`
	select * from v_total_programs;
	`
	var totalProgram int
	rows, err = d.db.QueryContext(ctx, query)
	if err != nil {

		d.logger.Error("database error",
			"operation", "query v_total_programs",
			"type", "exec",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}

	for rows.Next() {
		if err = rows.Scan(&totalProgram); err != nil {

			d.logger.Error("database error",
				"operation", "row.Next v_total_programs",
				"type", "exec",
				"err", err,
			)

			return nil, repoutils.HandleRepoErr(err)
		}
	}

	query =
		`
	select * from v_active_enrollments;
	`

	var activeEnrollments int
	rows, err = d.db.QueryContext(ctx, query)
	if err != nil {

		d.logger.Error("database error",
			"operation", "query v_active_enrollments",
			"type", "exec",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}

	for rows.Next() {
		if err = rows.Scan(&activeEnrollments); err != nil {

			d.logger.Error("database error",
				"operation", "row.Next v_active_enrollments",
				"type", "exec",
				"err", err,
			)

			return nil, repoutils.HandleRepoErr(err)
		}
	}

	query =
		`
	select * from v_programs_ending_soon;
	`

	var programEndingSoon []entity.ProgramEndingSoon
	rows, err = d.db.QueryContext(ctx, query)
	if err != nil {

		d.logger.Error("database error",
			"operation", "query v_programs_ending_soon",
			"type", "exec",
			"err", err,
		)

		return nil, repoutils.HandleRepoErr(err)
	}

	for rows.Next() {
		var program entity.ProgramEndingSoon
		if err = rows.Scan(
			&program.NameProfEducation,
			&program.EndDate,
			&program.TotalListeners,
		); err != nil {

			d.logger.Error("database error",
				"operation", "row.Next v_programs_ending_soon",
				"type", "exec",
				"err", err,
			)

			return nil, repoutils.HandleRepoErr(err)
		}
		programEndingSoon = append(programEndingSoon, program)
	}

	data := entity.UserDashBoard{
		TotalListener:     totalListener,
		TotalProgram:      totalProgram,
		ActiveEnrollments: activeEnrollments,
		ProgramEndingSoon: programEndingSoon,
	}

	return &data, nil

}

func (d *DashboardRepo) AdminDashboard(ctx context.Context) (*entity.AdminDashboard, error) {
	query :=
		`
	select * from admin_active_sessions 
	`

	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, repoutils.HandleRepoErr(err)
	}

	defer rows.Close()

	var pgStats []entity.PgStat

	for rows.Next() {
		var pgStat entity.PgStat
		if err = rows.Scan(
			&pgStat.Pid,
			&pgStat.ApplicationName,
			&pgStat.ClientAddr,
			&pgStat.ClientPort,
			&pgStat.State,
			&pgStat.QueryStart,
		); err != nil {
			return nil, repoutils.HandleRepoErr(err)
		}
		pgStats = append(pgStats, pgStat)
	}

	query =
		`
	select * from admin_db_stats
	`

	rows, err = d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, repoutils.HandleRepoErr(err)
	}

	defer rows.Close()

	var adminStat entity.AdminStat
	for rows.Next() {
		if err = rows.Scan(
			&adminStat.DatabaseName,
			&adminStat.Totalsize,
			&adminStat.ActiveConnection,
			&adminStat.CommittedTx,
			&adminStat.RolledbackTx,
			&adminStat.DiskBlockRead,
			&adminStat.BufferHits,
		); err != nil {
			return nil, repoutils.HandleRepoErr(err)
		}
	}

	query =
		`
	select * from role
	`

	rows, err = d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, repoutils.HandleRepoErr(err)
	}

	defer rows.Close()

	var roles []entity.Role
	for rows.Next() {
		var role entity.Role
		if err = rows.Scan(
			&role.ID,
			&role.Role,
		); err != nil {
			return nil, repoutils.HandleRepoErr(err)
		}
		roles = append(roles, role)
	}

	data := entity.AdminDashboard{
		PgStat:    pgStats,
		AdminStat: adminStat,
		Role:      roles,
	}

	return &data, nil
}
