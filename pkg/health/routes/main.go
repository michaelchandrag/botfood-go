package health

import (
	"net/http"

	log "github.com/michaelchandrag/botfood-go/lib/log"
	"github.com/gin-gonic/gin"
)

func SetupRoutes (r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		log.Logger.Info("API HEALTH")
		c.JSON(http.StatusOK, gin.H{
  			"message": "pong",
		})
	})
}