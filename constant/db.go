package constant

const (

	/****************************************** table name ****************************************/

	TableUser     = "user"
	TableGroup    = "group"
	TableNotice   = "notice"
	TableTemplate = "template"

	/****************************************** user ****************************************/

	UserDeleteStatus   = -10
	UserUnFollowStatus = 0
	UserFollowStatus   = 5

	/****************************************** notice ****************************************/

	NoticeDeleteStatus = -10
	NoticeExpireStatus = -1
	NoticePubStatus    = 5

	/****************************************** group ****************************************/

	GroupDelStatus    = -10
	GroupCommonStatus = 5

	/****************************************** feedback ****************************************/

	FeedbackUnReadStatus = 0
	FeedbackReadedStatus = 1

	/****************************************** redis ****************************************/

	// redis 字段
	RedisDefaultExpire     = 3600 * 24 * 7 // 7天
	RedisDefaultRandExpire = 3600 * 24     // 1天

	RedisNoticeFileUnHandNames = "notice:file:name:%s" // notice:file:name:<notice id>

	RedisUserInfo = "user:info:%s" // format: user:info:<unionid>
	//  format: user:notice:week:<week start timestamp>:<unionid>
	// list value: notice id
	RedisUserWeekNotice = "user:notice:week:%d:%s"
	RedisUserWeek       = "user:notice:week:%d:*"

	RedisGroupInfo        = "group:info:%s"   // format: group:info:<_id>
	RedisGroupQrcode      = "group:qrcode:%s" // format: group:qrcode:<code>
	RedisGroupCodePool    = "group:code:pool" // 列表存储圈子code, 每次存储 100 个
	RedisGroupCodePoolNum = 100
	RedisGroupCodeNextNum = "group:code:next_num" // 记录圈子code转换数字，初始 1000
	RedisGroupInitNextNum = 1000

	RedisWeixinAccessToken = "weixin:access_token"

	RedisNoticeZipFileIsUpdate = "notice:file:zip:%s:is_update" // 1 更新过, 0 没有
	RedisNoticeZipFileURL      = "notice:file:zip:%s"           // notice:file:zip:<notice id>， 过期时间30天
	RedisNoticeZipFileExpire   = 3600 * 24 * 30
)
