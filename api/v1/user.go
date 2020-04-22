package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"project/models"
)

type UserAccountCheckData struct {
	Exist bool `json:"exist"`
}

// UserAccountCheck godoc
// @Summary 账号检测
// @Description 检查账号是否存在
// @Tags 用户
// @Accept  json
// @Produce  json
// @Success 200 {object} UserAccountCheckData "desc"
// @Router /user/account/check [get]
func UserAccountCheck(c *gin.Context) {
	result := &Result{}
	defer c.JSON(http.StatusOK, result)

	var wxOpenID string
	wxOpenID = c.Request.Header.Get("X-WXOpenID")
	if wxOpenID == "" {
		cookie, err := c.Request.Cookie("wid")
		if err == nil {
			wxOpenID = cookie.Value
		}
	}

	exist, err := models.UserExistByWXOpenID(wxOpenID)
	if err != nil {
		log.Error().Err(err).Msg("")
		result.Code = CodeErrSystem
		result.Msg = MsgErrSystem
		return
	}

	var data UserAccountCheckData
	data.Exist = exist

	result.Msg = MsgSuccess
	result.Code = CodeSuccess
	result.Data = data
}

type UserLoginParam struct {
	LoginType   string `form:"login_type" json:"login_type"`
	Username    string `form:"username" json:"username"`
	Password    string `form:"password" json:"password"`
	PhoneNumber string `form:"phone_number" json:"phone_number"`
	Code        string `form:"code" json:"code"`
}

type UserLoginData struct {
	Token string `json:"token"`
}
type UserInfoData struct {
	ID              uint64 `json:"id"`
	RealName        string `json:"real_name"`
	IDCard          string `json:"id_card"`
	Passport        string `json:"passport"`
	Province        string `json:"province"`
	City            string `json:"city"`
	District        string `json:"district"`
	Addr            string `json:"addr"`
	CompanyName     string `json:"company_name"`
	CompanyProvince string `json:"company_province"`
	CompanyCity     string `json:"company_city"`
	CompanyDistrict string `json:"company_district"`
	CompanyAddr     string `json:"company_addr"`
	CompanyPark     string `json:"company_park"`
	IsFirstSubmit   bool   `json:"is_first_submit"`
}

func UserInfo(c *gin.Context) {
	result := &Result{}
	defer c.JSON(http.StatusOK, result)

	userID := ContextUserID(c)

	user, err := models.UserGet(userID)
	if err != nil {
		log.Error().Err(err).Msg("")
		result.Code = CodeErrSystem
		result.Msg = MsgErrSystem
		return
	}

	var data UserInfoData
	data.ID = user.ID
	data.RealName = user.RealName
	data.IDCard = user.IDCard
	data.Passport = user.Passport
	data.Province = user.Province
	data.City = user.City
	data.District = user.District
	data.Addr = user.Addr
	data.CompanyName = user.CompanyName
	data.CompanyProvince = user.CompanyProvince
	data.CompanyCity = user.CompanyCity
	data.CompanyDistrict = user.CompanyDistrict
	data.CompanyAddr = user.CompanyAddr
	data.CompanyPark = user.CompanyPark
	data.IsFirstSubmit = user.IsFirstSubmit

	result.Data = data
	result.Code = CodeSuccess
	result.Msg = MsgSuccess
}