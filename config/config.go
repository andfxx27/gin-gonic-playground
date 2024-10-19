package config

import (
	"github.com/joeshaw/envdecode"
	"github.com/rs/zerolog/log"
)

type Conf struct {
	Server  ConfServer
	JWT     ConfJWT
	Caching ConfCaching
	DB      ConfDB
}

type ConfServer struct {
	ApplicationName string `env:"APPLICATION_NAME,required"`
	Port            int    `env:"SERVER_PORT,required"`
	ClientPort      int    `env:"CLIENT_PORT,required"`
}

type ConfJWT struct {
	JWTSecret string `env:"JWT_SECRET,required"`
}

type ConfCaching struct {
	Addr     string `env:"CACHING_ADDR,required"`
	Password string `env:"CACHING_PASS"`
}

type ConfDB struct {
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT,required"`
	Username string `env:"DB_USER,required"`
	Password string `env:"DB_PASS,required"`
	DBName   string `env:"DB_NAME,required"`
}

func NewConf() *Conf {
	var conf Conf
	if err := envdecode.StrictDecode(&conf); err != nil {
		log.Fatal().Err(err).Msg("Error decoding environment variable into config.")
	}

	return &conf
}
