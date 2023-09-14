# 获取微信 AccessToken

## 接口

`http://wxtoken.xxx.cn/wxserver/get_accesstoken`

## 请求参数：

| 参数 | 必选 | 说明 |
| :-: | :-: | :-: |
| wxapp_id | 是 | 微信 appId |
| secret | 是 | 微信 appsecret |
| access_token | 否 | 若传则认为该 accessToken 已过期，重新获取新的 accessToken |

## 返回信息

| 参数 | 说明 |
| :-: | :-: |
| code | 0:获取成功 1000:参数错误 1001:数据错误 1002:服务发送异常 1003:redis错误|
| data | access_token:获取到的凭证 expires_in:凭证有效时间，单位：秒 |
| msg | 信息 |
| trace | 跟踪信息 |

```json
{
    "code": 0,
    "data": {
        "access_token": "15_XIuITClkEG4Nm_cQKpZnhW0utv-uGR_1U2XFTO2877L0vHZFLoNMQgFMcV7sFUEHbm2PBRedVSomXXuLBf02cJx2n-seZKE5RxkNO02eWu3b8qiUIsgP3cPrg3MFNiBpmjL_KI5jYA0uhcCaMIBiCEADZS",
        "expires_in": 3640
    },
    "msg": "",
    "trace": {
        "Trace": "27c89ad809a70148",
        "Span": "06dfcfd6e6342343"
    }
}
```

## 使用指南

高频次场景
    client启动的时候请求 http://wxtoken.xxx.cn/wxserver/get_accesstoken?wxapp_id=&secret= 获取中控服最新token，并保存到全局变量，同时设定一个定时器在 expires_in 秒之后再次请求，以此循环；
    注意：expires_in 无需自行留出时间间隔，中控服已有考虑；

低频次场景
    每次都通过中控服获取token，不关心过期时间；
