package config

import (
	"time"
)

type ApiGwConfig struct {
	// 应用ID，不能为空
	ClientId string

	// 应用密钥，不能为空
	ClientSecret string

	// 客户端超时时间，不能为空
	ClientTimeOut time.Duration

	DialTimeout time.Duration

	DialKeepAlive time.Duration

	MaxIdleConns int

	IdleConnTimeout time.Duration

	TLSHandshakeTimeout time.Duration

	ExpectContinueTimeout time.Duration

	DisableKeepAlives bool
}

func (config *ApiGwConfig) GetClientId() string {
	return config.ClientId
}

func (config *ApiGwConfig) SetClientId(clientId string) {
	config.ClientId = clientId
}

func (config *ApiGwConfig) GetClientSecret() string {
	return config.ClientSecret
}

func (config *ApiGwConfig) SetClientSecret(clientSecret string) {
	config.ClientSecret = clientSecret
}

func (config *ApiGwConfig) GetClientTimeOut() time.Duration {
	return config.ClientTimeOut
}

func (config *ApiGwConfig) SetClientTimeOut(clientTimeOut time.Duration) {
	config.ClientTimeOut = clientTimeOut
}

func (config *ApiGwConfig) GetDialTimeout() time.Duration {
	return config.DialTimeout
}

func (config *ApiGwConfig) SetDialTimeout(dialTimeout time.Duration) {
	config.DialTimeout = dialTimeout
}

func (config *ApiGwConfig) GetDialKeepAlive() time.Duration {
	return config.DialKeepAlive
}

func (config *ApiGwConfig) SetDialKeepAlive(dialKeepAlive time.Duration) {
	config.DialKeepAlive = dialKeepAlive
}

func (config *ApiGwConfig) GetMaxIdleConns() int {
	return config.MaxIdleConns
}

func (config *ApiGwConfig) SetMaxIdleConns(maxIdleConns int) {
	config.MaxIdleConns = maxIdleConns
}

func (config *ApiGwConfig) GetIdleConnTimeout() time.Duration {
	return config.IdleConnTimeout
}

func (config *ApiGwConfig) SetIdleConnTimeout(idleConnTimeout time.Duration) {
	config.IdleConnTimeout = idleConnTimeout
}

func (config *ApiGwConfig) GetTLSHandshakeTimeout() time.Duration {
	return config.TLSHandshakeTimeout
}

func (config *ApiGwConfig) SetTLSHandshakeTimeout(tLSHandshakeTimeout time.Duration) {
	config.TLSHandshakeTimeout = tLSHandshakeTimeout
}

func (config *ApiGwConfig) GetExpectContinueTimeout() time.Duration {
	return config.ExpectContinueTimeout
}

func (config *ApiGwConfig) SetExpectContinueTimeout(expectContinueTimeout time.Duration) {
	config.ExpectContinueTimeout = expectContinueTimeout
}

func (config *ApiGwConfig) GetDisableKeepAlives() bool {
	return config.DisableKeepAlives
}

func (config *ApiGwConfig) SetDisableKeepAlives(disableKeepAlives bool) {
	config.DisableKeepAlives = disableKeepAlives
}
