# RAGLite Go SDK

RAGLite 的官方 Go SDK，提供简洁、类型安全的 API 来访问 RAGLite 服务。

## 特性

- 🎯 **简洁易用** - 流畅的 API 设计，符合 Go 语言习惯
- 🔒 **类型安全** - 完整的类型定义，编译时错误检查
- 🧩 **模块化** - 按服务分离的清晰架构
- 🚀 **高性能** - 支持并发操作，HTTP 连接复用
- 🛠️ **可扩展** - 灵活的配置选项，易于定制
- 📝 **完整文档** - 详细的代码注释和示例

## 安装

```bash
go get github.com/chaitin/raglite-go-sdk
```

## 快速开始

```go
package main

import (
    "context"
    "fmt"
    "log"

    sdk "github.com/chaitin/raglite-go-sdk"
)

func main() {
    client, err := sdk.NewClient(
        "http://localhost:8080",
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // 健康检查
    health, err := client.Health.Check(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Service: %s, Status: %s\n", health.Service, health.Status)

    // 创建数据集
    dataset, err := client.Datasets.Create(ctx, &sdk.CreateDatasetRequest{
        Name:        "我的知识库",
        Description: "技术文档知识库",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Dataset created: %s\n", dataset.Name)

    // 搜索
    results, err := client.Search.Search(ctx, &sdk.SearchRequest{
        Query:     "如何使用 RAGLite",
        DatasetID: dataset.ID,
        TopK:      10,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Found %d results\n", len(results.Results))

    // 问答
    answer, err := client.QA.Ask(ctx, &sdk.QARequest{
        Query:     "RAGLite 的主要功能是什么？",
        DatasetID: dataset.ID,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Answer: %s\n", answer.Answer)
}
```

## 核心功能

### 1. 客户端配置

```go
import "time"

// 基础配置
client, _ := sdk.NewClient("http://localhost:8080")

// 使用 API Key 认证
client, _ := sdk.NewClient(
    "http://localhost:8080",
    sdk.WithAPIKey("your-api-key-here"),
)

// 自定义超时
client, _ := sdk.NewClient(
    "http://localhost:8080",
    sdk.WithTimeout(60 * time.Second),
)

// 组合使用多个选项
client, _ := sdk.NewClient(
    "http://localhost:8080",
    sdk.WithAPIKey("your-api-key-here"),
    sdk.WithTimeout(60 * time.Second),
)

// 自定义 HTTP 客户端
httpClient := &http.Client{
    Timeout: 30 * time.Second,
    Transport: customTransport,
}
client, _ := sdk.NewClient(
    "http://localhost:8080",
    sdk.WithHTTPClient(httpClient),
)
```

### 2. AI 模型管理

```go
// 创建模型
model, err := client.Models.Create(ctx, &sdk.CreateModelRequest{
    Name:      "GPT-4",
    ModelType: "chat",
    Provider:  "openai",
    ModelName: "gpt-4",
    Config: map[string]interface{}{
        "api_key":     "your-api-key",
        "temperature": 0.7,
    },
    IsDefault: true,
})

// 列出模型
models, err := client.Models.List(ctx, &sdk.ListModelsRequest{
    ModelType: "chat",
    Provider:  "openai",
})

// 获取模型详情
model, err := client.Models.Get(ctx, modelID)

// 更新模型
newName := "GPT-4 Turbo"
model, err := client.Models.Update(ctx, modelID, &sdk.UpdateModelRequest{
    Name: &newName,
})

// 删除模型
err := client.Models.Delete(ctx, modelID)

// 检查模型配置
check, err := client.Models.Check(ctx, &sdk.CheckModelRequest{
    Provider:  "openai",
    ModelName: "gpt-4",
    Config: map[string]interface{}{
        "api_key": "your-api-key",
    },
})
```

### 3. 数据集管理

```go
// 创建数据集
dataset, err := client.Datasets.Create(ctx, &sdk.CreateDatasetRequest{
    Name:        "技术文档",
    Description: "公司技术文档知识库",
    Config: map[string]interface{}{
        "chunk_size":    512,
        "chunk_overlap": 50,
    },
})

// 列出数据集
datasets, err := client.Datasets.List(ctx, &sdk.ListDatasetsRequest{
    Status: "active",
})

// 获取数据集详情
dataset, err := client.Datasets.Get(ctx, datasetID)

// 更新数据集
dataset, err := client.Datasets.Update(ctx, datasetID, &sdk.UpdateDatasetRequest{
    Description: ptrString("更新后的描述"),
})

// 获取统计信息
stats, err := client.Datasets.GetStats(ctx, datasetID)
fmt.Printf("总文档数: %d, 已完成: %d\n", stats.TotalDocuments, stats.CompletedDocs)

// 删除数据集
err := client.Datasets.Delete(ctx, datasetID)
```

### 4. 文档管理

```go
import (
    "os"
    "strings"
)

// 上传文档（从文件）
file, _ := os.Open("document.md")
defer file.Close()

resp, err := client.Documents.Upload(ctx, &sdk.UploadDocumentRequest{
    DatasetID: datasetID,
    File:      file,
    Filename:  "document.md",
    Tags:      []string{"技术", "文档"},
    Metadata: map[string]interface{}{
        "author":  "John Doe",
        "version": "1.0",
    },
})

// 上传文档（从字符串）
content := "# 标题\n文档内容..."
resp, err := client.Documents.Upload(ctx, &sdk.UploadDocumentRequest{
    DatasetID: datasetID,
    File:      strings.NewReader(content),
    Filename:  "doc.md",
})

// 列出文档
docs, err := client.Documents.List(ctx, &sdk.ListDocumentsRequest{
    DatasetID: datasetID,
})

// 获取文档详情
doc, err := client.Documents.Get(ctx, datasetID, documentID)

// 删除文档
err := client.Documents.Delete(ctx, datasetID, documentID)

// 批量删除文档
err := client.Documents.BatchDelete(ctx, &sdk.BatchDeleteDocumentsRequest{
    DatasetID:   datasetID,
    DocumentIDs: []string{docID1, docID2, docID3},
})
```

### 5. 搜索

```go
// 基础搜索
results, err := client.Search.Search(ctx, &sdk.SearchRequest{
    Query:     "RAGLite 是什么",
    DatasetID: datasetID,
    TopK:      10,
})

// 高级搜索（带过滤）
results, err := client.Search.Search(ctx, &sdk.SearchRequest{
    Query:               "机器学习算法",
    DatasetID:           datasetID,
    TopK:                20,
    RetrievalMode:       "smart",  // full | smart
    SimilarityThreshold: 0.7,
    Tags:                []string{"AI", "算法"},
    Metadata: map[string]interface{}{
        "category": "research",
    },
})

// 处理搜索结果
for i, result := range results.Results {
    fmt.Printf("%d. [Score: %.3f] %s\n", i+1, result.Score, result.DocumentTitle)
    fmt.Printf("   Section: %s\n", result.SectionTitle)
    fmt.Printf("   Content: %s\n", result.Content)
}
```

### 6. 问答

```go
// 基础问答
answer, err := client.QA.Ask(ctx, &sdk.QARequest{
    Query:     "什么是 RAG？",
    DatasetID: datasetID,
})
fmt.Printf("Answer: %s\n", answer.Answer)

// 高级问答
answer, err := client.QA.Ask(ctx, &sdk.QARequest{
    Query:               "RAG 系统如何工作？",
    DatasetID:           datasetID,
    TopK:                5,
    RetrievalMode:       "smart",
    SimilarityThreshold: 0.8,
})

// 查看引用的上下文
fmt.Printf("Answer: %s\n", answer.Answer)
fmt.Printf("Referenced %d contexts:\n", len(answer.Context))
for _, ctx := range answer.Context {
    fmt.Printf("  - %s (Score: %.3f)\n", ctx.DocumentTitle, ctx.Score)
}
```

### 7. 生成

```go
// 直接生成答案（不检索，需要手动提供上下文）
result, err := client.Generate.Generate(ctx, &sdk.GenerateRequest{
    Query:     "请总结以下内容",
    Context:   "上下文内容...",
    DatasetID: datasetID,
})
fmt.Printf("Generated: %s\n", result.Answer)
```

## 错误处理

SDK 提供了类型化的错误处理：

```go
dataset, err := client.Datasets.Get(ctx, "invalid-id")
if err != nil {
    // 类型断言获取 API 错误
    if apiErr, ok := err.(*sdk.APIError); ok {
        // 检查错误类型
        switch {
        case apiErr.IsNotFound():
            fmt.Println("Dataset not found")
        case apiErr.IsBadRequest():
            fmt.Println("Invalid request")
        case apiErr.IsServerError():
            fmt.Println("Server error")
        default:
            fmt.Printf("API error: %s (status: %d)\n", apiErr.Message, apiErr.StatusCode)
        }
    } else {
        // 其他错误（网络错误等）
        fmt.Printf("Request failed: %v\n", err)
    }
}
```

## 并发操作

SDK 是并发安全的，可以在多个 goroutine 中使用：

```go
// 并发上传多个文档
var wg sync.WaitGroup
documents := []string{"doc1.md", "doc2.md", "doc3.md"}

for _, filename := range documents {
    wg.Add(1)
    go func(fn string) {
        defer wg.Done()
        
        file, _ := os.Open(fn)
        defer file.Close()
        
        _, err := client.Documents.Upload(ctx, &sdk.UploadDocumentRequest{
            DatasetID: datasetID,
            File:      file,
            Filename:  fn,
        })
        if err != nil {
            log.Printf("Failed to upload %s: %v", fn, err)
        }
    }(filename)
}

wg.Wait()
```

## 完整示例

查看 `examples/` 目录获取更多示例：

- `examples/basic/` - 基础使用示例
- `examples/advanced/` - 高级功能示例

## API 文档

### 服务列表

- **Models** - AI 模型管理
  - `Create()`, `List()`, `Get()`, `Update()`, `Delete()`
  - `ListProviderModels()`, `Check()`

- **Datasets** - 数据集管理
  - `Create()`, `List()`, `Get()`, `Update()`, `Delete()`
  - `GetStats()`

- **Documents** - 文档管理
  - `Upload()`, `List()`, `Get()`, `Delete()`

- **Search** - 搜索服务
  - `Search()`

- **QA** - 问答服务
  - `Ask()`

- **Generate** - 生成服务
  - `Generate()`

- **Health** - 健康检查
  - `Check()`

## 最佳实践

### 1. 使用 Context

始终传递 context 以支持超时和取消：

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

results, err := client.Search.Search(ctx, req)
```

### 2. 复用客户端

客户端是并发安全的，应该复用：

```go
// 好的做法
var globalClient *sdk.Client

func init() {
    globalClient, _ = sdk.NewClient("http://localhost:8080")
}

// 避免每次请求都创建新客户端
func badPractice() {
    client, _ := sdk.NewClient("http://localhost:8080") // ❌
    // ...
}
```

### 3. 错误处理

总是检查错误并适当处理：

```go
result, err := client.Search.Search(ctx, req)
if err != nil {
    // 记录日志
    log.Printf("Search failed: %v", err)
    
    // 根据错误类型决定是否重试
    if apiErr, ok := err.(*sdk.APIError); ok {
        if apiErr.IsServerError() {
            // 可以重试
        }
    }
    
    return err
}
```

### 4. 资源清理

及时清理不需要的资源：

```go
// 删除测试数据集
defer func() {
    if err := client.Datasets.Delete(ctx, testDatasetID); err != nil {
        log.Printf("Failed to cleanup: %v", err)
    }
}()
```

## 配置选项

| 选项 | 说明 | 默认值 |
|------|------|--------|
| `WithTimeout()` | 设置请求超时时间 | 30s |
| `WithHTTPClient()` | 使用自定义 HTTP 客户端 | 默认客户端 |
| `WithTransport()` | 设置自定义 Transport | 默认 Transport |

## 常见问题

### Q: 如何设置代理？

```go
transport := &http.Transport{
    Proxy: http.ProxyURL(proxyURL),
}
client, _ := sdk.NewClient(
    baseURL,
    sdk.WithTransport(transport),
)
```

### Q: 如何启用 TLS？

```go
transport := &http.Transport{
    TLSClientConfig: &tls.Config{
        // TLS 配置
    },
}
client, _ := sdk.NewClient(
    "https://your-server",
    sdk.WithTransport(transport),
)
```

### Q: 如何添加请求日志？

```go
// 自定义 RoundTripper 添加日志
type LoggingTransport struct {
    Transport http.RoundTripper
}

func (t *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    log.Printf("Request: %s %s", req.Method, req.URL)
    resp, err := t.Transport.RoundTrip(req)
    if err == nil {
        log.Printf("Response: %d", resp.StatusCode)
    }
    return resp, err
}

client, _ := sdk.NewClient(
    baseURL,
    sdk.WithTransport(&LoggingTransport{
        Transport: http.DefaultTransport,
    }),
)
```

## 贡献

欢迎贡献代码、报告问题或提出建议！

## 许可证

MIT License

