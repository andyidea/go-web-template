package admin

import (
	"net/http"
	"project/models"
	"project/util"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
)

type RegParam struct {
	Username        string `form:"username" json:"username" valid:"Required;Match(/^[a-zA-Z0-9_-]{4,16}$/)"`
	Password        string `form:"password" json:"password" valid:"Required"`
	PasswordConfirm string `form:"password_confirm" json:"password_confirm" valid:"Required"`
}

type RegData struct {
	Token string `json:"token"`
}

func Reg(c *gin.Context) {
	result := &Result{}
	defer c.JSON(http.StatusOK, result)

	var param RegParam
	err := util.GinBind(c, &param)
	if err != nil {
		log.Warn().Err(err).Msg("")
		result.Code = CodeErrParams
		result.Msg = err.Error()
		return
	}

	var exist = true
	_, err = models.AdminByUsername(param.Username)
	if gorm.IsRecordNotFoundError(err) {
		exist = false
	} else if err != nil {
		log.Error().Err(err).Msg("")
		result.Code = CodeErrSystem
		result.Msg = MsgErrSystem
		return
	}

	if exist {
		result.Code = CodeErrLogic
		result.Msg = "账号已存在"
		return
	}

	if param.Password != param.PasswordConfirm {
		result.Code = CodeErrLogic
		result.Msg = "两次输入的密码不一致"
		return
	}

	var admin models.Admin
	admin.Username = param.Username
	admin.Salt = strconv.FormatInt(time.Now().Unix(), 10)
	admin.Password = util.EncryptPassword(param.Password, admin.Salt)

	err = admin.Add()
	if err != nil {
		log.Error().Err(err).Msg("")
		result.Code = CodeErrSystem
		result.Msg = MsgErrSystem
		return
	}

	token, err := admin.GenToken()
	if err != nil {
		log.Error().Err(err).Msg("")
		result.Code = CodeErrSystem
		result.Msg = MsgErrSystem
		return
	}
	c.SetCookie("Admin-Auth-Token", token, 9999999, "/", c.Request.Host, false, false)

	var data RegData
	data.Token = token
	result.Data = data
	result.Code = CodeSuccess
	result.Msg = MsgSuccess

}

func Login(c *gin.Context) {
	result := &Result{}
	defer c.JSON(http.StatusOK, result)

	var param AdminUserLoginParam
	err := util.GinBindJSON(c, &param)
	if err != nil {
		log.Warn().Err(err).Msg("")
		result.Code = CodeErrParams
		result.Msg = err.Error()
		return
	}

	username := param.Username

	admin, err := models.AdminByUsername(username)
	if err == gorm.ErrRecordNotFound {
		result.Code = CodeErrLogic
		result.Msg = "邮箱不存在或者密码错误"
		return
	}
	if err != nil {
		log.Error().Err(err).Msg("")
		result.Code = CodeErrSystem
		result.Msg = MsgErrSystem
		return
	}

	var pw = util.EncryptPassword(param.Password, admin.Salt)
	if admin.Password != pw {
		result.Code = CodeErrLogic
		result.Msg = "用户名不存在或者密码错误"
		return
	}

	//只有超级用户可以登录
	if !admin.IsSuper {
		result.Code = CodeErrLogic
		result.Msg = "无权限"
		return
	}

	var token string
	token, err = admin.GenToken()
	if err != nil {
		log.Error().Err(err).Msg("")
		result.Code = CodeErrSystem
		result.Msg = MsgErrSystem
		return
	}
	var data AdminUserLoginData
	data.Token = token
	result.Data = data
	result.Code = CodeSuccess
	result.Msg = MsgSuccess

	models.AdminOplogAdd(&models.AdminOplog{
		AdminUserID: admin.ID,
		OpType:      models.OpTypeLogin,
		Content:     "登录了系统",
		IP:          c.ClientIP(),
	})

}

func Logout(c *gin.Context) {
	result := &Result{}
	defer c.JSON(http.StatusOK, result)

	adminID := ContextAdminID(c)

	models.AdminOplogAdd(&models.AdminOplog{
		AdminUserID: adminID,
		OpType:      models.OpTyoeLogout,
		Content:     "登出了系统",
		IP:          c.ClientIP(),
	})

	result.Code = CodeSuccess
	result.Msg = MsgSuccess
}

type InfoGetData struct {
	RealName string `json:"real_name"`
}

func InfoGet(c *gin.Context) {
	result := &Result{}
	defer c.JSON(http.StatusOK, result)

	adminID := ContextAdminID(c)

	models.AdminGet(adminID)
	admin, err := models.AdminGet(adminID)
	if err != nil {
		log.Error().Err(err).Msg("")
		result.Code = CodeErrSystem
		result.Msg = MsgErrSystem
		return
	}

	var data InfoGetData

	data.RealName = admin.RealName

	result.Data = data
	result.Code = CodeSuccess
	result.Msg = MsgSuccess
}

type InfoUpdateParam struct {
	RealName string `form:"real_name" json:"real_name"`
}

func InfoUpdate(c *gin.Context) {
	result := &Result{}
	defer c.JSON(http.StatusOK, result)

	var param InfoUpdateParam
	err := util.GinBind(c, &param)
	if err != nil {
		log.Warn().Err(err).Msg("")
		result.Code = CodeErrParams
		result.Msg = err.Error()
		return
	}

	adminID := ContextAdminID(c)

	admin, err := models.AdminGet(adminID)
	if err != nil {
		log.Error().Err(err).Msg("")
		result.Code = CodeErrSystem
		result.Msg = MsgErrSystem
		return
	}

	if len(param.RealName) > 0 {
		admin.RealName = param.RealName
	}

	err = admin.Update()
	if err != nil {
		log.Error().Err(err).Msg("")
		result.Code = CodeErrSystem
		result.Msg = MsgErrSystem
		return
	}

	result.Code = CodeSuccess
	result.Msg = MsgSuccess

	models.AdminOplogAdd(&models.AdminOplog{
		AdminUserID: admin.ID,
		OpType:      "更新信息",
		Content:     "更新信息",
		Request:     util.JsonMarshal(param),
		IP:          c.ClientIP(),
	})
}
