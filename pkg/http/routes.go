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
}
