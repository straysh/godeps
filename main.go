// main.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Node 节点结构
type Node struct {
	Name     string
	Children []*Node
}

// 获取直接依赖
func getDirectDeps(workdir string) map[string]bool {
	cmd := exec.Command("go", "list", "-mod=readonly", "-m", "-f", "{{if not .Indirect}}{{.Path}}{{end}}", "all")
	cmd.Dir = workdir
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running go list for direct deps:", err)
		os.Exit(1)
	}

	directDeps := make(map[string]bool)
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			directDeps[line] = true
		}
	}
	return directDeps
}

// 获取间接依赖
func getIndirectDeps(workdir string) map[string]bool {
	cmd := exec.Command("go", "list", "-mod=readonly", "-m", "-f", "{{if .Indirect}}{{.Path}}{{end}}", "all")
	cmd.Dir = workdir
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running go list for indirect deps:", err)
		os.Exit(1)
	}

	indirectDeps := make(map[string]bool)
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			indirectDeps[line] = true
		}
	}
	return indirectDeps
}

// 解析 go mod graph 输出
func parseGraph(workdir string) (map[string][]string, map[string][]string) {
	cmd := exec.Command("go", "mod", "graph")
	cmd.Dir = workdir
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running go mod graph:", err)
		os.Exit(1)
	}

	deps := make(map[string][]string)
	reverse := make(map[string][]string)

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 2 {
			parent := parts[0]
			child := parts[1]
			deps[parent] = append(deps[parent], child)
			reverse[child] = append(reverse[child], parent)
		}
	}
	return deps, reverse
}

// 递归打印依赖链
func printDeps(node string, deps map[string][]string, prefix, childPrefix string, visited map[string]bool, color bool, indirectDeps map[string]bool, maxDepth int, currentDepth int) {
	nodeDisplay := node
	if color && indirectDeps[extractPackageName(node)] {
		nodeDisplay = node + " (Indirect)"
	}
	fmt.Println(prefix + nodeDisplay)

	if visited[node] {
		return
	}
	visited[node] = true

	// 检查深度限制
	if maxDepth > 0 && currentDepth >= maxDepth {
		return
	}

	children := deps[node]
	for i, child := range children {
		last := i == len(children)-1
		newPrefix := childPrefix + "├── "
		newChildPrefix := childPrefix + "│   "
		if last {
			newPrefix = childPrefix + "└── "
			newChildPrefix = childPrefix + "    "
		}
		printDeps(child, deps, newPrefix, newChildPrefix, visited, color, indirectDeps, maxDepth, currentDepth+1)
	}
}

// 递归打印反向依赖链
func printReverseDeps(node string, reverse map[string][]string, prefix, childPrefix string, visited map[string]bool, color bool, indirectDeps map[string]bool, maxDepth int, currentDepth int) {
	nodeDisplay := node
	if color && indirectDeps[extractPackageName(node)] {
		nodeDisplay = node + " (Indirect)"
	}
	fmt.Println(prefix + nodeDisplay)

	if visited[node] {
		return
	}
	visited[node] = true

	// 检查深度限制
	if maxDepth > 0 && currentDepth >= maxDepth {
		return
	}

	parents := reverse[node]
	for i, parent := range parents {
		last := i == len(parents)-1
		newPrefix := childPrefix + "├── "
		newChildPrefix := childPrefix + "│   "
		if last {
			newPrefix = childPrefix + "└── "
			newChildPrefix = childPrefix + "    "
		}
		printReverseDeps(parent, reverse, newPrefix, newChildPrefix, visited, color, indirectDeps, maxDepth, currentDepth+1)
	}
}

// 提取包名（去除版本号）
func extractPackageName(fullName string) string {
	parts := strings.Split(fullName, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return fullName
}

func getRootModule(workdir string) string {
	cmd := exec.Command("go", "list", "-m")
	cmd.Dir = workdir
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error running go list -m:", err)
		os.Exit(1)
	}
	return strings.TrimSpace(string(out))
}

// 查找匹配的包名（支持不带版本号的搜索）
func findMatchingPackages(searchTarget string, deps map[string][]string) []string {
	var matches []string
	for pkg := range deps {
		if extractPackageName(pkg) == searchTarget || pkg == searchTarget {
			matches = append(matches, pkg)
		}
	}
	return matches
}

// 查找匹配的包名（在reverse映射中）
func findMatchingPackagesReverse(searchTarget string, reverse map[string][]string) []string {
	var matches []string
	for pkg := range reverse {
		if extractPackageName(pkg) == searchTarget || pkg == searchTarget {
			matches = append(matches, pkg)
		}
	}
	return matches
}

func main() {
	var (
		path   = flag.String("path", "./", "path to a project")
		search = flag.String("search", "", "search for dependencies of package_name")
		color  = flag.Bool("color", false, "add (Indirect) marker for indirect dependencies")
		depth  = flag.Int("depth", 0, "maximum depth of dependency tree (0 means unlimited, only works when search is empty)")
	)
	flag.Parse()

	deps, reverse := parseGraph(*path)
	var indirectDeps map[string]bool
	if *color {
		indirectDeps = getIndirectDeps(*path)
	}

	if len(*search) == 0 {
		root := getRootModule(*path)
		fmt.Println("project package_name: ", root)
		visited := make(map[string]bool)
		printDeps(root, deps, "", "", visited, *color, indirectDeps, *depth, 0)
		return
	}

	target := *search
	fmt.Printf("# %s 包的依赖链路\n", target)

	// 查找匹配的包名（用于依赖链路）
	matches := findMatchingPackages(target, deps)
	visited := make(map[string]bool)
	for _, match := range matches {
		if len(matches) > 1 {
			fmt.Printf("## %s\n", match)
		}
		printDeps(match, deps, "", "", visited, *color, indirectDeps, 0, 0)
	}
	if len(matches) == 0 {
		fmt.Println("(没有找到依赖)")
	}

	fmt.Println()
	fmt.Printf("# 依赖 %s 包的链路\n", target)

	// 查找匹配的包名（用于反向依赖链路）
	reverseMatches := findMatchingPackagesReverse(target, reverse)
	visited = make(map[string]bool)
	for _, match := range reverseMatches {
		if len(reverseMatches) > 1 {
			fmt.Printf("## %s\n", match)
		}
		printReverseDeps(match, reverse, "", "", visited, *color, indirectDeps, 0, 0)
	}
	if len(reverseMatches) == 0 {
		fmt.Println("(没有找到反向依赖)")
	}
}
