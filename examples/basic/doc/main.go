package main

import (
	"context"
	"fmt"
	"log"

	sdk "github.com/chaitin/raglite-go-sdk"
)

type DocumentMetadata struct {
	GroupIDs []int `json:"group_ids"`
}

type Document struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	DatasetID   string           `json:"dataset_id"`
	Status      string           `json:"status"`
	ProgressMsg string           `json:"progress_msg"`
	MetaData    DocumentMetadata `json:"meta_data"`
	Tags        []string         `json:"tags"`
}

func main() {
	// 创建客户端
	client, err := sdk.NewClient("http://localhost:5050")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	docs, err := client.Documents.List(ctx, &sdk.ListDocumentsRequest{
		DatasetID: "9e2f3ae9-6624-4c52-ba11-f88b847987f6",
		PageSize:  1,
	})
	if err != nil {
		log.Fatalf("Failed to list documents: %v", err)
	}
	var documents []Document
	for _, item := range docs.Documents {
		doc := Document{
			ID:          item.ID,
			Name:        item.Title,
			DatasetID:   item.DatasetID,
			Status:      item.Status,
			ProgressMsg: item.ProgressMsg,
		}
		doc.MetaData = sdk.Decode[DocumentMetadata](item.Metadata.Data)
		documents = append(documents, doc)
	}
	fmt.Printf("%+v\n", documents)
}
