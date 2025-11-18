package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

// DocumentsService 文档管理服务
type DocumentsService struct {
	client *Client
}

// UploadDocumentRequest 上传文档请求
type UploadDocumentRequest struct {
	DatasetID  string
	DocumentID string
	File       io.Reader
	Filename   string
	Tags       []string
	Metadata   map[string]interface{}
}

// UploadDocumentResponse 上传文档响应
type UploadDocumentResponse struct {
	DocumentID string `json:"document_id"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	Filename   string `json:"filename"`
	Title      string `json:"title"`
	Size       int64  `json:"size"`
}

// ListDocumentsRequest 列表查询请求
type ListDocumentsRequest struct {
	DatasetID string
}

// ListDocumentsResponse 列表响应
type ListDocumentsResponse struct {
	Documents []Document `json:"documents"`
	Total     int64      `json:"total"`
	Page      int        `json:"page"`
	PageSize  int        `json:"page_size"`
}

// BatchDeleteDocumentsRequest 批量删除文档请求
type BatchDeleteDocumentsRequest struct {
	DatasetID   string   `json:"-"`
	DocumentIDs []string `json:"document_ids"`
}

// Upload 上传文档
func (s *DocumentsService) Upload(ctx context.Context, req *UploadDocumentRequest) (*UploadDocumentResponse, error) {
	// 创建 multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件
	part, err := writer.CreateFormFile("file", req.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, req.File); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	// 添加 tags
	if len(req.Tags) > 0 {
		tagsJSON, err := json.Marshal(req.Tags)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal tags: %w", err)
		}
		if err := writer.WriteField("tags", string(tagsJSON)); err != nil {
			return nil, fmt.Errorf("failed to write tags field: %w", err)
		}
	}

	// 添加 metadata
	if req.Metadata != nil {
		metadataJSON, err := json.Marshal(req.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata: %w", err)
		}
		if err := writer.WriteField("metadata", string(metadataJSON)); err != nil {
			return nil, fmt.Errorf("failed to write metadata field: %w", err)
		}
	}

	// 添加可选的 document_id（用于更新）
	if req.DocumentID != "" {
		if err := writer.WriteField("document_id", req.DocumentID); err != nil {
			return nil, fmt.Errorf("failed to write document_id field: %w", err)
		}
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// 发送请求
	path := fmt.Sprintf("/api/v1/datasets/%s/documents", req.DatasetID)
	fullURL := s.client.baseURL + path

	httpReq, err := http.NewRequestWithContext(ctx, "POST", fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.client.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 检查状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr APIError
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			return nil, &APIError{
				StatusCode: resp.StatusCode,
				Message:    string(respBody),
			}
		}
		apiErr.StatusCode = resp.StatusCode
		return nil, &apiErr
	}

	// 解析响应
	var apiResp struct {
		Success bool                   `json:"success"`
		Message string                 `json:"message"`
		Data    UploadDocumentResponse `json:"data"`
	}
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &apiResp.Data, nil
}

// List 列出文档
func (s *DocumentsService) List(ctx context.Context, req *ListDocumentsRequest) (*ListDocumentsResponse, error) {
	var result ListDocumentsResponse
	path := fmt.Sprintf("/api/v1/datasets/%s/documents", req.DatasetID)
	err := s.client.do(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Get 获取文档详情
func (s *DocumentsService) Get(ctx context.Context, datasetID, documentID string) (*Document, error) {
	var result Document
	path := fmt.Sprintf("/api/v1/datasets/%s/documents/%s", datasetID, documentID)
	err := s.client.do(ctx, "GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete 删除文档
func (s *DocumentsService) Delete(ctx context.Context, datasetID, documentID string) error {
	path := fmt.Sprintf("/api/v1/datasets/%s/documents/%s", datasetID, documentID)
	return s.client.do(ctx, "DELETE", path, nil, nil)
}

// BatchDelete 批量删除文档
func (s *DocumentsService) BatchDelete(ctx context.Context, req *BatchDeleteDocumentsRequest) error {
	if len(req.DocumentIDs) == 0 {
		return nil
	}

	path := fmt.Sprintf("/api/v1/datasets/%s/documents/batch-delete", req.DatasetID)

	body := map[string]interface{}{
		"document_ids": req.DocumentIDs,
	}

	return s.client.do(ctx, "POST", path, body, nil)
}

// UpdateDocumentRequest 更新文档请求
type UpdateDocumentRequest struct {
	DatasetID  string                 `json:"-"`
	DocumentID string                 `json:"-"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Tags       []string               `json:"tags,omitempty"`
}

// Update 更新文档的 metadata 和 tags
func (s *DocumentsService) Update(ctx context.Context, req *UpdateDocumentRequest) (*Document, error) {
	var result Document
	path := fmt.Sprintf("/api/v1/datasets/%s/documents/%s", req.DatasetID, req.DocumentID)

	body := make(map[string]interface{})
	if req.Metadata != nil {
		body["metadata"] = req.Metadata
	}
	if req.Tags != nil {
		body["tags"] = req.Tags
	}

	err := s.client.do(ctx, "PATCH", path, body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
