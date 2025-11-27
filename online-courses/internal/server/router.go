package server

import (
	"log/slog"
	"online-courses/internal/server/http/handlers"
	"online-courses/internal/server/http/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(server *gin.Engine,
	listenerHandler *handlers.ListenerHandler,
	divisionsHandler *handlers.DivisionsEducationHandler,
	educationTypeHandler *handlers.EducationTypeHandler,
	levelEducationHandler *handlers.LevelEducationHandler,
	programEducationHandler *handlers.ProgramEducationHandler,
	enrollmentHandler *handlers.EnrollmentListenerHandler,
	dashboardHandler *handlers.DashboardHandler,
	backupHandler *handlers.BackupHandler,
	reportHandler *handlers.ReportHandler,
	contractorHandler *handlers.ContractHandler,
	legalEntityHandler *handlers.LegalEntityHandler,
	Logger *slog.Logger) {

	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:8081"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := server.Group("api/v1")

	{

		api.Use(middleware.LoggerMiddleware(Logger))
		api.Use(middleware.AuthMiddleware())

		listener := api.Group("/listener")
		{

			listener.POST("/", middleware.RoleProtecteMiddleware("worker"), listenerHandler.CreateListenerHandler)
			listener.DELETE("/:id", middleware.RoleProtecteMiddleware("worker"), listenerHandler.DeleteListenerHandler)
			listener.GET("/", middleware.RoleProtecteMiddleware("worker"), listenerHandler.GetListener) // :page
			listener.GET("/details/:id", middleware.RoleProtecteMiddleware("worker"), listenerHandler.GetFullListener)
			listener.PUT("/:id", middleware.RoleProtecteMiddleware("worker"), listenerHandler.UpdateListener)
		}

		divisions := api.Group("/divisions")
		{
			divisions.POST("/", middleware.RoleProtecteMiddleware("worker"), divisionsHandler.CreateDivision)
			divisions.GET("/", middleware.RoleProtecteMiddleware("worker"), divisionsHandler.ReadDivisions) // :page
			// divisions.GET("/:id", middleware.RoleProtecteMiddleware("worker"), divisionsHandler.ReadDivisionsByID)
			divisions.PUT("/:id", middleware.RoleProtecteMiddleware("worker"), divisionsHandler.UpdateDivision)
			divisions.DELETE("/:id", middleware.RoleProtecteMiddleware("worker"), divisionsHandler.DeleteDivision)
		}

		educationType := api.Group("/educationtype")
		{
			educationType.POST("/", middleware.RoleProtecteMiddleware("worker"), educationTypeHandler.CreateEducationType)
			educationType.GET("/", middleware.RoleProtecteMiddleware("worker"), educationTypeHandler.ReadEducationType) // :page
			// educationType.GET("/:id", middleware.RoleProtecteMiddleware("worker"), educationTypeHandler.ReadEducationTypeByID)
			educationType.PUT("/:id", middleware.RoleProtecteMiddleware("worker"), educationTypeHandler.UpdateEducationType)
			educationType.DELETE("/:id", middleware.RoleProtecteMiddleware("worker"), educationTypeHandler.DeleteEducationtype)
		}

		levelEducation := api.Group("/leveleducation")
		{
			levelEducation.GET("/", middleware.RoleProtecteMiddleware("worker"), levelEducationHandler.GetLevelEducations)
		}

		programEducation := api.Group("/programeducation")
		{
			programEducation.POST("/", middleware.RoleProtecteMiddleware("worker"), programEducationHandler.CreateProgram)
			programEducation.GET("/", middleware.RoleProtecteMiddleware("worker"), programEducationHandler.ReadProgram) // :page
			programEducation.GET("/:id", middleware.RoleProtecteMiddleware("worker"), programEducationHandler.ReadByID)
			programEducation.PUT("/:id", middleware.RoleProtecteMiddleware("worker"), programEducationHandler.UpdateProgram)
			programEducation.DELETE("/:id", middleware.RoleProtecteMiddleware("worker"), programEducationHandler.DeleteProgram)
		}

		enrollment := api.Group("/enrollment")
		{
			enrollment.POST("/", middleware.RoleProtecteMiddleware("worker"), enrollmentHandler.CreateEnrollment)
			enrollment.GET("/", middleware.RoleProtecteMiddleware("worker"), enrollmentHandler.ReadEnrollment) // :page
			enrollment.PUT("/:id_listener/:id_program", middleware.RoleProtecteMiddleware("worker"), enrollmentHandler.UpdateEnrollment)
			enrollment.DELETE("/:id_listener/:id_program", middleware.RoleProtecteMiddleware("worker"), enrollmentHandler.DeleteEnrollment)
			enrollment.GET("/details/:id", middleware.RoleProtecteMiddleware("worker"), enrollmentHandler.ReadDetailListener)
			enrollment.GET("/:id", middleware.RoleProtecteMiddleware("worker"), enrollmentHandler.ReadByProgram) // :page
			enrollment.GET("/application-data-card", middleware.RoleProtecteMiddleware("worker"), enrollmentHandler.GetInfoToCreateCard)
		}

		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("/user", middleware.RoleProtecteMiddleware("worker", "accountant"), dashboardHandler.UserDashboard)
			dashboard.GET("/admin", middleware.RoleProtecteMiddleware("admin"), dashboardHandler.AdminDashboard)
		}

		backupDB := api.Group("/backup")
		{
			backupDB.GET("/", middleware.RoleProtecteMiddleware("admin"), backupHandler.BackupDB)
			backupDB.GET("/restore", middleware.RoleProtecteMiddleware("admin"), backupHandler.RestoreDB)
			backupDB.GET("/all", middleware.RoleProtecteMiddleware("admin"), backupHandler.AllBackup)
		}

		report := api.Group("/report")
		{
			report.GET("/period", middleware.RoleProtecteMiddleware("accountant"), reportHandler.ReportPeriod)
			report.GET("/expensive", middleware.RoleProtecteMiddleware("accountant"), reportHandler.MostExpensiveProgram)
		}

		contractor := api.Group("/contractor")
		{

			contractor.POST("/:id", middleware.RoleProtecteMiddleware("worker"), contractorHandler.CreateContract)
			contractor.DELETE("/:id", middleware.RoleProtecteMiddleware("worker"), contractorHandler.Delete)
			contractor.PUT("/:id", middleware.RoleProtecteMiddleware("worker"), contractorHandler.UpdateContractor)
		}

		legalEntity := api.Group("/legalentity")
		{

			legalEntity.POST("/", middleware.RoleProtecteMiddleware("worker"), legalEntityHandler.CreateLegalEntity)
			legalEntity.DELETE("/:id", middleware.RoleProtecteMiddleware("worker"), legalEntityHandler.DeleteLegalEntity)
			legalEntity.GET("/", middleware.RoleProtecteMiddleware("worker"), legalEntityHandler.ReadLegalEntity) // :page
			legalEntity.GET("/details/:id", middleware.RoleProtecteMiddleware("worker"), legalEntityHandler.ReadFullData)
			legalEntity.PUT("/:id", middleware.RoleProtecteMiddleware("worker"), legalEntityHandler.UpdateLegalEntity)
		}

	}
}
