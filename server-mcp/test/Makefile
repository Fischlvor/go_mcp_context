# Makefile for running tests and generating coverage reports
# Usage:
#   make test              - Run all tests
#   make test TEST=TestName - Run specific test
#   make coverage          - Generate coverage report
#   make clean             - Clean up generated files

# 变量定义
PROJECT_ROOT := ..
TEST_DIR := .
UNIT_TEST_DIR := $(TEST_DIR)/unit
INTEGRATION_TEST_DIR := $(TEST_DIR)/integration
COVERAGE_FILE := coverage.out
TEST_LOG := test_log.txt
COVERAGE_REPORT := all_functions_coverage.txt
COVERAGE_PKG := go-mcp-context/internal/service/...

# 测试超时时间
TIMEOUT := 5m

# 默认目标 - 显示帮助信息
.DEFAULT_GOAL := help

# 一键测试 + 覆盖率（清理 + 运行所有测试 + 生成覆盖率）
.PHONY: all
all: clean test coverage
	@echo ""
	@echo "✅ 测试和覆盖率报告生成完成！"
	@echo "📊 查看覆盖率: make show-coverage"
	@echo "📝 查看日志: make show-log"

# 运行所有测试（单元测试 + 集成测试）
.PHONY: test
test:
	@echo "=== 开始运行所有测试（单元 + 集成） ==="
	@echo "所有测试时间: $$(date '+%Y-%m-%d %H:%M:%S')" > $(TEST_LOG)
	@echo "" >> $(TEST_LOG)
	@cd $(TEST_DIR) && stdbuf -oL -eL go test ../test/... \
		-v \
		-coverprofile=$(COVERAGE_FILE) \
		-coverpkg=$(COVERAGE_PKG) \
		-timeout $(TIMEOUT) \
		2>&1 | tee -a $(TEST_LOG)
	@echo "" >> $(TEST_LOG)
	@echo "=== 所有测试完成 ===" | tee -a $(TEST_LOG)

# 运行单元测试
.PHONY: test-unit
test-unit:
	@echo "=== 开始运行测试 ==="
	@echo "测试时间: $$(date '+%Y-%m-%d %H:%M:%S')" > $(TEST_LOG)
	@echo "测试参数: TEST=$(TEST)" >> $(TEST_LOG)
	@echo "" >> $(TEST_LOG)
ifdef TEST
	@echo "运行指定测试: $(TEST)"
	@echo "运行指定测试: $(TEST)" >> $(TEST_LOG)
	@echo "========================================" >> $(TEST_LOG)
	@cd $(TEST_DIR) && stdbuf -oL -eL go test ../test/unit/... \
		-v \
		-run $(TEST) \
		-coverprofile=$(COVERAGE_FILE) \
		-coverpkg=$(COVERAGE_PKG) \
		-timeout $(TIMEOUT) 2>&1 | tee -a $(TEST_LOG); \
	TEST_EXIT=$$?; \
	echo "========================================" >> $(TEST_LOG); \
	exit $$TEST_EXIT
else
	@echo "运行所有单元测试..."
	@echo "运行所有单元测试" >> $(TEST_LOG)
	@echo "========================================" >> $(TEST_LOG)
	@cd $(TEST_DIR) && stdbuf -oL -eL go test ../test/unit/... \
		-v \
		-coverprofile=$(COVERAGE_FILE) \
		-coverpkg=$(COVERAGE_PKG) \
		-timeout $(TIMEOUT) 2>&1 | tee -a $(TEST_LOG); \
	TEST_EXIT=$$?; \
	echo "========================================" >> $(TEST_LOG); \
	exit $$TEST_EXIT
endif
	@echo "" >> $(TEST_LOG)
	@echo "=== 测试完成 ===" | tee -a $(TEST_LOG)
	@echo "测试日志已保存到: $(TEST_LOG)"

# 运行集成测试
.PHONY: test-integration
test-integration:
	@echo "=== 开始运行集成测试 ==="
	@echo "集成测试时间: $$(date '+%Y-%m-%d %H:%M:%S')" > $(TEST_LOG)
	@echo "" >> $(TEST_LOG)
	@cd $(TEST_DIR) && stdbuf -oL -eL go test ../test/integration/... \
		-v \
		-coverprofile=coverage_integration.out \
		-coverpkg=$(COVERAGE_PKG) \
		-timeout $(TIMEOUT) \
		2>&1 | tee -a $(TEST_LOG)
	@echo "" >> $(TEST_LOG)
	@echo "=== 集成测试完成 ===" | tee -a $(TEST_LOG)


# 生成覆盖率报告
.PHONY: coverage
coverage:
	@if [ ! -f $(COVERAGE_FILE) ]; then \
		echo "错误: 覆盖率文件不存在，请先运行 'make test'"; \
		exit 1; \
	fi
	@echo "=== 生成覆盖率报告 ==="
	@echo "生成时间: $$(date '+%Y-%m-%d %H:%M:%S')" >> $(TEST_LOG)
	@echo "" >> $(TEST_LOG)
	@echo "正在生成覆盖率报告..."
	@go tool cover -func=$(COVERAGE_FILE) | \
		grep "go-mcp-context/internal/service" | \
		awk '{gsub(/.*\//, "", $$1); printf "%-6s %-50s %s\n", $$3, $$1, $$2}' | \
		sort -n > $(COVERAGE_REPORT)
	@echo "" >> $(COVERAGE_REPORT)
	@echo "========================================" >> $(COVERAGE_REPORT)
	@go tool cover -func=$(COVERAGE_FILE) | tail -1 >> $(COVERAGE_REPORT)
	@echo "覆盖率报告已生成: $(COVERAGE_REPORT)"
	@echo "" >> $(TEST_LOG)
	@echo "=== 覆盖率统计 ===" | tee -a $(TEST_LOG)
	@go tool cover -func=$(COVERAGE_FILE) | tail -1 | tee -a $(TEST_LOG)
	@echo "" | tee -a $(TEST_LOG)
	@echo "覆盖率为0的函数: $$(grep '^0.0%' $(COVERAGE_REPORT) | wc -l) 个" | tee -a $(TEST_LOG)
	@echo "" | tee -a $(TEST_LOG)
	@echo "前20个最低覆盖率函数:" | tee -a $(TEST_LOG)
	@head -20 $(COVERAGE_REPORT) | tee -a $(TEST_LOG)

# 生成 HTML 覆盖率报告
.PHONY: coverage-html
coverage-html: coverage
	@echo "=== 生成 HTML 覆盖率报告 ==="
	@go tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "HTML 覆盖率报告已生成: coverage.html"
	@echo "在浏览器中打开: file://$(shell pwd)/coverage.html"

# 查看覆盖率报告
.PHONY: show-coverage
show-coverage:
	@if [ ! -f $(COVERAGE_REPORT) ]; then \
		echo "错误: 覆盖率报告不存在，请先运行 'make coverage'"; \
		exit 1; \
	fi
	@echo "=== 覆盖率报告 ==="
	@echo ""
	@go tool cover -func=$(COVERAGE_FILE) | tail -1
	@echo ""
	@echo "覆盖率为0的函数: $$(grep '^0.0%' $(COVERAGE_REPORT) | wc -l) 个"
	@echo ""
	@echo "所有函数覆盖率（按覆盖率从低到高排序）:"
	@cat $(COVERAGE_REPORT)

# 查看测试日志
.PHONY: show-log
show-log:
	@if [ ! -f $(TEST_LOG) ]; then \
		echo "错误: 测试日志不存在，请先运行测试"; \
		exit 1; \
	fi
	@cat $(TEST_LOG)

# 查看失败的测试
.PHONY: show-failures
show-failures:
	@if [ ! -f $(TEST_LOG) ]; then \
		echo "错误: 测试日志不存在，请先运行测试"; \
		exit 1; \
	fi
	@echo "=== 失败的测试 ==="
	@grep -E "(FAIL|panic|error)" $(TEST_LOG) || echo "没有发现失败的测试"

# 快速测试（不生成覆盖率）
.PHONY: test-quick
test-quick:
	@echo "=== 快速测试（无覆盖率） ==="
ifdef TEST
	@cd $(PROJECT_ROOT) && go test ./test/unit/... -v -run $(TEST) -timeout $(TIMEOUT)
else
	@cd $(PROJECT_ROOT) && go test ./test/unit/... -v -timeout $(TIMEOUT)
endif

# 测试特定文件
.PHONY: test-file
test-file:
ifndef FILE
	@echo "错误: 请指定文件，例如: make test-file FILE=processor"
	@exit 1
endif
	@echo "=== 测试文件: $(FILE) ==="
	@cd $(TEST_DIR) && go test ./unit/$(FILE)_test.go ./unit/setup_test.go -v -timeout $(TIMEOUT)

# 清理生成的文件
.PHONY: clean
clean:
	@echo "=== 清理生成的文件 ==="
	@rm -f $(COVERAGE_FILE)
	@rm -f coverage_integration.out
	@rm -f coverage.html
	@rm -f $(TEST_LOG)
	@rm -f $(COVERAGE_REPORT)
	@echo "清理完成"

# 清理覆盖率报告（保留测试日志）
.PHONY: clean-coverage
clean-coverage:
	@echo "=== 清理覆盖率文件 ==="
	@rm -f $(COVERAGE_FILE)
	@rm -f coverage_integration.out
	@rm -f coverage.html
	@rm -f $(COVERAGE_REPORT)
	@echo "清理完成"

# 帮助信息
.PHONY: help
help:
	@echo "╔════════════════════════════════════════════════════════════════╗"
	@echo "║           测试 Makefile 使用指南                                ║"
	@echo "╚════════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "⚡ 一键命令（推荐）:"
	@echo "  ┌─────────────────────────────────────────────────────────────┐"
	@echo "  │ make all               清理 + 单元测试 + 覆盖率              │"
	@echo "  └─────────────────────────────────────────────────────────────┘"
	@echo ""
	@echo "📋 测试命令（按推荐顺序）:"
	@echo "  ┌─────────────────────────────────────────────────────────────┐"
	@echo "  │ make test              运行所有测试（单元 + 集成）           │"
	@echo "  │ make test-unit         仅运行单元测试                        │"
	@echo "  │ make test-integration  仅运行集成测试（需要外部服务）        │"
	@echo "  └─────────────────────────────────────────────────────────────┘"
	@echo ""
	@echo "🎯 指定测试:"
	@echo "  make test-unit TEST=TestName          运行指定的单元测试"
	@echo "  make test-unit TEST=Test_Processor_*  运行匹配的测试"
	@echo "  make test-file FILE=processor         测试指定文件（processor_test.go）"
	@echo "  make test-quick                       快速测试（不生成覆盖率）"
	@echo ""
	@echo "📊 覆盖率报告:"
	@echo "  make coverage          生成覆盖率报告（all_functions_coverage.txt）"
	@echo "  make coverage-html     生成 HTML 覆盖率报告（可在浏览器查看）"
	@echo "  make show-coverage     查看覆盖率统计和函数列表"
	@echo ""
	@echo "📝 日志查看:"
	@echo "  make show-log          查看完整测试日志（test_log.txt）"
	@echo "  make show-failures     仅查看失败的测试"
	@echo ""
	@echo "🧹 清理:"
	@echo "  make clean             清理所有生成的文件"
	@echo "  make clean-coverage    仅清理覆盖率文件（保留日志）"
	@echo ""
	@echo "💡 常用工作流:"
	@echo "  ┌─────────────────────────────────────────────────────────────┐"
	@echo "  │ 1. 日常开发（最常用）:                                       │"
	@echo "  │    make all                                                 │"
	@echo "  │    # 自动清理 + 运行单元测试 + 生成覆盖率报告                │"
	@echo "  │                                                              │"
	@echo "  │ 2. 调试失败的测试:                                           │"
	@echo "  │    make test-unit TEST=TestName                             │"
	@echo "  │    make show-failures                                       │"
	@echo "  │                                                              │"
	@echo "  │ 3. 提交前完整测试:                                           │"
	@echo "  │    make clean && make test && make coverage                 │"
	@echo "  │    # 运行所有测试（包括集成测试）                            │"
	@echo "  │                                                              │"
	@echo "  │ 4. 查看详细覆盖率:                                           │"
	@echo "  │    make coverage-html                                       │"
	@echo "  │    # 在浏览器打开 test/coverage.html                        │"
	@echo "  └─────────────────────────────────────────────────────────────┘"
	@echo ""
	@echo "📁 生成的文件:"
	@echo "  coverage.out                  Go 覆盖率原始数据"
	@echo "  test_log.txt                  测试详细日志（包含所有输出）"
	@echo "  all_functions_coverage.txt    函数覆盖率列表（按覆盖率排序）"
	@echo "  coverage.html                 HTML 覆盖率报告（可选）"
	@echo ""
	@echo "❓ 更多帮助: 查看 test/README.md"
