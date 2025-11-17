package sdk

import (
	"encoding/json"
	"time"
)

// APIResponse 统一的 API 响应结构
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ListResponse 列表响应
type ListResponse struct {
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

// AIModel AI 模型
type AIModel struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	ModelType    string                 `json:"model_type"`
	Provider     string                 `json:"provider"`
	ModelName    string                 `json:"model_name"`
	Config       map[string]interface{} `json:"config"`
	Capabilities map[string]interface{} `json:"capabilities"`
	Status       string                 `json:"status"`
	IsDefault    bool                   `json:"is_default"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// Dataset 数据集
type Dataset struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	DenseModelID    string                 `json:"dense_model_id"`
	SparseModelID   *string                `json:"sparse_model_id,omitempty"`
	AnalysisModelID *string                `json:"analysis_model_id,omitempty"`
	RerankerModelID *string                `json:"reranker_model_id,omitempty"`
	VisionModelID   *string                `json:"vision_model_id,omitempty"`
	Config          map[string]interface{} `json:"config"`
	Status          string                 `json:"status"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// Document 文档
type Document struct {
	ID           string                 `json:"id"`
	DatasetID    string                 `json:"dataset_id"`
	Title        string                 `json:"title"`
	Filename     string                 `json:"filename"`
	FilePath     string                 `json:"file_path"`
	FileHash     string                 `json:"file_hash"`
	FileSize     int64                  `json:"file_size"`
	Metadata     map[string]interface{} `json:"metadata"`
	Tags         []string               `json:"tags"`
	Status       string                 `json:"status"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// SearchResult 搜索结果
type SearchResult struct {
	ChunkID       string                 `json:"chunk_id"`
	DocumentID    string                 `json:"document_id"`
	DocumentTitle string                 `json:"document_title"`
	SectionTitle  string                 `json:"section_title"`
	Content       string                 `json:"content"`
	Score         float64                `json:"score"`
	Metadata      map[string]interface{} `json:"metadata"`
	Tags          []string               `json:"tags"`
	Highlights    []string               `json:"highlights,omitempty"`
}

// DatasetStats 数据集统计
type DatasetStats struct {
	Dataset        Dataset `json:"dataset"`
	TotalDocuments int64   `json:"total_documents"`
	PendingDocs    int64   `json:"pending_docs"`
	ProcessingDocs int64   `json:"processing_docs"`
	CompletedDocs  int64   `json:"completed_docs"`
	FailedDocs     int64   `json:"failed_docs"`
	TotalFileSize  int64   `json:"total_file_size"`
}

// ProviderModel 供应商模型
type ProviderModel struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Capabilities []string `json:"capabilities"`
}

// CheckModelResponse 检查模型响应
type CheckModelResponse struct {
	Valid     bool        `json:"valid"`
	Error     string      `json:"error,omitempty"`
	ModelInfo interface{} `json:"model_info,omitempty"`
}

// JSON 辅助类型，用于处理 JSON 字段
type JSON struct {
	Data interface{}
}

// MarshalJSON 实现 json.Marshaler
func (j JSON) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Data)
}

// UnmarshalJSON 实现 json.Unmarshaler
func (j *JSON) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &j.Data)
}
