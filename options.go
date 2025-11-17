package sdk

import (
	"net/http"
	"time"
)

// Option 客户端配置选项
type Option func(*Client)

// WithHTTPClient 设置自定义 HTTP 客户端
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithTimeout 设置请求超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithTransport 设置自定义 Transport
func WithTransport(transport http.RoundTripper) Option {
	return func(c *Client) {
		c.httpClient.Transport = transport
	}
}
