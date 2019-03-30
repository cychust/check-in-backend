package controller

import (
	"check-in-backend/constant"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"net/http"
)

// jwt
func getJWTUserID(c echo.Context) string {
	return c.Get(constant.JWTContextKey).(*jwt.Token).Claims.(jwt.MapClaims)["user_id"].(string)
}

// ErrorRes ErrorResponse
type ErrorRes struct {
	Status int    `json:"status"`
	ErrMsg string `json:"err_msg"`
}

// DataRes DataResponse
type DataRes struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

// RetError response error, wrong response
func retError(c echo.Context, status int, errMsg string) error {
	// return c.JSON(code, ErrorRes{
	// 	Status: status,
	// 	ErrMsg: errMsg,
	// })
	return c.JSON(http.StatusBadGateway, ErrorRes{
		Status: status,
		ErrMsg: errMsg,
	})
}

// RetData response data, correct response
func retData(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, DataRes{
		Status: 200,
		Data:   data,
	})
}

