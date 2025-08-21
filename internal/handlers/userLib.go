package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/financial_tracer/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func PostJWT(c *gin.Context, secretKey string, id int) (ResponseJSONUser, error) {
	accsessToken, err := JWTAccsessToken(secretKey, id)
	if err != nil {
		return ResponseJSONUser{}, err
	}

	refreshToken, err := JWTRefreshToken(secretKey, id)
	if err != nil {
		return ResponseJSONUser{}, err
	}

	return ResponseJSONUser{
		AccsessToken: accsessToken,
		RefreshToken: refreshToken,
	}, nil
}

func JWTAccsessToken(secretKey string, id int) (string, error) {
	const op = "handlers.JWTAccsessToken"

	payload := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 48).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return t, nil
}

func JWTRefreshToken(secretKey string, id int) (string, error) {
	const op = "handlers.JWTAccsessToken"

	payload := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 148).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return t, nil
}

func RegistrationError(c *gin.Context, op string, err error) {
	if errors.Is(err, domain.ErrorDuplicated) {
		c.JSON(http.StatusConflict, gin.H{"error": "error duplicated email"})
		return
	}
	if errors.Is(err, domain.ErrorEmail) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error email"})
		return
	}
	if errors.Is(err, domain.ErrorPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error password"})
		return
	}
	if errors.Is(err, domain.ErrorNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not found"})
		return
	}

}
