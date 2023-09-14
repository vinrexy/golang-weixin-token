package utils

import (
	"github.com/liangdas/mqant/utils/uuid"
	"strings"
	"crypto/md5"
	"io"
	"fmt"
)

func GenerateUUID() string {
	return strings.Replace(uuid.Rand().Hex(), "-", "", -1)
}
// MD5 md5加密
func MD5(data string) string {
	t := md5.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

//HasIntElem 是否有元素
func HasIntElem(arr *[]int, ele int) bool {
	size := len(*arr)
	for i := 0; i < size; i++ {
		if (*arr)[i] == ele {
			return true
		}
	}
	return false
}
//HasStrElem 是否有元素
func HasStrElem(arr *[]string, ele string) bool {
	size := len(*arr)
	for i := 0; i < size; i++ {
		if (*arr)[i] == ele {
			return true
		}
	}
	return false
}
//Signature 签名字符串 md5(method+headers+uri+secret).toUpper()
func Signature(method, headers, url, secret string) string {
	return strings.ToUpper(MD5(method + headers + url + secret))
}