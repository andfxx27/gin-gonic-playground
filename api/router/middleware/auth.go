package middleware

import (
	"github.com/andfxx27/gin-gonic-playground/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

func AuthorizedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.Request.Header.Get("Authorization")
		if authorizationHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": http.StatusText(http.StatusUnauthorized),
			})
			return
		}

		if !strings.HasPrefix(authorizationHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": http.StatusText(http.StatusUnauthorized),
			})
			return
		}

		tokenString := authorizationHeader[7:]
		claims, err := util.ParseJWTWithClaims(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": http.StatusText(http.StatusUnauthorized),
			})
			return
		}

		c.Set("subject", claims["sub"])

		log.Info().Msg("Passing auth middleware")
		c.Next()
	}
}
