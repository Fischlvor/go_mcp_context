package config

// Chunker 文档分块配置
type Chunker struct {
	ChunkSize int `json:"chunk_size" yaml:"chunk_size"` // 分块大小（tokens）
	Overlap   int `json:"overlap" yaml:"overlap"`       // 重叠大小（tokens）
}
