package server

import (
	"fmt"
	"net/http"
	"github.com/liangdas/mqant/module/base"
	"github.com/liangdas/mqant/module"
	"github.com/liangdas/mqant/conf"
	"server/controllers"
	"github.com/liangdas/mqant/log"
	"tools"
)

var Module = func() *Server {
	server := new(Server)
	return server
}

type Server struct {
	basemodule.BaseModule
	StaticPath	string
	Port		int
	RedisUrl	string
}

type RedisConf struct {
	Address  string		`json:"Address"`
	Password string		`json:"Password"`
	Pool     PoolConf	`json:"Pool"`
}

type PoolConf struct {
	MaxIdle     int		`json:"MaxIdle"`
	MaxActive   int		`json:"MaxActive"`
	IdleTimeout int		`json:"IdleTimeout"`
}

// GetType 获取模块类型标识
func (self *Server) GetType() string {
	//很关键,需要与配置文件中的Module配置对应
	return "WXAccessToken"
}

//Version 获取Web模块版本号
func (self *Server) Version() string {
	//可以在监控时了解代码版本
	return "1.0.0"
}

//OnInit Web模块初始化方法
func (self *Server) OnInit(app module.App, settings *conf.ModuleSettings) {
	self.BaseModule.OnInit(self, app, settings)
	self.StaticPath = self.App.GetSettings().Settings["StaticPath"].(string)
	self.Port = int(settings.Settings["Port"].(float64))
	self.RedisUrl = app.GetSettings().Settings["RedisUrl"].(string)
	tools.RedisUrl = self.RedisUrl
}

//Run Web模块启动方法
func (self *Server) Run(closeSig chan bool) {
	//这里如果出现异常请检查8080端口是否已经被占用
	e := server.SetupRouter()
	e.Static("/static", self.StaticPath)
	go func() {
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", self.Port)))
		log.Info("", http.ListenAndServe(fmt.Sprintf(":%v", self.Port), nil))
	}()

	<-closeSig
	e.Close()
	log.Info("webapp server Shutting down...")
}

//OnDestroy Web模块注销方法
func (self *Server) OnDestroy() {
	//一定别忘了关闭RPC
	self.GetServer().OnDestroy()
}
