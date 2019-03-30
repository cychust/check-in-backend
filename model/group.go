package model

import (
	"check-in-backend/constant"
	"check-in-backend/model/db"
	"check-in-backend/util"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"sync"
)

type Group struct {
	ID bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	// 状态: -10 表示解散状态, 5 表示正常状态
	Status int    `bson:"status" json:"status"`
	Code   string `bson:"code" json:"code"` // 圈子code -> 邀请码, unique
	//AvatarURL  string   `bson:"avatar_url" json:"avatar_url"`   // 群头像
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

func CreateGroup(unionid, title string) (string, error) {

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
		//AvatarURL:  avatarURL,
		Title:     title,
		OwnerID:   unionid,
		PersonNum: 1,
	}
	err := insertGroups(group)

	if err != nil {
		return code, err
	}

	go func() {
		AddUserOwnGroup(unionid, group.ID.Hex())
		qrcode, _ := CreateQrcodeByGroupCode(code)
		if qrcode.Ticket != "" {
			qrcode.Ticket = fmt.Sprintf(constant.URLQrcodeTicket, qrcode.Ticket)
			setRedisGroupQrcode(code, qrcode.URL, qrcode.Ticket)
		}
	}()
	return code, err
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
