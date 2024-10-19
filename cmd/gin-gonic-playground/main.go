package main

import (
	"context"
	"fmt"
	"github.com/andfxx27/gin-gonic-playground/api/router"
	"github.com/andfxx27/gin-gonic-playground/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const fmtDbString = "postgres://%s:%s@%s:%d/%s"

func main() {
	conf := config.NewConf()

	aggregateWriter := config.NewMultiLevelWriter()
	log.Logger = zerolog.New(aggregateWriter).With().Timestamp().Logger()

	ctx := context.Background()
	dbConnString := fmt.Sprintf(fmtDbString, conf.DB.Username, conf.DB.Password, conf.DB.Host, conf.DB.Port, conf.DB.DBName)
	dbConnPool := config.NewDBConnection(dbConnString, ctx)
	redisClient := config.NewRedisConnection(conf.Caching.Addr, conf.Caching.Password)

	defer config.CloseDBConnection(dbConnPool, ctx)

	r := gin.Default()
	r.Use(gin.LoggerWithWriter(aggregateWriter))
	router.RegisterRoutes(r, dbConnPool, redisClient)

	log.Fatal().Err(r.Run(fmt.Sprintf("localhost:%d", conf.Server.Port)))
}
