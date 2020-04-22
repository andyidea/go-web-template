package main

import (
	"fmt"
	"os"
	"project/cache"
	"project/conf"
	"project/models"
	"project/util"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	//config init
	fmt.Println("init config.")
	err := conf.Init()
	if err != nil {
		panic(err)
	}

	//logger init
	fmt.Println("init logger.")
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	if conf.Config.Debug {
		log.Logger = log.Logger.With().Caller().Logger()
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	//cache init
	if conf.Config.Redis.Enable {
		cache.Init()
	}

	//db init
	fmt.Println("init models.")
	err = models.Init()
	if err != nil {
		panic(err)
	}

	var data = map[string]interface{}{
		"admin_id":   1,
		"created_at": time.Now().Unix(),
	}

	token, _ := util.GenerateToken(data, "test_admin_token_secret")
	println(token)
}
