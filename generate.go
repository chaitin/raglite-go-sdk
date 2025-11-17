package sdk

import "context"

// GenerateService 生成服务
type GenerateService struct {
	client *Client
}

// GenerateRequest 生成请求
type GenerateRequest struct {
	Query     string `json:"query"`
	Context   string `json:"context"`
	DatasetID string `json:"dataset_id"`
}

// GenerateResponse 生成响应
type GenerateResponse struct {
	Answer string `json:"answer"`
}

// Generate 生成答案（不检索，直接生成）
func (s *GenerateService) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	var result GenerateResponse
	err := s.client.do(ctx, "POST", "/api/v1/generate", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
