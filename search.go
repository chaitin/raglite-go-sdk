package sdk

import "context"

// SearchService 搜索服务
type SearchService struct {
	client *Client
}

// RetrieveRequest 召回请求
type RetrieveRequest struct {
	Query               string                 `json:"query"`
	DatasetID           string                 `json:"dataset_id"`
	TopK                int                    `json:"top_k,omitempty"`
	RetrievalMode       string                 `json:"retrieval_mode,omitempty"` // full | smart
	SimilarityThreshold float64                `json:"similarity_threshold,omitempty"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
	Tags                []string               `json:"tags,omitempty"`
	ChatHistory         []ChatMessage          `json:"chat_history,omitempty"`
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Query     string         `json:"query"`
	Results   []SearchResult `json:"results"`
	Total     int            `json:"total"`
	LatencyMs int64          `json:"latency_ms"`
}

// Search 执行搜索
func (s *SearchService) Retrieve(ctx context.Context, req *RetrieveRequest) (*SearchResponse, error) {
	// 设置默认值
	if req.TopK <= 0 {
		req.TopK = 10
	}

	var result SearchResponse
	err := s.client.do(ctx, "POST", "/api/v1/search", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
