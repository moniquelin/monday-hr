package api

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (app *Application) createToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		claims)

	tokenString, err := token.SignedString([]byte(app.Config.Jwt.Secret))
	if err != nil {
		return "Unable to sign token", err
	}

	return tokenString, nil
}

func (app *Application) verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return app.Config.Jwt.Secret, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
