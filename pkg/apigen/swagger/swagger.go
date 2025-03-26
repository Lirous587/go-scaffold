package swagger

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"scaffold/pkg/apigen"
	"scaffold/pkg/config"
	"strings"
)

// Swagger Swagger文档生成器
type Swagger struct {
	doc OpenAPIDoc
}

// New 创建新的Swagger文档生成器
func New() *Swagger {
	cfg := config.Cfg.Swagger
	if !cfg.Enabled {
		return nil
	}

	// 初始化OpenAPI文档
	doc := OpenAPIDoc{
		OpenAPI: "3.0.0",
		Info: Info{
			Title:       cfg.Info.Title,
			Description: cfg.Info.Description,
			Version:     cfg.Info.Version,
		},
		Paths:      make(map[string]PathItem),
		Components: Components{Schemas: make(map[string]Schema)},
	}

	// 添加联系人信息
	if cfg.Info.Contact.Name != "" || cfg.Info.Contact.Email != "" || cfg.Info.Contact.URL != "" {
		contactInfo := map[string]interface{}{
			"name":  cfg.Info.Contact.Name,
			"email": cfg.Info.Contact.Email,
			"url":   cfg.Info.Contact.URL,
		}
		// 移除空字段
		for k, v := range contactInfo {
			if v == "" {
				delete(contactInfo, k)
			}
		}
		if len(contactInfo) > 0 {
			contactBytes, _ := json.Marshal(contactInfo)
			json.Unmarshal(contactBytes, &doc.Info.Contact)
		}
	}

	// 添加许可证信息
	if cfg.Info.License.Name != "" || cfg.Info.License.URL != "" {
		licenseInfo := map[string]interface{}{
			"name": cfg.Info.License.Name,
			"url":  cfg.Info.License.URL,
		}
		// 移除空字段
		for k, v := range licenseInfo {
			if v == "" {
				delete(licenseInfo, k)
			}
		}
		if len(licenseInfo) > 0 {
			licenseBytes, _ := json.Marshal(licenseInfo)
			json.Unmarshal(licenseBytes, &doc.Info.License)
		}
	}

	// 添加服务条款
	if cfg.Info.TermsOfService != "" {
		doc.Info.TermsOfService = cfg.Info.TermsOfService
	}

	// 添加服务器信息
	doc.Servers = make([]Server, 0, len(cfg.Servers))
	for _, server := range cfg.Servers {
		doc.Servers = append(doc.Servers, Server{
			URL:         server.URL,
			Description: server.Description,
		})
	}

	// 添加安全定义
	if len(cfg.SecurityDefinitions) > 0 {
		doc.Components.SecuritySchemes = make(map[string]interface{})
		for name, scheme := range cfg.SecurityDefinitions {
			securityScheme := map[string]interface{}{
				"type":        scheme.Type,
				"description": scheme.Description,
			}

			// 根据类型添加特定字段
			switch scheme.Type {
			case "apiKey":
				securityScheme["name"] = scheme.Name
				securityScheme["in"] = scheme.In
			case "http":
				securityScheme["scheme"] = scheme.Scheme
				if scheme.BearerFormat != "" {
					securityScheme["bearerFormat"] = scheme.BearerFormat
				}
			}

			// 移除空字段
			for k, v := range securityScheme {
				if v == "" {
					delete(securityScheme, k)
				}
			}

			doc.Components.SecuritySchemes[name] = securityScheme
		}
	}

	// 添加默认安全要求
	if len(cfg.Security) > 0 {
		doc.Security = make([]map[string][]string, 0)
		for _, secName := range cfg.Security {
			doc.Security = append(doc.Security, map[string][]string{
				secName: {},
			})
		}
	}

	swagger := Swagger{}
	swagger.doc = doc
	return &swagger
}

// GenerateDocs 为API控制器生成Swagger文档
func (s *Swagger) GenerateDocs(pathPrefix string, apiInfos []apigen.ApiInfo) {
	for _, apiInfo := range apiInfos {
		// 生成操作对象
		operation := s.generateOperation(&apiInfo)

		// 规范化路径
		normPath := path.Join("/", pathPrefix, apiInfo.RouteInfo.Path)

		// 添加操作到路径
		httpMethod := strings.ToLower(apiInfo.RouteInfo.Method)
		if s.doc.Paths[normPath] == nil {
			s.doc.Paths[normPath] = make(PathItem)
		}
		s.doc.Paths[normPath][httpMethod] = operation
	}
}

// Save 保存Swagger文档到文件
func (s *Swagger) Save() error {
	if s == nil {
		return nil
	}

	fmt.Println("开始生成Swagger文档...")

	// 保存文档到文件
	docsDir := "./docs"
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(s.doc, "", "  ")
	if err != nil {
		return err
	}

	swaggerFile := filepath.Join(docsDir, "swagger.json")
	if err := os.WriteFile(swaggerFile, jsonData, 0644); err != nil {
		return err
	}

	fmt.Println("生成Swagger文档完成")
	return nil
}
