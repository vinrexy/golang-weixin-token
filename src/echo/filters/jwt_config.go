package filters

import (
	"github.com/labstack/echo/middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"fmt"
)

type JwtCustomClaims struct {
	UserId 	int64 `json:"user_id"`
	jwt.StandardClaims
}

type JwtUserClaims struct {
	UserId 	string 	`json:"user_id"`
	IP 		string 	`json:"ip"`
	TS 		int64 	`json:"ts"`
	Sign	string	`json:"sign"`
	jwt.StandardClaims
}

var (
	HS256_KEY=[]byte("secret")

	RSA_Public_Key=[]byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAo9INPwxdAG6DmDx7dybc
Bd6EOulCXOjv9a8H8dY8IRO/0QpuIpbDx0XrjRlcFfRcp7YAwaUkraMnk5Q1VKgM
DYj2Q6w80k7y2nJGyzVBzHAYTqkw7iyamb5HPWirnSheaDl6NINZ/zgzHrKvTU1R
MOL5Gm7yo4VwoFbcIcSDdTzIWHI+TUuv4hjiqgwVc5r6+B7K+lhfUraOnlUbJtFz
4LakWGQrpuXG7TLQOogBe+YR0eQtjVqICdDXXP51Ypt5ovaP13NSN3KLKEK7zrWJ
SLN5/nHDnfkoZdRjV0fsk//bkIwLshDi6cbQherm0nvkJd0pLtyRnxRAeR7WPizu
1wIDAQAB
-----END PUBLIC KEY-----`)
	RSA_Private_Key=[]byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAo9INPwxdAG6DmDx7dybcBd6EOulCXOjv9a8H8dY8IRO/0Qpu
IpbDx0XrjRlcFfRcp7YAwaUkraMnk5Q1VKgMDYj2Q6w80k7y2nJGyzVBzHAYTqkw
7iyamb5HPWirnSheaDl6NINZ/zgzHrKvTU1RMOL5Gm7yo4VwoFbcIcSDdTzIWHI+
TUuv4hjiqgwVc5r6+B7K+lhfUraOnlUbJtFz4LakWGQrpuXG7TLQOogBe+YR0eQt
jVqICdDXXP51Ypt5ovaP13NSN3KLKEK7zrWJSLN5/nHDnfkoZdRjV0fsk//bkIwL
shDi6cbQherm0nvkJd0pLtyRnxRAeR7WPizu1wIDAQABAoIBAQCfYtgWdp3ptJx+
OqJZbEp3v2Zhtt8lbFFDdTdCmRDJkeB3rzE2n/79W87w6jKI/cZEOjUEMvu7oNj5
oiI9Kn8HkDEh6GsIp11rIfI09az+DjXoGn8LzAPi/5lqavNFpagsuXdnrkCaqwA5
ptoeWNJcwQhiKn4SkNridYAZrovEPH7wvtW1tNwQ5KO4NnxeTXrXIA0YMMiIjNYX
jsOT7dKARZ+4DsklOXpUeysZPVAI0qYQGGhSiToQ+GjLKJcrzgWqurtPMCqs8v33
STaRTSpMjVIvDstnfpuJA8/7ac2PiflLP7TpGZDDqOAI7Zt9cnzgDQrLzEl2IsIr
57nYLfxRAoGBANlaVvprjdM9UeBEdUxVQY2tKQbAhsED3KSi0nOhdj9NWi0PDKGN
bbJvEW7NkQbtYjUfb2JznUDiwS9Wp4CN7tEwQhnMSp8Yo7eowWIYtWGPKb+o1Mnu
9CSsPxXP9ZCraTXdK+qTC8yxCWHpuBGg7+aED2AwxaA3eKO4PVnU/D5JAoGBAMDy
+mrup7qVH+Vk3pjjiDlA8YB4UJLZJl7cmLIJfvniy3spgIS1HW0WxQc76NI7F2ep
JNPbcUYqBKkbQ3ypAQMydiBJTZ73Gipd3Ca5XLQX/DqboOh1by34jAdMtJxn8SKQ
8anmhEhJQCQt50+7KRhhhaqrPAX0g0TZaXySZkQfAoGAA9cqzkX0PZVJyxKql+yx
udUjcnEYcHSnA2m1GkHyGvA89arcaEZdd9eqkTCkrWCoaZPinfS5BJp9G18Gmqjn
XV7i7B3F+8WtruMWd6tEGTM0Y6SSDfdg7Pz2KGaCSkodE8ySqBRtEvLV3ZsJm5Yi
ZwpSUzrJYylXwlzRCLNQubECgYAlqUenx52FlcX8CIxKW18jjcGVyeYwQ6Jxsa08
Uw4tyE7fY2JqhM+Rk3gxyUfQgSg4W5OMprCdeWYfe+rYUkSYUykrdCNqe+DnlBp8
lIG7xVK+PdJSjVl+J51tb1Nxk/hFPvVsrEn1shaK+UrFDUsgLyjf/zxgDTHyJl2o
qwq7EQKBgEX/C48LBBRiOSqYqF5QM27ZS+UbkQNDNaIBaCQrxxmohOVSfyt3B6qA
Y8QdXutAWsGX/yC2K6FHLvJ9WVQordWAVy028HFZvGJ+DkaOLYDQtWRKCtgGcXgz
g7/2JPiVsVQPrdBSmOECDJzLoHhbQsU8lcN1+bLT+5WRF6NyL3K5
-----END RSA PRIVATE KEY-----
`)
)

var jwtOptions map[string]interface{}

func skipRoute(ctx echo.Context) bool {
	if jwtOptions == nil {
		return false
	}
	return false
}

func GetJWT(ctx echo.Context) (*JwtCustomClaims,error) {
	if ctx==nil{
		return nil,errors.New("echo.Context is nil")
	}
	jwttoken:=ctx.Get("JwtToken")
	jwt_token,ok:=jwttoken.(*jwt.Token)
	if !ok{
		return nil,errors.Errorf("JwtToken is not *jwt.Token %v",jwt_token)
	}
	return jwt_token.Claims.(*JwtCustomClaims),nil
}

func GetSDKJWT(ctx echo.Context) (*JwtUserClaims,error) {
	if ctx==nil{
		return nil,errors.New("echo.Context is nil")
	}
	jwttoken:=ctx.Get("JwtToken")
	jwt_token,ok:=jwttoken.(*jwt.Token)
	if !ok{
		return nil,errors.Errorf("JwtToken is not *jwt.Token %v",jwt_token)
	}
	return jwt_token.Claims.(*JwtUserClaims),nil
}

func GetJWTConfig(options map[string]interface{}) middleware.JWTConfig {
	jwtOptions = options
	return middleware.JWTConfig{
		Skipper: skipRoute,
		Claims:     &JwtCustomClaims{},
		ContextKey:"JwtToken",
		SigningKey:HS256_KEY,
	}
}

func GetJWTRSAConfig(options map[string]interface{}) (middleware.JWTConfig) {
	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(RSA_Public_Key);
	if  err != nil {
		panic(fmt.Sprintf("Unable to parse RSA public key: %v", err))
	}
	jwtOptions = options
	return middleware.JWTConfig{
		Skipper: skipRoute,
		Claims:     &JwtUserClaims{},
		ContextKey:"JwtToken",
		SigningKey: rsaPublicKey,
		SigningMethod:jwt.SigningMethodRS512.Alg(),
	}
}