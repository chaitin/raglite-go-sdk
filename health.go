package sdk

import "context"

// HealthService 健康检查服务
type HealthService struct {
	client *Client
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

// Check 健康检查
func (s *HealthService) Check(ctx context.Context) (*HealthResponse, error) {
	var result HealthResponse
	err := s.client.do(ctx, "GET", "/health", nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
