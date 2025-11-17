package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/raglite/raglite/sdk"
)

func main() {
	// 创建客户端
	client, err := sdk.NewClient("http://localhost:8080")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// 示例 1: 健康检查
	fmt.Println("=== 健康检查 ===")
	health, err := client.Health.Check(ctx)
	if err != nil {
		log.Printf("Health check failed: %v", err)
	} else {
		fmt.Printf("Service: %s, Status: %s\n\n", health.Service, health.Status)
	}

	// 示例 2: 创建 AI 模型
	fmt.Println("=== 创建 AI 模型 ===")
	model, err := client.Models.Create(ctx, &sdk.CreateModelRequest{
		Name:      "OpenAI GPT-4",
		ModelType: "chat",
		Provider:  "openai",
		ModelName: "gpt-4",
		Config: map[string]interface{}{
			"api_key":     os.Getenv("OPENAI_API_KEY"),
			"temperature": 0.7,
		},
		IsDefault: true,
	})
	if err != nil {
		log.Printf("Failed to create model: %v", err)
	} else {
		fmt.Printf("Model created: %s (ID: %s)\n\n", model.Name, model.ID)
	}

	// 示例 3: 列出所有模型
	fmt.Println("=== 列出所有模型 ===")
	models, err := client.Models.List(ctx, &sdk.ListModelsRequest{
		ModelType: "chat",
	})
	if err != nil {
		log.Printf("Failed to list models: %v", err)
	} else {
		fmt.Printf("Found %d models\n", models.Total)
		for _, m := range models.Models {
			fmt.Printf("  - %s (%s/%s) [%s]\n", m.Name, m.Provider, m.ModelName, m.Status)
		}
		fmt.Println()
	}

	// 示例 4: 创建数据集
	fmt.Println("=== 创建数据集 ===")
	dataset, err := client.Datasets.Create(ctx, &sdk.CreateDatasetRequest{
		Name:        "技术文档",
		Description: "公司技术文档知识库",
		Config: map[string]interface{}{
			"chunk_size":    512,
			"chunk_overlap": 50,
		},
	})
	if err != nil {
		log.Printf("Failed to create dataset: %v", err)
		return
	}
	fmt.Printf("Dataset created: %s (ID: %s)\n\n", dataset.Name, dataset.ID)

	// 示例 5: 上传文档
	fmt.Println("=== 上传文档 ===")
	content := `# RAGLite 使用指南

## 简介

RAGLite 是一个轻量级的 RAG (Retrieval-Augmented Generation) 系统。

## 功能特性

- 文档管理
- 向量检索
- 智能问答
`
	uploadResp, err := client.Documents.Upload(ctx, &sdk.UploadDocumentRequest{
		DatasetID: dataset.ID,
		File:      strings.NewReader(content),
		Filename:  "guide.md",
		Tags:      []string{"文档", "指南"},
		Metadata: map[string]interface{}{
			"author":  "RAGLite Team",
			"version": "1.0",
		},
	})
	if err != nil {
		log.Printf("Failed to upload document: %v", err)
	} else {
		fmt.Printf("Document uploaded: %s (Status: %s)\n\n", uploadResp.Filename, uploadResp.Status)
	}

	// 示例 6: 搜索
	fmt.Println("=== 搜索文档 ===")
	searchResp, err := client.Search.Search(ctx, &sdk.SearchRequest{
		Query:     "RAGLite 有哪些功能",
		DatasetID: dataset.ID,
		TopK:      5,
	})
	if err != nil {
		log.Printf("Failed to search: %v", err)
	} else {
		fmt.Printf("Found %d results (in %dms)\n", searchResp.Total, searchResp.LatencyMs)
		for i, result := range searchResp.Results {
			fmt.Printf("  %d. %s (Score: %.3f)\n", i+1, result.DocumentTitle, result.Score)
			fmt.Printf("     %s\n", truncate(result.Content, 100))
		}
		fmt.Println()
	}

	// 示例 7: 问答
	fmt.Println("=== 智能问答 ===")
	qaResp, err := client.QA.Ask(ctx, &sdk.QARequest{
		Query:     "RAGLite 的主要功能是什么？",
		DatasetID: dataset.ID,
		TopK:      3,
	})
	if err != nil {
		log.Printf("Failed to ask question: %v", err)
	} else {
		fmt.Printf("答案: %s\n", qaResp.Answer)
		fmt.Printf("参考了 %d 个上下文\n\n", len(qaResp.Context))
	}

	// 示例 8: 获取数据集统计
	fmt.Println("=== 数据集统计 ===")
	stats, err := client.Datasets.GetStats(ctx, dataset.ID)
	if err != nil {
		log.Printf("Failed to get stats: %v", err)
	} else {
		fmt.Printf("总文档数: %d\n", stats.TotalDocuments)
		fmt.Printf("  - 待处理: %d\n", stats.PendingDocs)
		fmt.Printf("  - 处理中: %d\n", stats.ProcessingDocs)
		fmt.Printf("  - 已完成: %d\n", stats.CompletedDocs)
		fmt.Printf("  - 失败: %d\n", stats.FailedDocs)
		fmt.Printf("总文件大小: %d bytes\n", stats.TotalFileSize)
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
