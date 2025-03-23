package apigen

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

// OpenAPIDoc OpenAPI文档结构
type OpenAPIDoc struct {
	OpenAPI    string              `json:"openapi"`
	Info       Info                `json:"info"`
	Servers    []Server            `json:"servers"`
	Paths      map[string]PathItem `json:"paths"`
	Components Components          `json:"components"`
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type Server struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type PathItem map[string]Operation

type Operation struct {
    Tags        []string              `json:"tags,omitempty"`
    Summary     string                `json:"summary,omitempty"`
    Description string                `json:"description,omitempty"`
    OperationID string                `json:"operationId,omitempty"`
    Deprecated  bool                  `json:"deprecated,omitempty"`
    Parameters  []Parameter           `json:"parameters,omitempty"`
    RequestBody *RequestBody          `json:"requestBody,omitempty"`
    Responses   map[string]Response   `json:"responses"`
    Consumes    []string              `json:"consumes,omitempty"`
    Produces    []string              `json:"produces,omitempty"`
    Security    []map[string][]string `json:"security,omitempty"`
}

type Parameter struct {
	Name        string `json:"name"`
	In          string `json:"in"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required"`
	Schema      Schema `json:"schema"`

	// 可以直接添加到参数级别的属性
	Example interface{} `json:"example,omitempty"`
	Enum    []string    `json:"enum,omitempty"`
}

type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Required    bool                 `json:"required"`
	Content     map[string]MediaType `json:"content"`
}

type MediaType struct {
	Schema Schema `json:"schema"`
}

type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

type Schema struct {
	Type        string            `json:"type,omitempty"`
	Format      string            `json:"format,omitempty"`
	Description string            `json:"description,omitempty"`
	Properties  map[string]Schema `json:"properties,omitempty"`
	Items       *Schema           `json:"items,omitempty"`
	Required    []string          `json:"required,omitempty"`
	Ref         string            `json:"$ref,omitempty"`
	Example     interface{}       `json:"example,omitempty"`
	Enum        []string          `json:"enum,omitempty"`
	Default     interface{}       `json:"default,omitempty"`
	MinLength   *int              `json:"minLength,omitempty"`
	MaxLength   *int              `json:"maxLength,omitempty"`
	Minimum     *float64          `json:"minimum,omitempty"`
	Maximum     *float64          `json:"maximum,omitempty"`
	Pattern     string            `json:"pattern,omitempty"`
}

type Components struct {
	Schemas map[string]Schema `json:"schemas"`
}

// GenerateSwaggerDocs 为所有API生成Swagger文档
func GenerateSwaggerDocs(controllers ...interface{}) error {
	fmt.Println("开始生成Swagger文档...")

	// 初始化OpenAPI文档
	doc := OpenAPIDoc{
		OpenAPI: "3.0.0",
		Info: Info{
			Title:       "API服务框架",
			Description: "自动生成的API文档",
			Version:     "1.0",
		},
		Servers: []Server{
			{URL: "http://localhost:8080", Description: "开发服务器"},
		},
		Paths: make(map[string]PathItem),
		Components: Components{
			Schemas: make(map[string]Schema),
		},
	}

	// 遍历所有控制器
	for _, controller := range controllers {
		if err := generateControllerDocs(&doc, controller); err != nil {
			return err
		}
	}

	// 保存文档到文件
	docsDir := "./docs"
	if err := os.MkdirAll(docsDir, 0755); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}

	swaggerFile := filepath.Join(docsDir, "swagger.json")
	if err := os.WriteFile(swaggerFile, jsonData, 0644); err != nil {
		return err
	}

	fmt.Printf("Swagger文档生成完成: %s\n", swaggerFile)
	return nil
}

// 为控制器生成文档
func generateControllerDocs(doc *OpenAPIDoc, controller interface{}) error {
	controllerValue := reflect.ValueOf(controller)
	controllerType := controllerValue.Type()

	// 遍历控制器的所有方法
	for i := 0; i < controllerType.NumMethod(); i++ {
		methodInfo := controllerType.Method(i)

		// 验证并获取API信息
		apiInfo := ValidateMethod(controller, methodInfo)
		if !apiInfo.Valid {
			continue
		}

		// 提取参数信息
		params := ExtractParameters(apiInfo.ReqType)

		// 添加路由到Swagger
		operation := Operation{
			Tags:        apiInfo.RouteInfo.Tags,
			Summary:     apiInfo.RouteInfo.Summary,
			Description: apiInfo.RouteInfo.Description,
			OperationID: apiInfo.RouteInfo.OperationID,
			Deprecated:  apiInfo.RouteInfo.Deprecated,
			Consumes:    apiInfo.RouteInfo.Consumes,
			Produces:    apiInfo.RouteInfo.Produces,
			Responses:   make(map[string]Response),
		}

		// 添加路径参数
		for _, param := range params {
			if param.In == "path" {
				operation.Parameters = append(operation.Parameters, Parameter{
					Name:        param.Name,
					In:          param.In,
					Description: param.Description,
					Required:    true, // 路径参数总是必需的
					Schema: Schema{
						Type: convertGoTypeToSwagger(param.Type),
					},
				})
			}
		}

		// 处理请求体
		hasBodyParams := false
		for _, param := range params {
			if param.In == "body" {
				hasBodyParams = true
				break
			}
		}

		// 如果有请求体参数，添加请求体
		if hasBodyParams {
			// 将请求类型添加到组件schemas
			reqSchema := generateModelSchema(apiInfo.ReqType)
			schemaName := apiInfo.ReqType.Name()
			doc.Components.Schemas[schemaName] = reqSchema

			// 添加请求体引用
			operation.RequestBody = &RequestBody{
				Description: "请求参数",
				Required:    true,
				Content: map[string]MediaType{
					"application/json": {
						Schema: Schema{
							Ref: "#/components/schemas/" + schemaName,
						},
					},
				},
			}
		}

		// 添加响应
		// 将响应类型添加到组件schemas
		resSchema := generateModelSchema(apiInfo.ResType)
		resSchemaName := apiInfo.ResType.Name()
		doc.Components.Schemas[resSchemaName] = resSchema

		// 成功响应
		operation.Responses["200"] = Response{
			Description: "成功响应",
			Content: map[string]MediaType{
				"application/json": {
					Schema: Schema{
						Ref: "#/components/schemas/" + resSchemaName,
					},
				},
			},
		}

		// 错误响应
		operation.Responses["400"] = Response{
			Description: "请求参数错误",
			Content: map[string]MediaType{
				"application/json": {
					Schema: Schema{
						Type: "object",
						Properties: map[string]Schema{
							"code":    {Type: "integer"},
							"message": {Type: "string"},
							"errors":  {Type: "object"},
						},
					},
				},
			},
		}

		operation.Responses["500"] = Response{
			Description: "服务器内部错误",
			Content: map[string]MediaType{
				"application/json": {
					Schema: Schema{
						Type: "object",
						Properties: map[string]Schema{
							"code":    {Type: "integer"},
							"message": {Type: "string"},
							"error":   {Type: "string"},
						},
					},
				},
			},
		}

		// 规范化路径格式
		path := apiInfo.RouteInfo.Path
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}

		// 添加操作到路径
		httpMethod := strings.ToLower(apiInfo.RouteInfo.Method)
		if doc.Paths[path] == nil {
			doc.Paths[path] = make(PathItem)
		}
		doc.Paths[path][httpMethod] = operation
	}

	return nil
}

// 生成模型Schema
func generateModelSchema(typ reflect.Type) Schema {
	schema := Schema{
		Type:       "object",
		Properties: make(map[string]Schema),
	}

	// 所需字段列表
	required := []string{}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// 跳过Meta字段
		if field.Name == "Meta" {
			continue
		}

		// 获取JSON字段名
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// 提取字段名
		fieldName := strings.Split(jsonTag, ",")[0]

		// 创建属性Schema
		propSchema := Schema{
			Type: convertGoTypeToSwagger(field.Type.String()),
		}

		// 添加描述
		if dcTag := field.Tag.Get("dc"); dcTag != "" {
			propSchema.Description = dcTag
		} else {
			propSchema.Description = fmt.Sprintf("%s字段", field.Name)
		}

		// 添加格式
		if formatTag := field.Tag.Get("format"); formatTag != "" {
			propSchema.Format = formatTag
		}

		// 添加示例值
		if exampleTag := field.Tag.Get("example"); exampleTag != "" {
			propSchema.Example = parseExampleValue(exampleTag, field.Type)
		}

		// 添加枚举值
		if enumTag := field.Tag.Get("enum"); enumTag != "" {
			propSchema.Enum = strings.Split(enumTag, ",")
		}

		// 添加默认值
		if defaultTag := field.Tag.Get("default"); defaultTag != "" {
			propSchema.Default = parseExampleValue(defaultTag, field.Type)
		}

		// 添加字符串长度限制
		if minLengthTag := field.Tag.Get("minLength"); minLengthTag != "" {
			if val, err := strconv.Atoi(minLengthTag); err == nil {
				propSchema.MinLength = &val
			}
		}

		if maxLengthTag := field.Tag.Get("maxLength"); maxLengthTag != "" {
			if val, err := strconv.Atoi(maxLengthTag); err == nil {
				propSchema.MaxLength = &val
			}
		}

		// 添加数值范围限制
		if minTag := field.Tag.Get("min"); minTag != "" {
			if val, err := strconv.ParseFloat(minTag, 64); err == nil {
				propSchema.Minimum = &val
			}
		}

		if maxTag := field.Tag.Get("max"); maxTag != "" {
			if val, err := strconv.ParseFloat(maxTag, 64); err == nil {
				propSchema.Maximum = &val
			}
		}

		// 添加正则表达式模式
		if patternTag := field.Tag.Get("pattern"); patternTag != "" {
			propSchema.Pattern = patternTag
		}

		// 添加属性到schema
		schema.Properties[fieldName] = propSchema

		// 如果字段是必填的，添加到required列表
		if strings.Contains(field.Tag.Get("binding"), "required") {
			required = append(required, fieldName)
		}
	}

	if len(required) > 0 {
		schema.Required = required
	}

	return schema
}

// 将Go类型转换为Swagger类型
func convertGoTypeToSwagger(goType string) string {
	switch {
	case strings.Contains(goType, "string"):
		return "string"
	case strings.Contains(goType, "int"), strings.Contains(goType, "uint"):
		return "integer"
	case strings.Contains(goType, "float"):
		return "number"
	case strings.Contains(goType, "bool"):
		return "boolean"
	case strings.Contains(goType, "time.Time"):
		return "string" // 时间类型在OpenAPI中用字符串表示
	case strings.Contains(goType, "[]"):
		return "array"
	default:
		return "object"
	}
}
