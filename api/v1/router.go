package v1

import (
	"fmt"
	"net/http"
	"project/conf"
	"project/util"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func RouteWarp(r *gin.Engine) {
	v1 := r.Group("api/v1")
	{
		//user
		user := v1.Group("/user")
		{
			user.GET("account/check", UserAccountCheck)                      //账号是否存在
			user.GET("info", Authenticator(), UserInfo)                      //用户信息
		}

	}
}

func respwithlog(c *gin.Context, code int, result interface{}) {
	log.Info().Int("code", code).Interface("result", result).Msg("response info")
	c.JSON(code, result)
}

func ContextUserID(c *gin.Context) uint64 {
	value, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return uint64(value.(float64))
}

func ContextWxOpenID(c *gin.Context) string {
	value, exists := c.Get("wxopenid")
	if !exists {
		return ""
	}
	str := fmt.Sprintf("%v", value)
	return str
}

func AuthVerify(verify bool) gin.HandlerFunc {
	return func(c *gin.Context) {
			var token string
			token = c.Request.Header.Get("X-Auth-Token")
			if token == "" {
				token = c.Request.Header.Get("x-auth-token")
			}
			if token == "" {
				cookie, err := c.Request.Cookie("Auth-Token")
				if err != nil {
					log.Warn().Msg(err.Error())
					if !verify {
						c.Next()
						return
					}
					c.SetCookie("Auth-Token", "", 0, "/", "hongkou.vsattech.com", false, false)
					c.AbortWithStatusJSON(http.StatusOK, gin.H{
						"code": CodeFailedAuthVerify,
						"msg":  MsgFailedAuthVerify,
						"data": map[string]interface{}{},
					})
					return
				}

				token = cookie.Value
			}

			jwtT, err := util.ParseToken(token, conf.Config.TokenSecret)
			if err != nil {
				log.Warn().Err(err).Msg("")
				if !verify {
					c.Next()
					return
				}

				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"code": CodeFailedAuthVerify,
					"msg":  MsgFailedAuthVerify,
					"data": map[string]interface{}{},
				})
				return
			}

			claims := jwtT.Claims.(jwt.MapClaims)

			create := claims["created_at"]
			createdAt := int64(create.(float64))
			nowUnix := time.Now().Unix()
			if nowUnix-createdAt >= 10*24*60*60 {
				//过期
				if !verify {
					c.Next()
					return
				}
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"code": CodeTokenExceed,
					"msg":  MsgTokenExceed,
					"data": map[string]interface{}{},
				})
				return
			}

			c.Set("user_id", claims["user_id"])

			//10分钟刷新一下token
			//if nowUnix-createdAt >= 60*10 {
			//	userID, _ := UserID(c)
			//	user, err := models.UserGet(userID)
			//	if err != nil {
			//		if err == gorm.ErrRecordNotFound {
			//			c.AbortWithStatusJSON(http.StatusOK, gin.H{
			//				"code": CodeFailedAuthVerify,
			//				"msg":  MsgFailedAuthVerify,
			//				"data": map[string]interface{}{},
			//			})
			//			return
			//		}
			//		log.Warn().Msg(err.Error())
			//		if !verify {
			//			c.Next()
			//			return
			//		}
			//		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			//			"code": CodeErrSystem,
			//			"msg":  MsgErrSystem,
			//			"data": map[string]interface{}{},
			//		})
			//		return
			//	}
			//
			//	token, err := user.GenToken("")
			//	if err != nil {
			//		log.Warn().Msg(err.Error())
			//		if !verify {
			//			c.Next()
			//			return
			//		}
			//		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			//			"code": CodeErrSystem,
			//			"msg":  MsgErrSystem,
			//			"data": map[string]interface{}{},
			//		})
			//		return
			//	}
			//
			//	c.Header("Token-Refresh", token)
			//}

		c.Next()
	}
}

//Authenticator 身份验证
func Authenticator() gin.HandlerFunc {
	return AuthVerify(true)
}