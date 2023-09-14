package filters

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/labstack/echo"
	"github.com/liangdas/mqant/log"
	"strings"
	"time"
)

type logger struct {

}

func (this *logger) Write(p []byte) (n int, err error) {

	return 0, nil
}

func GetLogMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			req := ctx.Request()
			res := ctx.Response()
			start := time.Now()
			defer func() {
				if p := recover(); p != nil {
					log.TError(ctx.Get("trace").(log.TraceSpan), "panic recover! %v %v %v p: %v", req.Method, res.Status, req.URL.Path, p)
				}
			}()
			if err = next(ctx); err != nil {
				log.TError(ctx.Get("trace").(log.TraceSpan), "Error: %v", err.Error())
				ctx.Error(err)
			}
			q, _ := GetQueryJsonData(ctx)
			log.TInfo(ctx.Get("trace").(log.TraceSpan), "%v %v %v %v [%v] %v in %v", req.Method, res.Status, req.URL.Path, req.Proto, ctx.RealIP(), q, time.Since(start))
			return
		}
	}
}

func GetQueryJsonData(ctx echo.Context) (string, error) {
	req := ctx.Request()
	if req.ContentLength == 0 {
		if req.Method == echo.GET || req.Method == echo.DELETE {
			b, err := json.Marshal(ctx.QueryParams())
			return string(b), err
		}
		return "{}", nil
	}
	ctype := req.Header.Get(echo.HeaderContentType)
	switch {
	case strings.HasPrefix(ctype, echo.MIMEApplicationJSON):
		i := map[string]interface{}{}
		if err := json.NewDecoder(req.Body).Decode(&i); err != nil {
			if ute, ok := err.(*json.UnmarshalTypeError); ok {
				return "", fmt.Errorf("unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)
			} else if se, ok := err.(*json.SyntaxError); ok {
				return "", fmt.Errorf("syntax error: offset=%v, error=%v", se.Offset, se.Error())
			} else {
				return "", err
			}
		}
		b, err := json.Marshal(i)
		return string(b), err
	case strings.HasPrefix(ctype, echo.MIMEApplicationXML), strings.HasPrefix(ctype, echo.MIMETextXML):
		i := map[string]interface{}{}
		if err := xml.NewDecoder(req.Body).Decode(&i); err != nil {
			if ute, ok := err.(*xml.UnsupportedTypeError); ok {
				return "", fmt.Errorf("unsupported type error: type=%v, error=%v", ute.Type, ute.Error())
			} else if se, ok := err.(*xml.SyntaxError); ok {
				return "", fmt.Errorf("syntax error: line=%v, error=%v", se.Line, se.Error())
			} else {
				return "", err
			}
		}
		b, err := json.Marshal(i)
		return string(b), err
	case strings.HasPrefix(ctype, echo.MIMEApplicationForm), strings.HasPrefix(ctype, echo.MIMEMultipartForm):
		params, err := ctx.FormParams()
		if err != nil {
			return "", err
		}
		b, err := json.Marshal(params)
		return string(b), err
	default:
		return "", fmt.Errorf("unsupported media type")
	}
}
