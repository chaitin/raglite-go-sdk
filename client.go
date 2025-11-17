package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client RAGLite SDK 客户端
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client

	// Services
	Models    *ModelsService
	Datasets  *DatasetsService
	Documents *DocumentsService
	Search    *SearchService
	QA        *QAService
	Generate  *GenerateService
	Health    *HealthService
}

// NewClient 创建新的 SDK 客户端
func NewClient(baseURL string, opts ...Option) (*Client, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("baseURL is required")
	}

	c := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// 应用选项
	for _, opt := range opts {
		opt(c)
	}

	// 初始化各个服务
	c.Models = &ModelsService{client: c}
	c.Datasets = &DatasetsService{client: c}
	c.Documents = &DocumentsService{client: c}
	c.Search = &SearchService{client: c}
	c.QA = &QAService{client: c}
	c.Generate = &GenerateService{client: c}
	c.Health = &HealthService{client: c}

	return c, nil
}

// do 执行 HTTP 请求
func (c *Client) do(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}

	fullURL := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// 设置 API Key
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// 检查 HTTP 状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr APIError
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			return &APIError{
				StatusCode: resp.StatusCode,
				Message:    string(respBody),
			}
		}
		apiErr.StatusCode = resp.StatusCode
		return &apiErr
	}

	// 解析响应
	if result != nil {
		var apiResp APIResponse
		apiResp.Data = result
		if err := json.Unmarshal(respBody, &apiResp); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}

		if !apiResp.Success {
			return &APIError{
				StatusCode: resp.StatusCode,
				Message:    apiResp.Message,
			}
		}
	}

	return nil
}

// buildURL 构建带查询参数的 URL
func (c *Client) buildURL(path string, params map[string]string) string {
	if len(params) == 0 {
		return path
	}

	u, _ := url.Parse(path)
	q := u.Query()
	for k, v := range params {
		if v != "" {
			q.Set(k, v)
		}
	}
	u.RawQuery = q.Encode()
	return u.String()
}
