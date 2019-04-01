package util

import (
	"check-in-backend/config"
	"errors"
	"math/rand"
	"time"

	timeUtil "github.com/jinzhu/now"

	"github.com/imroc/req"
	"github.com/json-iterator/go"
	"gopkg.in/gomail.v2"
)

var (
	baseStr    string = "0123456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	baseStrLen uint64
	base       []byte = []byte(baseStr)
	baseMap    map[byte]int
	BaseLen    int = 4 // 生成的基础长度
)

func init() {
	InitBaseMap()
	timeUtil.WeekStartDay = time.Monday
}

// JSONStructToMap convert struct to map
func JSONStructToMap(obj interface{}) map[string]interface{} {
	jsonBytes, _ := jsoniter.Marshal(obj)
	var data map[string]interface{}
	jsoniter.Unmarshal(jsonBytes, &data)
	return data
}

func GetFilterMap(obj interface{}, filter []string) map[string]interface{} {
	jsonBytes, _ := jsoniter.Marshal(obj)
	var data map[string]interface{}
	jsoniter.Unmarshal(jsonBytes, &data)

	resData := map[string]interface{}{}
	for _, f := range filter {
		if v, ok := data[f]; ok {
			resData[f] = v
		}
	}
	return resData
}

// BindGetJSONData bind the json data of method GET
// body must be a point
func BindGetJSONData(url string, param req.Param, body interface{}) error {
	r, err := req.Get(url, param)
	if err != nil {
		return err
	}
	err = r.ToJSON(body)
	if err != nil {
		return err
	}
	return nil
}

func SendEmail(name, subject, content string, emailTos []string) {
	m := gomail.NewMessage()
	emailInfo := config.Conf.EmailInfo
	m.SetAddressHeader("From", emailInfo.From, name) // 发件人

	// 收件人
	m.SetHeader("To",
		emailTos...,
	)
	m.SetHeader("Subject", subject) // 主题
	m.SetBody("text/html", content) // 正文

	d := gomail.NewPlainDialer(emailInfo.Host, 465, emailInfo.From, emailInfo.AuthCode) // 发送邮件服务器、端口、发件人账号、发件人密码
	d.DialAndSend(m)
}

func GetNowTimestamp() int64 {
	return time.Now().UnixNano() / 1000000
}

func GetNextAvatarURL() string {
	var urls []string
	urls = make([]string, 10)
	urls[0] = "https://raw.githubusercontent.com/cychust/check-in/master/backgrounds/1.jpeg"
	urls[1] = "https://raw.githubusercontent.com/cychust/check-in/master/backgrounds/1.jpeg"
	urls[2] = "https://raw.githubusercontent.com/cychust/check-in/master/backgrounds/1.jpeg"
	return urls[rand.Int()%2]
}

func GetNextDayEndTimestamp() int64 {
	return timeUtil.EndOfDay().Add(time.Hour * 24).UnixNano() / 1000000
}

func GetNextWeekEndTimestamp() int64 {
	return timeUtil.EndOfWeek().Add(time.Hour * 24 * 7).UnixNano() / 1000000
}

//
//func GetWeekStartTimestamp(t time.Time) int64 {
//	return now.New(t).BeginningOfWeek().UnixNano() / 1000000
//}
//
//func GetWeekEndTimestamp(t time.Time) int64 {
//	return now.New(t).EndOfWeek().UnixNano() / 1000000
//}

func InitBaseMap() {
	baseStrLen = uint64(len(baseStr))
	baseMap = make(map[byte]int)
	for i, v := range base {
		baseMap[v] = i
	}
}

func Base34(n uint64) []byte {
	var mod uint64
	l := []byte{}
	for n != 0 {
		mod = n % baseStrLen
		n = n / baseStrLen
		l = append(l, base[int(mod)])
	}

	lLen := len(l)
	if lLen < BaseLen {
		for i := lLen; i < BaseLen; i++ {
			l = append(l, base[0])
		}
	}
	return l
}

func Base34ToNum(str []byte) (uint64, error) {
	if len(str) == 0 {
		return 0, errors.New("parameter is nil or empty")
	}
	res := uint64(0)
	b := uint64(1)
	for i := 0; i < len(str); i++ {
		v, ok := baseMap[str[i]]
		if !ok {
			return 0, errors.New("character is not base")
		}
		res += b * uint64(v)
		b *= baseStrLen
	}
	return res, nil
}
