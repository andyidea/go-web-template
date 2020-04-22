package conf

import (
	"github.com/koding/multiconfig"
)

type RedisConf struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	MaxConn  int    `json:"maxconn"`
	Prefix   string `json:"prefix"`
	Password string `json:"password"`
	DB       int    `json:"db"`
	Enable   bool   `json:"enable"`
}

var Config = struct {
	Debug  bool `json:"debug"`
	Listen struct {
		Host string `json:"host" required:"true"`
		Port string `json:"port" required:"true"`
	} `json:"listen"`
	Database struct {
		Host         string `json:"host" required:"true"`
		Port         string `json:"port" required:"true"`
		User         string `json:"user" required:"true"`
		Dbname       string `json:"dbname" required:"true"`
		Password     string `json:"password" required:"true"`
		Sslmode      string `json:"sslmode" required:"true"`
		MaxIdleConns int    `json:"max_idle_conns" required:"true"`
		MaxOpenConns int    `json:"max_open_conns" required:"true"`
	} `json:"database"`
	Redis RedisConf `json:"redis" required:"true"`
	TokenSecret      string `json:"token_secret" required:"true"`
	AdminTokenSecret string `json:"admin_token_secret" required:"true"`
}{}

func Init() error {
	m := multiconfig.NewWithPath("config.json")
	err := m.Load(&Config)
	return err
}
