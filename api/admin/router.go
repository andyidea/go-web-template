package admin

import (
	"net/http"
	"project/conf"
	"project/models"
	"project/proto"
	"project/util"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RouteWarp(r *gin.Engine) {
	admin := r.Group("api/admin")
	{
		//身份校验
		admin.Use(AdminAuthenticator())
		//adminuser
		admin.POST("login", Login)   //登录
		admin.POST("reg", Reg)       //注册
		admin.POST("logout", Logout) //登出
		admin.GET("info/get", InfoGet)
		admin.POST("info/update", InfoUpdate)

		//oplog
		admin.GET("oplog/list", OplogList)
	}
}

func ContextAdminID(c *gin.Context) uint64 {
	value, exists := c.Get("admin_id")
	if !exists {
		return 0
	}
	return uint64(value.(float64))
}

func ContextOrganID(c *gin.Context) uint64 {
	value, exists := c.Get("organ_id")
	if !exists {
		return 0
	}
	v, _ := strconv.ParseUint(value.(string), 10, 64)
	return v
}

func ContextRoles(c *gin.Context) string {
	value, exists := c.Get("roles")
	if !exists {
		return ""
	}
	return value.(string)
}

//AuthVerify 身份验证
func AdminAuthenticator() gin.HandlerFunc {
	return func(c *gin.Context) {
		//跳过验证 白名单
		if strings.Contains(c.Request.URL.Path, "/api/admin/login") ||
			strings.Contains(c.Request.URL.Path, "/api/admin/reg") ||
			strings.Contains(c.Request.URL.Path, "/sms/code/send") {
			c.Next()
			return
		}
		var token string
		token = c.Request.Header.Get("X-Admin-Auth-Token")
		if token == "" {
			cookie, err := c.Request.Cookie("Admin-Auth-Token")
			if err != nil {
				log.Warn().Msg(err.Error())
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"code": CodeFailedAuthVerify,
					"msg":  MsgFailedAuthVerify,
					"data": map[string]interface{}{},
				})
				return
			}

			token = cookie.Value
		}

		jwtT, err := util.ParseToken(token, conf.Config.AdminTokenSecret)
		if err != nil {
			log.Warn().Err(err).Msg("")
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": CodeFailedAuthVerify,
				"msg":  MsgFailedAuthVerify,
				"data": map[string]interface{}{},
			})
			return
		}

		claims := jwtT.Claims.(jwt.MapClaims)
		c.Set("admin_id", claims["admin_id"])
		c.Set("roles", claims["roles"])
		cookie, err := c.Request.Cookie("oid")
		if err == nil {
			c.Set("organ_id", cookie.Value)
		}

		c.Next()
	}
}

type OplogListParam struct {
	proto.Pagination
}

type OplogListItem struct {
	ID        uint64 `json:"id"`
	ShowName  string `json:"show_name"`
	OpType    string `json:"op_type"`
	Content   string `json:"content"`
	Request   string `json:"request"`
	IP        string `json:"ip"`
	CreatedAt int64  `json:"created_at"`
}

type OplogListData struct {
	Count int             `json:"count"`
	Items []OplogListItem `json:"items"`
}

func OplogList(c *gin.Context) {
	result := &Result{}
	defer c.JSON(http.StatusOK, &result)

	var param OplogListParam
	err := util.GinBind(c, &param)
	if err != nil {
		log.Warn().Err(err).Msg("")
		result.Code = CodeErrParams
		result.Msg = MsgErrParams
		return
	}

	organID := ContextOrganID(c)

	count, ops, err := models.AdminOplogs(param.Pagination, organID)
	if err != nil {
		log.Warn().Err(err).Msg("")
		result.Code = CodeErrSystem
		result.Msg = MsgErrSystem
		return
	}

	var data OplogListData
	data.Items = make([]OplogListItem, 0)

	for _, ao := range ops {
		var item OplogListItem
		item.ID = ao.ID
		item.ShowName = ao.Cellphone
		item.OpType = ao.OpType
		item.Content = ao.Content
		item.Request = ao.Request
		item.IP = ao.IP
		item.CreatedAt = ao.CreatedAt.Unix()

		data.Items = append(data.Items, item)
	}

	data.Count = count

	result.Data = data
	result.Code = CodeSuccess
	result.Msg = MsgSuccess
}

//func permCheck(perms []string) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		rolestr, _ := c.Get("roles")
//
//		roles := strings.Split(rolestr.(string), ",")
//
//		var pass = false
//		for _, role := range roles {
//			if role == "admin" {
//				pass = true
//				break
//			}
//
//			if util.StringsIn(perms, role) {
//				pass = true
//				break
//			}
//
//		}
//
//		if !pass {
//			c.AbortWithStatusJSON(http.StatusOK, gin.H{
//				"code": CodeNoPerm,
//				"msg":  MsgNoPerm,
//				"data": map[string]interface{}{},
//			})
//			return
//		}
//
//		c.Next()
//	}
//}
