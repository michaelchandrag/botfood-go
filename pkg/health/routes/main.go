package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes (r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
  			"message": "pong",
		})
	})
}