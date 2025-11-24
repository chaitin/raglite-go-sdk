package sdk

import (
	"encoding/json"
	"time"
)

// Ptr 返回值的指针，方便创建指针字段
func Ptr[T any](v T) *T {
	return &v
}

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
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	ModelType    string            `json:"model_type"`
	Provider     string            `json:"provider"`
	ModelName    string            `json:"model_name"`
	Config       AIModelConfig     `json:"config"`
	Capabilities ModelCapabilities `json:"capabilities"`
	Status       string            `json:"status"`
	IsDefault    bool              `json:"is_default"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// AIModelConfig AI 模型配置
type AIModelConfig struct {
	// API 配置
	APIKey     string `json:"api_key,omitempty"`
	APIBase    string `json:"api_base,omitempty"`
	APIVersion string `json:"api_version,omitempty"`
	OrgID      string `json:"org_id,omitempty"`

	// 模型参数（Chat Model 专用）
	Temperature      *float64 `json:"temperature,omitempty"`
	MaxTokens        *int     `json:"max_tokens,omitempty"`
	TopP             *float64 `json:"top_p,omitempty"`
	FrequencyPenalty *float64 `json:"frequency_penalty,omitempty"`
	PresencePenalty  *float64 `json:"presence_penalty,omitempty"`

	// 其他配置
	Timeout         int                    `json:"timeout,omitempty"`
	MaxRetries      int                    `json:"max_retries,omitempty"`
	CustomHeaders   map[string]string      `json:"custom_headers,omitempty"`
	ExtraParameters map[string]interface{} `json:"extra_parameters,omitempty"`
}

// ModelCapabilities 模型能力参数
type ModelCapabilities struct {
	// Embedding 模型能力
	VectorDimension *int  `json:"vector_dimension,omitempty"`
	MaxInputTokens  *int  `json:"max_input_tokens,omitempty"`
	SupportsBatch   *bool `json:"supports_batch,omitempty"`

	// Chat Model 能力
	SupportsStreaming *bool `json:"supports_streaming,omitempty"`
	SupportsFunctions *bool `json:"supports_functions,omitempty"`
	ContextWindow     *int  `json:"context_window,omitempty"`
	MaxOutputTokens   *int  `json:"max_output_tokens,omitempty"`

	// 通用能力
	CostPer1KTokens *float64 `json:"cost_per_1k_tokens,omitempty"`
}

// Dataset 数据集
type Dataset struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	DenseModelID    string        `json:"dense_model_id"`
	SparseModelID   *string       `json:"sparse_model_id,omitempty"`
	AnalysisModelID *string       `json:"analysis_model_id,omitempty"`
	RerankerModelID *string       `json:"reranker_model_id,omitempty"`
	VisionModelID   *string       `json:"vision_model_id,omitempty"`
	Config          DatasetConfig `json:"config"`
	Status          string        `json:"status"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
}

// DatasetConfig 数据集配置
type DatasetConfig struct {
	ChunkSize    int                    `json:"chunk_size"`             // 块大小（tokens）
	ChunkOverlap int                    `json:"chunk_overlap"`          // 块重叠（tokens）
	IndexParams  map[string]interface{} `json:"index_params,omitempty"` // 索引配置（可选）
}

// Document 文档
type Document struct {
	ID          string                 `json:"id"`
	DatasetID   string                 `json:"dataset_id"`
	Title       string                 `json:"title"`
	Filename    string                 `json:"filename"`
	FilePath    string                 `json:"file_path"`
	FileHash    string                 `json:"file_hash"`
	FileSize    int64                  `json:"file_size"`
	Metadata    map[string]interface{} `json:"metadata"`
	Tags        []string               `json:"tags"`
	Status      string                 `json:"status"`
	ProgressMsg string                 `json:"progress_msg,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
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
