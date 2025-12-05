package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	sdk "github.com/chaitin/raglite-go-sdk"
)

func main() {
	// 创建带自定义配置的客户端
	client, err := sdk.NewClient(
		"http://localhost:5050",
		sdk.WithTimeout(60*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// 高级示例 1: 使用 Upsert 创建或更新模型
	fmt.Println("=== 使用 Upsert 创建或更新模型 ===")
	upsertResp, err := client.Models.Upsert(ctx, &sdk.UpsertModelRequest{
		Name:      "GPT-4",
		ModelType: "chat",
		Provider:  "openai",
		ModelName: "gpt-4",
		Config: sdk.AIModelConfig{
			APIKey:  os.Getenv("OPENAI_API_KEY"),
			APIBase: "https://api.openai.com/v1",
		},
		Capabilities: sdk.ModelCapabilities{
			ContextWindow: sdk.Ptr(8192),
		},
	})
	if err != nil {
		log.Printf("Failed to upsert model: %v", err)
	} else {
		fmt.Printf("Model %s: %s (ID: %s)\n", upsertResp.Action, upsertResp.Model.Name, upsertResp.Model.ID)
	}
	fmt.Println()

	// 高级示例 2: 检查模型配置
	fmt.Println("=== 检查模型配置 ===")
	checkResp, err := client.Models.Check(ctx, &sdk.CheckModelRequest{
		Provider:  "openai",
		ModelName: "gpt-4",
		Config: sdk.AIModelConfig{
			APIKey: os.Getenv("OPENAI_API_KEY"),
		},
	})
	if err != nil {
		log.Printf("Failed to check model: %v", err)
	} else {
		if checkResp.Valid {
			fmt.Println("Model configuration is valid")
		} else {
			fmt.Printf("Model configuration is invalid: %s\n", checkResp.Error)
		}
	}
	fmt.Println()

	// 高级示例 3: 创建多模型数据集
	fmt.Println("=== 创建多模型数据集 ===")

	// 先创建必要的模型
	denseModel, err := createOrGetModel(ctx, client, "dense", "text-embedding-3-small")
	if err != nil {
		log.Fatalf("Failed to get dense model: %v", err)
	}

	sparseModel, err := createOrGetModel(ctx, client, "sparse", "bm25")
	if err != nil {
		log.Printf("Failed to get sparse model: %v", err)
	}

	rerankerModel, err := createOrGetModel(ctx, client, "reranker", "bge-reranker-v2-m3")
	if err != nil {
		log.Printf("Failed to get reranker model: %v", err)
	}

	// 创建数据集，使用多个模型
	dataset, err := client.Datasets.Create(ctx, &sdk.CreateDatasetRequest{
		Name:            "高级知识库",
		Description:     "使用混合检索和重排序的高级知识库",
		DenseModelID:    &denseModel.ID,
		SparseModelID:   sdk.Ptr(sparseModel.ID),
		RerankerModelID: sdk.Ptr(rerankerModel.ID),
		Config: sdk.DatasetConfig{
			ChunkSize:    512,
			ChunkOverlap: 100,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create dataset: %v", err)
	}
	fmt.Printf("Dataset created: %s (ID: %s)\n\n", dataset.Name, dataset.ID)

	// 高级示例 4: 带过滤的搜索
	fmt.Println("=== 带过滤的高级搜索 ===")
	searchResp, err := client.Search.Retrieve(ctx, &sdk.RetrieveRequest{
		Query:               "人工智能",
		DatasetID:           dataset.ID,
		TopK:                10,
		RetrievalMode:       "smart", // 使用智能检索模式
		SimilarityThreshold: 0.7,     // 只返回相似度 > 0.7 的结果
		Tags:                []string{"技术", "AI"},
		Metadata: map[string]interface{}{
			"category": "research",
		},
	})
	if err != nil {
		log.Printf("Failed to search: %v", err)
	} else {
		fmt.Printf("Found %d results\n", len(searchResp.Results))
		for i, result := range searchResp.Results {
			fmt.Printf("%d. [Score: %.3f] %s\n", i+1, result.Score, result.DocumentTitle)
		}
	}
	fmt.Println()

	// 高级示例 5: 批量操作
	fmt.Println("=== 批量上传文档 ===")
	documents := []struct {
		filename string
		content  string
	}{
		{"doc1.md", "# 文档 1\n这是第一个文档"},
		{"doc2.md", "# 文档 2\n这是第二个文档"},
		{"doc3.md", "# 文档 3\n这是第三个文档"},
	}

	for _, doc := range documents {
		// 这里可以使用 goroutine 并发上传
		go func(filename, content string) {
			_, err := client.Documents.Upload(ctx, &sdk.UploadDocumentRequest{
				DatasetID: dataset.ID,
				File:      nil, // 实际使用时应该传入文件内容
				Filename:  filename,
			})
			if err != nil {
				log.Printf("Failed to upload %s: %v", filename, err)
			} else {
				fmt.Printf("Uploaded: %s\n", filename)
			}
		}(doc.filename, doc.content)
	}

	// 等待上传完成
	time.Sleep(2 * time.Second)
	fmt.Println()

	// 高级示例 6: 错误处理
	fmt.Println("=== 错误处理示例 ===")
	_, err = client.Datasets.Get(ctx, "non-existent-id")
	if err != nil {
		if apiErr, ok := err.(*sdk.APIError); ok {
			if apiErr.IsNotFound() {
				fmt.Println("Dataset not found (404)")
			} else if apiErr.IsBadRequest() {
				fmt.Println("Bad request (400)")
			} else if apiErr.IsServerError() {
				fmt.Println("Server error (5xx)")
			} else {
				fmt.Printf("API error: %v\n", apiErr)
			}
		} else {
			fmt.Printf("Other error: %v\n", err)
		}
	}
	fmt.Println()

	// 高级示例 7: 更新模型配置
	fmt.Println("=== 更新数据集配置 ===")
	updatedDataset, err := client.Datasets.Update(ctx, dataset.ID, &sdk.UpdateDatasetRequest{
		Status: sdk.Ptr("active"),
		Config: sdk.Ptr(sdk.DatasetConfig{
			ChunkSize:    1024,
			ChunkOverlap: 100,
		}),
	})
	if err != nil {
		log.Printf("Failed to update dataset: %v", err)
	} else {
		fmt.Printf("Dataset updated: %s\n", updatedDataset.Name)
	}
}

func createOrGetModel(ctx context.Context, client *sdk.Client, modelType, modelName string) (*sdk.AIModel, error) {
	// 使用 Upsert 方法自动处理创建或更新
	provider := "openai"
	apiBase := "https://api.openai.com/v1"
	if modelType == "sparse" {
		provider = "local"
		apiBase = "http://localhost:9200"
	} else if modelType == "reranker" {
		provider = "local"
		apiBase = "http://localhost:8001"
	}

	result, err := client.Models.Upsert(ctx, &sdk.UpsertModelRequest{
		Name:      fmt.Sprintf("%s-%s", modelType, modelName),
		ModelType: modelType,
		Provider:  provider,
		ModelName: modelName,
		Config: sdk.AIModelConfig{
			APIKey:  os.Getenv("OPENAI_API_KEY"),
			APIBase: apiBase,
		},
	})
	if err != nil {
		return nil, err
	}

	return &result.Model, nil
}
