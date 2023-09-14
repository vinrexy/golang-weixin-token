package server

import (
	"echo/filters"
	"net"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func RegisterRouters(e *echo.Echo) {
	server := new(ServerController)
	server.registerRouter(e.Group("/wxserver"))
}

func registerFilter(e *echo.Echo) {
	// middleware
	e.Pre(filters.SetLogTrace())
	e.Use(filters.GetLogMiddleware())
	e.Use(middleware.Recover())
}

type customValidator struct {
}

func (cv *customValidator) Validate(i interface{}) error {
	if _, err := govalidator.ValidateStruct(i); err != nil {
		return err
	}
	return nil
}

// SetupRouter 启动路由，添加中间件
func SetupRouter() *echo.Echo {
	e := echo.New()
	e.Server.ReadTimeout = 5 * time.Second
	e.Server.WriteTimeout = 10 * time.Second
	e.Server.SetKeepAlivesEnabled(false) //禁用KeepAlives 客户端连接不稳定很耗资源
	//在go1.3之后，提供了一个ConnState的hook，我们能通过这个来获取到对应的connection，这样在服务结束的时候我们就能够close掉这个connection了。该hook会在如下几种ConnState状态的时候调用。
	//
	//StateNew：新的连接，并且马上准备发送请求了
	//StateActive：表明一个connection已经接收到一个或者多个字节的请求数据，在server调用实际的handler之前调用hook。
	//StateIdle：表明一个connection已经处理完成一次请求，但因为是keepalived的，所以不会close，继续等待下一次请求。
	//StateHijacked：表明外部调用了hijack，最终状态。
	//StateClosed：表明connection已经结束掉了，最终状态。
	//通常，我们不会进入hijacked的状态（如果是websocket就得考虑了），所以一个可能的hook函数如下，
	//参考http://rcrowley.org/talks/gophercon-2014.html
	e.Server.ConnState = func(conn net.Conn, state http.ConnState) {
		switch state {
		case http.StateNew:

		case http.StateActive:

		case http.StateIdle:

		case http.StateHijacked, http.StateClosed:

		}
	}
	govalidator.TagMap["alphaversion"] = govalidator.Validator(func(str string) bool {
		return str == "alpha" || str == "beta" || str == "rc" || str == "release"
	})
	e.Validator = &customValidator{}
	registerFilter(e)
	RegisterRouters(e)
	return e
}
