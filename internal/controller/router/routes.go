package routes

import (
	"fmt"
	docs "github.com/backstagefood/video-processor-uploader/docs/http"
	"github.com/backstagefood/video-processor-uploader/internal/controller/handlers"
	"github.com/backstagefood/video-processor-uploader/pkg/adapter"
	"github.com/backstagefood/video-processor-uploader/utils"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func NewRouter(connectionManager adapter.ConnectionManager) *gin.Engine {
	r := gin.Default()
	initSwagger()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// Grupo de rotas /api com middleware
	apiGroup := r.Group("/v1")
	apiGroup.Use(apiAuthMiddleware())
	{
		videoHandler := handlers.NewVideoHandler(connectionManager.GetMessageProducer(), connectionManager.GetBucketConn())
		apiGroup.POST("/upload", videoHandler.HandleVideoUpload)
	}

	// outros
	r.GET("/info", handlers.HandleInfo)
	r.GET("/health", handlers.HandleHealth)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// logger
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s %s %d | %10s | %15s | %-7s \"%s\"\n",
			param.TimeStamp.Format("2006/01/02 15:04:05"),
			"INFO",
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
		)
	}))

	return r
}

// Middleware para rotas /api/*
func apiAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userEmail := c.GetHeader("X-User-Email")
		if userEmail == "" {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Header X-User-Email é obrigatório",
			})
			return
		}
		c.Set("user_email", userEmail)
		c.Next()
	}
}

func initSwagger() {
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Version = handlers.Version
	docs.SwaggerInfo.Host = utils.GetEnvVarOrDefault("SWAGGER_HOST", "localhost:8080")
}
