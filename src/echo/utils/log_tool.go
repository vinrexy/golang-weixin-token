package utils

import (
	"github.com/labstack/echo"
	"github.com/liangdas/mqant/log"
)
// Logger 日志对象
type Logger struct {
	trace log.TraceSpan
}
// Debugf 调试格式化打印
func (logger *Logger) Debugf(format string, a ... interface{})  {
	log.TDebug(logger.trace, format, a...)
}
// Errorf 错误格式化打印
func (logger *Logger) Errorf(format string, a ... interface{})  {
	log.TError(logger.trace, format, a...)
}
// Warnf 警告格式化打印
func (logger *Logger) Warnf(format string, a ... interface{})  {
	log.TWarning(logger.trace, format, a...)
}
// Infof 重要信息格式化打印
func (logger *Logger) Infof(format string, a ... interface{})  {
	log.TInfo(logger.trace, format, a...)
}
// GetLogger 获取 Logger 对象
func GetLogger(ctx echo.Context) *Logger {
	return &Logger{
		trace: ctx.Get("trace").(log.TraceSpan),
	}
}