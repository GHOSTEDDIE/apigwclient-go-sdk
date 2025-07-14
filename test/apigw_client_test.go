package test

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/GHOSTEDDIE/apigwclient-go-sdk/client"
	"github.com/GHOSTEDDIE/apigwclient-go-sdk/config"
	"github.com/GHOSTEDDIE/apigwclient-go-sdk/model"
	"github.com/GHOSTEDDIE/apigwclient-go-sdk/util"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

// 根据配置构建http client的transport
func TestApiGwConfig(t *testing.T) {
	gwConfig := buildApiGwConfig()
	fmt.Println(gwConfig.ClientTimeOut.String())
}

func buildApiGwConfig() *config.ApiGwConfig {
	gwConfig := &config.ApiGwConfig{
		ClientId:              getClientId(),
		ClientSecret:          getClientSecret(),
		ClientTimeOut:         3 * time.Second,
		DialTimeout:           15 * time.Second,
		DialKeepAlive:         15 * time.Second,
		MaxIdleConns:          100,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     true,
	}
	return gwConfig
}

func getClientId() string {
	return "236514"
}

func getClientSecret() string {
	return "3eaf929860ac418d6dcd4de903913787"
}

// 自定义http client的transport
func TestApiGwConfigCustomer(t *testing.T) {
	gwConfig := config.ApiGwConfig{
		ClientId:      getClientId(),
		ClientSecret:  getClientSecret(),
		ClientTimeOut: 60 * time.Second,
	}
	fmt.Printf(gwConfig.ClientTimeOut.String())
	transport := buildTransport()
	fmt.Println(transport.IdleConnTimeout.String())
}

func buildTransport() http.Transport {
	transport := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   15 * time.Second,
			KeepAlive: 15 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     true,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	return transport
}

func TestInitApiGwClient(t *testing.T) {
	gwClient := client.InitApiGwClient(buildApiGwConfig())
	fmt.Println(gwClient.Config.MaxIdleConns)
}

func TestInitApiGwClientByCustomer(t *testing.T) {
	gwClient := client.InitApiGwClientByCustomer(buildApiGwConfig(), buildTransport())
	fmt.Println(gwClient.Config.MaxIdleConns)
}

func getHost() string {
	return "https://api.kingdee.com"
}

func TestApiGwRequestInit(t *testing.T) {
	gwRequest := &model.ApiGwRequest{Host: getHost(), Path: "/jdy/test/get", Method: "GET", ContentType: "application/json; charset=utf-8"}
	fmt.Println(gwRequest.ContentType)

	gwRequest = &model.ApiGwRequest{}
	gwRequest.Host = getHost()
	gwRequest.Path = "/jdy/test/get"
	gwRequest.SetBodyJson([]byte(""))
	gwRequest.SetBodyText([]byte(""))
	gwRequest.AddHeader("x-test", "test")
	gwRequest.AddSignHeaders("x-test", "test")

	header := make(map[string]string, 0)
	header["x-test"] = "test"
	gwRequest.AddHeaders(header)
}

func TestApiGwSendRequest(t *testing.T) {
	gwClient := client.InitApiGwClientByCustomer(buildApiGwConfig(), buildTransport())
	gwRequest := &model.ApiGwRequest{Host: getHost(), Path: "/jdyconnector/app_management/kingdee_auth_token", Method: "GET", ContentType: "application/json; charset=utf-8"}
	response, err := gwClient.SyncSend(gwRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.StatusCode)
	fmt.Println(string(response.Body))
}

func TestApiGwSendRequestByQuery(t *testing.T) {
	gwClient := client.InitApiGwClientByCustomer(buildApiGwConfig(), buildTransport())
	gwRequest := &model.ApiGwRequest{Host: getHost(), Path: "/jdyconnector/app_management/kingdee_auth_token", Method: "GET", ContentType: "application/json; charset=utf-8"}
	appKey := "IDFdwQGf"
	appSignature := util.Base64AndSha256HMAC(appKey, gwClient.Config.ClientSecret)
	query := make(map[string]string, 0)
	query["app_key"] = appKey
	query["app_signature"] = appSignature
	gwRequest.Query = query
	response, err := gwClient.SyncSend(gwRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.StatusCode)
	fmt.Println(string(response.Body))
}

func getToken() string {
	return "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHQiOnsiYWNjb3VudElkIjoiMTY0MTM2MDk5MzAxMDMyNDI1MCIsImdyb3VwTmFtZSI6Im5zLXByZTEiLCJhcHBfa2V5IjoiZnVscGNIVDciLCJ0ZW5hbnRJZCI6Ijc5ODY2OTM5OTU4IiwidXNlck5hbWUiOiJraW5nZGVldGVzdDAwMDAifSwiZ3JwIjoibnMtcHJlMSIsImV4cCI6MTY3MzQ1MzE5MCwiYWlkIjoiMTY0MTM2MDk5MzAxMDMyNDI1MCIsImlhdCI6MTY3MzM2Njc5MH0.GJADrTJBaspA_ukJpmUDpUNf67ZAcvDOtTiKXv-l6oo"
}

func TestApiGwSendRequestParam(t *testing.T) {
	gwClient := client.InitApiGwClientByCustomer(buildApiGwConfig(), buildTransport())
	gwRequest := &model.ApiGwRequest{Host: getHost(), Path: "/jdy/v2/bd/material", Method: "GET", ContentType: "application/json; charset=utf-8"}
	//gwRequest.AddHeader("app-token", getToken())
	query := make(map[string]string, 0)
	query["page"] = "1"
	query["page_size"] = "10"
	query["enable"] = "1"
	query["parent"] = "1307741726651203584"
	gwRequest.Query = query
	response, err := gwClient.SyncSend(gwRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.StatusCode)
	fmt.Println(string(response.Body))
}

// JSON格式请求
func TestApiGwSendRequestJsonBody(t *testing.T) {
	gwClient := client.InitApiGwClientByCustomer(buildApiGwConfig(), buildTransport())
	gwRequest := &model.ApiGwRequest{Host: getHost(), Path: "/jdy/v2/ls/member", Method: "POST", ContentType: "application/json; charset=utf-8"}
	bodyStr := "{\"birthday\":\"2022-08-01\",\"end_date\":\"2022-01-01\",\"custom_no\":\"1\",\"join_way\":\"1\",\"period\":\"1\",\"address\":\"410\",\"shop\":\"850693217107061760\",\"level\":\"1\",\"detail_address\":\"XXX\",\"mobile\":\"13333333333\",\"referrere_mp\":\"1\",\"remark\":\"备注\",\"label\":null,\"source\":\"1\",\"mbcard_id\":\"1\",\"is_delete\":\"0\",\"sale_shop\":\"1\",\"enable\":\"1\",\"name\":\"会员A\",\"sex_group\":\"1\",\"customer\":\"001\",\"start_date\":\"2022-01-01\"}"
	gwRequest.SetBodyJson([]byte(bodyStr))
	response, err := gwClient.SyncSend(gwRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.StatusCode)
	fmt.Println(string(response.Body))
}

// PUT请求
func TestApiGwSendRequestPut(t *testing.T) {
	gwClient := client.InitApiGwClientByCustomer(buildApiGwConfig(), buildTransport())
	gwRequest := &model.ApiGwRequest{Host: getHost(), Path: "/jdy/test/put", Method: "PUT", ContentType: "application/json; charset=utf-8"}
	query := make(map[string]string, 0)
	query["outerInstanceId"] = "1307741726651203584"
	gwRequest.Query = query
	response, err := gwClient.SyncSend(gwRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.StatusCode)
	fmt.Println(string(response.Body))
}

// DELETE请求
func TestApiGwSendRequestDelete(t *testing.T) {
	gwClient := client.InitApiGwClientByCustomer(buildApiGwConfig(), buildTransport())
	gwRequest := &model.ApiGwRequest{Host: getHost(), Path: "/jdy/test/delete", Method: "DELETE", ContentType: "application/json; charset=utf-8"}
	query := make(map[string]string, 0)
	query["outerInstanceId"] = "1307741726651203584"
	gwRequest.Query = query
	response, err := gwClient.SyncSend(gwRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.StatusCode)
	fmt.Println(string(response.Body))
}

// 异步回调对象
type ApiCallBackTest struct {
	Name string
}

func (callBack ApiCallBackTest) OnFailure(request *model.ApiGwRequest, err error) {
	fmt.Println(request.Path)
	fmt.Println(err)
}

func (callBack ApiCallBackTest) OnResponse(request *model.ApiGwRequest, result *model.ApiResult) {
	fmt.Println(request.Path)
	fmt.Println(string(result.Body))
}

// 异步发起请求
func TestApiGwAsyncSendRequest(t *testing.T) {
	gwClient := client.InitApiGwClientByCustomer(buildApiGwConfig(), buildTransport())
	gwRequest := &model.ApiGwRequest{Host: getHost(), Path: "/jdy/test/get", Method: "GET", ContentType: "application/json; charset=utf-8"}
	query := make(map[string]string, 0)
	query["outerInstanceId"] = "1307741726651203584"
	gwRequest.Query = query
	var test ApiCallBackTest
	test.Name = ""
	var callBack model.ApiCallBack
	callBack = test
	gwClient.AsyncSend(gwRequest, callBack)
	time.Sleep(10 * time.Minute)
}

// 文件上传请求
func TestApiGwUploadFileRequest(t *testing.T) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", "F:\\go.txt")
	if err != nil {
		fmt.Println("error writing to buffer")
	}
	//打开文件句柄操作
	fh, err := os.Open("F:\\go.txt")
	if err != nil {
		fmt.Println("error opening file")
	}
	defer fh.Close()
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
	}
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	gwClient := client.InitApiGwClientByCustomer(buildApiGwConfig(), buildTransport())
	gwRequest := &model.ApiGwRequest{Host: getHost(), Path: "/jdy/test/attachment_upload", Method: "POST", ContentType: "application/json; charset=utf-8"}
	/* query := make(map[string]string, 0)
	   query["client_id"] = ""
	   query["uid"] = ""
	   gwRequest.Query = query
	   gwRequest.AddHeader("groupName", "")
	   gwRequest.AddHeader("accountId", "")*/
	gwRequest.ContentType = contentType
	gwRequest.Body = bodyBuf.Bytes()
	response, err := gwClient.SyncSend(gwRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(response.StatusCode)
	fmt.Println(string(response.Body))
}
