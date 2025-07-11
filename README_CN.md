# Go 依赖分析工具 (godeps)

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.18-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

中文版本 | [English](README.md)

一个用于分析 Go 模块依赖关系的命令行工具。它提供清晰的树状结构视图来展示项目的依赖关系，包括直接依赖和间接依赖。

## ✨ 功能特性

- **完整依赖树**: 以树状结构显示完整的依赖链
- **反向依赖查找**: 查找哪些包依赖于特定模块
- **间接依赖标记**: 使用 `(Indirect)` 标记突出显示间接依赖
- **深度控制**: 限制依赖树的深度以提高可读性
- **清晰输出**: 美观的树状结构可视化，具有适当的缩进

## 📦 安装

### 从源码编译

```bash
git clone https://github.com/straysh/godeps.git
cd godeps
go build -o godeps main.go
```

### 直接安装

```bash
go install github.com/straysh/godeps@latest
```

## 🚀 使用方法

### 基本用法

```bash
# 显示当前项目的完整依赖树
./godeps

# 显示指定项目的依赖树
./godeps --path=/path/to/your/project

# 搜索特定包的依赖关系
./godeps --path=/path/to/your/project --search=github.com/gin-gonic/gin

# 显示带有间接依赖标记的依赖树
./godeps --path=/path/to/your/project --color

# 限制依赖树的深度
./godeps --path=/path/to/your/project --depth=2
```

## 📋 命令行选项

| 参数 | 类型 | 默认值 | 描述 |
|------|------|---------|-------------|
| `--path` | string | `./` | Go 项目目录的路径 |
| `--search` | string | `""` | 搜索特定包的依赖关系 |
| `--color` | bool | `false` | 为间接依赖添加 `(Indirect)` 标记 |
| `--depth` | int | `0` | 依赖树的最大深度 (0 = 无限制，仅在 search 为空时有效) |

## 📖 使用示例

### 1. 显示完整依赖树

```bash
./godeps --path=/path/to/project
```

**输出:**
```
project package_name: github.com/example/app
├── github.com/gorilla/mux@v1.8.0
│   └── github.com/gorilla/context@v1.1.1
└── golang.org/x/text@v0.3.2
    └── golang.org/x/tools@v0.0.0-20180917221912-90fa682c2a6e
```

### 2. 搜索特定包

```bash
./godeps --path=/path/to/project --search=github.com/gin-gonic/gin
```

**输出:**
```
# github.com/gin-gonic/gin包名 的依赖链路
github.com/gin-gonic/gin@v1.9.1
├── github.com/gin-contrib/sse@v0.1.0
├── github.com/go-playground/validator/v10@v10.14.0
└── github.com/json-iterator/go@v1.1.12

# 依赖 github.com/gin-gonic/gin包名 的链路
github.com/example/app
└── github.com/gin-gonic/gin@v1.9.1
```

### 3. 显示间接依赖

```bash
./godeps --path=/path/to/project --color
```

**输出:**
```
project package_name: github.com/example/app
├── github.com/gorilla/mux@v1.8.0
│   └── github.com/gorilla/context@v1.1.1 (Indirect)
└── golang.org/x/text@v0.3.2 (Indirect)
    └── golang.org/x/tools@v0.0.0-20180917221912-90fa682c2a6e (Indirect)
```

### 4. 限制依赖深度

```bash
./godeps --path=/path/to/project --depth=1
```

**输出:**
```
project package_name: github.com/example/app
├── github.com/gorilla/mux@v1.8.0
└── golang.org/x/text@v0.3.2
```

## 🔧 工作原理

此工具利用 Go 的内置模块系统命令:

- **所有依赖**: `go mod graph`
- **直接依赖**: `go list -mod=readonly -m -f '{{if not .Indirect}}{{.Path}}{{end}}' all`
- **间接依赖**: `go list -mod=readonly -m -f '{{if .Indirect}}{{.Path}}{{end}}' all`

## 🤝 贡献

欢迎贡献！请随时提交 Pull Request。对于重大更改，请先打开 issue 讨论您想要更改的内容。

1. Fork 仓库
2. 创建您的功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交您的更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开 Pull Request

## 📝 许可证

此项目使用 MIT 许可证 - 详情请参阅 [LICENSE](LICENSE) 文件。

## 🙏 致谢

- 基于 Go 模块系统的强大功能构建
- 受到对更好的依赖可视化工具需求的启发

---

**愉快地分析依赖关系！** 🎉 