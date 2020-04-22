package models

import (
	"fmt"
	"project/conf"
	"project/util"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	db *gorm.DB
)

func IsNoRowsInResultError(err error) bool {
	if err != nil && strings.Contains(err.Error(), "no rows in result set") {
		return true
	}
	return false
}

type Model struct {
	ID        uint64     `gorm:"primary_key;AUTO_INCREMENT" json:"id"` //用户ID
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func Init() error {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		conf.Config.Database.Host,
		conf.Config.Database.Port,
		conf.Config.Database.User,
		conf.Config.Database.Dbname,
		conf.Config.Database.Password,
		conf.Config.Database.Sslmode)
	db, err = gorm.Open("postgres", connStr)
	if err != nil {
		return err
	}
	if conf.Config.Debug {
		db.LogMode(true)
	}

	if err = db.AutoMigrate(&User{},  &Admin{},
		&AdminOplog{}).Error; nil != err {
		return err
	}

	db.DB().SetMaxIdleConns(conf.Config.Database.MaxIdleConns)
	db.DB().SetMaxOpenConns(conf.Config.Database.MaxOpenConns)

	createInitAdmin()

	return nil
}

func CloseDB() error {
	return db.Close()
}

func GetDB() *gorm.DB {
	return db
}

func createInitAdmin() {
	count, err := AdminCount()
	if err != nil {
		log.Warn().Err(err).Msg("")
		return
	}

	if count <= 0 {
		admin := Admin{}
		admin.Username = "admin"
		admin.Salt = "123456"
		admin.Password = util.EncryptPassword(util.MD5("admin2020"), admin.Salt)
		admin.IsSuper = true
		admin.Roles = "admin"
		_ = admin.Add()
	}

}

//Offset 获取数据库查询的offset
func Offset(page, limit int) int {
	if page <= 0 {
		return -1 //cancel offset
	}

	return (page - 1) * limit
}

func Limit(limit int) int {
	if limit <= 0 {
		return -1
	}

	return limit
}
