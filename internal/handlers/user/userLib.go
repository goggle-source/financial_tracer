package userHandlers

import (
	"fmt"
	"time"

	"github.com/financial_tracer/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

// ResponseJWTUser represents a jwt model
type ResponseJWTUser struct {
	RefreshToken string `json:"refresh_token,omitempty"`
	AccsessToken string `json:"access_token"`
}

type Claims struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func PostJWT(c *gin.Context, secretKey string, id uint, name string) (ResponseJWTUser, error) {
	accsessToken, err := JWTAccessToken(secretKey, id, name)
	if err != nil {
		return ResponseJWTUser{}, fmt.Errorf("error create access token: %s", err)
	}

	refreshToken, err := JWTRefreshToken(secretKey, id, name)
	if err != nil {
		return ResponseJWTUser{}, fmt.Errorf("error create refresh token: %s", err)
	}

	return ResponseJWTUser{
		AccsessToken: accsessToken,
		RefreshToken: refreshToken,
	}, nil
}

func JWTAccessToken(secretKey string, id uint, name string) (string, error) {
	const op = "handlers.JWTAccsessToken"

	payload := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 48).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return t, nil
}

func JWTRefreshToken(secretKey string, id uint, name string) (string, error) {
	const op = "handlers.JWTRefreshToken"

	payload := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 148).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return t, nil
}

func CheckAccess(refreshToken string, secretKey string, log *logrus.Logger) (uint, string, error) {
	const op = "handlers.CheckAccess"
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%s: %w", op, domain.ErrorInternal)
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		log.WithFields(logrus.Fields{
			"op":  op,
			"err": "error parse token",
		}).Error("error parse token")
		return 0, "", fmt.Errorf("%s, not valid method: %w", op, err)
	}

	if !token.Valid {
		log.WithFields(logrus.Fields{
			"op":  op,
			"err": "error valid token",
		}).Error("error valid token")
		return 0, "", fmt.Errorf("%s not valid token", op)
	}

	tokenClaims, ok := token.Claims.(*Claims)
	if !ok {
		log.WithFields(logrus.Fields{
			"op":  op,
			"err": "error claims",
		}).Error("error claims")
		return 0, "", fmt.Errorf("%s error convert token in claim", op)
	}

	if tokenClaims.ExpiresAt.Unix() < time.Now().Unix() {
		log.WithFields(logrus.Fields{
			"op":  op,
			"err": "error expired token",
		}).Error("error expired token")
		return 0, "", fmt.Errorf("%s the deadline has ended", op)
	}

	if tokenClaims.ID == "" {
		log.WithFields(logrus.Fields{
			"op":  op,
			"err": "error id",
		}).Error("error id")
		return 0, "", fmt.Errorf("%s error id token", op)
	}

	return tokenClaims.Id, tokenClaims.Name, nil
}
