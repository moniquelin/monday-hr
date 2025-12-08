package api

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// createToken creates JWT tokens for user login
func (app *Application) createToken(userId int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": userId,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		claims)

	tokenString, err := token.SignedString([]byte(app.Config.Jwt.Secret))
	if err != nil {
		return "Unable to sign token", err
	}

	return tokenString, nil
}
