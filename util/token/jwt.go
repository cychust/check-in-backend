package token

import (
	"check-in-backend/config"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// GetJWTInfo 获取 payload
func GetJWTInfo(c echo.Context) jwt.MapClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims
}

func GetJWTToken(data map[string]interface{}) string {
	t := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	claims := t.Claims.(jwt.MapClaims)
	for key, value := range data {
		claims[key] = value
	}
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	jwtToken, _ := t.SignedString([]byte(config.Conf.Security.Secret))

	return jwtToken
}

func GetJWTTokenWithClaims(claims jwt.Claims) string {
	t := jwt.New(jwt.SigningMethodHS256)
	t.Claims = claims
	jwtToken, _ := t.SignedString([]byte(config.Conf.Security.Secret))
	return jwtToken
}
