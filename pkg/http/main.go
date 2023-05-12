package http

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/michaelchandrag/botfood-go/pkg/handlers"
	utils "github.com/michaelchandrag/botfood-go/utils"

	"github.com/gin-contrib/cors"
	gin "github.com/gin-gonic/gin"
)

type ServerHTTP struct {
	DB     *sqlx.DB
	Router *gin.Engine
}

const (
	GIN_DEBUG_MODE_RELEASE = "release"
	GIN_DEBUG_MODE_DEBUG   = "debug"
)

func ServeHTTP(port string, db *sqlx.DB) error {

	serverHTTP := ServerHTTP{
		DB:     db,
		Router: gin.Default(),
	}

	debugMode := utils.GetEnv("BOTFOOD_APP_ENV", "dev")
	if debugMode == "prod" {
		gin.SetMode(GIN_DEBUG_MODE_RELEASE)
	} else {
		gin.SetMode(GIN_DEBUG_MODE_DEBUG)
	}
	serverHTTP.setupHTTPRequestCORS()
	serverHTTP.Router.Use(buildHTTPResponseCORS())
	serverHTTP.Router.Use(gin.Recovery())

	serverHTTP.Router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{})
	})

	// build http handler first, to inject inside routes
	handler := handlers.NewHTTPHandler(serverHTTP.DB)

	serverHTTP.registerRoutes(handler)
	if port == "" {
		port = "8080"
	}
	serverHTTP.Router.Run(fmt.Sprintf(":%s", port))

	return nil
}

// incoming request
func (serverHTTP *ServerHTTP) setupHTTPRequestCORS() {
	serverHTTP.Router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "DELETE", "PUT", "PATCH"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
		AllowOrigins:     []string{"*"},
		MaxAge:           86400,
	}))
}

// sending response
func buildHTTPResponseCORS() gin.HandlerFunc {
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
