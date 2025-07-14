package model

import (
	"apigwclient-go-sdk/util"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

// 网关请求对象
type ApiGwRequest struct {
	Body []byte

	Host string

	Path string

	Method string

	Headers map[string]string

	Query map[string]string

	// 需要签名的请求头属性
	SignHeaders []string

	ContentType string

	// 请求时间戳，默认系统生成
	TimeStamp string

	// 请求随机数，默认系统生成
	Nonce string
}

func (request *ApiGwRequest) AddHeaders(headers map[string]string) {
	for param, value := range headers {
		request.AddHeader(param, value)
	}
}

func (request *ApiGwRequest) AddHeader(param string, value string) {
	if request.Headers == nil {
		request.Headers = make(map[string]string, 0)
	}
	request.Headers[param] = value
}

func (request *ApiGwRequest) AddSignHeaders(param string, value string) {
	request.SignHeaders = append(request.SignHeaders, param)
	request.AddHeader(param, value)
}

func (request *ApiGwRequest) SetBodyJson(body []byte) {
	request.ContentType = "application/json; charset=utf-8"
	request.Body = body
}

func (request *ApiGwRequest) SetBodyText(body []byte) {
	request.ContentType = "application/text; charset=utf-8"
	request.Body = body
}

// 获取随机数
func (request *ApiGwRequest) GetNonce() string {
	if request.Nonce == "" {
		request.Nonce = util.GetApiGwNonce()
	}
	return request.Nonce
}

func (request *ApiGwRequest) GetTimeStamp() string {
	if request.TimeStamp == "" {
		request.TimeStamp = fmt.Sprint(time.Now().Unix())
	}
	return request.TimeStamp
}

// 获取请求路径的url编码
func (request *ApiGwRequest) GetPathEncode() string {
	return util.UrlEncoder(request.Path)
}

// 获取请求路径参数的url编码，两次url编码
func (request *ApiGwRequest) GetQueryEncode() string {
	if len(request.Query) == 0 {
		return ""
	}
	rawQueryString := request.GetEncoderQueryString()
	queryStrings := strings.Split(rawQueryString, "&")
	result := make([]string, 0)
	for i := 0; i < len(queryStrings); i++ {
		param := queryStrings[i]
		index := strings.Index(param, "=")
		if index >= 1 {
			encodeParam := util.UrlEncoder(param[0:index])
			encodeValue := util.UrlEncoder(param[index+1:])
			result = append(result, fmt.Sprintf("%s=%s", encodeParam, encodeValue))
		}
	}
	return strings.Join(result, "&")
}

// 获取请求路径参数的url编码，一次url编码
func (request *ApiGwRequest) GetEncoderQueryString() string {
	result := make([]string, 0)
	for param, value := range request.Query {
		encodeParam := util.UrlEncoder(param)
		encodeValue := util.UrlEncoder(value)
		result = append(result, fmt.Sprintf("%s=%s", encodeParam, encodeValue))
	}
	sort.Strings(result)
	return strings.Join(result, "&")
}

// 回调对象
type ApiCallBack interface {
	OnFailure(request *ApiGwRequest, err error)

	OnResponse(request *ApiGwRequest, result *ApiResult)
}

// 网关返回结果对象
type ApiResult struct {
	Status string

	StatusCode int

	Body []byte

	Header http.Header

	Request *http.Request
}
