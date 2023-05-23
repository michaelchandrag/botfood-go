package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/michaelchandrag/botfood-go/infrastructures/database"
	health_service "github.com/michaelchandrag/botfood-go/pkg/modules/health/services"
	me_service "github.com/michaelchandrag/botfood-go/pkg/modules/me/services"
	open_api_service "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/services"
	report_service "github.com/michaelchandrag/botfood-go/pkg/modules/report/services"
)

type Handler struct {
	c              *gin.Context
	dao            database.MainDB
	healthService  health_service.HealthService
	reportService  report_service.ReportService
	meService      me_service.MeService
	openApiService open_api_service.OpenApiService
}

func NewHTTPHandler(mainDB *sqlx.DB) Handler {

	db := database.NewDB(mainDB)

	var healthService health_service.HealthService
	var reportService report_service.ReportService
	var meService me_service.MeService
	var openApiService open_api_service.OpenApiService

	healthService = health_service.RegisterHealthService()
	reportService = report_service.RegisterReportService(db)
	meService = me_service.RegisterMeService(db)
	openApiService = open_api_service.RegisterOpenApiService(db)

	return Handler{
		dao:            db,
		healthService:  healthService,
		reportService:  reportService,
		meService:      meService,
		openApiService: openApiService,
	}
}
