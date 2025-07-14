package client

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/GHOSTEDDIE/apigwclient-go-sdk/config"
	"github.com/GHOSTEDDIE/apigwclient-go-sdk/model"
	"github.com/GHOSTEDDIE/apigwclient-go-sdk/util"
	"io/ioutil"
	"net"
	"net/http"
	"sort"
	"strings"
)

var DefaultSignHeaders = []string{"X-Api-Nonce", "X-Api-TimeStamp"}

type ApiGwClient struct {
	Config *config.ApiGwConfig

	Transport http.Transport
}

func (client *ApiGwClient) GetConfig() *config.ApiGwConfig {
	return client.Config
}

func InitApiGwClient(config *config.ApiGwConfig) *ApiGwClient {
	transport := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   config.DialTimeout,
			KeepAlive: config.DialKeepAlive,
		}).DialContext,
		MaxIdleConns:          config.GetMaxIdleConns(),
		IdleConnTimeout:       config.GetIdleConnTimeout(),
		TLSHandshakeTimeout:   config.GetTLSHandshakeTimeout(),
		ExpectContinueTimeout: config.GetExpectContinueTimeout(),
		DisableKeepAlives:     config.GetDisableKeepAlives(),
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	return &ApiGwClient{Config: config, Transport: transport}
}

func InitApiGwClientByCustomer(config *config.ApiGwConfig, transport http.Transport) *ApiGwClient {
	return &ApiGwClient{Config: config, Transport: transport}
}

func (client *ApiGwClient) sortSignHeaders(request *model.ApiGwRequest) []string {
	var signHeaders []string
	if len(request.SignHeaders) > 0 {
		signHeaders = append(signHeaders, request.SignHeaders...)
		signHeaders = append(signHeaders, DefaultSignHeaders...)
		sort.Strings(signHeaders)
	} else {
		signHeaders = DefaultSignHeaders
	}
	return signHeaders
}

func (client *ApiGwClient) buildSignature(request *model.ApiGwRequest, signHeaders []string) string {
	bufSign := &strings.Builder{}
	bufSign.WriteString(request.Method)
	bufSign.WriteString("\n")
	bufSign.WriteString(request.GetPathEncode())
	bufSign.WriteString("\n")
	bufSign.WriteString(request.GetQueryEncode())
	bufSign.WriteString("\n")
	for _, value := range signHeaders {
		bufSign.WriteString(strings.ToLower(value))
		bufSign.WriteString(":")
		bufSign.WriteString(request.Headers[value])
		bufSign.WriteString("\n")
	}
	return util.Base64AndSha256HMAC(bufSign.String(), client.Config.ClientSecret)
}

func (client *ApiGwClient) setSignature(request *model.ApiGwRequest) {
	signHeaders := client.sortSignHeaders(request)
	request.AddHeader("X-Api-ClientID", client.Config.ClientId)
	request.AddHeader("X-Api-Auth-Version", "2.0")
	request.AddHeader("X-Api-TimeStamp", request.GetTimeStamp())
	request.AddHeader("X-Api-Nonce", request.GetNonce())
	request.AddHeader("X-Api-SignHeaders", strings.Join(signHeaders, ","))
	request.AddHeader("X-Api-Signature", client.buildSignature(request, signHeaders))
	request.SignHeaders = signHeaders
}

func (client *ApiGwClient) getHttpClient() *http.Client {
	httpClient := &http.Client{Transport: &client.Transport}
	httpClient.Timeout = client.Config.ClientTimeOut
	return httpClient
}

// 同步请求
func (client *ApiGwClient) SyncSend(request *model.ApiGwRequest) (*model.ApiResult, error) {
	client.setSignature(request)
	return client.send(request)
}

// 异步请求
func (client *ApiGwClient) AsyncSend(request *model.ApiGwRequest, callBack model.ApiCallBack) {
	client.setSignature(request)
	go client.asyncSend(request, callBack)
}

func (client *ApiGwClient) asyncSend(request *model.ApiGwRequest, callBack model.ApiCallBack) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("async http send error", err)
		}
	}()
	response, err := client.send(request)
	if err != nil {
		callBack.OnFailure(request, err)
	} else {
		callBack.OnResponse(request, response)
	}
}

func (client *ApiGwClient) send(request *model.ApiGwRequest) (*model.ApiResult, error) {
	httpReq, errReq := buildHttpRequest(request)
	if errReq != nil {
		fmt.Println("build http request error ", errReq)
		return nil, errReq
	}
	resp, err := client.getHttpClient().Do(httpReq)
	if err != nil {
		fmt.Println("do http error ", err)
		return nil, err
	}

	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("http response read error ", err)
		return nil, err
	}
	result := &model.ApiResult{Status: resp.Status, StatusCode: resp.StatusCode, Body: respData, Header: resp.Header, Request: resp.Request}
	return result, nil
}

func buildHttpRequest(request *model.ApiGwRequest) (*http.Request, error) {
	httpReq, errReq := http.NewRequest(request.Method, fmt.Sprintf("%s%s", request.Host, request.Path), bytes.NewBuffer(request.Body))
	if errReq != nil {
		fmt.Println("create request error ", errReq)
		return httpReq, errReq
	}
	for k, v := range request.Headers {
		httpReq.Header[k] = []string{v}
	}
	httpReq.Header.Add("Content-type", request.ContentType)

	httpReq.URL.RawQuery = request.GetEncoderQueryString()

	return httpReq, errReq
}
