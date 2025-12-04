package flag

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

// 定义 CLI 标志
var (
	sqlFlag = &cli.BoolFlag{
		Name:  "sql",
		Usage: "Initializes the structure of the PostgreSQL database tables.",
	}
)

// Run 执行基于命令行标志的相应操作
func Run(c *cli.Context) {
	if c.NumFlags() > 1 {
		fmt.Println("Only one command can be specified")
		os.Exit(1)
	}

	switch {
	case c.Bool(sqlFlag.Name):
		if err := SQL(); err != nil {
			fmt.Printf("Failed to create table structure: %v\n", err)
		} else {
			fmt.Println("Successfully created table structure")
		}
	default:
		fmt.Println("Unknown command")
	}
}

// NewApp 创建 CLI 应用程序
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "Go MCP Context"
	app.Flags = []cli.Flag{
		sqlFlag,
	}
	app.Action = Run
	return app
}

// InitFlag 初始化并运行 CLI 应用程序
func InitFlag() {
	if len(os.Args) > 1 {
		app := NewApp()
		if err := app.Run(os.Args); err != nil {
			fmt.Printf("Application error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}
}
