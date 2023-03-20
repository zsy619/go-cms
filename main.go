package main

import (
	"fmt"
	"log"
	"os"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/session"
	"github.com/kardianos/service"
)

var globalSessions *session.Manager

// Session初始化
func InitSession() {
	sessionConfig := &session.ManagerConfig{
		CookieName:      "gosessionid",
		EnableSetCookie: true,
		Gclifetime:      7200,
		Maxlifetime:     7200,
		Secure:          false,
		CookieLifeTime:  7200,
		ProviderConfig:  "./session",
	}

	globalSessions, _ = session.NewManager("memory", sessionConfig)
	go globalSessions.GC()
}

func init() {
	err := os.MkdirAll("./logs", os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = logs.SetLogger(logs.AdapterFile, `{"filename":"logs/project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
	if err != nil {
		panic(err)
	}
	logs.Info("hello cms!")
	web.BConfig.WebConfig.Session.SessionOn = true
	// InitSession()
}

var serviceConfig = &service.Config{
	Name:        "cms",
	DisplayName: "内容管理",
	Description: "内容管理",
}

func main() {
	// 构建服务对象
	prog := &Program{}
	s, err := service.New(prog, serviceConfig)
	if err != nil {
		log.Fatal(err)
	}

	// 用于记录系统日志
	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		err = s.Run()
		if err != nil {
			logger.Error(err)
		}
		return
	}

	cmd := os.Args[1]

	if cmd == "install" {
		err = s.Install()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("安装成功")
	}
	if cmd == "uninstall" {
		err = s.Uninstall()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("卸载成功")
	}
}

type Program struct{}

func (p *Program) Start(s service.Service) error {
	log.Println("开始服务")
	go p.run()
	return nil
}

func (p *Program) Stop(s service.Service) error {
	log.Println("停止服务")
	return nil
}

func (p *Program) run() { // 此处编写具体的服务代码
	// beego过滤器 https://beego.me/docs/mvc/controller/filter.md
	// if global.AllowLogin {
	// 	// sso过滤器
	// 	web.InsertFilter("/sso/*", web.BeforeRouter, routers.FilterSSO)
	// 	web.InsertFilter("/mkt/*", web.BeforeRouter, routers.FilterSSO)
	// 	web.InsertFilter("/cms/*", web.BeforeRouter, routers.FilterSSO)
	// 	// mobile过滤器
	// 	web.InsertFilter("/mobile/x/*", web.BeforeRouter, routers.FilterMobile)
	// }

	web.BConfig.CopyRequestBody = true

	web.Run()
}

func RefreshToken() {
	// oauth := mp.NewOAuth2(&models.MpConfig{
	// 	AppId:          global.WxMpConfig.AppId,
	// 	AppSecret:      global.WxMpConfig.AppSecret,
	// 	EncodingAesKey: global.WxMpConfig.EncodingAesKey,
	// 	Token:          global.WxMpConfig.Token,
	// })
	// oauth.RefreshToken()
}
