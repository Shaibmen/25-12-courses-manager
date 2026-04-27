package main

import (
	"log"
	"online-courses/internal/config"
	"online-courses/internal/database/postgres"
	"online-courses/internal/logger"
	"online-courses/internal/repo/pg"
	repoutils "online-courses/internal/repo/pg/repo_utils"
	"online-courses/internal/server"
	"online-courses/internal/server/http/handlers"
	"online-courses/internal/service"
	"online-courses/internal/validate"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	log.Print("Запуск API")

	cfg := config.MustLoadConfig()

	Logger := logger.MustInitLogger(cfg.Logger)
	Logger.Info("логгер инициализирован без ошибок")

	db := postgres.MustNewConnectionPostgresSQL(cfg.DB_STRING_CONN, Logger)

	validate.InitValid()
	repoutils.InitRepoLogger(Logger)

	educationListenerRepo := pg.NewEducationListenerRepo(db, Logger)
	listenerRepo := pg.NewListenerRepo(db, Logger)
	passportRepo := pg.NewPassportRepo(db, Logger)
	placeWorkRepo := pg.NewPlaceWorkRepo(db, Logger)
	registrationAddressRepo := pg.NewRegistrationAddressRepo(db, Logger)
	divisionsRepo := pg.NewDivisionsEducationRepo(db, Logger)
	educationTypeRepo := pg.NewEducationTypeRepo(db, Logger)
	levelEducationRepo := pg.NewLevelEducationRepo(db, Logger)
	programEducationRepo := pg.NewProgramEducationRepo(db, Logger)
	enrollmentRepo := pg.NewEnrollmentListenerRepo(db, Logger)
	dashboardRepo := pg.NewDashboardRepo(db, Logger)
	backupRepo := pg.NewBackupRepo(db, cfg.DBUSER, cfg.DBPASSWORD, cfg.DBNAME, cfg.DBHOST, Logger)
	reportRepo := pg.NewReportRepo(db, Logger)
	// procedureRepo := pg.NewProcedureRepo(db, Logger)
	contractorRepo := pg.NewContractorRepo(db, Logger)
	legalEntityRepo := pg.NewLegalEntity(db, Logger)
	executerRepo := pg.NewExecutorRepo(db, Logger)
	documentRepo := pg.NewDocumentsRepo(db, Logger)
	groupRepo := pg.NewGroupDB(db, Logger)
	graphicRepo := pg.NewGraphicsRepo(db, Logger)

	// go func() {
	// 	for {
	// 		procedureRepo.DeactivationNoValidEnrollment()
	// 		procedureRepo.ShuffleLevelEducation()
	// 		procedureRepo.ShuffleProgram()
	//
	// 		time.Sleep(1 * time.Hour)
	// 	}
	//
	// }()

	if _, err := os.Stat("reports/excel/"); os.IsNotExist(err) {
		os.MkdirAll("reports/excel/", 0755)
	}

	if _, err := os.Stat("backups/pg/"); os.IsNotExist(err) {
		os.MkdirAll("backups/pg/", 0755)
	}

	listenerSevice := service.NewListenerService(db, passportRepo, educationListenerRepo, placeWorkRepo, registrationAddressRepo, listenerRepo)
	divisionsService := service.NewDivisionsEducationService(divisionsRepo)
	educationTypeService := service.NewEducationTypeService(educationTypeRepo)
	levelEducationService := service.NewLevelEducationService(levelEducationRepo)
	programEducationService := service.NewProgramEducationService(programEducationRepo)
	enrollmentService := service.NewEnrollmentListenerService(enrollmentRepo)
	dashboardService := service.NewDashboardService(dashboardRepo)
	backupService := service.NewBackupService(backupRepo)
	reportService := service.NewRepostService(reportRepo)
	contractorService := service.NewContractorService(db, contractorRepo, listenerRepo, passportRepo, registrationAddressRepo)
	legalEntityService := service.NewLegalEntityService(db, legalEntityRepo, registrationAddressRepo, listenerRepo)
	executerService := service.NewExecutorService(executerRepo)
	documentService := service.NewDocumentService(documentRepo)
	groupService := service.NewGroupService(groupRepo)
	graphicService := service.NewGraphicsService(graphicRepo)

	listenerHanlder := handlers.NewListenerHandler(listenerSevice)
	divisionsHandler := handlers.NewDivisionsEducationHandler(divisionsService)
	educationTypeHandler := handlers.NewEducationTypeHandler(educationTypeService)
	levelEducationHandler := handlers.NewLevelEducationHandler(levelEducationService)
	programEducationHandler := handlers.NewProgramEducationHandler(programEducationService)
	enrollmentHandler := handlers.NewEnrollmentListenerHandler(enrollmentService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	backupHandler := handlers.NewBackupHandler(backupService)
	reportHandler := handlers.NewReportHandler(reportService)
	contractorHandler := handlers.NewContractHandler(contractorService)
	legalEntityHandler := handlers.NewLegalEntityHandler(legalEntityService)
	executerHandler := handlers.NewExecutorHandler(executerService)
	documentHandler := handlers.NewDocumentHandler(documentService, cfg)
	groupHandler := handlers.NewGroupHandler(groupService)
	graphicHandler := handlers.NewGraphicsHandler(graphicService)
	ыфсcanDiplomHandler := handlers.NewScanDiplom(cfg)

	r := gin.Default()

	server.SetupRoutes(
		r,
		listenerHanlder,
		divisionsHandler,
		educationTypeHandler,
		levelEducationHandler,
		programEducationHandler,
		enrollmentHandler,
		dashboardHandler,
		backupHandler,
		reportHandler,
		contractorHandler,
		legalEntityHandler,
		executerHandler,
		documentHandler,
		groupHandler,
		graphicHandler,
		ыфсcanDiplomHandler,
		cfg,
		Logger)

	r.Run(":8080")
}
