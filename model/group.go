package model

import (
	"check-in-backend/constant"
	"check-in-backend/model/db"
	"check-in-backend/util"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"sync"
)

type Group struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	// 状态: -10 表示解散状态, 5 表示正常状态
	Status     int      `bson:"status" json:"status"`
	Code       string   `bson:"code" json:"code"`               // 圈子code -> 邀请码, unique
	AvatarURL  string   `bson:"avatar_url" json:"avatar_url"`   // 群头像
	Title      string   `bson:"title" json:"title"`             // 圈子title
	CreateTime int64    `bson:"create_time" json:"create_time"` // 创建时间
	OwnerID    string   `bson:"owner_id" json:"owner_id"`       // unionid 注：以下三种身份不会重复，如：members中不会有owner
	Managers   []string `bson:"managers" json:"managers"`       // 管理员
	Members    []string `bson:"members" json:"members"`         // 成员
	PersonNum  int      `bson:"person_num" json:"person_num"`   // 总人数：1 + 管理员人数 + 成员人数
}

var (
	groupCodeNextNumMutex sync.Mutex
)

func init() {
	InitGroupCodeNextNum()
}
func InitGroupCodeNextNum() {
	cntrl := db.NewRedisDBCntlr()
	defer cntrl.Close()

	nextNum, _ := cntrl.GETInt64(constant.RedisGroupCodeNextNum)
	if nextNum == 0 {
		cntrl.SET(constant.RedisGroupCodeNextNum, constant.RedisGroupInitNextNum)
	}
}

///***get group code****** 待读
func getGroupCode() (string, error) {
	cntrl := db.NewRedisDBCntlr()
	defer cntrl.Close()

	len, _ := cntrl.LLEN(constant.RedisGroupCodePool)
	if len == 0 {
		groupCodeNextNumMutex.Lock()
		nextNum, _ := cntrl.GETInt64(constant.RedisGroupCodeNextNum)
		if nextNum == 0 {
			return "", errors.New("next_num wrong")
		}
		codePool := make([]interface{}, constant.RedisGroupCodePoolNum)
		for i := 0; i < constant.RedisGroupCodePoolNum; i++ {
			nextNum++
			codePool[i] = string(util.Base34(uint64(nextNum)))
		}
		cntrl.RPUSH(constant.RedisGroupCodePool, codePool...)
		cntrl.SET(constant.RedisGroupCodeNextNum, nextNum)
		groupCodeNextNumMutex.Unlock()
	}

	return cntrl.LPOP(constant.RedisGroupCodePool)
}

func CreateGroup(openid, title string) (string, error) {

	var code string
	for i := 0; i < 5; i++ {
		code, _ = getGroupCode()
		if code != "" {
			break
		}
	}
	group := Group{
		ID:         bson.NewObjectId(),
		Status:     constant.GroupCommonStatus,
		CreateTime: util.GetNowTimestamp(),
		Code:       code,
		AvatarURL:  util.GetNextAvatarURL(),
		Title:      title,
		OwnerID:    openid, //todo
		PersonNum:  1,
	}
	err := insertGroups(group)

	if err != nil {
		return code, err
	}

	go func() {
		AddUserOwnGroup(openid, group.ID.Hex())
		//qrcode, _ := CreateQrcodeByGroupCode(code)
		//if qrcode.Ticket != "" {
		//	qrcode.Ticket = fmt.Sprintf(constant.URLQrcodeTicket, qrcode.Ticket)
		//	setRedisGroupQrcode(code, qrcode.URL, qrcode.Ticket)
		//}
	}()
	return code, err
}

//解散圈子
func DelGroupOwner(groupID, ownerID string) error {
	if !bson.IsObjectIdHex(groupID) {
		return constant.ErrorIDFormatWrong
	}
	cntrl := db.NewCloneMgoDBCntlr()
	defer cntrl.Close()
	groupTable := cntrl.GetTable(constant.TableGroup)
	query := bson.M{
		"_id":      bson.ObjectIdHex(groupID),
		"owner_id": ownerID,
		"status": bson.M{
			"$gte": constant.GroupCommonStatus,
		},
	}
	group := Group{}
	selector := bson.M{
		"owner_id": 1,
		"managers": 1,
		"members":  1,
	}
	log.Println(groupID)
	log.Println(ownerID)
	err := groupTable.FindId(bson.ObjectIdHex(groupID)).Select(selector).One(&group)
	if err != nil {
		log.Println(err)
		return err
	}
	update := bson.M{
		"$set": bson.M{
			"status": constant.GroupDelStatus,
		},
	}

	err = groupTable.Update(query, update)
	if err != nil {
		writeLog("group.go","1","ssss",err)
		return err
	}

	// 创建者、管理员、成员 更新
	userTable := cntrl.GetTable(constant.TableUser)
	query = bson.M{
		"openid": ownerID,
	}
	update = bson.M{
		"$pull": bson.M{
			"own_groups": groupID,
		},
	}
	userTable.Update(query, update)

	query = bson.M{
		"openid": bson.M{
			"$in": group.Managers,
		},
	}
	update = bson.M{
		"$pull": bson.M{
			"manage_groups": groupID,
		},
	}
	userTable.UpdateAll(query, update)

	query = bson.M{
		"openid": bson.M{
			"$in": group.Members,
		},
	}
	update = bson.M{
		"$pull": bson.M{
			"join_groups": groupID,
		},
	}
	userTable.UpdateAll(query, update)
	return nil
}

func GetGroup(id string) (Group, error) {
	if !bson.IsObjectIdHex(id) {
		return Group{}, constant.ErrorIDFormatWrong
	}
	query := bson.M{
		"_id": bson.ObjectIdHex(id),
		"status": bson.M{
			"$gte": constant.GroupCommonStatus,
		},
	}
	return findGroup(query, DefaultSelector)
}

func findGroup(query, selectField interface{}) (Group, error) {
	data := Group{}
	cntrl := db.NewCopyMgoDBCntlr()
	defer cntrl.Close()
	table := cntrl.GetTable(constant.TableGroup)
	err := table.Find(query).Select(selectField).One(&data)
	return data, err
}

func insertGroups(docs ...interface{}) error {
	return insertDocs(constant.TableGroup, docs...)
}

func getRedisGroupQrcode(code string) (map[string]string, error) {
	cntrl := db.NewRedisDBCntlr()
	defer cntrl.Close()
	key := fmt.Sprintf(constant.RedisGroupQrcode, code)
	return cntrl.HGETALL(key)
}

func JoinGroup(code, unionid string) error {
	return groupAction(code, unionid, false)
}

func groupAction(code, unionid string, isJoin bool) error {
	query := bson.M{
		"code": code,
		"status": bson.M{
			"$gte": constant.GroupCommonStatus,
		},
		"owner_id": bson.M{
			"$ne": unionid,
		},
		"managers": bson.M{
			"$nin": []string{unionid},
		},
	}
	group, err := findGroup(query, bson.M{"code": 1})
	if err != nil {
		return err
	}
	var update bson.M
	if isJoin {
		update = bson.M{
			"$addToSet": bson.M{
				"members": unionid,
			},
			"$inc": bson.M{
				"person_num": 1,
			},
		}
	} else {
		update = bson.M{
			"$pull": bson.M{
				"members": unionid,
			},
			"$inc": bson.M{
				"person_num": -1,
			},
		}
	}
	err = updateGroup(query, update)
	if err != nil {
		return err
	}
	query = bson.M{
		"unionid": unionid,
	}

	if isJoin {
		update = bson.M{
			"$addToSet": bson.M{
				"join_groups": group.ID.Hex(),
			},
		}
	} else {
		update = bson.M{
			"$pull": bson.M{
				"join_groups": group.ID.Hex(),
			},
		}
	}
	return updateUser(query, update)
}

func updateGroup(query, update interface{}) error {
	return updateDoc(constant.TableGroup, query, update)
}

func setRedisGroupQrcode(code, url, ticket string) error {
	cntrl := db.NewRedisDBCntlr()
	defer cntrl.Close()
	key := fmt.Sprintf(constant.RedisGroupQrcode, code)
	cntrl.HMSET(key, "url", url, "ticket", ticket)
	_, err := cntrl.EXPIRE(key, 3600*24*30)
	return err
}

func GetRedisGroupInfos(ids []string) ([]map[string]string, error) {
	cntrl := db.NewRedisDBCntlr()
	defer cntrl.Close()

	res := make([]map[string]string, len(ids))
	for i, id := range ids {
		key := fmt.Sprint(constant.RedisGroupInfo, id)
		data, err := cntrl.HGETALL(key)
		if len(data) == 0 || err != nil {
			group, err := setRedisGroupInfo(id)
			if err != nil {
				return res, err
			}
			data["title"] = group.Title
			data["code"] = group.Code
			data["owner_id"] = group.OwnerID
			data["person_num"] = string(group.PersonNum)
			data["avatar_url"] = group.AvatarURL

			userInfo, _ := GetRedisUserInfo(group.OwnerID)
			log.Println(userInfo)
			data["owner_name"] = userInfo["nickname"]
		} else {
			userInfo, _ := GetRedisUserInfo(data["owner_id"])
			data["owner_name"] = userInfo["nickname"]
		}
		data["id"] = id

		res[i] = data
	}
	return res, nil
}
func setRedisGroupInfo(id string) (group Group, err error) {
	query := bson.M{
		"_id": bson.ObjectIdHex(id),
		"status": bson.M{
			"$gte": constant.GroupCommonStatus,
		},
	}
	selector := bson.M{
		"title":      1,
		"code":       1,
		"owner_id":   1, // string   `bson:"owner_id" json:"owner_id"` // unionid 注：以下三种身份不会重复，如：members中不会有owner
		"person_num": 1, // int      `bson:"person_num" json:"person_num"`   // 总人数：1 + 管理员人数 + 成员人数
		"avatar_url": 1,
	}
	group, err = findGroup(query, selector)
	if err != nil {
		return
	}
	key := fmt.Sprint(constant.RedisGroupInfo, group.ID.Hex())

	args := []interface{}{
		"title",
		group.Title,
		"code",
		group.Code,
		"owner_id",
		group.OwnerID,
		"person_num",
		group.PersonNum,
		"avatar_url",
		group.AvatarURL,
	}
	cntrl := db.NewRedisDBCntlr()
	defer cntrl.Close()
	_, err = cntrl.HMSET(key, args...)
	cntrl.EXPIRE(key, getRedisDefaultExpire())
	return
}
