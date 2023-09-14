package filters

import (
	"github.com/labstack/echo"
	"fmt"
	"runtime"
	"github.com/liangdas/mqant/log"
)
var (
	StackSize=         4 << 10 // 4 KB
	DisableStackAll=   false
	DisablePrintStack= false
)
func GetRecoverMiddleware() echo.MiddlewareFunc  {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, StackSize)
					length := runtime.Stack(stack, !DisableStackAll)
					if !DisablePrintStack {
						log.TError(ctx.Get("trace").(log.TraceSpan),"[PANIC RECOVER] %v %s\n", err, stack[:length])
					}
					ctx.Error(err)
				}
			}()
			return next(ctx)
		}
	}
}
