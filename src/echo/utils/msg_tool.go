package utils

import (
	"fmt"
	"strings"

	"github.com/labstack/echo"
	"github.com/liangdas/mqant/log"
)

type RetMsg struct {
	Code  int           `json:"code"`
	Data  interface{}   `json:"data"`
	Msg   string        `json:"msg"`
	Trace log.TraceSpan `json:"trace"`
	ctx   echo.Context  `json:"-"`
}

func NewRetMsg(ctx echo.Context) *RetMsg {
	var trace log.TraceSpan = nil
	if ctx != nil {
		trace = ctx.Get("trace").(log.TraceSpan)
	}
	return &RetMsg{
		ctx:   ctx,
		Trace: trace,
	}
}

func (ret *RetMsg) GetTrace() log.TraceSpan {
	return ret.Trace
}

func getRequestID(ctx echo.Context) string {
	if ctx == nil {
		return ""
	}
	return ctx.Get("trace").(log.TraceSpan).TraceId()
}

func (ret *RetMsg) PackError(code int, msg ...interface{}) {
	ret.Code = code
	if len(msg) == 0 {
		ret.Msg = getErrorMsg(code)
	} else {
		f := strings.Repeat("%+v ", len(msg))
		ret.Msg = fmt.Sprintf(f, msg...)
	}
}

func (ret *RetMsg) Write(result []byte) {
	ret.ctx.Response().Write(result)
}

func (ret *RetMsg) PackResult(result interface{}) {
	ret.Data = result
}
