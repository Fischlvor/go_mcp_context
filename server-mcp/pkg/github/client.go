package github

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"go-mcp-context/pkg/global"
)

// Client GitHub API 客户端
type Client struct {
	httpClient *http.Client
	token      string
}

// NewClient 创建 GitHub 客户端
func NewClient() *Client {
	// 优先从环境变量读取 token，其次从配置文件
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		token = global.Config.GitHub.Token
	}

	// 创建 HTTP 客户端（支持代理）
	transport := &http.Transport{}

	// 使用配置文件中的代理（为空则不使用代理）
	if proxyURL := global.Config.GitHub.Proxy; proxyURL != "" {
		if proxy, err := url.Parse(proxyURL); err == nil {
			transport.Proxy = http.ProxyURL(proxy)
		}
	}

	return &Client{
		httpClient: &http.Client{
			Timeout:   60 * time.Second, // 下载文件可能需要更长时间
			Transport: transport,
		},
		token: token,
	}
}

// ==========================================
// API 请求
// ==========================================

// doRequest 执行 GitHub API 请求
func (c *Client) doRequest(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error: %d - %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

// ==========================================
// 仓库信息
// ==========================================

// GetRepoInfo 获取仓库信息
func (c *Client) GetRepoInfo(ctx context.Context, repo string) (*Repo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s", repo)
	data, err := c.doRequest(ctx, url)
	if err != nil {
		return nil, err
	}

	var repoInfo Repo
	if err := json.Unmarshal(data, &repoInfo); err != nil {
		return nil, err
	}

	return &repoInfo, nil
}

// ==========================================
// 版本获取（智能分页）
// ==========================================

// GetReleases 获取所有正式发布版本（智能分页）
func (c *Client) GetReleases(ctx context.Context, repo string) ([]string, error) {
	baseURL := fmt.Sprintf("https://api.github.com/repos/%s/releases?per_page=100", repo)

	// 先获取总页数（从 Link 响应头）
	req, _ := http.NewRequestWithContext(ctx, "HEAD", baseURL, nil)
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	// 解析最后一页
	lastPage := 1
	if linkHeader := resp.Header.Get("Link"); linkHeader != "" {
		re := regexp.MustCompile(`page=(\d+)>; rel="last"`)
		if matches := re.FindStringSubmatch(linkHeader); len(matches) > 1 {
			fmt.Sscanf(matches[1], "%d", &lastPage)
		}
	}

	// 并行获取所有页面
	type pageResult struct {
		page     int
		versions []string
		err      error
	}
	results := make(chan pageResult, lastPage)

	for page := 1; page <= lastPage; page++ {
		go func(p int) {
			url := fmt.Sprintf("%s&page=%d", baseURL, p)
			data, err := c.doRequest(ctx, url)
			if err != nil {
				results <- pageResult{page: p, err: err}
				return
			}

			var releases []Release
			if err := json.Unmarshal(data, &releases); err != nil {
				results <- pageResult{page: p, err: err}
				return
			}

			// 过滤正式版本
			var versions []string
			versionRe := regexp.MustCompile(`^v?\d+\.\d+(\.\d+)?$`)
			for _, r := range releases {
				if !r.Prerelease && !r.Draft && versionRe.MatchString(r.TagName) {
					versions = append(versions, r.TagName)
				}
			}
			results <- pageResult{page: p, versions: versions}
		}(page)
	}

	// 收集结果
	var allVersions []string
	for i := 0; i < lastPage; i++ {
		result := <-results
		if result.err != nil {
			return nil, result.err
		}
		allVersions = append(allVersions, result.versions...)
	}

	// 按版本号排序（倒序）
	sort.Slice(allVersions, func(i, j int) bool {
		return compareVersions(allVersions[i], allVersions[j]) > 0
	})

	return allVersions, nil
}

// GetMajorVersions 获取每个大版本的最新版本
// 如果只有一个大版本，则改为按 minor 版本分组
func (c *Client) GetMajorVersions(ctx context.Context, repo string, maxCount int) ([]string, error) {
	allVersions, err := c.GetReleases(ctx, repo)
	if err != nil {
		return nil, err
	}

	if len(allVersions) == 0 {
		return nil, nil
	}

	// 先统计有多少个大版本
	majorSet := make(map[string]bool)
	for _, v := range allVersions {
		majorSet[extractMajorVersion(v)] = true
	}
	uniqueMajors := len(majorSet)

	// 根据大版本数量决定分组策略
	var versionMap map[string]string
	if uniqueMajors == 1 {
		// 只有一个大版本，按 minor 版本分组
		versionMap = make(map[string]string)
		for _, v := range allVersions {
			minor := extractMinorVersion(v)
			if _, exists := versionMap[minor]; !exists {
				versionMap[minor] = v
			}
		}
	} else {
		// 多个大版本，按 major 版本分组
		versionMap = make(map[string]string)
		for _, v := range allVersions {
			major := extractMajorVersion(v)
			if _, exists := versionMap[major]; !exists {
				versionMap[major] = v
			}
		}
	}

	// 转为数组并排序
	var keys []string
	for key := range versionMap {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return compareVersions(keys[i], keys[j]) > 0
	})

	// 限制数量
	if maxCount > 0 && len(keys) > maxCount {
		keys = keys[:maxCount]
	}

	// 返回每个版本组的最新版本
	var result []string
	for _, key := range keys {
		result = append(result, versionMap[key])
	}

	return result, nil
}

// ==========================================
// 目录树
// ==========================================

// GetTree 获取目录树
func (c *Client) GetTree(ctx context.Context, repo, ref string) (*Tree, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/git/trees/%s?recursive=1", repo, ref)
	data, err := c.doRequest(ctx, url)
	if err != nil {
		return nil, err
	}

	var tree Tree
	if err := json.Unmarshal(data, &tree); err != nil {
		return nil, err
	}

	return &tree, nil
}

// FilterTree 过滤目录树
func (c *Client) FilterTree(tree *Tree, pathFilter string, excludes []string) []TreeItem {
	var filtered []TreeItem

	// 支持的文件扩展名（只支持 Markdown）
	docExtensions := map[string]bool{
		".md": true, ".mdx": true,
	}

	// 默认排除目录（参考 Context7 规则）
	// 排除：.github, test(s), dist, node_modules, vendor, fixtures, bench
	//       archive/archived/deprecated/legacy/old/outdated
	//       i18n 非英语目录
	excludeDirPatterns := []string{
		".github", "node_modules", "vendor", "dist",
		"test", "tests", "__tests__", "fixtures", "bench", "benchmark", "benchmarks",
		"archive", "archived", "deprecated", "legacy", "outdated",
		"i18n", "zh-cn", "zh-tw", "zh-hk", "zh_cn", "zh_tw",
	}

	// 排除文件名
	excludeFileNames := map[string]bool{
		"CHANGELOG.md": true, "CHANGELOG.mdx": true,
		"LICENSE.md": true, "LICENSE.mdx": true,
		"CODE_OF_CONDUCT.md": true, "CODE_OF_CONDUCT.mdx": true,
		"CONTRIBUTING.md": true, // 保留可能有用，但先排除
	}

	// 合并用户自定义排除
	allExcludeDirs := append(excludeDirPatterns, excludes...)

	for _, item := range tree.Tree {
		// 只处理文件
		if item.Type != "blob" {
			continue
		}

		// 路径过滤
		if pathFilter != "" && !strings.HasPrefix(item.Path, pathFilter) {
			continue
		}

		// 扩展名检查
		ext := strings.ToLower(filepath.Ext(item.Path))
		if !docExtensions[ext] {
			continue
		}

		// 排除文件名检查
		fileName := filepath.Base(item.Path)
		if excludeFileNames[fileName] {
			continue
		}

		// 排除目录检查（路径中包含排除模式）
		excluded := false
		pathLower := strings.ToLower(item.Path)
		for _, pattern := range allExcludeDirs {
			// 检查路径段是否匹配（避免误匹配，如 "test" 不应该匹配 "contest"）
			patternLower := strings.ToLower(pattern)
			if strings.Contains(pathLower, "/"+patternLower+"/") ||
				strings.HasPrefix(pathLower, patternLower+"/") ||
				strings.Contains(pathLower, "/"+patternLower) && strings.HasSuffix(pathLower, patternLower) {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}

		filtered = append(filtered, item)
	}

	return filtered
}

// ==========================================
// 文件下载
// ==========================================

// DownloadFile 下载单个文件
func (c *Client) DownloadFile(ctx context.Context, repo, ref, path string) ([]byte, error) {
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s", repo, ref, path)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// TarballFile 表示从 tarball 中提取的文件
type TarballFile struct {
	Path    string
	Content []byte
	Size    int64
}

// DownloadTarballFiles 下载 tarball 并流式提取指定文件
// 返回一个 channel，每提取一个文件就发送一次
func (c *Client) DownloadTarballFiles(ctx context.Context, repo, ref string, filter func(path string) bool) (<-chan TarballFile, <-chan error) {
	fileChan := make(chan TarballFile, 10)
	errChan := make(chan error, 1)

	go func() {
		defer close(fileChan)
		defer close(errChan)

		// tarball URL
		url := fmt.Sprintf("https://github.com/%s/archive/refs/tags/%s.tar.gz", repo, ref)

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			errChan <- err
			return
		}

		if c.token != "" {
			req.Header.Set("Authorization", "Bearer "+c.token)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			errChan <- err
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			errChan <- fmt.Errorf("failed to download tarball: %d", resp.StatusCode)
			return
		}

		// 流式解压 gzip
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			errChan <- fmt.Errorf("failed to create gzip reader: %w", err)
			return
		}
		defer gzReader.Close()

		// 流式读取 tar
		tarReader := tar.NewReader(gzReader)

		for {
			header, err := tarReader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				errChan <- fmt.Errorf("failed to read tar: %w", err)
				return
			}

			// 跳过目录
			if header.Typeflag != tar.TypeReg {
				continue
			}

			// 去掉第一层目录（repo-tag/）
			path := header.Name
			if idx := strings.Index(path, "/"); idx != -1 {
				path = path[idx+1:]
			}

			// 应用过滤器
			if !filter(path) {
				continue
			}

			// 读取文件内容（只有匹配的文件才读取）
			content, err := io.ReadAll(tarReader)
			if err != nil {
				continue // 跳过读取失败的文件
			}

			fileChan <- TarballFile{
				Path:    path,
				Content: content,
				Size:    header.Size,
			}
		}
	}()

	return fileChan, errChan
}

// ==========================================
// 辅助函数
// ==========================================

// extractMajorVersion 提取大版本号 (v1.x.x -> v1)
func extractMajorVersion(version string) string {
	version = strings.TrimPrefix(version, "v")
	parts := strings.Split(version, ".")
	if len(parts) > 0 {
		return "v" + parts[0]
	}
	return version
}

// extractMinorVersion 提取次版本号 (v1.10.x -> v1.10)
func extractMinorVersion(version string) string {
	version = strings.TrimPrefix(version, "v")
	parts := strings.Split(version, ".")
	if len(parts) >= 2 {
		return "v" + parts[0] + "." + parts[1]
	}
	if len(parts) == 1 {
		return "v" + parts[0]
	}
	return version
}

// compareVersions 比较版本号（返回 >0 表示 a > b）
func compareVersions(a, b string) int {
	a = strings.TrimPrefix(a, "v")
	b = strings.TrimPrefix(b, "v")

	partsA := strings.Split(a, ".")
	partsB := strings.Split(b, ".")

	for i := 0; i < len(partsA) && i < len(partsB); i++ {
		var numA, numB int
		fmt.Sscanf(partsA[i], "%d", &numA)
		fmt.Sscanf(partsB[i], "%d", &numB)
		if numA != numB {
			return numA - numB
		}
	}

	return len(partsA) - len(partsB)
}
