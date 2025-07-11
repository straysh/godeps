# Go Dependencies Analyzer (godeps)

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.18-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

[中文版本](README_CN.md) | English

A command-line tool for analyzing Go module dependencies. It provides a clear, tree-structured view of your project's dependency relationships, including both direct and indirect dependencies.

## ✨ Features

- **Complete Dependency Tree**: Display full dependency chains with tree-structured output
- **Reverse Dependency Lookup**: Find which packages depend on a specific module
- **Indirect Dependencies Marking**: Highlight indirect dependencies with `(Indirect)` marker
- **Depth Control**: Limit dependency tree depth for better readability
- **Clean Output**: Beautiful tree-structured visualization with proper indentation

## 📦 Installation

### From Source

```bash
git clone https://github.com/straysh/godeps.git
cd godeps
go build -o godeps main.go
```

### Direct Installation

```bash
go install github.com/straysh/godeps@latest
```

## 🚀 Usage

### Basic Usage

```bash
# Show complete dependency tree for current project
./godeps

# Show dependency tree for specific project
./godeps --path=/path/to/your/project

# Search for specific package dependencies
./godeps --path=/path/to/your/project --search=github.com/gin-gonic/gin

# Show dependencies with indirect marking
./godeps --path=/path/to/your/project --color

# Limit dependency tree depth
./godeps --path=/path/to/your/project --depth=2
```

## 📋 Command Line Options

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--path` | string | `./` | Path to the Go project directory |
| `--search` | string | `""` | Search for dependencies of specific package |
| `--color` | bool | `false` | Add `(Indirect)` marker for indirect dependencies |
| `--depth` | int | `0` | Maximum depth of dependency tree (0 = unlimited, only works when search is empty) |

## 📖 Examples

### 1. Show Complete Dependency Tree

```bash
./godeps --path=/path/to/project
```

**Output:**
```
project package_name: github.com/example/app
├── github.com/gorilla/mux@v1.8.0
│   └── github.com/gorilla/context@v1.1.1
└── golang.org/x/text@v0.3.2
    └── golang.org/x/tools@v0.0.0-20180917221912-90fa682c2a6e
```

### 2. Search for Specific Package

```bash
./godeps --path=/path/to/project --search=github.com/gin-gonic/gin
```

**Output:**
```
# github.com/gin-gonic/gin 包的依赖链路
github.com/gin-gonic/gin@v1.9.1
├── github.com/gin-contrib/sse@v0.1.0
├── github.com/go-playground/validator/v10@v10.14.0
└── github.com/json-iterator/go@v1.1.12

# 依赖 github.com/gin-gonic/gin 包的链路
github.com/example/app
└── github.com/gin-gonic/gin@v1.9.1
```

### 3. Show Indirect Dependencies

```bash
./godeps --path=/path/to/project --color
```

**Output:**
```
project package_name: github.com/example/app
├── github.com/gorilla/mux@v1.8.0
│   └── github.com/gorilla/context@v1.1.1 (Indirect)
└── golang.org/x/text@v0.3.2 (Indirect)
    └── golang.org/x/tools@v0.0.0-20180917221912-90fa682c2a6e (Indirect)
```

### 4. Limit Dependency Depth

```bash
./godeps --path=/path/to/project --depth=1
```

**Output:**
```
project package_name: github.com/example/app
├── github.com/gorilla/mux@v1.8.0
└── golang.org/x/text@v0.3.2
```

## 🔧 How It Works

This tool leverages Go's built-in module system commands:

- **All Dependencies**: `go mod graph`
- **Direct Dependencies**: `go list -mod=readonly -m -f '{{if not .Indirect}}{{.Path}}{{end}}' all`
- **Indirect Dependencies**: `go list -mod=readonly -m -f '{{if .Indirect}}{{.Path}}{{end}}' all`

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with the power of Go's module system
- Inspired by the need for better dependency visualization tools

---

**Happy Dependency Analyzing!** 🎉

