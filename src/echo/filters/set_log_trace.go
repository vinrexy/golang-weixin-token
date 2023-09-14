package filters

import (
	"github.com/labstack/echo"
	"github.com/liangdas/mqant/log"
)
// SetLogTrace 设置每个请求对日志跟踪对象
func SetLogTrace() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			req := ctx.Request()
			res := ctx.Response()
			id := req.Header.Get(echo.HeaderXRequestID)
			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			if id == "" {
				trace := log.CreateRootTrace()
				ctx.Set("trace", trace)
				req.Header.Set(echo.HeaderXRequestID, trace.TraceId())
			}else{
				trace := log.CreateTrace(id, log.CreateRootTrace().SpanId())
				ctx.Set("trace", trace)
			}
			return next(ctx)
		}
	}
}