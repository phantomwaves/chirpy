package main

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (cfg *apiConfig) parseJwtFromHeader(req *http.Request) (*jwt.Token, error) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("missing Authorization header")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("invalid Authorization header")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (cfg *apiConfig) generateToken(expiration, id int) (string, error) {
	var exp time.Time
	if expiration <= 0 || expiration > 86400 {
		exp = time.Now().UTC().Add(24 * time.Hour)
	} else {
		exp = time.Now().UTC().Add(time.Second * time.Duration(expiration))
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		Issuer:    "chirpy",
		Subject:   strconv.Itoa(id),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(exp),
	})

	s, err := t.SignedString([]byte(cfg.jwtSecret))

	return s, err
}
