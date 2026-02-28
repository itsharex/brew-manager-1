# 🍺 Brew Manager

> **基于 Wails + Vue 3 构建的 macOS 原生质感 Homebrew 管理工具**

[![Wails](https://img.shields.io/badge/Built%20with-Wails-red.svg)](https://wails.io/)
[![Vue](https://img.shields.io/badge/Frontend-Vue%203-brightgreen.svg)](https://vuejs.org/)
[![Go](https://img.shields.io/badge/Backend-Go-blue.svg)](https://go.dev/)
[![Platform](https://img.shields.io/badge/Platform-macOS-lightgrey.svg)]()

`Brew Manager` 是一个为 macOS 用户打造的轻量级 Homebrew 图形界面工具。它不仅能让你一眼看清系统安装的所有 Formulae 和 Casks，还能像原生系统服务一样管理终端工具的启动与停止。

---

## ✨ 核心特性

- 🖥️ **原生视觉体验**：采用 macOS Frameless 窗口设计，完美支持 **Vibrancy（毛玻璃）** 效果。
- ⚡ **服务一键管理**：支持对 `brew services` 进行图形化操作（启动/停止）。
- 🔍 **丝滑搜索**：实时过滤海量软件清单，瞬间找到目标。
- 🔄 **自动同步**：每 10 秒自动更新软件运行状态，无需手动刷新。
- 📦 **双列表展示**：清晰区分终端工具 (Formulae) 与桌面应用 (Casks)。

### ✨ 视觉与体验 (UX & UI)
* **Native macOS Feel**: 深度集成系统毛玻璃（Vibrancy）效果，支持暗色模式。
* **Fluid Animations**: 优化的 Toast 提示与列表交互动效，反馈更自然。
* **右上角通知系统**: 采用类 macOS 通知中心的交互设计，不干扰核心操作。

### ✨ 功能特性
* **视觉进化**：自动获取 Homebrew Cask 应用图标，界面更直观。
* **原生转换引擎**：集成 macOS 系统级 `sips` 技术，实现 `.icns` 到 Web 兼容格式的零延迟转换。
* **智能映射算法**：内置名称映射与模糊匹配逻辑，解决包名（如 `iterm2`）与应用名（如 `iTerm.app`）不一致的问题。
* **高性能加载**：采用异步并发提取机制，并配合进程级隔离的临时缓存，确保图标加载不卡顿、不错位。

---

## 📸 界面预览

| 概览 | 搜索与操作 |
| :--- | :--- |
| ![Main Window](https://cdn.jsdelivr.net/gh/mycherish/imgCloud/img/20260301002735681.png) | ![Search](https://cdn.jsdelivr.net/gh/mycherish/imgCloud/img/20260301002838262.png) |


---

## 🛠️ 技术栈

| 维度 | 技术 |
| :--- | :--- |
| **框架** | [Wails v2](https://wails.io/) (Go + Webview) |
| **前端** | Vue 3 (Composition API) + Vite |
| **样式** | CSS3 Glassmorphism + Native macOS Titlebar Inset |
| **后端** | Go (executing brew commands) |

---

## 🚀 快速开始

### 1. 前置要求
确保你的 Mac 已安装：
- [Go](https://go.dev/dl/) (1.20+)
- [Node.js](https://nodejs.org/) & [NPM](https://www.npmjs.com/)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation)
- [Homebrew](https://brew.sh/)

### 2. 开发模式
```bash
# 克隆仓库
git clone https://github.com/mycherish/brew-manager.git

# 进入目录
cd brew-manager

# 启动开发服务器
wails dev
```

## 💾 安装与使用 (Installation)

1. 前往 [Releases](https://github.com/mycherish/brew-manager/releases) 页面下载最新的 `.dmg` 文件。
2. 打开 `.dmg` 并将 `Brew-Manager` 拖入 **Applications** 文件夹。

### ⚠️ 解决“无法验证开发者”问题
由于本应用未进行 Apple 开发者签名，首次打开时可能会提示“无法验证开发者”或“应用已损坏”。请执行以下操作：

1. **右键点击**应用程序文件夹中的 `Brew-Manager`，选择 **“打开”**，在弹出的对话框中再次点击 **“打开”**。
2. 如果依然无法运行，请打开终端执行以下命令：
   ```bash
   sudo xattr -rd com.apple.quarantine /Applications/Brew-Manager.app
   ```
3. 输入开机密码后即可正常使用。