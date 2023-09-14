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
	e.Server.SetKeepAlivesEnabled(false)
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
