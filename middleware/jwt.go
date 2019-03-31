package mid



import (
	"check-in-backend/config"
	"check-in-backend/constant"
	jwt "github.com/dgrijalva/jwt-go"
"github.com/labstack/echo"
"github.com/labstack/echo/middleware"
)

// DefaultJWTConfig jwt配置
var DefaultJWTConfig = middleware.JWTConfig{
	SigningKey:  []byte(config.Conf.Security.Secret),
	TokenLookup: "header:" + echo.HeaderAuthorization,
	AuthScheme:  "Bearer",
	Claims:      jwt.MapClaims{},
	ContextKey:  constant.JWTContextKey,
	Skipper:     Skipper,
}

// Skipper 过滤
func Skipper(c echo.Context) bool {
	if c.Path() == "/api/v1/login" {
		return true
	}
	return false
}

// CustomJWTConfig custom jwt config
func CustomJWTConfig(skipperPaths []string, authScheme string) middleware.JWTConfig {
	if authScheme == "" {
		authScheme = "Bearer"
	}
	return middleware.JWTConfig{
		SigningKey:  []byte(config.Conf.Security.Secret),
		TokenLookup: "header:" + echo.HeaderAuthorization,
		AuthScheme:  authScheme,
		Claims:      jwt.MapClaims{},
		ContextKey:  constant.JWTContextKey,
		Skipper:     CustomSkipper(skipperPaths),
	}
}

// CustomSkipper custom skipper
func CustomSkipper(skipperPaths []string) func(c echo.Context) bool {
	return func(c echo.Context) bool {
		for _, skipperPath := range skipperPaths {
			if c.Path() == skipperPath {
				return true
			}
		}
		return false
	}
}

