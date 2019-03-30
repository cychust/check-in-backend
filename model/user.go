package model

import (
	"check-in-backend/model/db"
	"check-in-backend/util"
	"gopkg.in/mgo.v2/bson"
	"log"
	"phs-mp-develop/src/constant"
)

type User struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	// WeixinUserInfo
	Openid   string `bson:"openid" json:"openid"`     // openid
	Unionid  string `bson:"unionid" json:"user_id"`   // unionid
	Nickname string `bson:"nickname" json:"nickname"` // 用户昵称
	Gender   int    `bson:"gender" json:"gender"`     // 性别 0：未知、1：男、2：女
	Province string `bson:"province" json:"province"` // 省份
	City     string `bson:"city" json:"city"`         // 城市
	Country  string `bson:"country" json:"country"`   // 国家
	//AvatarURL string `bson:"avatar_url" json:"avatar_url"` // 用户头像
	Language string `bson:"language" json:"language"` // 语言

	OwnGroups    []string `bson:"own_groups" json:"own_groups"`
	ManageGroups []string `bson:"manage_groups" json:"manage_groups"`
	JoinGroups   []string `bson:"join_groups" json:"join_groups"`
}

func GetUser(unionid string) (User, error) {
	query := bson.M{
		"unionid": unionid,
	}

	selector := bson.M{
		"_id":           0,
		"openid":        0,
		"unionid":       0,
		"own_groups":    0,
		"manage_groups": 0,
		"join_groups":   0,
	}
	data := User{}
	cntrl := db.NewCloneMgoDBCntlr()
	defer cntrl.Close()

	table := cntrl.GetTable("user")
	err := table.Find(query).Select(selector).One(&data)
	return data, err
}

func CreateUser(userInfo *util.DecryptUserInfo) error {
	//if userInfo.UnionID == "" {
	//	return constant.ErrorIDFormatWrong
	//}
	//query := bson.M{
	//	"unionid": userInfo.UnionID,
	//}
	//user := User{}
	//seletor := bson.M{
	//	"unionid":  1,
	//	"nickname": 1,
	//}
	cntrl := db.NewCopyMgoDBCntlr()
	defer cntrl.Close()

	table := cntrl.GetTable(constant.TableUser)

	user := User{
		ID:       bson.NewObjectId(),
		Openid:   userInfo.OpenID,
		Unionid:  userInfo.UnionID,
		Nickname: userInfo.NickName,
		Gender:   userInfo.Gender,
		Language: userInfo.Language,
		Country:  userInfo.Country,
		City:     userInfo.City,
		Province: userInfo.Province,
	}
	log.Println("aaaaaaa")
	eil := table.Insert(user)
	log.Println(eil)
	return eil

}

func AddUserOwnGroup(unionid, id string) error {
	if !bson.IsObjectIdHex(id) {
		return constant.ErrorIDFormatWrong
	}

	cntrl := db.NewCloneMgoDBCntlr()
	defer cntrl.Close()
	table := cntrl.GetTable(constant.TableUser)
	query := bson.M{
		"unionid": unionid,
	}
	update := bson.M{
		"$addToSet": bson.M{
			"own_groups": id,
		},
	}
	return table.Update(query, update)
}

func updateUser(query, update interface{}) error {
	return insertDoc(constant.TableUser, query, update)
}
