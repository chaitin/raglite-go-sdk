package sdk

import (
	"context"
	"fmt"
)

// DatasetsService 数据集管理服务
type DatasetsService struct {
	client *Client
}

// CreateDatasetRequest 创建数据集请求
type CreateDatasetRequest struct {
	Name            string        `json:"name"`
	Description     string        `json:"description,omitempty"`
	DenseModelID    *string       `json:"dense_model_id,omitempty"`
	SparseModelID   *string       `json:"sparse_model_id,omitempty"`
	AnalysisModelID *string       `json:"analysis_model_id,omitempty"`
	RerankerModelID *string       `json:"reranker_model_id,omitempty"`
	VisionModelID   *string       `json:"vision_model_id,omitempty"`
	Config          DatasetConfig `json:"config,omitempty"`
}

// UpdateDatasetRequest 更新数据集请求
type UpdateDatasetRequest struct {
	Name            *string        `json:"name,omitempty"`
	Description     *string        `json:"description,omitempty"`
	DenseModelID    *string        `json:"dense_model_id,omitempty"`
	SparseModelID   *string        `json:"sparse_model_id,omitempty"`
	AnalysisModelID *string        `json:"analysis_model_id,omitempty"`
	RerankerModelID *string        `json:"reranker_model_id,omitempty"`
	VisionModelID   *string        `json:"vision_model_id,omitempty"`
	Config          *DatasetConfig `json:"config,omitempty"`
	Status          *string        `json:"status,omitempty"`
}

// ListDatasetsRequest 列表查询请求
type ListDatasetsRequest struct {
	Status string
}

// ListDatasetsResponse 列表响应
type ListDatasetsResponse struct {
	Datasets []Dataset `json:"datasets"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
}

// Create 创建数据集
func (s *DatasetsService) Create(ctx context.Context, req *CreateDatasetRequest) (*Dataset, error) {
	var result Dataset
	err := s.client.do(ctx, "POST", "/api/v1/datasets", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// List 列出数据集
func (s *DatasetsService) List(ctx context.Context, req *ListDatasetsRequest) (*ListDatasetsResponse, error) {
	params := make(map[string]string)
	if req != nil && req.Status != "" {
		params["status"] = req.Status
	}

	path := s.client.buildURL("/api/v1/datasets", params)
	var result ListDatasetsResponse
	err := s.client.do(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Get 获取数据集详情
func (s *DatasetsService) Get(ctx context.Context, datasetID string) (*Dataset, error) {
	var result Dataset
	path := fmt.Sprintf("/api/v1/datasets/%s", datasetID)
	err := s.client.do(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Update 更新数据集
func (s *DatasetsService) Update(ctx context.Context, datasetID string, req *UpdateDatasetRequest) (*Dataset, error) {
	var result Dataset
	path := fmt.Sprintf("/api/v1/datasets/%s", datasetID)
	err := s.client.do(ctx, "PUT", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete 删除数据集
func (s *DatasetsService) Delete(ctx context.Context, datasetID string) error {
	path := fmt.Sprintf("/api/v1/datasets/%s", datasetID)
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

// GetStats 获取数据集统计信息
func (s *DatasetsService) GetStats(ctx context.Context, datasetID string) (*DatasetStats, error) {
	var result DatasetStats
	path := fmt.Sprintf("/api/v1/datasets/%s/stats", datasetID)
	err := s.client.do(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
