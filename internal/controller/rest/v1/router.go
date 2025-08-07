// Package v1 implements routing paths. Each services in own file.
package v1

import (
	"net/http"

	_ "github.com/Homyakadze14/AFFARM_tz/docs"
	"github.com/Homyakadze14/AFFARM_tz/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Swagger spec:
// @title       TestSote
// @description RestAPI for test site
// @version     1.0
// @host        localhost:8080
// @BasePath    /api/v1
func NewRouter(handler *gin.Engine, h *usecase.CryptocurrencyService) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Set cors
	corsConf := cors.DefaultConfig()
	corsConf.AllowOrigins = []string{"http://localhost:5173", "http://147.45.235.14:5173"}
	corsConf.AllowCredentials = true
	handler.Use(cors.New(corsConf))

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	handler.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Routers
	g := handler.Group("/api/v1")
	{
		NewHellotRoutes(g, h)
	}
}
