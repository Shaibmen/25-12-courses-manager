package pg

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"online-courses/internal/database"
	"online-courses/internal/domain/entity"
	repoutils "online-courses/internal/repo/pg/repo_utils"

	"github.com/google/uuid"
)

type DocumentsRepo struct {
	repo   database.DB
	logger *slog.Logger
}

func NewDocumentsRepo(repo database.DB, logger *slog.Logger) *DocumentsRepo {
	return &DocumentsRepo{repo: repo, logger: logger}
}

func (d *DocumentsRepo) PrepareDataDocument(ctx context.Context, idListener, idProgram, idExecutor uuid.UUID) (*entity.FullDocument, error) {
	query :=
		`
		
	select 
	
	l.first_name, l.second_name, l.middle_name, l.date_of_birth, l.snils, l.contact_phone, l.email,
	pas.place_birth, pas.citizenship, pas.gender, pas.seria, pas.number, pas.passport_given, pas.date_given, pas.code,
	r.mail_index, r.region, r.city, r.street, r.house, r.building, r.apartment,
	educ.diplom_seria, educ.diplom_number, educ.city as educ_city, educ.region as educ_region, educ.educational_institution , educ.speciality, lvl.education,
	pw.name_company, pw.job_title, pw.all_experience, pw.job_title_experience,
	p.name_prof_education, p.time_education, d.divisions,  et.type_name,
	e.start_date, e.end_date, e.current_price, e.is_active, e.group_number, e.type_of_retraining
	
	from enrollmentlistener as e
	inner join listener l on e.id_listener = l.id_listener
	left join passport pas on l.id_passport = pas.id_passport
	inner join registrationaddress r on l.id_regaddress = r.id_regaddress
	left join educationlistener educ on l.id_educationlistener = educ.id_educationlistener
	left join leveleducation lvl on educ.level_education = lvl.id_leveleducation
	left join placework pw on l.id_placework = pw.id_placework
	inner join programeducation p on e.id_programeducation = p.id_programeducation
	inner join educationtypes et on p.id_educationtype = et.id_educationtype
	inner join divisionseducation d on p.id_divisionseducation = d.id_divisionseducation
	left join contractor c on l.id_contractor = c.id_contractor
	left join passport pasc on c.id_passport = pasc.id_passport
	left join registrationaddress rc on c.id_regaddress = rc.id_regaddress
	where e.id_listener = $1 and e.id_programeducation = $2;
	
	
	`

	info := entity.PersonalCardInfo{}
	err := d.repo.GetContext(ctx, &info, query, idListener, idProgram)
	repoutils.HandleRepoErr(errors.New("формируется юрик или ошибка поиска слушателя"))

	queryExecutor :=
		`
	select e.status, e.first_name, e.second_name, e.middle_name
	from executor as e
	where e.id_executor = $1
	`

	executor := entity.Executor{}
	err = d.repo.GetContext(ctx, &executor, queryExecutor, idExecutor)
	if err != nil {
		fmt.Println("1234")
		return nil, repoutils.HandleRepoErr(err)
	}

	queryContractor :=
		`
	select 
	c.first_name, c.second_name, c.middle_name, c.contact_phone, c.email,
	pasc.place_birth, pasc.citizenship, pasc.gender, pasc.seria, pasc.number, pasc.passport_given, pasc.date_given, pasc.code,
    rc.mail_index, rc.region, rc.city, rc.street, rc.house, rc.building, rc.apartment
    from enrollmentlistener as e
    inner join listener l on e.id_listener = l.id_listener
    left join contractor c on l.id_contractor = c.id_contractor
	left join passport pasc on c.id_passport = pasc.id_passport
	left join registrationaddress rc on c.id_regaddress = rc.id_regaddress
	where e.id_listener = $1 and e.id_programeducation = $2;
	`
	contractor := entity.ContractorDoc{}
	err = d.repo.GetContext(ctx, &contractor, queryContractor, idListener, idProgram)
	repoutils.HandleRepoErr(errors.New("формируется юрик или нет контрактора"))

	prepareInfo := entity.FullDocument{
		PersonalInfo: info,
		Executor:     executor,
		Contractor:   contractor,
	}

	return &prepareInfo, nil
}
