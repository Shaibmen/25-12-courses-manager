package pg

import (
	"context"
	"log/slog"
	"online-courses/internal/database"
	"time"
)

type ProcedureRepo struct {
	repo   database.DB
	logger *slog.Logger
}

func NewProcedureRepo(db database.DB, logger *slog.Logger) *ProcedureRepo {
	return &ProcedureRepo{repo: db, logger: logger}
}

func (p *ProcedureRepo) DeactivationNoValidEnrollment() {
	query :=
		`
		call deactivate_finished_enrollments()
	`

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := p.repo.ExecContext(ctx, query)
	if err != nil {
		p.logger.Error("деактивация записей на курс не сработала", "err", err)
	}
}

func (p ProcedureRepo) ShuffleProgram() {
	query :=

		`
	call shuffle_program_order()
	`
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := p.repo.ExecContext(ctx, query)
	if err != nil {
		p.logger.Error("обновление программ обучения не сработало", "err", err)
	}
}

func (p ProcedureRepo) ShuffleLevelEducation() {
	query :=

		`
	call shuffle_education_type()
	`
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := p.repo.ExecContext(ctx, query)
	if err != nil {
		p.logger.Error("обновление типов обучения не сработало", "err", err)
	}
}
