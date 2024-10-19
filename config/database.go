package config

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func NewDBConnection(dbConnString string, ctx context.Context) *pgxpool.Pool {
	connPool, err := pgxpool.New(ctx, dbConnString)
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to connect to database")
	}

	log.Info().Msg("Connected to database")

	return connPool
}

func CloseDBConnection(connPool *pgxpool.Pool, ctx context.Context) {
	connPool.Close()
	log.Info().Msg("Closed database connection")
}
