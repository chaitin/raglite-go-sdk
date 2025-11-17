package sdk

import "context"

// QAService 问答服务
type QAService struct {
	client *Client
}

// QARequest 问答请求
type QARequest struct {
	Query               string  `json:"query"`
	DatasetID           string  `json:"dataset_id"`
	TopK                int     `json:"top_k,omitempty"`
	RetrievalMode       string  `json:"retrieval_mode,omitempty"` // full | smart
	Stream              bool    `json:"stream,omitempty"`
	SimilarityThreshold float64 `json:"similarity_threshold,omitempty"`
}

// QAResponse 问答响应
type QAResponse struct {
	Answer  string         `json:"answer"`
	Context []SearchResult `json:"context"`
}

// Ask 提出问题并获取答案
func (s *QAService) Ask(ctx context.Context, req *QARequest) (*QAResponse, error) {
	// 设置默认值
	if req.TopK <= 0 {
		req.TopK = 10
	}

	var result QAResponse
	err := s.client.do(ctx, "POST", "/api/v1/qa", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
