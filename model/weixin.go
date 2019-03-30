package model

import (
	"check-in-backend/config"
	"check-in-backend/constant"
	"check-in-backend/util"
	"fmt"
	"github.com/imroc/req"
)

type WeixinTokenRes struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}

type WeixinSessRes struct {
	Unionid    string `json:"unionid"`
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
}

type QrcodeParam struct {
	ExpireSeconds int              `json:"expire_seconds"`
	ActionName    string           `json:"action_name"`
	ActionInfo    QrcodeActionInfo `json:"action_info"`
}

type QrcodeActionInfo struct {
	Scene QrcodeScene `json:"scene"`
}

type QrcodeScene struct {
	SceneStr string `json:"scene_str"`
	SceneId  int32  `json:"scene_id"`
}

type QrcodeRes struct {
	ExpireSeconds int    `json:"expire_seconds,omitempty"`
	Ticket        string `json:"ticket"`
	URL           string `json:"url"`
}

type WeixinWebTokenRes struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Unionid      string `json:"unionid"`
	Scope        string `json:"scope"`
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
}

type WeixinWebUserInfo struct {
	Openid     string   `json:"openid" bson:"openid"`
	Unionid    string   `json:"unionid" bson:"unionid"`
	Nickname   string   `json:"nickname" bson:"nickname"`
	Sex        int      `json:"sex" bson:"sex"`
	Province   string   `json:"province" bson:"province"`
	City       string   `json:"city" bson:"city"`
	Country    string   `json:"country" bson:"country"`
	HeadImgURL string   `json:"headimgurl" bson:"headimgurl"`
	Privilege  []string `json:"privilege" bson:"privilege"`
	Errcode    int      `json:"errcode" bson:"errcode"`
	Errmsg     string   `json:"errmsg" bson:"errmsg"`
}

func GetWeixinWebAccessToken(code string) (WeixinWebTokenRes, error) {
	conf := config.Conf
	param := req.Param{
		"appid":      conf.WebWechat.AppID,
		"secret":     conf.WebWechat.AppSecret,
		"code":       code,
		"grant_type": "authorization_code",
	}

	weixinTokenRes := WeixinWebTokenRes{}
	err := util.BindGetJSONData(constant.WechatOpenToken, param, &weixinTokenRes)
	return weixinTokenRes, err
}

func GetWeixinWebUserInfo(accessToken, openid string) (*util.DecryptUserInfo, error) {
	userInfo, err := getWeixinWebUserInfo(accessToken, openid)
	if err != nil {
		return nil, err
	}

	return &util.DecryptUserInfo{
		OpenID:    userInfo.Openid,
		UnionID:   userInfo.Unionid,
		NickName:  userInfo.Nickname,
		Gender:    userInfo.Sex,
		City:      userInfo.City,
		Province:  userInfo.Province,
		Country:   userInfo.Country,
		AvatarURL: userInfo.HeadImgURL,
	}, nil
}

func getWeixinWebUserInfo(accessToken, openid string) (WeixinWebUserInfo, error) {
	param := req.Param{
		"access_token": accessToken,
		"openid":       openid,
	}
	userInfo := WeixinWebUserInfo{}
	err := util.BindGetJSONData(constant.WechatOpenUserInfo, param, &userInfo)
	if err != nil {
		return userInfo, err
	}

	if userInfo.HeadImgURL == "" {
		userInfo.HeadImgURL = constant.ImgDefaultWeixinUserHead
	}
	return userInfo, nil
}

func GetWeixinSession(code string) (WeixinSessRes, error) {
	data := WeixinSessRes{}
	appInfo := config.Conf.Wechat
	param := req.Param{
		"appid":      appInfo.AppID,
		"secret":     appInfo.AppSecret,
		"js_code":    code,
		"grant_type": "authorization_code",
	}
	url := constant.WechatSessionURIPrefix
	err := util.BindGetJSONData(url, param, &data)
	return data, err
}

func DecryptWeixinEncryptedData(sessionKey, encryptedData, iv string) (*util.DecryptUserInfo, error) {
	pc := util.NewWXBizDataCrypt(config.Conf.Wechat.AppID, sessionKey)
	return pc.Decrypt(encryptedData, iv)
}

func getAccessToken() (WeixinTokenRes, error) {
	data := WeixinTokenRes{}
	appInfo := config.Conf.Wechat
	param := req.Param{
		"appid":      appInfo.AppID,
		"secret":     appInfo.AppSecret,
		"grant_type": "client_credential",
	}
	url := constant.WechatTokenURIPrefix
	err := util.BindGetJSONData(url, param, &data)
	return data, err
}

// CreateQrcodeByGroupCode 创建二维码
func CreateQrcodeByGroupCode(code string) (QrcodeRes, error) {
	resData := QrcodeRes{}

	qrcodeMap, err := getRedisGroupQrcode(code)
	if err == nil && qrcodeMap["ticket"] != "" {
		resData.URL = qrcodeMap["url"]
		resData.Ticket = qrcodeMap["ticket"]
		return resData, nil
	}
	str := fmt.Sprintf(constant.WechatScanCodeJoinPhsMPGroup, code)
	reqData := QrcodeParam{
		ExpireSeconds: 3600 * 24 * 30,
		ActionName:    "QR_STR_SCENE",
		ActionInfo: QrcodeActionInfo{
			Scene: QrcodeScene{
				SceneStr: str,
			},
		},
	}

	r, err := req.Post(constant.URLCreateQrcode, req.BodyJSON(&reqData))
	if err != nil {
		return resData, err
	}
	err = r.ToJSON(&resData)
	if resData.Ticket != "" {
		resData.Ticket = fmt.Sprintf(constant.URLQrcodeTicket, resData.Ticket)
		go setRedisGroupQrcode(code, resData.URL, resData.Ticket)
	}
	return resData, err
}

/****************************************** weixin redis action ****************************************/

func UpdateRedisAccessToken() error {
	data, err := getAccessToken()
	if err != nil || data.Errcode != 0 {
		return errors.New(data.Errmsg)
	}
	return updateRedisAccessToken(data.AccessToken, data.ExpiresIn)
}

func updateRedisAccessToken(accessToken string, expire int64) error {
	cntlr := db.NewRedisDBCntlr()
	defer cntlr.Close()

	key := constant.RedisWeixinAccessToken
	_, err := cntlr.SETEX(key, expire, accessToken)
	return err
}

func getRedisAccessToken() (string, error) {
	cntlr := db.NewRedisDBCntlr()
	defer cntlr.Close()

	key := constant.RedisWeixinAccessToken
	return cntlr.GET(key)
}
