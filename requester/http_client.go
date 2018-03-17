package requester

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"time"
)

var (
	// TLSConfig tls 链接配置
	TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
)

// HTTPClient http client
type HTTPClient struct {
	http.Client
	UserAgent string
}

// NewHTTPClient 返回 HTTPClient 的指针,
// 预设了一些配置
func NewHTTPClient() *HTTPClient {
	jar, _ := cookiejar.New(nil)
	return &HTTPClient{
		Client: http.Client{
			Transport: &http.Transport{
				Proxy:                 http.ProxyFromEnvironment,
				DialContext:           dialContext,
				Dial:                  dial,
				DialTLS:               dialTLS,
				TLSClientConfig:       TLSConfig,
				TLSHandshakeTimeout:   10 * time.Second,
				DisableKeepAlives:     false,
				DisableCompression:    false, // gzip
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
				ExpectContinueTimeout: 10 * time.Second,
			},
			Jar:     jar,
			Timeout: 30 * time.Second,
		},
		UserAgent: UserAgent,
	}
}

// SetCookiejar 设置cookie
func (h *HTTPClient) SetCookiejar(c *cookiejar.Jar) {
	if c != nil {
		h.Jar = c
		return
	}
	h.ResetCookiejar()
}

// ResetCookiejar 清空cookie
func (h *HTTPClient) ResetCookiejar() {
	h.Jar, _ = cookiejar.New(nil)
}

// SetHTTPSecure 是否启用https 安全检测,默认不检测
func (h *HTTPClient) SetHTTPSecure(b bool) {
	h.Transport.(*http.Transport).TLSClientConfig.InsecureSkipVerify = !b
}

// SetGzip 是否开启Gzip
func (h *HTTPClient) SetGzip(b bool) {
	h.Transport.(*http.Transport).DisableCompression = !b
}

// SetResponseHeaderTimeout 设置目标服务器响应超时时间
func (h *HTTPClient) SetResponseHeaderTimeout(t time.Duration) {
	h.Transport.(*http.Transport).ResponseHeaderTimeout = t
}

// SetTimeout 设置http 请求超时时间,默认为30s
func (h *HTTPClient) SetTimeout(t time.Duration) {
	h.Timeout = t
}
