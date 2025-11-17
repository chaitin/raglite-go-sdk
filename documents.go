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
	DatasetID string
	File      io.Reader
	Filename  string
	Tags      []string
	Metadata  map[string]interface{}
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
