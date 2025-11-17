package sdk

import "fmt"

// APIError API 错误
type APIError struct {
	StatusCode int
	Message    string
}

// Error 实现 error 接口
func (e *APIError) Error() string {
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}

// IsNotFound 是否为 404 错误
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == 404
}

// IsBadRequest 是否为 400 错误
func (e *APIError) IsBadRequest() bool {
	return e.StatusCode == 400
}

// IsServerError 是否为服务器错误（5xx）
func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}
