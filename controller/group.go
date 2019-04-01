package controller

import (
	"check-in-backend/constant"
	"check-in-backend/controller/param"
	"check-in-backend/model"
	"github.com/labstack/echo"
	"log"
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
 * @api {post} /api/v1/group CreateGroup
 * @apiVersion 1.0.0
 * @apiName CreateGroup
 * @apiGroup  group
 * @apiUse CreateGroup
 */

func CreateGroup(c echo.Context) (error) {
	data := param.TitleParm{}
	err := c.Bind(&data)
	log.Println(data)
	if err != nil {
		WriteGroupLog("CreateGroup", constant.ErrorMsgParamWrong, err)
		return retError(c, http.StatusBadRequest, constant.ErrorMsgParamWrong)
	}
	if data.Title == "" {
		WriteGroupLog("CreateGroup", constant.ErrorMsgParamWrong, err)
		return retError(c, http.StatusBadRequest, constant.ErrorMsgParamWrong)
	}
	userId := getJWTUserID(c) //openid
	code, err := model.CreateGroup(userId, data.Title)
	if err != nil {
		writeLog("group.go", "CreateGroup", "创建群组失败", err)
		return retError(c, http.StatusBadRequest, err.Error())
	}
	resData := map[string]string{
		"code": code,
	}
	log.Println("aaaaa")
	return retData(c, resData)

}

/**
 * @apiDefine GetGroups  GetGroups
 * @apiDescription    获取所有群组信息
 *
 * @apiSuccess {Number} status=200 状态码
 * @apiSuccess {Object} data 正确返回数据
 *
 * @apiSuccessExample Success-Response:
 *         HTTP/1.1 200 OK
 *     {
 *       "status": 200,
 *       "data": {
 *           "own_groups": [{
 *               "id": String,
 *               "title": "软件1601",
 *               "code": "唯一群code，同时也是邀请码",
 *             }]
 *           "manage_groups": [{
 *               "id": String,
 *               "title": "软件1601",
 *               "code": "唯一群code，同时也是邀请码",
 *             }]
 *           "join_groups": [{
 *               "id": String,
 *               "title": "软件1601",
 *               "code": "唯一群code，同时也是邀请码",
 *             }]
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
 * @api {get} /api/group/list GetGroup
 * @apiVersion 1.0.0
 * @apiName GetGroups
 * @apiGroup group
 * @apiUse GetGroups
 */

func GetGroups(c echo.Context) error {
	userID := getJWTUserID(c)
	ownGroups, manageGroups, joinGroups, err := findGroupInfosByUserID(userID)
	if err != nil {
		writeIndexLog("GetGroup", "查询群组信息错误", err)
		return retError(c, http.StatusBadGateway, "查询群组信息错误")
	}
	resData := map[string]interface{}{
		"own_groups":     ownGroups,
		"manager_groups": manageGroups,
		"join_groups":    joinGroups,
	}
	return retData(c, resData)
}

/**
 * @apiDefine DeleteGroup  DeleteGroup
 * @apiDescription     解散圈子
 *
 * @apiSuccess {Number} status=200 状态码
 * @apiSuccess {Object} data 正确返回数据
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "status": 200,
 *       "data": {
 *          "group_id":"group_id",
 *			"owner_id":"owner_id"
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
 * @api {post} /api/v1/group/delete DeleteGroup
 * @apiVersion 1.0.0
 * @apiName DeleteGroup
 * @apiGroup group
 * @apiUse DeleteGroup
 */

func DeleteGroups(c echo.Context) error {
	data := param.DeleteParam{}
	err := c.Bind(&data)
	if err != nil {
		writeIndexLog("DeleteGroups", constant.ErrorMsgParamWrong, err)
		return retError(c, http.StatusBadRequest, constant.ErrorMsgParamWrong)
	}
	err = model.DelGroupOwner(data.GroupID, data.OwnerId)
	if err != nil {
		writeIndexLog("DeleteGroups", "删除圈子错误", err)
		return retError(c, http.StatusBadRequest, "删除圈子错误")
	}
	return retData(c, "")
}

func findGroupInfosByUserID(userID string) ([]map[string]string, []map[string]string, []map[string]string, error) {
	ownGroups, manageGroups, joinGroups, err := model.FindGroupsByUserID(userID)
	if err != nil {
		return nil, nil, nil, err
	}
	groups := append(append(append([]string{}, ownGroups...), manageGroups...), joinGroups...)
	groupInfos, err := model.GetRedisGroupInfos(groups)
	if err != nil {
		return nil, nil, nil, err
	}
	ownGroupsLen := len(ownGroups)
	manageGroupsLen := len(manageGroups)
	return groupInfos[:ownGroupsLen], groupInfos[ownGroupsLen : ownGroupsLen+manageGroupsLen], groupInfos[ownGroupsLen+manageGroupsLen:], nil
}

func WriteGroupLog(funcName, errMsg string, err error) {
	writeLog("group.go", funcName, errMsg, err)
}
