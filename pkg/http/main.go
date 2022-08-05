package http

import (
	"fmt"

	healthRoutes "github.com/michaelchandrag/botfood-go/pkg/health/routes"
	utils "github.com/michaelchandrag/botfood-go/utils"

	gin "github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

type ServerHTTP struct {
	Router *gin.Engine
}

const (
	GIN_DEBUG_MODE_RELEASE = "release"
	GIN_DEBUG_MODE_DEBUG = "debug"
)

func ServeHTTP (port string) error {

	serverHTTP := ServerHTTP{}
	serverHTTP.Router = gin.Default()

	debugMode := utils.GetEnv("BOTFOOD_APP_ENV", "dev")
	if debugMode == "prod" {
		gin.SetMode(GIN_DEBUG_MODE_RELEASE)
	} else {
		gin.SetMode(GIN_DEBUG_MODE_DEBUG)
	}
	serverHTTP.Router.Use(CORSMiddleware())
	serverHTTP.Router.Use(gin.Recovery())
	serverHTTP.initRoutes()
	if port == "" {
		port = "8080"
	}
	serverHTTP.Router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "DELETE", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
		AllowOrigins:     []string{"*"},
		MaxAge:           86400,
	}))
	serverHTTP.Router.Run(fmt.Sprintf(":%s", port))

	return nil
}

func (serverHTTP *ServerHTTP) initRoutes () {
	healthRoutes.SetupRoutes(serverHTTP.Router)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}