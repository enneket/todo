# Todo App

一个基于 Go (Wails) 和 Vue 3 构建的现代、隐私优先、跨平台桌面待办事项应用。

[English Documentation](README.md)

## ✨ 特性

- **跨平台支持**: 支持 Windows, Linux 和 macOS。
- **现代 UI/UX**: 使用 Vue 3, Tailwind CSS 和 Phosphor Icons 构建。
- **隐私优先**: 所有数据通过 SQLite 本地存储。
- **强大的后端**: 基于 Go 1.24+ 和 Wails v2 开发。
- **开发友好**: 通过 Makefile 提供统一的工作流。

## 🛠 技术栈

- **后端**: Go 1.24+, Wails v2.11+, SQLite (`modernc.org/sqlite`)
- **前端**: Vue 3.5+, Pinia, Tailwind CSS 3.3+, Vite 5+, TypeScript
- **通信**: HTTP REST API (localhost:8081)

## 🚀 快速开始

### 前置要求

- **Go**: v1.24 或更高版本
- **Node.js**: v18 或更高版本
- **Wails CLI**: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`
- **Make**: 用于运行统一命令 (Windows 用户可使用 WSL 或安装 Make 工具)

### 安装

1. 克隆仓库:
   ```bash
   git clone https://github.com/yourusername/todo-app.git
   cd todo-app
   ```

2. 安装依赖:
   ```bash
   make setup
   ```

### 开发

以开发模式启动应用（支持热重载）:

```bash
make dev
```

### 构建

构建生产环境应用:

```bash
make build
```

构建产物位于 `build/bin/` 目录。

### 跨平台构建

- **Windows**: `make build-windows`
- **Linux**: `make build-linux`
- **macOS**: `make build-darwin`

## 🧪 测试与质量保证

- **运行所有检查**: `make check` (Lint + Test + Build)
- **单元测试**: `make test`
- **E2E 测试**: `make test-e2e` (需要应用正在运行)
- **代码检查 (Lint)**: `make lint`
- **代码格式化**: `make format`

## 📂 项目结构

```
├── backend/      # Go 后端代码 (Server, DB, Service)
├── frontend/     # Vue 前端代码 (Stores, Components)
├── docs/         # 项目文档
├── build/        # 构建产物
├── main.go       # 应用入口
└── Makefile      # 统一开发命令
```

## 📚 文档

更多详细文档请参考 `docs/` 目录:

- [架构设计](docs/ARCHITECTURE.md)
- [代码模式](docs/CODE_PATTERNS.md)
- [API 文档](docs/API_DOCUMENTATION.md)

## 📄 许可证

MIT License
