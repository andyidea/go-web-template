package admin

import "project/proto"

// Result represents HTTP response body.
type Result struct {
	Code   int         `json:"code"` // return code, 0 for succ
	Msg    string      `json:"msg"`  // message
	Data   interface{} `json:"data"` // data object
	Detail string      `json:"detail"`
}

type PagaParam struct {
	proto.Pagination
}

const (
	CodeSuccess          = 0
	CodeErrSystem        = 10001 //系统错误
	CodeErrParams        = 10002 //参数错误
	CodeErrLogic         = 10003 //逻辑错误
	CodeFailedAuthVerify = 10004 //身份验证失败
	CodeRecordExists     = 10005 //记录已存在
	CodeNoPerm           = 10006 //没有权限
)

const (
	MsgSuccess          = "Success"
	MsgErrSystem        = "系统错误"
	MsgErrParams        = "参数错误"
	MsgFailedAuthVerify = "身份验证失败"
	MsgRecordExists     = "记录已存在"
	MsgUserNoExists     = "用户不存在"
	MsgNoPerm           = "没有权限"
)

type AdminUserLoginParam struct {
	Username string `form:"username" binding:"required" json:"username"`
	Password string `form:"password" binding:"required" json:"password"`
}

type AdminUserOplogListItem struct {
	ID        uint64 `json:"id"`
	Email     string `json:"email"`
	OpType    string `json:"op_type"`
	Content   string `json:"content"`
	Request   string `json:"request"`
	IP        string `json:"ip"`
	CreatedAt int64  `json:"created_at"`
}

type AdminUserOplogListData struct {
	Count int                      `json:"count"`
	Items []AdminUserOplogListItem `json:"items"`
}

type AdminUserLoginData struct {
	Token string `json:"token"`
}

type AdminUserInfoData struct {
	ID     uint64   `json:"id"`
	Email  string   `json:"email"`
	Avatar string   `json:"avatar"`
	Roles  []string `json:"roles"`
}

type AdminUserOplogListParam struct {
	proto.Pagination
}

type AdminUserListParam struct {
	proto.Pagination
}

type AdminUserSearchListParam struct {
	proto.Pagination
	RealName  string `form:"real_name"`
	Cellphone string `form:"cellphone"`
	Email     string `form:"email"`
}

type AdminUserListItem struct {
	ID        uint64 `json:"id"`
	Email     string `json:"email"`
	RealName  string `json:"real_name"`
	Cellphone string `json:"cellphone"`
	Roles     string `json:"roles"`
	CreatedAt int64  `json:"created_at"`
}

type AdminUserListData struct {
	Count int                 `json:"count"`
	Items []AdminUserListItem `json:"items"`
}

type OrganBindParam struct {
	OrganID     uint64 `form:"organ_id"`
	AdminUserID uint64 `form:"admin_user_id"`
	IsClearBind bool   `form:"is_clear_bind"`
}
