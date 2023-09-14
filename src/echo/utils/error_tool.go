package utils

func CatchException(f func(string))  {
	if err := recover(); err != nil {
		switch r := err.(type) {
		case string:
			f(r)
		case error:
			f(r.Error())
		}
	}
}

const ErrorOK  = 0

const (
	ErrorParams = iota + 1000
	ErrorData
	ErrorServer
	ErrorDB
)

var errorMsg = map[int]string{
	ErrorParams : "参数错误",
	ErrorData: "数据错误",
	ErrorServer: "服务发送异常",
	ErrorDB: "数据库错误",
}

func getErrorMsg(code int) string {
	 if msg, ok := errorMsg[code]; ok {
	 	return msg
	 }
	 return Int2String(code)
}
