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

/**
 * @apiDefine Login  Login
 * @apiDescription 登录以及创建用户和更新用户信息
 * @apiParamExample  {json} Request-Example:
 *     {
 *        "code":code,
*		  "userInfo": {
 *           "nickName": String, // 用户昵称
 *           "gender": Number, // 性别 0：未知、1：男、2：女
 *           "province": String, // 省份
 *           "city": String, // 城市
 *           "country": String, // 国家
 *           "avatarUrl": String, // 用户头像
 *           "language": String, // 用户的语言，简体中文为zh_CN
 *         }
 *       "rawData": String, // 不包括敏感信息的原始数据字符串，用于计算签名
 *       "signature": String, // 使用 sha1( rawData + sessionkey ) 得到字符串，用于校验用户信息
 *       "encryptedData": String, // 包括敏感数据在内的完整用户信息的加密数据
 *       "iv": String, // 加密算法的初始向量
 *     }
 *
 * @apiSuccess {Number} status=200 状态码
 * @apiSuccess {Object} data 正确返回数据
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "status": 200,
 *       "data": {
 *           "jwt_token": "jwt_token",
 *         }
 *     }
 * @apiError {Number} status 状态码
 * @apiError {String} err_msg 错误信息
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 401 Unauthorized
 *     {
 *       "status": 502,
 *       "err_msg": "BadGateWay"
 *     }
 */
/**
 * @api {post} /api/v1/login Login
 * @apiVersion 1.0.0
 * @apiName Login
 * @apiGroup Index
 * @apiUse Login
 */

func Login(c echo.Context) error {

	data := param.WeixinLoginData{}
	err := c.Bind(&data)
	if err != nil {
		writeIndexLog("Login", constant.ErrorMsgParamWrong, err)
		return retError(c, http.StatusBadRequest, constant.ErrorMsgParamWrong)
	}

	//todo
	weixinSessRes, err := model.GetWeixinSession(data.Code)
	if err != nil {
		writeIndexLog("GetWeiXinSession", constant.ErrorMsgParamWrong, err)
		return retError(c, http.StatusBadRequest, constant.ErrorMsgParamWrong)
	}

	var userInfo *util.DecryptUserInfo
	if weixinSessRes.Openid == "" {
		userInfo, err = model.DecryptWeixinEncryptedData(weixinSessRes.SessionKey, data.EncryptedData, data.Iv)
		if err != nil {
			writeIndexLog("DecryptWei", constant.ErrorMsgParamWrong, err)
			return retError(c, http.StatusBadRequest, constant.ErrorMsgParamWrong)
		}
	} else {
		userInfo = &util.DecryptUserInfo{
			//UnionID:   weixinSessRes.Unionid,
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
		"user_id": userInfo.OpenID, //todo
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
		"user_id": userInfo.OpenID, //todo
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
