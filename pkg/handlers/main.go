package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/michaelchandrag/botfood-go/infrastructures/database"
	health_service "github.com/michaelchandrag/botfood-go/pkg/modules/health/services"
	me_service "github.com/michaelchandrag/botfood-go/pkg/modules/me/services"
	consumer_service "github.com/michaelchandrag/botfood-go/pkg/modules/openapi/services"
	report_service "github.com/michaelchandrag/botfood-go/pkg/modules/report/services"
)

type Handler struct {
	c               *gin.Context
	dao             database.MainDB
	healthService   health_service.HealthService
	reportService   report_service.ReportService
	meService       me_service.MeService
	consumerService consumer_service.ConsumerService
}

func NewHTTPHandler(mainDB *sqlx.DB) Handler {

	db := database.NewDB(mainDB)

	var healthService health_service.HealthService
	var reportService report_service.ReportService
	var meService me_service.MeService
	var consumerService consumer_service.ConsumerService

	healthService = health_service.RegisterHealthService()
	reportService = report_service.RegisterReportService(db)
	meService = me_service.RegisterMeService(db)
	consumerService = consumer_service.RegisterConsumerService(db)

	return Handler{
		dao:             db,
		healthService:   healthService,
		reportService:   reportService,
		meService:       meService,
		consumerService: consumerService,
	}
}
