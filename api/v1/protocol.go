package v1

// Result represents HTTP response body.
type Result struct {
	Code int         `json:"code"` // return code, 0 for succ
	Msg  string      `json:"msg"`  // message
	Data interface{} `json:"data"` // data object
}

const (
	CodeSuccess             = 0
	CodeErrSystem           = 10001 //系统错误
	CodeErrParams           = 10002 //参数错误
	CodeErrLogic            = 10003 //逻辑错误
	CodeFailedAuthVerify    = 10004 //身份验证失败
	CodeTokenExceed         = 10005 //身份过期
	CodeWXNoPhone           = 10006 //微信未绑定手机
	CodeNotPassage          = 10007 //无通行资格
	CodeLackDailyHealthInfo = 10008 //缺少每日健康信息
	CodeIllHealth           = 10009 //体温异常或有异常状况
	CodeNotAllowResumption  = 10010 //不允许复工或未申请复工
	CodeNotFullDay          = 10011 //未满14天
)

const (
	MsgSuccess             = "Success"
	MsgFailed              = "Failed"
	MsgErrSystem           = "系统错误"
	MsgNotPassage          = "无通行资格"
	MsgErrParams           = "参数错误"
	MsgFailedAuthVerify    = "身份验证失败或者已过期，请退出重新登录"
	MsgTokenExceed         = "身份过期，请重新登录"
	MsgWXNoPhone           = "微信未绑定手机"
	MsgShiminNoIDCard      = "请在随申办市民云实名认证后再进入该系统"
	MsgDayConfirm          = "该用户返沪未满14天，是否确认通行" //7天确认
	MsgLackDailyHealthInfo = "检测到您没有填报当日健康信息,请填写后再次申请通行"
	MsgCodeIllHealth       = "检测到您体温异常或有异常状况,暂时无法通行"
	MsgNotAllowResumption  = "检测到您所在的企业未被批准复工"
	MsgNotFullDay          = "检测到您从外地返回未满14天"
)

type SmsCodeSendParam struct {
	PhoneNumber string `form:"phone_number" json:"phone_number" valid:"Required"`
	SmsType     string `form:"sms_type" json:"sms_type" valid:"Required"`
}
