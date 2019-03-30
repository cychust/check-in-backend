package constant

const (

	/****************************************** wechat ****************************************/

	// https://api.weixin.qq.com/sns/jscode2session?appid=APPID&secret=SECRET&js_code=JSCODE&grant_type=authorization_code
	WechatSessionURIPrefix = "https://api.weixin.qq.com/sns/jscode2session"
	// format https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=APPID&secret=APPSECRET
	WechatTokenURIPrefix = "https://api.weixin.qq.com/cgi-bin/token"
	// https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?access_token=ACCESS_TOKEN
	WechatTemplateSendURIPrefix = "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send"
	WechatDefaultHeadImgURL     = "http://image.upick.hustonline.net/user/wechat-default-headimgurl.jpg"
	WechatOpenCode              = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login#wechat_redirect"
	WechatOpenToken             = "https://api.weixin.qq.com/sns/oauth2/access_token"
	WechatOpenUserInfo          = "https://api.weixin.qq.com/sns/userinfo"

	URLCreateQrcode = "https://wechat.hustonline.net/api/v1/qrcode"
	URLQrcodeTicket = "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=%s"

	URLBingYanIsFollow     = "https://wechat.hustonline.net/api/v1/user/status/follow"
	URLBingYanSendTemplate = "https://wechat.hustonline.net/api/v1/msg/template/list/action/send"
)
