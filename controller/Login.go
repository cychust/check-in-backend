package controller

import (
	"check-in-backend/constant"
	"check-in-backend/model/db"
	"check-in-backend/model"
	"check-in-backend/util"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

/**

 */

type PostLoginData struct {
	UnionID  string `json:"union_id" query:"union_id"`
	NickName string `json:"nickname" query:"nickname"`
}

func Login(c echo.Context) error {
	//data := model.User{}
	//err := c.Bind(&data)
	//if err != nil {
	//return retError(c, http.StatusBadRequest, constant.ErrorMsgParamWrong)
	//}json

	data := PostLoginData{}

	err := c.Bind(&data)

	if err != nil {
		return retError(c, http.StatusBadRequest, constant.ErrorMsgParamWrong)
	}

	log.Println(data)
	query := bson.M{
		"unionid":  data.UnionID,
		"nickname": data.NickName,
	}

	cntrl := db.NewCloneMgoDBCntlr()
	defer cntrl.Close()
	table := cntrl.GetTable("user")
	user := model.User{}

	err = table.Find(query).One(&user)

	if err != nil {
		return retError(c, http.StatusBadGateway, "查找用户错误")
	}
	log.Println(user)

	userInfo := &util.DecryptUserInfo{
		OpenID:   "1",
		UnionID:  "1",   //  注意是unionId!
		NickName: "cyc", // 注意是nickName!
		Gender:   1,
		City:     "henan",
		Province: "henam",
		Country:  "china",
		Language: "chinese",
	}
	err = model.CreateUser(userInfo)
	if err != nil {
		return retError(c, http.StatusBadGateway, "创建用户错误")
	}
	return retData(c, "ccc")
}
