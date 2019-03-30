package controller

import (
	"check-in-backend/constant"
	"check-in-backend/controller/param"
	"check-in-backend/model"
	"github.com/labstack/echo"
	"net/http"
)


/**
 * @apiDefine CreateGroup  CreateGroup
 * @apiDescription 创建圈子
 *
 * @apiParam {String} title  圈子title
 *
 * @apiParamExample  {json} Request-Example:
 *     {
 *       "title":"今日英语"
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
 *           "code": String,
 *         }
 *     }
 * @apiError {Number} status 状态码
 * @apiError {String} err_msg 错误信息
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 401 Unauthorized
 *     {
 *       "status": 401,
 *       "err_msg": "Unauthorized"
 *     }
 */
/**
 * @api {post} /api/group CreateGroup
 * @apiVersion 1.0.0
 * @apiName CreateGroup
 * @apiGroup  group
 * @apiUse CreateGroup
 */

func CreateGroup(c echo.Context) (error) {
	data := param.TitleParm{}
	err := c.Bind(&data)
	if err != nil {
		WriteGroupLog("CreateGroup", constant.ErrorMsgParamWrong, err)
		return retError(c, http.StatusBadRequest, constant.ErrorMsgParamWrong)
	}
	if data.Title == "" {
		WriteGroupLog("CreateGroup", constant.ErrorMsgParamWrong, err)
		return retError(c, http.StatusBadRequest, constant.ErrorMsgParamWrong)
	}
	userId := getJWTUserID(c)
	code, err := model.CreateGroup(userId, data.Title)
	if err != nil {
		writeLog("group.go", "CreateGroup", "创建群组失败", err)
		return retError(c, http.StatusBadRequest, err.Error())
	}
	resData := map[string]string{
		"code": code,
	}
	return retData(c, resData)

}

func WriteGroupLog(funcName, errMsg string, err error) {
	writeLog("group.go", funcName, errMsg, err)
}
