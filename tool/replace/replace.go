package main

import (
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const oldModule = "scaffold"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: ./replace.exe <newModule>")
		os.Exit(1)
	}

	if len(os.Args) > 2 {
		fmt.Println("传递参数过多")
		os.Exit(1)
	}

	newModule := os.Args[1]

	replaceGoMod("../../go.mod", oldModule, newModule)

	if err := filepath.Walk("../../", func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".go") {
			return nil
		}
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil
		}
		changed := false
		for _, imp := range node.Imports {
			val := strings.Trim(imp.Path.Value, `"`)
			if strings.HasPrefix(val, oldModule+"/") {
				imp.Path.Value = `"` + strings.Replace(val, oldModule, newModule, 1) + `"`
				changed = true
			}
		}
		if changed {
			f, _ := os.Create(path)
			defer f.Close()
			if err := printer.Fprint(f, fset, node); err != nil {
				log.Printf("Error writing file:err:%v", err)
				return err
			}
		}
		return nil
	}); err != nil {
		panic(err)
	}
}

func replaceGoMod(goModPath, oldModule, newModule string) {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return
	}
	content := strings.Replace(string(data), oldModule, newModule, 1)
	_ = os.WriteFile(goModPath, []byte(content), 0644)
}
