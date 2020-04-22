package models

import (
	"project/conf"
	"project/util"
)

type Admin struct {
	Model
	Username  string `gorm:"type:varchar(128);"`      //用户名
	Email     string `gorm:"type:varchar(128);"`      //邮箱
	Cellphone string `gorm:"type:varchar(128);index"` //手机
	Password  string `gorm:"type:varchar(128)"`       //密码
	Salt      string `gorm:"type:varchar(128)"`       //盐值
	RealName  string `gorm:"type:varchar(128)"`       //真实姓名
	IsSuper   bool   //是不是超级用户
	Roles     string `gorm:"type:text"` //权限数组
}

func (au *Admin) GenToken() (string, error) {
	var data = map[string]interface{}{
		"admin_id":   au.ID,
		"username":   au.Username,
		"email":      au.Email,
		"cellphone":  au.Cellphone,
		"roles":      au.Roles,
		"created_at": au.CreatedAt.Unix(),
	}

	return util.GenerateToken(data, conf.Config.AdminTokenSecret)
}

func (a *Admin) Add() error {
	return db.Model(a).Create(a).Error
}

func (a *Admin) Update() error {
	return db.Model(a).Save(a).Error
}

func AdminGet(id uint64) (*Admin, error) {
	var admin Admin
	err := db.First(&admin, "id = ?", id).Error
	return &admin, err
}

func AdminByUsername(username string) (*Admin, error) {
	var admin Admin
	err := db.First(&admin, "username = ?", username).Error
	return &admin, err
}

func AdminByCellphone(cellphone string) (*Admin, error) {
	var admin Admin
	err := db.First(&admin, "cellphone = ?", cellphone).Error
	return &admin, err
}

func AdminIsExist(id uint64) (bool, error) {
	var count int64
	err := db.Model(&Admin{}).Where("id = ?", id).Count(&count).Error
	if count == 0 {
		return false, nil
	}
	return true, err
}

func AdminExistByCellphone(cellphone string) (bool, error) {
	var count int64
	err := db.Model(&Admin{}).Where("cellphone = ?", cellphone).Count(&count).Error
	if count == 0 {
		return false, nil
	}
	return true, err
}

func AdminCount() (int, error) {
	var count int
	err := db.Model(&Admin{}).Count(&count).Error
	return count, err
}
