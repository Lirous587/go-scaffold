package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"html/template"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic("加载环境变量失败")
	}
}

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

	// 查找go.mod所在路径
	module := ""
	cwd, err := os.Getwd()
	if err != nil {
		panic("os.Getwd() 调用失败")
	}
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

	data := map[string]string{
		"Domain":      domainLower,
		"DomainTitle": domainTitle,
		"Module":      module,
	}

	tempDir := filepath.Join("template_v" + *version)
	goModDir := filepath.Dir(goModPath)

	hasGenTemplates := false
	// 是否创建了模板目录
	outBase := filepath.Join(goModDir, "internal", domainLower)
	if _, err := os.Stat(outBase); os.IsNotExist(err) {
		if err := os.MkdirAll(outBase, 0755); err != nil {
			fmt.Printf("创建目录失败:%v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("%s 模板目录已存在，跳过生成\n", domainLower)
		hasGenTemplates = true
	}

	// 生成模板文件
	if !hasGenTemplates {
		err = filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(d.Name(), ".tmpl") {
				return nil
			}

			relPath, _ := filepath.Rel(tempDir, path)
			// 替换文件名中的 model 为 domainLower
			outFile := strings.ReplaceAll(strings.TrimSuffix(relPath, ".tmpl"), "model", domainLower) + ".go"
			outPath := filepath.Join(outBase, outFile)
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
			rollBackTemplate(outBase)
			os.Exit(1)
		}
	}

	hasGenCodes := false
	// 是否创建了codes
	codePath := filepath.Join(goModDir, "internal", "common", "reskit", "codes", domainLower+".go")
	if _, err := os.Stat(codePath); os.IsNotExist(err) {
		f, err := os.Create(codePath)
		if err != nil {
			fmt.Printf("创建code文件失败:%v\n", err)
			rollBackTemplate(outBase)
			os.Exit(1)
		}
		defer f.Close()
	} else {
		fmt.Printf("%s code文件已存在，跳过生成\n", domainLower)
		hasGenCodes = true
	}
	// 生成codes
	if !hasGenCodes {
		tmplPath := filepath.Join(goModDir, "tool", "gen", "codes.tmpl")
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			fmt.Printf("解析codes模板失败: %v\n", err)
			rollBackTemplate(outBase)
			os.Exit(1)
		}
		f, err := os.Create(codePath)
		if err != nil {
			fmt.Printf("创建code文件失败: %v\n", err)
			rollBackTemplate(outBase)
			os.Exit(1)
		}
		defer f.Close()
		if err := tmpl.Execute(f, data); err != nil {
			fmt.Printf("渲染codes模板失败: %v\n", err)
			rollBackCode(codePath)
			rollBackTemplate(outBase)
			os.Exit(1)
		}
		fmt.Printf("生成成功: %s\n", codePath)
	}

	// 创建基础表
	if err := createTable(domainLower); err != nil {
		rollBackCode(codePath)
		rollBackTemplate(outBase)
		panic(err)
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

// 创建模板回滚
func rollBackTemplate(outBase string) {
	if err := os.RemoveAll(outBase); err != nil {
		fmt.Printf("回滚失败，无法删除目录: %v\n", err)
	} else {
		fmt.Printf("已回滚并删除目录: %s\n", outBase)
	}
}

// 创建code回滚
func rollBackCode(codePath string) {
	if err := os.Remove(codePath); err != nil {
		fmt.Printf("回滚失败，无法删除code文件: %v\n", err)
	} else {
		fmt.Printf("已回滚并删除code文件: %s\n", codePath)
	}
}

// 建表
func createTable(domainLower string) error {
	cwd, err := os.Getwd()
	if err != nil {
		panic("os.Getwd() 调用失败")
	}

	goModPath, _ := findGoModPath(cwd)
	goModDir := filepath.Dir(goModPath)
	sqlPath := filepath.Join(goModDir, "tool", "gen", "ddl.sql")
	content, err := os.ReadFile(sqlPath)
	if err != nil {
		fmt.Printf("读取建表 SQL 失败: %v\n", err)
		return err
	}
	sqlStr := strings.ReplaceAll(string(content), "$Domain", domainLower)

	requiredKeys := []string{
		"PSQL_HOST", "PSQL_PORT", "PSQL_USERNAME",
		"PSQL_PASSWORD", "PSQL_DB_NAME", "PSQL_SSL_MODE",
	}
	envs := make(map[string]string)
	for _, key := range requiredKeys {
		val := os.Getenv(key)
		if val == "" {
			return fmt.Errorf("环境变量 %s 未设置", key)
		}
		envs[key] = val
	}
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		envs["PSQL_HOST"],
		envs["PSQL_PORT"],
		envs["PSQL_USERNAME"],
		envs["PSQL_PASSWORD"],
		envs["PSQL_DB_NAME"],
		envs["PSQL_SSL_MODE"],
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 判断表是否已存在
	var exists bool
	checkSQL := `SELECT EXISTS (
		SELECT 1 FROM information_schema.tables 
		WHERE table_schema = 'public' AND table_name = $1
	)`
	if err := db.QueryRow(checkSQL, domainLower).Scan(&exists); err != nil {
		return fmt.Errorf("检查表是否存在失败: %v", err)
	}
	if exists {
		fmt.Printf("表 %s 已存在，跳过建表\n", domainLower)
		return nil
	}

	if _, err := db.Exec(sqlStr); err != nil {
		return fmt.Errorf("执行建表 SQL 失败: %v", err)
	}
	fmt.Println("建表成功")

	cmd := exec.Command("sqlboiler", "psql")
	cmd.Dir = goModDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("执行 sqlboiler 失败: %v", err)
	}
	fmt.Println("sqlboiler 执行完成")

	return nil
}
