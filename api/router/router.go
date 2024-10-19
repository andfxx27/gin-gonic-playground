package router

import (
	"github.com/andfxx27/gin-gonic-playground/api/resource/health"
	"github.com/andfxx27/gin-gonic-playground/api/resource/member"
	"github.com/andfxx27/gin-gonic-playground/api/router/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(r *gin.Engine, dbConnPool *pgxpool.Pool, redisClient *redis.Client) {
	mainGroup := r.Group("/api/v1")

	// Route health
	healthHandler := health.NewHandler(dbConnPool)
	healthGroup := mainGroup.Group("/health")
	healthGroup.GET("/health-check", healthHandler.HealthCheck)

	// Route member
	memberRepository := member.NewRepositorier(dbConnPool)
	memberHandler := member.NewHandler(memberRepository, redisClient)
	memberGroup := mainGroup.Group("/member")
	memberGroup.POST("/sign-up", memberHandler.SignUp)
	memberGroup.POST("/sign-in", memberHandler.SignIn)

	memberGroup.Use(middleware.AuthorizedMiddleware())
	memberGroup.GET("/profile", memberHandler.GetProfile)
}
