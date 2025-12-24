package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"go-mcp-context/internal/model/database"
	"go-mcp-context/internal/model/request"
	"go-mcp-context/internal/model/response"
	"go-mcp-context/pkg/global"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	APIKeyPrefix      = "mcpsk-"
	MaxAPIKeysPerUser = 5
)

type ApiKeyService struct{}

// Create 创建 API Key
func (s *ApiKeyService) Create(userUUID string, req *request.APIKeyCreate) (*response.APIKeyCreateResponse, error) {
	// 检查用户已有的 API Key 数量
	var count int64
	if err := global.DB.Model(&database.APIKey{}).
		Where("user_uuid = ? AND deleted_at IS NULL", userUUID).
		Count(&count).Error; err != nil {
		global.Log.Error("查询 API Key 数量失败", zap.Error(err))
		return nil, errors.New("查询失败")
	}

	if count >= MaxAPIKeysPerUser {
		return nil, errors.New("已达到最大 API Key 数量限制（5 个）")
	}

	// 生成 API Key：mcpsk-<UUID v4>
	uuidV4, err := uuid.NewV4()
	if err != nil {
		global.Log.Error("生成 UUID 失败", zap.Error(err))
		return nil, errors.New("生成 API Key 失败")
	}

	apiKey := APIKeyPrefix + uuidV4.String()
	tokenSuffix := apiKey[len(apiKey)-4:] // 后 4 位

	// 计算 SHA256 哈希
	hash := sha256.Sum256([]byte(apiKey))
	tokenHash := hex.EncodeToString(hash[:])

	// 保存到数据库
	token := &database.APIKey{
		UserUUID:    userUUID,
		TokenHash:   tokenHash,
		TokenSuffix: tokenSuffix,
		Name:        req.Name,
	}

	if err := global.DB.Create(token).Error; err != nil {
		global.Log.Error("创建 API Key 失败", zap.Error(err))
		return nil, errors.New("创建失败")
	}

	global.Log.Info("API Key 创建成功",
		zap.String("user_uuid", userUUID),
		zap.Uint("token_id", token.ID),
		zap.String("name", req.Name),
	)

	return &response.APIKeyCreateResponse{
		ID:          token.ID,
		Name:        token.Name,
		APIKey:      apiKey, // 仅此一次返回完整 key
		TokenSuffix: tokenSuffix,
		CreatedAt:   token.CreatedAt,
	}, nil
}

// List 获取用户的 API Key 列表
func (s *ApiKeyService) List(userUUID string) ([]response.APIKeyListItem, error) {
	var tokens []database.APIKey
	if err := global.DB.
		Where("user_uuid = ? AND deleted_at IS NULL", userUUID).
		Order("created_at DESC").
		Find(&tokens).Error; err != nil {
		global.Log.Error("查询 API Key 列表失败", zap.Error(err))
		return nil, errors.New("查询失败")
	}

	items := make([]response.APIKeyListItem, len(tokens))
	for i, t := range tokens {
		items[i] = response.APIKeyListItem{
			ID:          t.ID,
			Name:        t.Name,
			TokenSuffix: t.TokenSuffix,
			LastUsedAt:  t.LastUsedAt,
			CreatedAt:   t.CreatedAt,
		}
	}

	return items, nil
}

// Delete 删除 API Key（软删除）
func (s *ApiKeyService) Delete(userUUID string, id uint) error {
	result := global.DB.
		Where("id = ? AND user_uuid = ? AND deleted_at IS NULL", id, userUUID).
		Delete(&database.APIKey{})

	if result.Error != nil {
		global.Log.Error("删除 API Key 失败", zap.Error(result.Error))
		return errors.New("删除失败")
	}

	if result.RowsAffected == 0 {
		return errors.New("API Key 不存在或已删除")
	}

	global.Log.Info("API Key 删除成功",
		zap.String("user_uuid", userUUID),
		zap.Uint("token_id", id),
	)

	return nil
}

// ValidateAPIKey 验证 API Key，返回用户 UUID
func (s *ApiKeyService) ValidateAPIKey(apiKey string) (string, error) {
	// 检查前缀
	if len(apiKey) < len(APIKeyPrefix) || apiKey[:len(APIKeyPrefix)] != APIKeyPrefix {
		return "", errors.New("无效的 API Key 格式")
	}

	// 计算哈希
	hash := sha256.Sum256([]byte(apiKey))
	tokenHash := hex.EncodeToString(hash[:])

	// 查询数据库
	var token database.APIKey
	if err := global.DB.
		Where("token_hash = ? AND deleted_at IS NULL", tokenHash).
		First(&token).Error; err != nil {
		return "", errors.New("无效的 API Key")
	}

	// 更新使用次数和最后使用时间（异步，不阻塞请求）
	go func() {
		now := time.Now()
		global.DB.Model(&token).Updates(map[string]interface{}{
			"usage_count":  gorm.Expr("usage_count + 1"),
			"last_used_at": now,
		})
	}()

	return token.UserUUID, nil
}
