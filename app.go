package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

// App struct
type App struct {
	ctx context.Context
}

// 定义一个结构体，方便前端使用
type BrewPackage struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Status  string `json:"status"` // "started", "stopped", 或 "none" (不是服务)
}

type BrewData struct {
	Formulae []BrewPackage `json:"formulae"`
	Casks    []BrewPackage `json:"casks"`
}
type ServiceInfo struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// 操作结果返回
type ActionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func getBrewPath() string {
	// 检查 M1/M2 路径
	if _, err := os.Stat("/opt/homebrew/bin/brew"); err == nil {
		return "/opt/homebrew/bin/brew"
	}
	// 检查 Intel 路径
	if _, err := os.Stat("/usr/local/bin/brew"); err == nil {
		return "/usr/local/bin/brew"
	}
	return "brew" // 保底方案
}

func (a *App) GetBrewData() BrewData {
	// 1. 获取所有服务状态
	services := make(map[string]string)
	serviceRaw, _ := exec.Command(getBrewPath(), "services", "info", "--all", "--json").Output()
	var serviceList []ServiceInfo
	json.Unmarshal(serviceRaw, &serviceList)
	for _, s := range serviceList {
		services[s.Name] = s.Status
	}

	// 2. 获取并包装数据
	return BrewData{
		Formulae: fetchWithStatus("--formula", services),
		Casks:    fetchWithStatus("--cask", services),
	}
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 重点：为 GUI 进程手动注入 PATH，否则执行 brew 命令会卡死或直接报错找不到命令
	path := os.Getenv("PATH")
	// 兼容 M1/M2/M3 和 Intel Mac，同时保留原有路径
	newPath := "/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:" + path
	os.Setenv("PATH", newPath)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// brew list 命令并返回字符串数组
func fetchWithStatus(flag string, serviceMap map[string]string) []BrewPackage {
	out, _ := exec.Command(getBrewPath(), "list", "--versions", flag).Output()
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var packages []BrewPackage

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			name := parts[0]
			version := strings.Join(parts[1:], " ")
			status, isService := serviceMap[name]
			if !isService {
				status = "none_tool" // 普通工具，没有服务状态，服务停止会显示 none
			}

			packages = append(packages, BrewPackage{
				Name:    name,
				Version: version,
				Status:  status,
			})
		}
	}
	// --- 排序逻辑开始 ---
	sort.Slice(packages, func(i, j int) bool {
		// 定义状态权重
		priority := map[string]int{
			"started":   1,
			"stopped":   2,
			"none":      3,
			"none_tool": 4,
		}

		// 如果状态不同，按优先级排（1最小，排最前）
		if packages[i].Status != packages[j].Status {
			return priority[packages[i].Status] < priority[packages[j].Status]
		}
		// 如果状态相同，按名称字母顺序排
		return packages[i].Name < packages[j].Name
	})
	// --- 排序逻辑结束 ---
	return packages
}

// StartService 启动指定的 Brew 服务
func (a *App) StartService(name string) ActionResponse {
	// 执行命令: brew services start <name>
	out, err := exec.Command(getBrewPath(), "services", "start", name).CombinedOutput()
	if err != nil {
		// 这样即使失败，也会把错误信息返回给前端的 alert
		return ActionResponse{
			Success: false,
			Message: fmt.Sprintf("后端执行失败: %v, 输出: %s", err, string(out)),
		}
	}
	return ActionResponse{Success: true, Message: "服务 " + name + "已启动"}
}

// StopService 停止指定的 Brew 服务
func (a *App) StopService(name string) ActionResponse {
	// 执行命令: brew services stop <name>
	out, err := exec.Command(getBrewPath(), "services", "stop", name).CombinedOutput()
	if err != nil {
		return ActionResponse{
			Success: false,
			Message: fmt.Sprintf("停止失败(%v): %s", err, string(out)),
		}
		// return ActionResponse{Success: false, Message: "停止失败：" + err.Error()}
	}
	return ActionResponse{Success: true, Message: "服务 " + name + "已停止"}
}

// GetAppIcon 获取应用的图标并转为 Base64
func (a *App) GetAppIcon(appName string) string {
	// 1. 建立手动映射表 (解决包名和 App 名完全对不上的情况)
	mismatchedNames := map[string]string{
		"iterm2":             "iTerm",
		"docker-desktop":     "Docker",
		"google-chrome":      "Google Chrome",
		"visual-studio-code": "Visual Studio Code",
		"dbeaver-community":  "DBeaver", // 针对 DBeaver 的映射
	}

	// 2. 生成可能的 App 搜索名称
	searchNames := []string{appName}
	if realName, ok := mismatchedNames[appName]; ok {
		searchNames = append(searchNames, realName)
	}

	// 补充常规变形：首字母大写、全大写
	searchNames = append(searchNames, capitalize(appName))
	searchNames = append(searchNames, strings.ToUpper(appName))

	var appPath string
	// --- 第一阶段：尝试所有已知可能的精确路径 ---
	for _, name := range searchNames {
		path := filepath.Join("/Applications", name+".app")
		if _, err := os.Stat(path); err == nil {
			appPath = path
			break
		}
	}

	// --- 第二阶段：如果还没找到，才进行谨慎的模糊匹配 ---
	if appPath == "" {
		// 如果还是找不到，尝试在 /Applications 目录下遍历一遍，看看有没有包含关键字的
		files, _ := os.ReadDir("/Applications")
		for _, f := range files {
			// 确保找到的是 .app 结尾的文件夹，且排除干扰项
			fileName := f.Name()
			if strings.HasSuffix(fileName, ".app") &&
				strings.Contains(strings.ToLower(fileName), strings.ToLower(appName)) {
				appPath = filepath.Join("/Applications", fileName)
				break
			}
		}
	}

	if appPath == "" {
		return ""
	}

	// 3. 读取 Info.plist 定位图标 (iTerm2 这里的图标名通常是 iTerm2.icns)
	iconFileName := "AppIcon.icns" // 默认
	plistPath := filepath.Join(appPath, "Contents", "Info.plist")
	if data, err := os.ReadFile(plistPath); err == nil {
		content := string(data)
		if idx := strings.Index(content, "CFBundleIconFile"); idx != -1 {
			sub := content[idx:]
			start := strings.Index(sub, "<string>") + 8
			end := strings.Index(sub, "</string>")
			if start > 7 && end > start {
				iconFileName = sub[start:end]
				if !strings.HasSuffix(iconFileName, ".icns") {
					iconFileName += ".icns"
				}
			}
		}
	}

	iconPath := filepath.Join(appPath, "Contents", "Resources", iconFileName)
	// 兜底方案：如果指定的图标不存在，扫描 Resources 下任何一个 .icns
	if _, err := os.Stat(iconPath); err != nil {
		resFiles, _ := os.ReadDir(filepath.Join(appPath, "Contents", "Resources"))
		for _, rf := range resFiles {
			if strings.HasSuffix(strings.ToLower(rf.Name()), ".icns") {
				iconPath = filepath.Join(appPath, "Contents", "Resources", rf.Name())
				break
			}
		}
	}

	// 4. 转换并返回 (带唯一 AppName 的临时文件防止图标错乱)
	tmpPng := filepath.Join(os.TempDir(), fmt.Sprintf("icon_%s.png", appName))
	// tmpPng := filepath.Join(os.TempDir(), fmt.Sprintf("brew_icon_%s.png", appName))
	defer os.Remove(tmpPng)

	cmd := exec.Command("sips", "-s", "format", "png", iconPath, "--out", tmpPng)
	if err := cmd.Run(); err != nil {
		return ""
	}

	imgData, _ := os.ReadFile(tmpPng)
	return base64.StdEncoding.EncodeToString(imgData)
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	// 将字符串转为 rune 切片以正确处理 Unicode 字符
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
