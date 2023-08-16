package http

import (
	"github.com/michaelchandrag/botfood-go/pkg/handlers"
)

func (serverHTTP *ServerHTTP) registerRoutes(h handlers.Handler) {

	v1 := serverHTTP.Router.Group("/api/v1")
	v1.GET("/health", h.GetHealthAction)
	v1.GET("/error", h.GetErrorAction)

	me := v1.Group("/me")
	me.Use(handlers.BasicMiddleware(h))
	me.GET("", h.GetMeAction)
	me.GET("/reviews", h.GetMeReviewsAction)

	report := v1.Group("/report")
	report.Use(handlers.BasicMiddleware(h))
	report.GET("/channel_report", h.GetChannelReportAction)
	report.GET("/promotion", h.GetPromotionReportAction)
	report.GET("/atp", h.GetATPReportAction)

	openApi := serverHTTP.Router.Group("/boa/v1")
	openApi.POST("/consume_message_queue", h.PostConsumeMessageQueueAction)

	openApiBranchChannel := openApi.Group("/branch_channel")
	openApiBranchChannel.Use(handlers.OpenApiMiddleware(h))
	openApiBranchChannel.GET("/list", h.GetOpenApiBranchChannelListAction)
	openApiBranchChannel.GET("/detail/:branch_channel_id", h.GetOpenApiBranchChannelDetailAction)

	openApiItem := openApi.Group("/item")
	openApiItem.Use(handlers.OpenApiMiddleware(h))
	openApiItem.GET("/list", h.GetOpenApiItemListAction)

	openApiVariant := openApi.Group("/variant")
	openApiVariant.Use(handlers.OpenApiMiddleware(h))
	openApiVariant.GET("/list", h.GetOpenApiVariantListAction)

	openApiReview := openApi.Group("/review")
	openApiReview.Use(handlers.OpenApiMiddleware(h))
	openApiReview.GET("/list", h.GetOpenApiReviewListAction)

	openApiReport := openApi.Group("/report")
	openApiReport.Use(handlers.OpenApiMiddleware(h))
	openApiReport.GET("/item_availability_report/list", h.GetOpenApiItemAvailabilityReportListAction)
	openApiReport.GET("/branch_channel_availability_report/list", h.GetOpenApiBranchChannelAvailabilityReportListAction)

}
