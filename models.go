package sdk

import (
	"context"
	"fmt"
)

// ModelsService AI 模型管理服务
type ModelsService struct {
	client *Client
}

// CreateRequest 创建模型请求
type CreateModelRequest struct {
	Name         string            `json:"name"`
	Description  string            `json:"description,omitempty"`
	ModelType    string            `json:"model_type"`
	Provider     string            `json:"provider"`
	ModelName    string            `json:"model_name"`
	Config       AIModelConfig     `json:"config,omitempty"`
	Capabilities ModelCapabilities `json:"capabilities,omitempty"`
	IsDefault    bool              `json:"is_default,omitempty"`
}

// UpdateModelRequest 更新模型请求
type UpdateModelRequest struct {
	Name         *string            `json:"name,omitempty"`
	Description  *string            `json:"description,omitempty"`
	Provider     *string            `json:"provider,omitempty"`
	ModelName    *string            `json:"model_name,omitempty"`
	Config       *AIModelConfig     `json:"config,omitempty"`
	Capabilities *ModelCapabilities `json:"capabilities,omitempty"`
	Status       *string            `json:"status,omitempty"`
	IsDefault    *bool              `json:"is_default,omitempty"`
	IsActive     *bool              `json:"is_active,omitempty"`
}

// ListModelsRequest 列表查询请求
type ListModelsRequest struct {
	ModelType string
	Provider  string
	Status    string
}

// ListModelsResponse 列表响应
type ListModelsResponse struct {
	Models []AIModel `json:"models"`
	Total  int64     `json:"total"`
}

// ListProviderModelsRequest 获取供应商支持的模型列表请求
type ListProviderModelsRequest struct {
	Provider string                 `json:"provider"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

// CheckModelRequest 检查模型配置请求
type CheckModelRequest struct {
	Provider  string        `json:"provider"`
	ModelName string        `json:"model_name"`
	Config    AIModelConfig `json:"config"`
}

// UpsertModelRequest 根据配置创建或更新模型请求
type UpsertModelRequest struct {
	Name         string            `json:"name,omitempty"`
	Description  string            `json:"description,omitempty"`
	ModelType    string            `json:"model_type"`
	Provider     string            `json:"provider"`
	ModelName    string            `json:"model_name"`
	Config       AIModelConfig     `json:"config"`
	Capabilities ModelCapabilities `json:"capabilities,omitempty"`
	IsDefault    bool              `json:"is_default,omitempty"`
	IsActive     bool              `json:"is_active,omitempty"`
}

// UpsertModelResponse Upsert 响应
type UpsertModelResponse struct {
	Action string  `json:"action"` // "created" 或 "updated"
	Model  AIModel `json:"model"`
}

// Create 创建 AI 模型
func (s *ModelsService) Create(ctx context.Context, req *CreateModelRequest) (*AIModel, error) {
	var result AIModel
	err := s.client.do(ctx, "POST", "/api/v1/models", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// List 列出 AI 模型
func (s *ModelsService) List(ctx context.Context, req *ListModelsRequest) (*ListModelsResponse, error) {
	params := make(map[string]string)
	if req != nil {
		if req.ModelType != "" {
			params["model_type"] = req.ModelType
		}
		if req.Provider != "" {
			params["provider"] = req.Provider
		}
		if req.Status != "" {
			params["status"] = req.Status
		}
	}

	path := s.client.buildURL("/api/v1/models", params)
	var result ListModelsResponse
	err := s.client.do(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Get 获取模型详情
func (s *ModelsService) Get(ctx context.Context, modelID string) (*AIModel, error) {
	var result AIModel
	path := fmt.Sprintf("/api/v1/models/%s", modelID)
	err := s.client.do(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Update 更新模型
func (s *ModelsService) Update(ctx context.Context, modelID string, req *UpdateModelRequest) (*AIModel, error) {
	var result AIModel
	path := fmt.Sprintf("/api/v1/models/%s", modelID)
	err := s.client.do(ctx, "PUT", path, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete 删除模型
func (s *ModelsService) Delete(ctx context.Context, modelID string) error {
	path := fmt.Sprintf("/api/v1/models/%s", modelID)
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

// ListProviderModels 获取供应商支持的模型列表
func (s *ModelsService) ListProviderModels(ctx context.Context, req *ListProviderModelsRequest) (interface{}, error) {
	var result interface{}
	err := s.client.do(ctx, "POST", "/api/v1/models/provider/supported", req, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Check 检查模型配置
func (s *ModelsService) Check(ctx context.Context, req *CheckModelRequest) (*CheckModelResponse, error) {
	var result CheckModelResponse
	err := s.client.do(ctx, "POST", "/api/v1/models/check", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Upsert 根据 API Base 和 Model Name 创建或更新模型
// 如果找到匹配的模型则更新，否则创建新模型
func (s *ModelsService) Upsert(ctx context.Context, req *UpsertModelRequest) (*UpsertModelResponse, error) {
	var result UpsertModelResponse
	err := s.client.do(ctx, "POST", "/api/v1/models/upsert", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
