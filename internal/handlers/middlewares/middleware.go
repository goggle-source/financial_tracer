package middlewares

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/financial_tracer/internal/handlers/api"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func Logging(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		log.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"time":   start,
		}).Info("request")

		c.Next()

		log.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"exp":    time.Since(start),
		}).Info("response")
	}
}

func JWToken(secretKey string, log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		const authPrefix = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, authPrefix) {
			log.Error("error Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ResponseUnauthorizedError("invalid Authorization header"))
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, authPrefix)
		tokenStr = strings.TrimSpace(tokenStr)

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Error("invalid signing method")
				return nil, fmt.Errorf("invalid signing method: %s", t.Header["alg"])
			}

			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			log.Error("invalid token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ResponseUnauthorizedError("invalid token"))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Error("error convert token in claims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ResponseUnauthorizedError("invalid token claims"))
			return
		}

		if exp, ok := claims["exp"]; ok {
			if expTime, ok := exp.(float64); ok {
				if time.Now().After(time.Unix(int64(expTime), 0)) {
					log.Error("invalid term token exp")
					c.AbortWithStatusJSON(http.StatusBadRequest, api.ResponseUnauthorizedError("invalid term token exp"))

					return
				}
			} else {
				log.Error("error exp format")
				c.AbortWithStatusJSON(http.StatusUnauthorized, api.ResponseUnauthorizedError("invalid exp format"))
			}
		} else {
			log.Error("error get exp")
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ResponseUnauthorizedError("invalid get exp"))
			return
		}
		idClaim, ok := claims["id"]
		if !ok {
			log.Error("error get userID")
			c.AbortWithStatusJSON(http.StatusBadRequest, api.ResponseUnauthorizedError("invalid get userID"))
			return
		}
		switch v := idClaim.(type) {
		case float64:
			userId := uint(v)
			c.Set("userID", userId)
			c.Next()
		case string:
			id, err := strconv.ParseUint(v, 10, 32)
			if err != nil {
				log.Error("error convert userID (string in uint)")
				c.AbortWithStatusJSON(http.StatusUnauthorized, api.ResponseUnauthorizedError("invalid userID format"))
			}
			userId := uint(id)
			c.Set("userID", userId)
			c.Next()
		default:
			log.Error("error convert userID")
			c.AbortWithStatusJSON(http.StatusUnauthorized, api.ResponseUnauthorizedError("the userID is not meet the requirements any format"))
		}
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
