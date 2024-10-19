package health

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"net/http"
)

func NewHandler(dbConnPool *pgxpool.Pool) Handler {
	return &handler{
		dbConnPool,
	}
}

type Handler interface {
	HealthCheck(c *gin.Context)
}

type handler struct {
	dbConnPool *pgxpool.Pool
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

	c.JSON(200, gin.H{
		"message": "Service OK",
	})
	return
}
