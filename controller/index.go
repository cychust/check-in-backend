package controller

import (
	"check-in-backend/constant"
	"check-in-backend/controller/param"
	"check-in-backend/model"
	"check-in-backend/util"
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

	writeIndexLog("login", "sss", err)

	//todo
	weixinSessRes, err := model.GetWeixinSession(data.Code)
	if err != nil {
		writeIndexLog("GetWeiXinSession", constant.ErrorMsgParamWrong, err)
		return retError(c, http.StatusBadRequest, constant.ErrorMsgParamWrong)
	}

	writeIndexLog("login", weixinSessRes.SessionKey, err)

	var userInfo *util.DecryptUserInfo
	if weixinSessRes.Unionid == "" {
		//userInfo, err = model.DecryptWeixinEncryptedData(weixinSessRes.SessionKey, data.EncryptedData, data.Iv)
		//if err != nil {
		//	writeIndexLog("DecryptWei", constant.ErrorMsgParamWrong, err)
		//	return retError(c, http.StatusBadRequest, constant.ErrorMsgParamWrong)
		//}
		userInfo = &util.DecryptUserInfo{
			UnionID:   "111111111111",
			OpenID:    "1111111111111",
			NickName:  data.UserInfo.Nickname,
			Gender:    data.UserInfo.Gender,
			Province:  data.UserInfo.Province,
			City:      data.UserInfo.City,
			Country:   data.UserInfo.Country,
			AvatarURL: data.UserInfo.AvatarURL,
			Language:  data.UserInfo.Language,
		}
	} else {
		userInfo = &util.DecryptUserInfo{
			UnionID:   weixinSessRes.Unionid,
			OpenID:    weixinSessRes.Openid,
			NickName:  data.UserInfo.Nickname,
			Gender:    data.UserInfo.Gender,
			Province:  data.UserInfo.Province,
			City:      data.UserInfo.City,
			Country:   data.UserInfo.Country,
			AvatarURL: data.UserInfo.AvatarURL,
			Language:  data.UserInfo.Language,
		}
	}

	err = model.CreateUser(userInfo)
	if err != nil {
		writeIndexLog("Login", "创建用户错误", err)
		return retError(c, http.StatusBadGateway, "创建用户错误")
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
