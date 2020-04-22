package mw

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//Cors 跨域
func Cors() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowOriginFunc = func(origin string) bool {
		return true
		//matched, err := regexp.MatchString(`((http(s?)://)?([:alnum:]+\.)?)(smm.cn|smmadmin.cn|anhuida.com|metal.com)`, origin)
		//if err != nil {
		//	log.Println("正则匹配域名错误:", err.Error())
		//	return false
		//}
		//
		//if matched == true {
		//	return true
		//}
		//return false
	}
	config.AllowCredentials = true
	config.AllowOrigins = []string{"http://localhost"}
	config.AllowHeaders = append(config.AllowHeaders, "X-Auth-Token", "X-Admin-Auth-Token")
	config.AllowMethods = append(config.AllowMethods, "OPTIONS", "GET", "POST")
	config.ExposeHeaders = append(config.ExposeHeaders, "Token-Refresh")

	return cors.New(config)
}
