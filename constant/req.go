package constant

const (
	ReqGroupUpdateOwnerType = iota + 1
	ReqGroupDelOwnerType
	ReqGroupSetManagerType
	ReqGroupUnSetManagerType
	ReqGroupDelManagerType
	ReqGroupDelMemberType
)

const (
	ReqNoticeGetAllType = iota + 1
	ReqNoticeGetByGroupIDType
	ReqNoticeGetFileType
)

const (
	ReqNoticeUpdate = iota
	ReqNoticeUpdateContentType
	ReqNoticeUpdateImgsType
	ReqNoticeUpdateTitleType
	ReqNoticeUpdateNoteType
	ReqNoticeUpdateNoticeTimeType
	ReqNoticeUpdateGroupIDType
)

const (
	ReqNoticeAddFileType = iota + 1
	ReqNoticeDelFileType
	ReqNoticeUpdateFileType
)

const (
	ReqNoticeGetGroupMembersType = iota
	ReqNoticeGetGroupInfoType
)

const (
	ReqNoticeGetNoticeType    = iota + 1
	ReqNoticeGetNoticeManType // 学委版
)
