package models

import (
	"github.com/jinzhu/gorm"
	"project/conf"
	"project/proto"
	"project/util"
	"time"
)

const (
	SexMale    = "male"    //男性
	SexFemale  = "female"  //女性
	SexSecrecy = "secrecy" //保密
)

//User 用户
type User struct {
	Model
	WXOpenID  string     `gorm:"type:varchar(128);index"`
	Email     string     `gorm:"type:varchar(128);"`      //邮箱
	Cellphone string     `gorm:"type:varchar(128);index"` //手机
	Username  string     `gorm:"type:varchar(128);index"` //用户名
	Password  string     `gorm:"type:varchar(128)"`       //密码
	Source    string     `gorm:"type:varchar(128)"`       //来源
	Salt      string     `gorm:"type:varchar(128)"`       //盐值
	RealName  string     `gorm:"type:varchar(128)"`       //真实姓名
	Profile   string     `gorm:"type:text"`               //简介
	Birthday  *time.Time `gorm:"type:date"`               //生日
	Country   string     `gorm:"type:varchar(255)"`       //国家

	//居住信息
	Province string `gorm:"type:varchar(255)"` //省
	City     string `gorm:"type:varchar(255)"` //市
	District string `gorm:"type:varchar(255)"` //区县
	Addr     string `gorm:"type:text"`         //联系地址
	IDCard   string `gorm:"type:varchar(255)"` //身份证
	Passport string `gorm:"type:varchar(255)"` //护照

	//公司信息
	CompanyName     string `gorm:"type:varchar(255)"` //公司名称
	CompanyProvince string `gorm:"type:varchar(255)"` //省
	CompanyCity     string `gorm:"type:varchar(255)"` //市
	CompanyDistrict string `gorm:"type:varchar(255)"` //区县
	CompanyAddr     string `gorm:"text"`              //公司详细地址
	CompanyPark     string `gorm:"type:varchar(255)"` //公司园区/楼宇

	IsFirstSubmit bool //是否首次提交
}

func (u *User) GenToken() (string, error) {
	var data = map[string]interface{}{
		"user_id":    u.ID,
		"username":   u.Username,
		"created_at": time.Now().Unix(),
	}

	return util.GenerateToken(data, conf.Config.TokenSecret)
}

func UserGet(userID uint64) (*User, error) {
	var user User
	err := db.First(&user, userID).Error
	return &user, err
}

func UserByUsername(username string) (*User, error) {
	var user User
	err := db.First(&user, "username = ?", username).Error
	return &user, err
}

func UserByCellphone(cellphone string) (*User, error) {
	var user User
	err := db.First(&user, "cellphone = ?", cellphone).Error
	return &user, err
}

func UserByCellphoneNull(cellphone string) (*User, error) {
	var user User
	var count int
	err := db.Model(&User{}).Where("cellphone = ?", cellphone).Count(&count).First(&user).Error
	if err != nil {
		if count == 0 {
			return &user, nil
		}
	}
	return &user, err
}

func UserByWXOpenID(wxOpenID string) (*User, error) {
	var user User
	err := db.First(&user, "wx_open_id = ?", wxOpenID).Error
	return &user, err
}

func UserByShiminUserID(shiminUserID uint64) (*User, error) {
	var user User
	err := db.First(&user, "shimin_user_id = ?", shiminUserID).Error
	return &user, err
}

func UserAdd(user *User) error {
	return db.Model(user).Create(user).Error
}

func UserUpdate(user *User) error {
	return db.Model(user).Save(user).Error
}

func UserExistByUsername(username string) (bool, error) {
	var count int
	err := db.Model(&User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count >= 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func UserExistByCellphone(cellphone string) (bool, error) {
	var count int
	err := db.Model(&User{}).Where("cellphone = ?", cellphone).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count >= 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func UserExistByUserID(userID uint64) (bool, error) {
	var count int
	err := db.Model(&User{}).Where("id = ?", userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count >= 1 {
		return true, nil
	} else {
		return false, nil
	}
}

func UserExistByWXOpenID(WXOpenID string) (bool, error) {
	var count int
	err := db.Model(&User{}).Where("wx_open_id = ?", WXOpenID).Count(&count).Error
	if err != nil {
		return false, err
	}
	if count >= 1 {
		return true, nil
	} else {
		return false, nil
	}
}

type UserListArgs struct {
	proto.Pagination
	RealName  string
	Cellphone string
}

type UserListArgsBySuper struct {
	proto.Pagination
	RealName  string
	Cellphone string
}

func UserList(args *UserListArgs) (result []User, count int, err error) {
	err = db.Model(&User{}).Count(&count).Order("id").Offset(Offset(args.Page, args.Limit)).Limit(args.Limit).Find(&result).Error
	if IsNoRowsInResultError(err) {
		err = nil
		return
	}
	return
}

func UserListBySuper(args *UserListArgsBySuper, result interface{}) (count int, err error) {
	query := db.Table("users").
		Order("id DESC")
	if len(args.Cellphone) > 0 {
		query = query.Where("cellphone like ?", "%"+args.Cellphone+"%")
	}
	if len(args.RealName) > 0 {
		query = query.Where("real_name like ?", "%"+args.RealName+"%")
	}
	query = query.Where("deleted_at is null")
	query = query.Select("CAST(extract(epoch FROM users.created_at)as int) as created_at,users.*").Count(&count).Offset(Offset(args.Page, args.Limit)).Limit(args.Limit).Scan(result)
	err = query.Error
	if gorm.IsRecordNotFoundError(err) {
		err = nil
		return
	}
	return
}

