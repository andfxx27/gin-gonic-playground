package util

import (
	"github.com/andfxx27/gin-gonic-playground/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
)

func CreateJWTWithClaims(claims jwt.MapClaims) (string, error) {
	conf := config.NewConf()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(conf.JWT.JWTSecret))
	if err != nil {
		log.Err(err).Msg("Error signing token.")
		return "", err
	}

	return tokenString, nil
}

func ParseJWTWithClaims(tokenString string) (jwt.MapClaims, error) {
	conf := config.NewConf()
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.JWT.JWTSecret), nil
	})
	if err != nil {
		log.Err(err).Msg("Error parsing token.")
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	log.Err(err).Msg("Invalid access token.")
	return nil, err
}
