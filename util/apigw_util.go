package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

// HMAC-SHA256算法后，Base64编码
func Base64AndSha256HMAC(param string, secret string) string {
	shaResult := Sha256HMAC(param, secret)
	return base64.StdEncoding.EncodeToString([]byte(shaResult))
}

// HMAC-SHA256算法，16进制格式
func Sha256HMAC(data string, secret string) string {
	sign := hmac.New(sha256.New, []byte(secret))
	sign.Write([]byte(data))
	return hex.EncodeToString(sign.Sum(nil))
}

// URL请求参数编码
func UrlEncoder(param string) string {
	return url.QueryEscape(param)
}

// 随机正整数
func GetApiGwNonce() string {
	build := &strings.Builder{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 10; i++ {
		build.WriteString(fmt.Sprintf("%d", r.Int31n(10)))
	}
	return build.String()
}
