package model

import (
	"check-in-backend/constant"
	"check-in-backend/model/db"
	"check-in-backend/util/log"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"time"
)

var (
	indexLogger            = log.GetLogger()
	DefaultSelector bson.M = bson.M{}
)

func getRedisDefaultExpire() int64 {
	rand.Seed(time.Now().UnixNano())
	return constant.RedisDefaultExpire + rand.Int63n(constant.RedisDefaultRandExpire)
}

func WriteLog(funcName, filename, errMsg string, err error) {
	writeLog(funcName, filename, errMsg, err)
}

func writeLog(filename, funcName, errMsg string, err error) {
	indexLogger.WithFields(logrus.Fields{
		"package":  "controller",
		"file":     filename,
		"function": funcName,
		"err":      err,
	}).Warn(errMsg)
}

//***********************db basic action************************//

func updateDoc(tableName string, query, update interface{}) error {
	cntrl := db.NewCloneMgoDBCntlr()
	defer cntrl.Close()
	table := cntrl.GetTable(tableName)
	return table.Update(query, update)
}

func updateDocs(tableName string, query, update interface{}) (interface{}, error) {
	cntrl := db.NewCloneMgoDBCntlr()
	defer cntrl.Close()
	table := cntrl.GetTable(tableName)
	return table.UpdateAll(query, update)
}

func insertDocs(tableName string, docs ...interface{}) error {
	cntrl := db.NewCloneMgoDBCntlr()
	defer cntrl.Close()
	table := cntrl.GetTable(tableName)
	return table.Insert(docs...)
}

func insertDoc(tableName string, docs ...interface{}) error {
	cntrl := db.NewCloneMgoDBCntlr()
	defer cntrl.Close()
	table := cntrl.GetTable(tableName)
	return table.Insert(docs)
}
