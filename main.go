package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"project/api/admin"
	v1 "project/api/v1"
	"project/cache"
	"project/conf"
	"project/models"
	"project/mw"
	"time"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "project/docs"
)

func init() {
	fmt.Println("Gin Version: ", gin.Version)
}

// @title GO WEB TEMPLATE API DOC
// @version 1.0
// @description 接口文档

// @host etong.vsattech.com
// @BasePath /backend/api/v1
// @query.collection.format multi
// @schemes https

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationurl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationurl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information

// @x-extension-openapi {"example": "value on a json format"}

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

	//web server
	fmt.Println("init server.")
	server := &http.Server{
		Addr:    conf.Config.Listen.Host + ":" + conf.Config.Listen.Port,
		Handler: newRouter(),
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		fmt.Println("receive interrupt signal")
		if err := server.Close(); err != nil {
			panic(err)
		}
	}()

	fmt.Println("launch server.")
	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			fmt.Println("Server closed under request")
		} else {
			panic(err)
		}
	}

	err = models.CloseDB()
	if err != nil {
		panic(err)
	}

	fmt.Println("quit")

}

func newRouter() *gin.Engine {
	if conf.Config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(ginzerolog.Logger("gin"))
	r.Use(gin.Recovery())
	r.Use(mw.Cors())

	r.GET("ping", func(c *gin.Context) { c.String(http.StatusOK, "pong") })
	// use ginSwagger middleware to serve the API docs
	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1.RouteWarp(r)
	admin.RouteWarp(r)

	return r
}