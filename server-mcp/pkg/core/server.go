package core

import (
	"fmt"

	"go-mcp-context/internal/initialize"
	"go-mcp-context/pkg/global"
)

// RunServer 启动服务器
func RunServer() {
	addr := fmt.Sprintf("%s:%d", global.Config.System.Host, global.Config.System.Port)
	router := initialize.InitRouter()

	fmt.Printf("Server running on %s\n", addr)
	if err := router.Run(addr); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
