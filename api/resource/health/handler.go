package health

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"net/http"
)

func NewHandler(dbConnPool *pgxpool.Pool, redisClient *redis.Client) Handler {
	return &handler{
		dbConnPool,
		redisClient,
	}
}

type Handler interface {
	HealthCheck(c *gin.Context)
}

type handler struct {
	dbConnPool  *pgxpool.Pool
	redisClient *redis.Client
}

func (h *handler) HealthCheck(c *gin.Context) {
	err := h.dbConnPool.Ping(c)
	if err != nil {
		log.Err(err).Msg("Service health check failed, error when pinging to db.")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "Service Not OK, connection to db unavailable",
		})
		return
	}

	err = h.redisClient.Ping(c).Err()
	if err != nil {
		log.Err(err).Msg("Service health check failed, error when pinging to redis client.")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"message": "Service Not OK, connection to redis unavailable",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Service OK",
	})
	return
}
