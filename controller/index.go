package controller

import (
	"check-in-backend/constant"
	"check-in-backend/controller/param"
	"check-in-backend/model"
	"check-in-backend/util/log"
	"check-in-backend/util/token"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
)

var indexLogger = log.GetLogger()

func Login(c echo.Context) error {
	data := param.WeixinLoginData{}
	err := c.Bind(&data)
	if err != nil {
		writeIndexLog("Login", constant.ErrorMsgParamWrong, err)
		return retError(c, http.StatusBadRequest, constant.ErrorMsgParamWrong)
	}
	//todo
}

func LoginWeb(c echo.Context) error {
	data := param.CodeParam{}
	err := c.Bind(&data)
	if err != nil {
		writeIndexLog("LoginWeb", constant.ErrorMsgParamWrong, err)
		return retError(c, http.StatusBadRequest, constant.ErrorMsgParamWrong)
	}
	code := data.Code

	weixinTokenRes, err := model.GetWeixinWebAccessToken(code)
	if err != nil {
		writeIndexLog("GetWeixinWebAccessToken", "code wrong", err)
		return retError(c, http.StatusBadRequest, "code wrong")
	}
	userInfo, err := model.GetWeixinWebUserInfo(weixinTokenRes.AccessToken, weixinTokenRes.Openid)
	if err != nil {
		writeIndexLog("GetWeixinWebInfo", "get userInfo err", err)
		return retError(c, http.StatusBadRequest, "get userInfo wrong")
	}

	err = model.CreateUser(userInfo)
	if err != nil {
		writeIndexLog("GetWeixinAccess", "get userInfo faild", err)
		return retError(c, http.StatusBadGateway, "get userInfo faild")
	}
	jwtAuth := map[string]interface{}{
		"user_id": userInfo.UnionID,
	}
	jwtToken := token.GetJWTToken(jwtAuth)
	resData := map[string]interface{}{
		"jwt_token": jwtToken,
	}
	return retData(c, resData)

}



func writeIndexLog(funcName, errMsg string, err error) {
	writeLog("index.go", funcName, errMsg, err)
}

func writeLog(filename, funcName, errMsg string, err error) {
	indexLogger.WithFields(logrus.Fields{
		"package":  "controller",
		"file":     filename,
		"function": funcName,
		"err":      err,
	}).Warn(errMsg)
}
