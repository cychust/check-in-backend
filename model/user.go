package model

import (
	"check-in-backend/constant"
	"check-in-backend/model/db"
	"check-in-backend/util"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
)

type User struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	// WeixinUserInfo
	Openid string `bson:"openid" json:"openid"` // openid
	//Unionid   string `bson:"unionid" json:"user_id"`       // unionid
	Nickname  string `bson:"nickname" json:"nickname"`     // 用户昵称
	Gender    int    `bson:"gender" json:"gender"`         // 性别 0：未知、1：男、2：女
	Province  string `bson:"province" json:"province"`     // 省份
	City      string `bson:"city" json:"city"`             // 城市
	Country   string `bson:"country" json:"country"`       // 国家
	AvatarURL string `bson:"avatar_url" json:"avatar_url"` // 用户头像
	Language  string `bson:"language" json:"language"`     // 语言

	OwnGroups    []string `bson:"own_groups" json:"own_groups"`
	ManageGroups []string `bson:"manage_groups" json:"manage_groups"`
	JoinGroups   []string `bson:"join_groups" json:"join_groups"`
}

func GetUser(unionid string) (User, error) {
	query := bson.M{
		"unionid": unionid,
	}

	selector := bson.M{
		"_id":    0,
		"openid": 0,
		//"unionid":       0,
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
	if userInfo.AvatarURL == "" {
		userInfo.AvatarURL = constant.WechatDefaultHeadImgURL
	}

	query := bson.M{
		//"unionid": userInfo.UnionID,
		"openid": userInfo.OpenID,
	}
	user := User{}
	selector := bson.M{
		"unionid":    1,
		"nickname":   1,
		"avatar_url": 1,
	}

	cntrl := db.NewCloneMgoDBCntlr()
	defer cntrl.Close()
	table := cntrl.GetTable(constant.TableUser)

	err := table.Find(query).Select(selector).One(&user)

	log.Println(user)

	if err != nil {
		user = User{
			ID:     bson.NewObjectId(),
			Openid: userInfo.OpenID,
			//Unionid:   userInfo.UnionID,
			Nickname:  userInfo.NickName,
			AvatarURL: userInfo.AvatarURL,
			Gender:    userInfo.Gender,
			Language:  userInfo.Language,
			Country:   userInfo.Country,
			City:      userInfo.City,
			Province:  userInfo.Province,
		}
		return table.Insert(user)
	}

	// update
	if userInfo.NickName != user.Nickname || userInfo.AvatarURL != user.AvatarURL {
		updateMap := map[string]interface{}{
			//"status":     status,
			"nickname":   userInfo.NickName,
			"avatar_url": userInfo.AvatarURL,
			"gender":     userInfo.Gender,
			"language":   userInfo.Language,
			"country":    userInfo.Country,
			"city":       userInfo.City,
			"province":   userInfo.Province,
		}
		if userInfo.OpenID != "" {
			updateMap["openid"] = userInfo.OpenID
		}
		update := bson.M{
			"$set": updateMap,
		}
		return table.Update(query, update)
	}
	return nil

}

func AddUserOwnGroup(openid, id string) error {
	if !bson.IsObjectIdHex(id) {
		return constant.ErrorIDFormatWrong
	}

	cntrl := db.NewCloneMgoDBCntlr()
	defer cntrl.Close()
	table := cntrl.GetTable(constant.TableUser)
	query := bson.M{
		//"unionid": unionid,
		"openid": openid,
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

func FindGroupsByUserID(openid string) (ownGroups, manageGroups, joinGroups []string, err error) {
	if openid == "" {
		err = constant.ErrorIDFormatWrong
		return
	}

	query := bson.M{
		"openid": openid,
	}
	selector := bson.M{
		"own_groups":    1,
		"manage_groups": 1,
		"join_groups":   1,
	}

	user, err := findUser(query, selector)
	ownGroups, manageGroups, joinGroups = user.OwnGroups, user.ManageGroups, user.JoinGroups
	return
}

func findUser(query, selector interface{}) (User, error) {
	data := User{}
	cntrl := db.NewCloneMgoDBCntlr()
	defer cntrl.Close()
	table := cntrl.GetTable(constant.TableUser)
	err := table.Find(query).Select(selector).One(&data)
	return data, err
}
func GetRedisUserInfos(openids []string) ([]map[string]string, error) {
	redisConn := db.NewRedisDBCntlr()
	defer redisConn.Close()

	resData := make([]map[string]string, len(openids))
	for i, openid := range openids {
		key := fmt.Sprintf(constant.RedisUserInfo, openid)
		userInfo, err := redisConn.HGETALL(key)
		if len(userInfo) == 0 || err != nil {
			user, _ := setRedisUserInfo(openid)
			userInfo = map[string]string{
				"nickname":   user.Nickname,
				"gender":     strconv.Itoa(user.Gender),
				"province":   user.Province,
				"city":       user.City,
				"country":    user.Country,
				"avatar_url": user.AvatarURL,
				"language":   user.Language,
			}
		}
		userInfo["user_id"] = openid
		resData[i] = userInfo
	}
	return resData, nil
}

func GetRedisUserInfo(openid string) (map[string]string, error) {
	redisConn := db.NewRedisDBCntlr()
	defer redisConn.Close()
	key := fmt.Sprintf(constant.RedisUserInfo, openid)
	userInfo, err := redisConn.HGETALL(key)
	if len(userInfo) == 0 || err != nil {
		user, _ := setRedisUserInfo(openid)
		userInfo = map[string]string{
			"nickname":   user.Nickname,
			"gender":     strconv.Itoa(user.Gender),
			"province":   user.Province,
			"city":       user.City,
			"country":    user.Country,
			"avatar_url": user.AvatarURL,
			"language":   user.Language,
		}
		log.Println("aaa")
		log.Println(userInfo)
	}
	userInfo["user_id"] = openid
	resData := userInfo
	return resData, nil
}

func setRedisUserInfo(openid string) (User, error) {
	query := bson.M{
		"openid": openid,
	}
	selector := bson.M{
		"nickname":   1,
		"gender":     1,
		"province":   1,
		"city":       1,
		"country":    1,
		"avatar_url": 1,
		"language":   1,
	}
	user, err := findUser(query, selector)
	if err != nil || user.Nickname == "" || user.AvatarURL == "" {
		return user, err
	}

	cntrl := db.NewRedisDBCntlr()
	defer cntrl.Close()

	key := fmt.Sprintf(constant.RedisUserInfo, openid)
	args := []interface{}{
		"nickname",
		user.Nickname,
		"gender",
		user.Gender,
		"province",
		user.Province,
		"city",
		user.City,
		"country",
		user.Country,
		"avatar_url",
		user.AvatarURL,
		"language",
		user.Language,
	}
	cntrl.HMSET(key, args...)
	cntrl.EXPIRE(key, getRedisDefaultExpire())
	return user, nil
}
