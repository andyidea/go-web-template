package models

import (
	"project/proto"

	"github.com/rs/zerolog/log"
)

const (
	OpTypeLogin                  = "登录"
	OpTyoeLogout                 = "登出"
	OpTypeAdminAdd               = "新增管理员"
	OpTypeAdminUpdate            = "更新管理员"
	OpTypeOrganAdd               = "新增机构"
	OpTypeOrganUpdate            = "编辑机构"
	OpTypeOrganDelete            = "删除机构"
	OpTypeCourseAdd              = "新建课程"
	OpTypeCourseUpdate           = "编辑课程"
	OpTypeCourseDelete           = "删除课程"
	OpTypePeriodAdd              = "新增课时"
	OpTypePeriodUpdate           = "编辑课时"
	OpTypePeriodDelete           = "删除课时"
	OpTypePeriodSetting          = "课时时间设置"
	OpTypeCoursewareAdd          = "新增课件"
	OpTypeCoursewareUpdate       = "编辑课件"
	OpTypeCoursewareDelete       = "删除课件"
	OpTypeAuthAdd                = "新增权限"
	OpTypeUserAdd                = "新增用户"
	OpTypeUserEdit               = "编辑用户"
	OpTypeUserUpdate             = "更新用户"
	OpTypeOrgCertification       = "认证企业"
	OpTypeOrgRejectCertification = "取消认证企业"
)

type AdminOplog struct {
	Model
	AdminUserID uint64
	OrganID     uint64
	OpType      string `gorm:"type:varchar(128)"`
	Content     string `gorm:"type:text"`
	Request     string `gorm:"type:text"`
	IP          string `gorm:"type:varchar(128)"`
	Cellphone   string `gorm:"-"`
}

func AdminOplogAdd(ao *AdminOplog) {
	err := db.Create(ao).Error
	if err != nil {
		log.Warn().Err(err).Msg("")
	}
}

func AdminOplogs(pag proto.Pagination, organID uint64) (count int, aos []AdminOplog, err error) {
	query := db.Model(&AdminOplog{})

	if organID > 0 {
		query = query.Where("organ_id = ?", organID)
	}
	err = query.Count(&count).
		Offset(Offset(pag.Page, pag.Limit)).
		Limit(pag.Limit).
		Select("admin_oplogs.*, admins.cellphone").
		Joins("left join admins on admin_oplogs.admin_user_id = admins.id").
		Find(&aos).Error
	return
}
