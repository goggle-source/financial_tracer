package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

const authPrefix = "Bearer "

func Logging(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		log.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		}).Info("request")

		c.Next()

		log.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"exp":    time.Since(start),
		}).Info("request")
	}
}

func JWToken(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, authPrefix) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "error Authorization header"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, authPrefix)
		tokenStr = strings.TrimSpace(tokenStr)

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method: %s", t.Header["alg"])
			}

			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "error token"})
			return
		}

		c.Set("jwtToken", token)

		c.Next()
	}
}
