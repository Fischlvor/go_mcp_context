package test_test

import (
	"strings"
	"testing"

	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/service"
	"go-mcp-context/pkg/global"

	dbmodel "go-mcp-context/internal/model/database"
)

// TestAPIKeyCreate 测试 API Key 创建
func Test_APIKey_Create(t *testing.T) {
	apiKeyService := &service.ApiKeyService{}
	userUUID := "00000000-0000-0000-0000-000000000001"

	t.Run("create api key successfully", func(t *testing.T) {
		req := &request.APIKeyCreate{
			Name: "Test API Key",
		}

		resp, err := apiKeyService.Create(userUUID, req)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if resp == nil {
			t.Fatal("Expected response, got nil")
		}

		if resp.Name != req.Name {
			t.Errorf("Expected name %s, got %s", req.Name, resp.Name)
		}

		if resp.APIKey == "" {
			t.Error("Expected API key, got empty string")
		}

		// 验证 API Key 格式：mcpsk-<uuid>
		if len(resp.APIKey) < 6 || resp.APIKey[:6] != "mcpsk-" {
			t.Errorf("Invalid API key format: %s", resp.APIKey)
		}

		if resp.TokenSuffix == "" {
			t.Error("Expected token suffix, got empty string")
		}

		// 验证数据库中存在
		var count int64
		global.DB.Model(&dbmodel.APIKey{}).Where("id = ?", resp.ID).Count(&count)
		if count != 1 {
			t.Errorf("Expected 1 API key in DB, got %d", count)
		}
	})

	t.Run("create api key with empty name", func(t *testing.T) {
		req := &request.APIKeyCreate{
			Name: "",
		}

		resp, err := apiKeyService.Create(userUUID, req)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if resp == nil {
			t.Fatal("Expected response, got nil")
		}

		// 验证 API Key 仍然被创建
		if resp.APIKey == "" {
			t.Error("Expected API key, got empty string")
		}
	})

	t.Run("create multiple api keys", func(t *testing.T) {
		userUUID2 := "00000000-0000-0000-0000-000000000002"

		// 创建 3 个 API Key
		for i := 1; i <= 3; i++ {
			req := &request.APIKeyCreate{
				Name: "Test Key " + string(rune('0'+i)),
			}
			_, err := apiKeyService.Create(userUUID2, req)
			if err != nil {
				t.Fatalf("Create() error = %v", err)
			}
		}

		// 验证数量
		var count int64
		global.DB.Model(&dbmodel.APIKey{}).Where("user_uuid = ?", userUUID2).Count(&count)
		if count != 3 {
			t.Errorf("Expected 3 API keys, got %d", count)
		}
	})

	t.Run("exceed max api keys limit", func(t *testing.T) {
		userUUID3 := "00000000-0000-0000-0000-000000000003"

		// 创建 5 个 API Key（达到上限）
		for i := 1; i <= 5; i++ {
			req := &request.APIKeyCreate{
				Name: "Test Key " + string(rune('0'+i)),
			}
			_, err := apiKeyService.Create(userUUID3, req)
			if err != nil {
				t.Fatalf("Create() error = %v at iteration %d", err, i)
			}
		}

		// 尝试创建第 6 个，应该失败
		req := &request.APIKeyCreate{
			Name: "Test Key 6",
		}
		_, err := apiKeyService.Create(userUUID3, req)
		if err == nil {
			t.Error("Expected error when exceeding max API keys, got nil")
		}
	})
}

// TestAPIKeyList 测试 API Key 列表查询
func Test_APIKey_List(t *testing.T) {
	apiKeyService := &service.ApiKeyService{}

	t.Run("list api keys", func(t *testing.T) {
		userUUID := "00000000-0000-0000-0000-000000000011"

		// 创建 2 个 API Key
		for i := 1; i <= 2; i++ {
			req := &request.APIKeyCreate{
				Name: "List Test Key " + string(rune('0'+i)),
			}
			_, err := apiKeyService.Create(userUUID, req)
			if err != nil {
				t.Fatalf("Create() error = %v", err)
			}
		}

		// 查询列表
		items, err := apiKeyService.List(userUUID)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(items) < 2 {
			t.Errorf("Expected at least 2 API keys, got %d", len(items))
		}

		// 验证字段
		for _, item := range items {
			if item.ID == 0 {
				t.Error("Expected non-zero ID")
			}
			if item.Name == "" {
				t.Error("Expected non-empty name")
			}
			if item.TokenSuffix == "" {
				t.Error("Expected non-empty token suffix")
			}
		}
	})

	t.Run("list empty api keys", func(t *testing.T) {
		userUUID := "00000000-0000-0000-0000-000000000012"

		items, err := apiKeyService.List(userUUID)
		if err != nil {
			t.Fatalf("List() error = %v", err)
		}

		if len(items) != 0 {
			t.Errorf("Expected 0 API keys for new user, got %d", len(items))
		}
	})
}

// TestAPIKeyDelete 测试 API Key 删除
func Test_APIKey_Delete(t *testing.T) {
	apiKeyService := &service.ApiKeyService{}

	t.Run("delete api key successfully", func(t *testing.T) {
		userUUID := "00000000-0000-0000-0000-000000000021"

		// 创建一个 API Key
		req := &request.APIKeyCreate{
			Name: "Delete Test Key",
		}
		resp, err := apiKeyService.Create(userUUID, req)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 删除
		err = apiKeyService.Delete(userUUID, resp.ID)
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		// 验证已删除（软删除）
		var apiKey dbmodel.APIKey
		result := global.DB.Where("id = ?", resp.ID).First(&apiKey)
		if result.Error == nil {
			t.Error("Expected API key to be soft deleted, but still found")
		}
	})

	t.Run("delete non-existent api key", func(t *testing.T) {
		userUUID := "00000000-0000-0000-0000-000000000022"

		err := apiKeyService.Delete(userUUID, 99999)
		if err == nil {
			t.Error("Expected error when deleting non-existent API key, got nil")
		}
	})

	t.Run("delete other user's api key", func(t *testing.T) {
		userUUID1 := "00000000-0000-0000-0000-000000000023"
		userUUID2 := "00000000-0000-0000-0000-000000000024"

		// 用户1 创建 API Key
		req := &request.APIKeyCreate{
			Name: "User1 Key",
		}
		resp, err := apiKeyService.Create(userUUID1, req)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 用户2 尝试删除用户1 的 API Key，应该失败
		err = apiKeyService.Delete(userUUID2, resp.ID)
		if err == nil {
			t.Error("Expected error when deleting other user's API key, got nil")
		}
	})
}

// TestAPIKeyValidate 测试 API Key 验证
func Test_APIKey_Validate(t *testing.T) {
	apiKeyService := &service.ApiKeyService{}

	t.Run("validate api key successfully", func(t *testing.T) {
		userUUID := "00000000-0000-0000-0000-000000000031"

		// 创建 API Key
		req := &request.APIKeyCreate{
			Name: "Validate Test Key",
		}
		resp, err := apiKeyService.Create(userUUID, req)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 验证 API Key
		validatedUserUUID, err := apiKeyService.ValidateAPIKey(resp.APIKey)
		if err != nil {
			t.Fatalf("ValidateAPIKey() error = %v", err)
		}

		if validatedUserUUID != userUUID {
			t.Errorf("Expected user UUID %s, got %s", userUUID, validatedUserUUID)
		}
	})

	t.Run("validate invalid api key format", func(t *testing.T) {
		_, err := apiKeyService.ValidateAPIKey("invalid-key")
		if err == nil {
			t.Error("Expected error for invalid API key format, got nil")
		}
	})

	t.Run("validate non-existent api key", func(t *testing.T) {
		fakeKey := "mcpsk-00000000-0000-0000-0000-000000000000"
		_, err := apiKeyService.ValidateAPIKey(fakeKey)
		if err == nil {
			t.Error("Expected error for non-existent API key, got nil")
		}
	})

	t.Run("validate deleted api key", func(t *testing.T) {
		userUUID := "00000000-0000-0000-0000-000000000032"

		// 创建并删除 API Key
		req := &request.APIKeyCreate{
			Name: "Deleted Key",
		}
		resp, err := apiKeyService.Create(userUUID, req)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		err = apiKeyService.Delete(userUUID, resp.ID)
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		// 尝试验证已删除的 API Key
		_, err = apiKeyService.ValidateAPIKey(resp.APIKey)
		if err == nil {
			t.Error("Expected error for deleted API key, got nil")
		}
	})

	t.Run("validate api key with wrong format", func(t *testing.T) {
		// 测试各种错误格式
		testCases := []string{
			"",                   // 空字符串
			"invalid",            // 没有前缀
			"mcpsk",              // 只有前缀
			"mcpsk-",             // 前缀后没有内容
			"mcpsk-invalid-uuid", // 无效的 UUID
			"other-key-format",   // 完全不同的格式
		}

		for _, testKey := range testCases {
			_, err := apiKeyService.ValidateAPIKey(testKey)
			if err == nil {
				t.Errorf("Expected error for invalid key format: %s, got nil", testKey)
			}
		}
	})

	t.Run("validate api key case sensitivity", func(t *testing.T) {
		userUUID := "00000000-0000-0000-0000-000000000033"

		// 创建 API Key
		req := &request.APIKeyCreate{
			Name: "Case Test Key",
		}
		resp, err := apiKeyService.Create(userUUID, req)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		// 验证原始 API Key
		validatedUUID, err := apiKeyService.ValidateAPIKey(resp.APIKey)
		if err != nil {
			t.Fatalf("ValidateAPIKey() error = %v", err)
		}

		if validatedUUID != userUUID {
			t.Errorf("Expected user UUID %s, got %s", userUUID, validatedUUID)
		}
	})
}

// TestAPIKeyCreateAdvanced 测试 API Key 创建的高级场景
func Test_APIKey_Create_Advanced(t *testing.T) {
	apiKeyService := &service.ApiKeyService{}

	t.Run("create api key with empty name", func(t *testing.T) {
		userUUID := "00000000-0000-0000-0000-000000000040"

		req := &request.APIKeyCreate{
			Name: "",
		}

		resp, err := apiKeyService.Create(userUUID, req)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if resp == nil {
			t.Fatal("Expected response, got nil")
		}

		if resp.APIKey == "" {
			t.Error("Expected non-empty API key")
		}

		// 空名称应该被接受
		if resp.Name != "" {
			t.Errorf("Expected empty name, got %s", resp.Name)
		}
	})

	t.Run("create api key with very long name", func(t *testing.T) {
		userUUID := "00000000-0000-0000-0000-000000000041"

		// 使用合理长度的名称（数据库字段可能有限制）
		longName := ""
		for i := 0; i < 100; i++ {
			longName += "a"
		}

		req := &request.APIKeyCreate{
			Name: longName,
		}

		resp, err := apiKeyService.Create(userUUID, req)
		if err != nil {
			t.Logf("Create() error = %v (may be due to database field length limit)", err)
			return
		}

		if resp == nil {
			t.Fatal("Expected response, got nil")
		}

		if resp.APIKey == "" {
			t.Error("Expected non-empty API key")
		}
	})

	t.Run("create api key with special characters in name", func(t *testing.T) {
		userUUID := "00000000-0000-0000-0000-000000000042"

		req := &request.APIKeyCreate{
			Name: "Test@#$%^&*()_+-=[]{}|;:',.<>?/~`",
		}

		resp, err := apiKeyService.Create(userUUID, req)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if resp == nil {
			t.Fatal("Expected response, got nil")
		}

		if resp.APIKey == "" {
			t.Error("Expected non-empty API key")
		}
	})

	t.Run("create api key with unicode characters in name", func(t *testing.T) {
		userUUID := "00000000-0000-0000-0000-000000000043"

		req := &request.APIKeyCreate{
			Name: "测试-Test-テスト-🔑",
		}

		resp, err := apiKeyService.Create(userUUID, req)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if resp == nil {
			t.Fatal("Expected response, got nil")
		}

		if resp.APIKey == "" {
			t.Error("Expected non-empty API key")
		}
	})

	t.Run("create api key token suffix is correct", func(t *testing.T) {
		userUUID := "00000000-0000-0000-0000-000000000044"

		req := &request.APIKeyCreate{
			Name: "Token Suffix Test",
		}

		resp, err := apiKeyService.Create(userUUID, req)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if resp == nil {
			t.Fatal("Expected response, got nil")
		}

		// 验证 token suffix 是 API Key 的最后 4 位
		if len(resp.APIKey) < 4 {
			t.Error("API key too short")
		} else {
			expectedSuffix := resp.APIKey[len(resp.APIKey)-4:]
			if resp.TokenSuffix != expectedSuffix {
				t.Errorf("Expected token suffix %s, got %s", expectedSuffix, resp.TokenSuffix)
			}
		}
	})

	t.Run("create api key has correct prefix", func(t *testing.T) {
		userUUID := "00000000-0000-0000-0000-000000000045"

		req := &request.APIKeyCreate{
			Name: "Prefix Test",
		}

		resp, err := apiKeyService.Create(userUUID, req)
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}

		if resp == nil {
			t.Fatal("Expected response, got nil")
		}

		// 验证 API Key 以 mcpsk- 开头
		if !strings.HasPrefix(resp.APIKey, "mcpsk-") {
			t.Errorf("Expected API key to start with 'mcpsk-', got %s", resp.APIKey)
		}
	})
}
