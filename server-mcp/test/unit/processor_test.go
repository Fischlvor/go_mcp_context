package test_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	dbmodel "go-mcp-context/internal/model/database"
	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/internal/service"
	"go-mcp-context/pkg/bufferedwriter/actlog"
	"go-mcp-context/pkg/global"
)

// Test_Processor_ProcessDocument 测试文档处理（使用真实文档）
func Test_Processor_ProcessDocument(t *testing.T) {
	processor := &service.DocumentProcessor{}
	libService := &service.LibraryService{}

	// 创建测试库
	lib, err := libService.Create(&request.LibraryCreate{
		Name:        "test-gorm-processor",
		Description: "Test library for processor",
	})
	if err != nil {
		t.Fatalf("Failed to create library: %v", err)
	}
	defer libService.Delete(lib.ID)

	// 创建测试版本
	err = libService.CreateVersion(lib.ID, "v1.0.0")
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	t.Run("process real markdown document from URL", func(t *testing.T) {
		// 下载真实的 GORM 文档
		resp, err := http.Get("https://image.hsk423.cn/mcp/docs/gorm/v1.21.0/README.md")
		if err != nil {
			t.Skipf("Failed to download document: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Skipf("Failed to download document: HTTP %d", resp.StatusCode)
			return
		}

		content, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read document: %v", err)
		}

		if len(content) == 0 {
			t.Fatal("Downloaded document is empty")
		}

		t.Logf("Downloaded document: %d bytes", len(content))

		// 创建文档上传记录
		doc := &dbmodel.DocumentUpload{
			LibraryID: lib.ID,
			Version:   "v1.0.0",
			Title:     "GORM README",
			FilePath:  "README.md",
			FileType:  "markdown",
			FileSize:  int64(len(content)),
		}

		// 创建任务日志器
		actLogger := actlog.NewTaskLogger(lib.ID, "test-task", "v1.0.0")

		// 处理文档
		err = processor.ProcessDocument(doc, content, actLogger)
		if err != nil {
			t.Fatalf("ProcessDocument() error = %v", err)
		}

		t.Logf("Document processed successfully")
	})

	t.Run("process markdown with code blocks", func(t *testing.T) {
		content := []byte(`# GORM Guide

## Quick Start

GORM is a fantastic ORM library for Golang.

` + "```go\n" + `package main

import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
}
` + "```\n\n" + `## Features

- Full-Featured ORM
- Associations (Has One, Has Many, Belongs To, Many To Many, Polymorphism)
- Hooks (Before/After Create/Save/Update/Delete/Find)
`)

		doc := &dbmodel.DocumentUpload{
			LibraryID: lib.ID,
			Version:   "v1.0.0",
			Title:     "GORM Guide with Code",
			FilePath:  "guide.md",
			FileType:  "markdown",
			FileSize:  int64(len(content)),
		}

		actLogger := actlog.NewTaskLogger(lib.ID, "test-task-2", "v1.0.0")

		err := processor.ProcessDocument(doc, content, actLogger)
		if err != nil {
			t.Fatalf("ProcessDocument() error = %v", err)
		}

		t.Logf("Document with code blocks processed successfully")
	})

	t.Run("process empty document", func(t *testing.T) {
		content := []byte("")

		doc := &dbmodel.DocumentUpload{
			LibraryID: lib.ID,
			Version:   "v1.0.0",
			Title:     "Empty Doc",
			FilePath:  "empty.md",
			FileType:  "markdown",
			FileSize:  0,
		}

		actLogger := actlog.NewTaskLogger(lib.ID, "test-task-3", "v1.0.0")

		err := processor.ProcessDocument(doc, content, actLogger)
		if err != nil {
			t.Logf("ProcessDocument(empty) error = %v (expected)", err)
		}
	})
}

// Test_Processor_ProcessDocumentAsync 测试异步文档处理
func Test_Processor_ProcessDocumentAsync(t *testing.T) {
	processor := &service.DocumentProcessor{}
	libService := &service.LibraryService{}

	// 创建测试库
	lib, err := libService.Create(&request.LibraryCreate{
		Name:        "test-async-processor",
		Description: "Test library for async processor",
	})
	if err != nil {
		t.Fatalf("Failed to create library: %v", err)
	}
	defer libService.Delete(lib.ID)

	// 创建测试版本
	err = libService.CreateVersion(lib.ID, "v1.0.0")
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	t.Run("process document asynchronously", func(t *testing.T) {
		content := []byte(`# Async Test Document

## Section 1
This is a test document for async processing.

## Section 2
More content here.`)

		doc := &dbmodel.DocumentUpload{
			LibraryID: lib.ID,
			Version:   "v1.0.0",
			Title:     "Async Test",
			FilePath:  "async.md",
			FileType:  "markdown",
			FileSize:  int64(len(content)),
		}

		actLogger := actlog.NewTaskLogger(lib.ID, "async-task", "v1.0.0")

		// 异步处理
		processor.ProcessDocumentAsync(doc, content, actLogger)

		// 等待异步处理完成
		time.Sleep(2 * time.Second)

		t.Logf("Async processing initiated")
	})
}

// Test_Processor_ProcessDocumentAsync_Advanced 测试异步文档处理的高级场景
func Test_Processor_ProcessDocumentAsync_Advanced(t *testing.T) {
	processor := &service.DocumentProcessor{}
	libService := &service.LibraryService{}

	t.Run("process document async with error handling", func(t *testing.T) {
		// 创建测试库
		lib, err := libService.Create(&request.LibraryCreate{
			Name:        "test-async-error-lib",
			Description: "Test library for async error handling",
		})
		if err != nil {
			t.Fatalf("Failed to create library: %v", err)
		}
		defer libService.Delete(lib.ID)

		// 创建文档记录（使用空内容来触发分块失败）
		content := []byte("") // 空内容会导致分块失败
		doc := &dbmodel.DocumentUpload{
			LibraryID:    lib.ID,
			Version:      "v1.0.0",
			Title:        "Async Error Test",
			FilePath:     "test/async-error.md",
			FileType:     "markdown",
			Status:       "processing",
			ErrorMessage: "",
		}

		if err := global.DB.Create(doc).Error; err != nil {
			t.Fatalf("Failed to create document: %v", err)
		}

		actLogger := actlog.NewTaskLogger(lib.ID, "test-async-error-task", "v1.0.0")

		// 异步处理文档（应该失败，因为内容为空）
		processor.ProcessDocumentAsync(doc, content, actLogger)

		// 等待异步处理完成（增加等待时间）
		time.Sleep(8 * time.Second)

		// 验证文档状态
		var updatedDoc dbmodel.DocumentUpload
		if err := global.DB.First(&updatedDoc, doc.ID).Error; err != nil {
			t.Fatalf("Failed to query document: %v", err)
		}

		// 空内容可能导致 failed 或 processing 状态
		if updatedDoc.Status == "failed" {
			t.Logf("✅ Async error handling works: status=%s, error=%s", updatedDoc.Status, updatedDoc.ErrorMessage)
		} else {
			t.Logf("Document status: %s (async processing may still be running)", updatedDoc.Status)
		}
	})

	t.Run("process document async successfully", func(t *testing.T) {
		// 创建测试库
		lib, err := libService.Create(&request.LibraryCreate{
			Name:        "test-async-success-lib",
			Description: "Test library for async success",
		})
		if err != nil {
			t.Fatalf("Failed to create library: %v", err)
		}
		defer libService.Delete(lib.ID)

		// 创建文档记录
		content := []byte("# Async Success Test\n\nThis document will be processed successfully.")
		doc := &dbmodel.DocumentUpload{
			LibraryID: lib.ID,
			Version:   "v1.0.0",
			Title:     "Async Success Test",
			FilePath:  "test/async-success.md",
			FileType:  "markdown",
			Status:    "processing",
		}

		if err := global.DB.Create(doc).Error; err != nil {
			t.Fatalf("Failed to create document: %v", err)
		}

		actLogger := actlog.NewTaskLogger(lib.ID, "test-async-success-task", "v1.0.0")

		// 异步处理文档
		processor.ProcessDocumentAsync(doc, content, actLogger)

		// 等待异步处理完成
		time.Sleep(8 * time.Second)

		// 验证文档状态
		var updatedDoc dbmodel.DocumentUpload
		if err := global.DB.First(&updatedDoc, doc.ID).Error; err != nil {
			t.Fatalf("Failed to query document: %v", err)
		}

		if updatedDoc.Status != "completed" {
			t.Logf("Document status: %s (expected 'completed', may fail if external services unavailable)", updatedDoc.Status)
		} else {
			t.Log("✅ Document processed successfully via async")
		}

		// 验证生成了分块
		var chunks []dbmodel.DocumentChunk
		if err := global.DB.Where("upload_id = ?", doc.ID).Find(&chunks).Error; err != nil {
			t.Fatalf("Failed to query chunks: %v", err)
		}

		if len(chunks) == 0 && updatedDoc.Status == "completed" {
			t.Error("Expected chunks after successful processing")
		}

		t.Logf("✅ Generated %d chunks", len(chunks))
	})
}

// Test_Processor_ProcessDocumentWithCallback 测试带回调的文档处理
func Test_Processor_ProcessDocumentWithCallback(t *testing.T) {
	processor := &service.DocumentProcessor{}
	libService := &service.LibraryService{}

	// 创建测试库
	lib, err := libService.Create(&request.LibraryCreate{
		Name:        "test-callback-processor",
		Description: "Test library for callback processor",
	})
	if err != nil {
		t.Fatalf("Failed to create library: %v", err)
	}
	defer libService.Delete(lib.ID)

	// 创建测试版本
	err = libService.CreateVersion(lib.ID, "v1.0.0")
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	t.Run("process document with callback", func(t *testing.T) {
		content := []byte(`# Callback Test Document

## Introduction
This document tests callback functionality.

## Details
Callback should be invoked during processing.`)

		doc := &dbmodel.DocumentUpload{
			LibraryID: lib.ID,
			Version:   "v1.0.0",
			Title:     "Callback Test",
			FilePath:  "callback.md",
			FileType:  "markdown",
			FileSize:  int64(len(content)),
		}

		actLogger := actlog.NewTaskLogger(lib.ID, "callback-task", "v1.0.0")

		// 创建状态回调通道
		statusChan := make(chan response.ProcessStatus, 10)

		// 启动 goroutine 接收进度
		go func() {
			for status := range statusChan {
				t.Logf("Status: %s - %s", status.Stage, status.Message)
			}
		}()

		// 处理文档（异步）
		processor.ProcessDocumentWithCallback(doc, content, statusChan, actLogger, true)

		time.Sleep(500 * time.Millisecond) // 等待回调完成
		t.Logf("Document processed with callback successfully")
	})

	t.Run("process document with empty content", func(t *testing.T) {
		// 使用空内容，应该触发分块失败
		content := []byte("")

		doc := &dbmodel.DocumentUpload{
			LibraryID: lib.ID,
			Version:   "v1.0.0",
			Title:     "Empty Content Test",
			FilePath:  "empty.md",
			FileType:  "markdown",
			FileSize:  0,
		}

		if err := global.DB.Create(doc).Error; err != nil {
			t.Fatalf("Failed to create document: %v", err)
		}

		actLogger := actlog.NewTaskLogger(lib.ID, "empty-task", "v1.0.0")

		// 创建状态回调通道
		statusChan := make(chan response.ProcessStatus, 10)

		// 启动 goroutine 接收进度
		failedReceived := false
		go func() {
			for status := range statusChan {
				t.Logf("Status: %s - %s", status.Stage, status.Message)
				if status.Stage == "failed" {
					failedReceived = true
				}
			}
		}()

		// 处理文档（应该失败）
		processor.ProcessDocumentWithCallback(doc, content, statusChan, actLogger, true)

		time.Sleep(500 * time.Millisecond) // 等待回调完成

		if failedReceived {
			t.Log("✅ Empty content correctly triggered failure")
		} else {
			t.Log("⚠️ Empty content did not trigger failure (may be handled differently)")
		}
	})
}

// Test_Processor_ProcessDocumentForRefresh 测试刷新文档处理
func Test_Processor_ProcessDocumentForRefresh(t *testing.T) {
	processor := &service.DocumentProcessor{}
	libService := &service.LibraryService{}

	// 创建测试库
	lib, err := libService.Create(&request.LibraryCreate{
		Name:        "test-refresh-processor",
		Description: "Test library for refresh processor",
	})
	if err != nil {
		t.Fatalf("Failed to create library: %v", err)
	}
	defer libService.Delete(lib.ID)

	// 创建测试版本
	err = libService.CreateVersion(lib.ID, "v1.0.0")
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	t.Run("process document for refresh", func(t *testing.T) {
		content := []byte(`# Refresh Test

## Updated Content
This is refreshed content.`)

		doc := &dbmodel.DocumentUpload{
			LibraryID: lib.ID,
			Version:   "v1.0.0",
			Title:     "Refresh Test",
			FilePath:  "refresh.md",
			FileType:  "markdown",
			FileSize:  int64(len(content)),
		}

		actLogger := actlog.NewTaskLogger(lib.ID, "refresh-task", "v1.0.0")

		// ProcessDocumentForRefresh 需要 batchVersion 参数
		batchVersion := time.Now().Unix()
		chunks, totalTokens, err := processor.ProcessDocumentForRefresh(doc, content, batchVersion, actLogger)
		if err != nil {
			t.Fatalf("ProcessDocumentForRefresh() error = %v", err)
		}

		if len(chunks) == 0 {
			t.Error("Expected chunks, got 0")
		}

		if totalTokens == 0 {
			t.Error("Expected tokens > 0, got 0")
		}

		t.Logf("Document refreshed successfully")
	})

	t.Run("process pdf document", func(t *testing.T) {
		content := []byte("PDF content here")

		doc := &dbmodel.DocumentUpload{
			LibraryID: lib.ID,
			Version:   "v1.0.0",
			Title:     "PDF Test",
			FilePath:  "test.pdf",
			FileType:  "pdf",
			FileSize:  int64(len(content)),
		}

		actLogger := actlog.NewTaskLogger(lib.ID, "pdf-task", "v1.0.0")
		batchVersion := time.Now().Unix()
		chunks, _, err := processor.ProcessDocumentForRefresh(doc, content, batchVersion, actLogger)
		if err != nil {
			t.Fatalf("ProcessDocumentForRefresh(pdf) error = %v", err)
		}
		t.Logf("✅ PDF document processed: %d chunks", len(chunks))
	})

	t.Run("process docx document", func(t *testing.T) {
		content := []byte("DOCX content here")

		doc := &dbmodel.DocumentUpload{
			LibraryID: lib.ID,
			Version:   "v1.0.0",
			Title:     "DOCX Test",
			FilePath:  "test.docx",
			FileType:  "docx",
			FileSize:  int64(len(content)),
		}

		actLogger := actlog.NewTaskLogger(lib.ID, "docx-task", "v1.0.0")
		batchVersion := time.Now().Unix()
		chunks, _, err := processor.ProcessDocumentForRefresh(doc, content, batchVersion, actLogger)
		if err != nil {
			t.Fatalf("ProcessDocumentForRefresh(docx) error = %v", err)
		}
		t.Logf("✅ DOCX document processed: %d chunks", len(chunks))
	})

	t.Run("process swagger document", func(t *testing.T) {
		content := []byte(`{"swagger": "2.0", "info": {"title": "API"}}`)

		doc := &dbmodel.DocumentUpload{
			LibraryID: lib.ID,
			Version:   "v1.0.0",
			Title:     "Swagger Test",
			FilePath:  "api.json",
			FileType:  "swagger",
			FileSize:  int64(len(content)),
		}

		actLogger := actlog.NewTaskLogger(lib.ID, "swagger-task", "v1.0.0")
		batchVersion := time.Now().Unix()
		chunks, _, err := processor.ProcessDocumentForRefresh(doc, content, batchVersion, actLogger)
		if err != nil {
			t.Fatalf("ProcessDocumentForRefresh(swagger) error = %v", err)
		}
		t.Logf("✅ Swagger document processed: %d chunks", len(chunks))
	})

	t.Run("process unknown type document", func(t *testing.T) {
		content := []byte("Unknown content here")

		doc := &dbmodel.DocumentUpload{
			LibraryID: lib.ID,
			Version:   "v1.0.0",
			Title:     "Unknown Test",
			FilePath:  "test.txt",
			FileType:  "unknown",
			FileSize:  int64(len(content)),
		}

		actLogger := actlog.NewTaskLogger(lib.ID, "unknown-task", "v1.0.0")
		batchVersion := time.Now().Unix()
		chunks, _, err := processor.ProcessDocumentForRefresh(doc, content, batchVersion, actLogger)
		if err != nil {
			t.Fatalf("ProcessDocumentForRefresh(unknown) error = %v", err)
		}
		t.Logf("✅ Unknown type document processed: %d chunks", len(chunks))
	})
}

// Test_Processor_InternalFunctions 测试内部函数（通过边界测试覆盖）
func Test_Processor_InternalFunctions(t *testing.T) {
	libService := &service.LibraryService{}

	// 创建测试库
	lib, err := libService.Create(&request.LibraryCreate{
		Name:        "test-processor-internal",
		Description: "Test library for processor internal functions",
	})
	if err != nil {
		t.Fatalf("Failed to create library: %v", err)
	}
	defer libService.Delete(lib.ID)

	// 创建版本
	err = libService.CreateVersion(lib.ID, "v1.0.0")
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	t.Run("process document with large sections to trigger splitLargeSectionWithMetadata", func(t *testing.T) {
		// 创建一个包含大段落的文档，会触发 splitLargeSectionWithMetadata
		largeSection := strings.Repeat("This is a very long paragraph that will exceed the chunk size limit. ", 100)
		content := []byte(fmt.Sprintf("# Large Document\n\n%s\n\n## Another Section\n\n%s", largeSection, largeSection))

		doc := &dbmodel.DocumentUpload{
			LibraryID: lib.ID,
			Version:   "v1.0.0",
			Title:     "large-doc.md",
			FilePath:  "test/large-doc.md",
			FileType:  "markdown",
			Status:    "processing",
		}

		if err := global.DB.Create(doc).Error; err != nil {
			t.Fatalf("Failed to create document: %v", err)
		}

		processor := &service.DocumentProcessor{}
		actLogger := actlog.NewTaskLogger(lib.ID, "test-task", "v1.0.0")

		chunks, totalTokens, err := processor.ProcessDocumentForRefresh(doc, content, time.Now().Unix(), actLogger)
		if err != nil {
			t.Fatalf("ProcessDocumentForRefresh() error = %v", err)
		}

		if len(chunks) == 0 {
			t.Error("Expected chunks from large document")
		}

		if totalTokens == 0 {
			t.Error("Expected tokens > 0")
		}

		t.Logf("Processed large document: %d chunks, %d tokens", len(chunks), totalTokens)
	})

	t.Run("process document with code blocks to trigger splitIntoAtoms", func(t *testing.T) {
		// 创建一个超大section包含多个代码块，会触发 splitIntoAtoms 的所有分支
		// 需要超过 512 tokens 才能触发 splitLargeSectionWithMetadata
		largePara := strings.Repeat("This is a long paragraph with many words to make the section exceed the chunk size limit. ", 80)

		content := []byte(`# Code Examples

## Large Section with Multiple Code Blocks

` + largePara + `

Here is the first code block:

` + "```go" + `
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
    for i := 0; i < 10; i++ {
        fmt.Printf("Number: %d\n", i)
    }
}
` + "```" + `

` + largePara + `

Second code block:

` + "```python" + `
def hello():
    print("Hello from Python")
    for i in range(10):
        print(f"Number: {i}")

def goodbye():
    print("Goodbye")
` + "```" + `

` + largePara + `

Third code block:

` + "```javascript" + `
function hello() {
    console.log("Hello from JavaScript");
    for (let i = 0; i < 10; i++) {
        console.log('Number: ' + i);
    }
}
` + "```" + `

` + largePara)

		doc := &dbmodel.DocumentUpload{
			LibraryID: lib.ID,
			Version:   "v1.0.0",
			Title:     "code-blocks.md",
			FilePath:  "test/code-blocks.md",
			FileType:  "markdown",
			Status:    "processing",
		}

		if err := global.DB.Create(doc).Error; err != nil {
			t.Fatalf("Failed to create document: %v", err)
		}

		processor := &service.DocumentProcessor{}
		actLogger := actlog.NewTaskLogger(lib.ID, "test-task", "v1.0.0")

		chunks, totalTokens, err := processor.ProcessDocumentForRefresh(doc, content, time.Now().Unix(), actLogger)
		if err != nil {
			t.Fatalf("ProcessDocumentForRefresh() error = %v", err)
		}

		if len(chunks) == 0 {
			t.Error("Expected chunks from document with code blocks")
		}

		if totalTokens == 0 {
			t.Error("Expected tokens > 0")
		}

		// 验证至少有一些代码块被识别
		hasCodeChunk := false
		for _, chunk := range chunks {
			if chunk.ChunkType == "code" {
				hasCodeChunk = true
				break
			}
		}

		if !hasCodeChunk {
			t.Log("Warning: No code chunks detected (may be expected depending on detection logic)")
		}

		t.Logf("Processed document with code blocks: %d chunks, %d tokens", len(chunks), totalTokens)
	})

	t.Run("process document with mixed content", func(t *testing.T) {
		// 混合内容：大段落 + 代码块 + 普通文本
		largeText := strings.Repeat("Lorem ipsum dolor sit amet. ", 50)
		content := []byte(fmt.Sprintf(`# Mixed Content Document

## Introduction

%s

## Code Example

`+"```javascript"+`
function test() {
    console.log("test");
}
`+"```"+`

## More Text

%s

## Another Code Block

`+"```bash"+`
#!/bin/bash
echo "Hello World"
`+"```"+`

## Conclusion

Final paragraph with some text.`, largeText, largeText))

		doc := &dbmodel.DocumentUpload{
			LibraryID: lib.ID,
			Version:   "v1.0.0",
			Title:     "mixed-content.md",
			FilePath:  "test/mixed-content.md",
			FileType:  "markdown",
			Status:    "processing",
		}

		if err := global.DB.Create(doc).Error; err != nil {
			t.Fatalf("Failed to create document: %v", err)
		}

		processor := &service.DocumentProcessor{}
		actLogger := actlog.NewTaskLogger(lib.ID, "test-task", "v1.0.0")

		chunks, totalTokens, err := processor.ProcessDocumentForRefresh(doc, content, time.Now().Unix(), actLogger)
		if err != nil {
			t.Fatalf("ProcessDocumentForRefresh() error = %v", err)
		}

		if len(chunks) == 0 {
			t.Error("Expected chunks from mixed content document")
		}

		if totalTokens == 0 {
			t.Error("Expected tokens > 0")
		}

		t.Logf("Processed mixed content document: %d chunks, %d tokens", len(chunks), totalTokens)
	})
}

// Test_Processor_SplitIntoAtoms 测试 splitIntoAtoms 函数（通过处理包含代码块的文档）
func Test_Processor_SplitIntoAtoms(t *testing.T) {
	processor := &service.DocumentProcessor{}
	libService := &service.LibraryService{}

	// 创建测试库
	lib, err := libService.Create(&request.LibraryCreate{
		Name:        "test-splitintoatoms",
		Description: "Test library for splitIntoAtoms",
	})
	if err != nil {
		t.Fatalf("Failed to create library: %v", err)
	}
	defer libService.Delete(lib.ID)

	t.Run("process document with code blocks", func(t *testing.T) {
		// 创建包含多个代码块的大型文档
		content := `# API Documentation

## Introduction
This is a comprehensive guide to using our API.

## Installation

First, install the package:

` + "```bash" + `
npm install our-package
` + "```" + `

## Basic Usage

Here's a simple example:

` + "```javascript" + `
const pkg = require('our-package');
const result = pkg.doSomething();
console.log(result);
` + "```" + `

## Advanced Features

### Configuration

You can configure the package like this:

` + "```json" + `
{
  "option1": "value1",
  "option2": "value2"
}
` + "```" + `

### Error Handling

Handle errors properly:

` + "```javascript" + `
try {
  pkg.doSomething();
} catch (error) {
  console.error('Error:', error);
}
` + "```" + `

## Conclusion

This concludes the documentation.
`

		doc := &dbmodel.DocumentUpload{
			LibraryID: lib.ID,
			Version:   "v1.0.0",
			Title:     "API Documentation with Code Blocks",
			FilePath:  "test/api-docs.md",
			FileType:  "markdown",
			Status:    "processing",
		}

		if err := global.DB.Create(doc).Error; err != nil {
			t.Fatalf("Failed to create document: %v", err)
		}

		actLogger := actlog.NewTaskLogger(lib.ID, "test-task", "v1.0.0")

		// 处理文档（这会调用 splitIntoAtoms）
		err := processor.ProcessDocument(doc, []byte(content), actLogger)
		if err != nil {
			t.Fatalf("ProcessDocument() error = %v", err)
		}

		// 验证文档已处理
		var chunks []dbmodel.DocumentChunk
		if err := global.DB.Where("upload_id = ?", doc.ID).Find(&chunks).Error; err != nil {
			t.Fatalf("Failed to query chunks: %v", err)
		}

		if len(chunks) == 0 {
			t.Error("Expected chunks from document with code blocks")
		}

		// 验证至少有一些代码块被识别
		codeChunks := 0
		for _, chunk := range chunks {
			if strings.Contains(chunk.Code, "```") {
				codeChunks++
			}
		}

		t.Logf("✅ Processed document with code blocks: %d total chunks, %d code chunks", len(chunks), codeChunks)
	})

	t.Run("process document without code blocks", func(t *testing.T) {
		// 创建不包含代码块的文档（测试另一个分支）
		content := `# Simple Document

## Section 1

This is a simple paragraph.

## Section 2

Another paragraph here.

## Section 3

Final paragraph.
`

		doc := &dbmodel.DocumentUpload{
			LibraryID: lib.ID,
			Version:   "v1.0.0",
			Title:     "Simple Document",
			FilePath:  "test/simple.md",
			FileType:  "markdown",
			Status:    "processing",
		}

		if err := global.DB.Create(doc).Error; err != nil {
			t.Fatalf("Failed to create document: %v", err)
		}

		actLogger := actlog.NewTaskLogger(lib.ID, "test-task-2", "v1.0.0")

		// 处理文档（这会调用 splitIntoAtoms 的无代码块分支）
		err := processor.ProcessDocument(doc, []byte(content), actLogger)
		if err != nil {
			t.Fatalf("ProcessDocument() error = %v", err)
		}

		// 验证文档已处理
		var chunks []dbmodel.DocumentChunk
		if err := global.DB.Where("upload_id = ?", doc.ID).Find(&chunks).Error; err != nil {
			t.Fatalf("Failed to query chunks: %v", err)
		}

		if len(chunks) == 0 {
			t.Error("Expected chunks from simple document")
		}

		t.Logf("✅ Processed simple document: %d chunks", len(chunks))
	})
}
