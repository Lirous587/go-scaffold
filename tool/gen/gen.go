package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var (
		model   = flag.String("m", "", "模块名称")
		version = flag.String("v", "1", "版本号（可选，默认1）")
		help    = flag.Bool("h", false, "显示帮助")
	)

	flag.Usage = func() {
		fmt.Println("用法: ./gen -m <model> [-v <version>]")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *model == "" {
		flag.Usage()
		os.Exit(1)
	}

	domainLower := strings.ToLower(*model)
	domainTitle := cases.Title(language.Und).String(domainLower)

	module := ""
	{
		cwd, _ := os.Getwd()
		goModPath, err := findGoModPath(cwd)
		if err != nil {
			fmt.Println("未找到 go.mod 文件")
			os.Exit(1)
		}
		module, err = getModuleName(goModPath)
		if err != nil {
			fmt.Println("无法解析 module 名称")
			os.Exit(1)
		}
	}

	data := map[string]string{
		"Domain":      domainLower,
		"DomainTitle": domainTitle,
		"Module":      module,
	}

	tempDir := filepath.Join("template_v" + *version)
	outBase := filepath.Join("../..", "internal", domainLower)

	// 脚本是否创建了目录

	if _, err := os.Stat(outBase); os.IsNotExist(err) {
		if err := os.MkdirAll(outBase, 0755); err != nil {
			fmt.Printf("创建目录失败:%v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("domain %s 已存在，跳过生成\n", domainLower)
		return
	}

	err := filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".tmpl") {
			return nil
		}
		fmt.Println(path)
		relPath, _ := filepath.Rel(tempDir, path)
		outPath := filepath.Join(outBase, strings.TrimSuffix(relPath, ".tmpl")) + ".go"
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}

		tmpl, err := template.ParseFiles(path)
		if err != nil {
			return err
		}
		f, err := os.Create(outPath)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := tmpl.Execute(f, data); err != nil {
			return err
		}
		fmt.Printf("生成成功:%s\n", outPath)
		return nil
	})

	if err != nil {
		fmt.Printf("生成失败:%v\n", err)
		os.RemoveAll(outBase)
		os.Exit(1)
	}
}

// 查找最近的 go.mod 文件路径
func findGoModPath(startDir string) (string, error) {
	dir := startDir
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return goModPath, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", os.ErrNotExist
}

// 解析 go.mod 文件获取 module 名称
func getModuleName(goModPath string) (string, error) {
	f, err := os.Open(goModPath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module")), nil
		}
	}
	return "", os.ErrNotExist
}
