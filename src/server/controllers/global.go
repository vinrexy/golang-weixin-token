package server

import (
	"echo/utils"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"tools"

	"github.com/gomodule/redigo/redis"
	"github.com/labstack/echo"
	"github.com/liangdas/mqant/log"
)

// 请求数据
type RequestData struct {
	WXAppId     string `json:"wxapp_id" query:"wxapp_id" form:"wxapp_id" valid:"required"`
	Secret      string `json:"secret" query:"secret" form:"secret" valid:"required"`
	AccessToken string `json:"access_token" query:"access_token" form:"access_token" valid:"required"`
}

type ResultData struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// 微信返回数据
type wxResultData struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	ErrCode     int64  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

type ServerController struct {
}

func (self *ServerController) registerRouter(g *echo.Group) {
	g.Any("/get_accesstoken", self.getWXAccessToken)
}

func (self *ServerController) getWXAccessToken(ctx echo.Context) error {
	logger := utils.GetLogger(ctx)
	req := new(RequestData)
	retmsg := utils.NewRetMsg(ctx)

	remoteAddr := ctx.RealIP()
	if err := ctx.Bind(req); err != nil {
		logger.Errorf("getWXAccessToken Exception: >> RemoteAddr:%v error:%v", remoteAddr, err.Error())
		retmsg.PackError(utils.ErrorData, err.Error())
		return ctx.JSON(http.StatusOK, retmsg)
	}

	if req.WXAppId == "" || req.Secret == "" {
		logger.Errorf("getWXAccessToken 参数错误: >> RemoteAddr:%v WXAppId=%v, secret=%v", remoteAddr, req.WXAppId, req.Secret)
		retmsg.PackError(utils.ErrorParams)
		return ctx.JSON(http.StatusOK, retmsg)
	}

	self.getAccessToken(req, retmsg)
	if retmsg.Code != 0 {
		logger.Errorf("getWXAccessToken Exception: >> RemoteAddr:%v WXAppId:%+v RetMsg:%+v", remoteAddr, req.WXAppId, retmsg)
	} else {
		logger.Infof("getWXAccessToken: >> RemoteAddr:%v ResponseData:%+v", remoteAddr, retmsg.Data)
	}
	return ctx.JSON(http.StatusOK, retmsg)
}

var wxAccessToken = "wx_access_token:%s:%s"
var wxTokenLock = "wx_token_lock:%s:%s"

func (self *ServerController) getAccessToken(req *RequestData, retMsg *utils.RetMsg) {
	var accessToken []any
	var err error

	accessToken, err = tools.GetAccessToken(fmt.Sprintf(wxAccessToken, req.WXAppId, req.Secret))
	if err == nil || err == redis.ErrNil {
		if accessToken != nil && req.AccessToken != accessToken[0] {
			retMsg.PackResult(&ResultData{
				accessToken[0].(string),
				accessToken[1].(int64) + int64(1+rand.Intn(50)),
			})
			return
		} else {
			if locked, err := tools.LockToken(fmt.Sprintf(wxTokenLock, req.WXAppId, req.Secret)); err == nil {
				if locked[0].(int64) == 1 {
					tools.DeleteKey(fmt.Sprintf(wxAccessToken, req.WXAppId, req.Secret))
					self.getAccessTokenFromWX(req.WXAppId, req.Secret, retMsg)
					return
				} else {
					reqNum := 1
				loop:
					if reqNum <= 3 {
						reqNum++
						//等待1s
						time.Sleep(time.Second * time.Duration(1))
						token, err := tools.GetAccessToken(fmt.Sprintf(wxAccessToken, req.WXAppId, req.Secret))
						if err != nil {
							goto loop
						}
						retMsg.PackResult(&ResultData{
							token[0].(string),
							token[1].(int64) + int64(1+rand.Intn(50)),
						})
						return
					} else {
						retMsg.PackError(utils.ErrorServer, "请求超时")
						return
					}
				}
			} else {
				retMsg.PackError(utils.ErrorDB, err.Error())
				return
			}
		}
	} else {
		retMsg.PackError(utils.ErrorDB, err.Error())
		return
	}
}

func (self *ServerController) getAccessTokenFromWX(appId string, secret string, retMsg *utils.RetMsg) {
	defer tools.DeleteKey(fmt.Sprintf(wxTokenLock, appId, secret))

	var url = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" +
		appId + "&secret=" + secret
	client := tools.SimpleHttpClient()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		retMsg.PackError(utils.ErrorServer, "make weixin request error:%v", err.Error())
		return
	}
	result := tools.MakeHttpRequest(client, req)
	if result.Err != nil {
		retMsg.PackError(utils.ErrorServer, "weixin request error:%v", result.Err)
		return
	} else {
		if result.StatusCode != http.StatusOK {
			retMsg.PackError(utils.ErrorServer, "weixin request error with http statusCode:%v", result.StatusCode)
			return
		}

		wxResult := new(wxResultData)
		if err := json.Unmarshal(result.Body, wxResult); err != nil {
			retMsg.PackError(utils.ErrorData, err.Error())
			return
		}

		if wxResult.ErrCode != 0 {
			retMsg.PackError(utils.ErrorData, fmt.Sprintf("wxerrcode:%v, wxerrmsg:{%v}", wxResult.ErrCode, wxResult.ErrMsg))
			return
		} else {
			log.TInfo(retMsg.Trace, "WXAppId:%v  Body:%v", appId, string(result.Body))
			tools.SetAccessToken(fmt.Sprintf(wxAccessToken, appId, secret), wxResult.AccessToken, wxResult.ExpiresIn-180)
			retMsg.PackResult(&ResultData{
				wxResult.AccessToken,
				wxResult.ExpiresIn - 180,
			})
			return
		}
	}
}
