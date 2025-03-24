package swagger

import (
	"fmt"
	"path"
	"reflect"
	"scaffold/pkg/apigen"
	"strconv"
	"strings"
)

// generateControllerDocs 为控制器生成API文档
func (s *Swagger) generateControllerDocs(pathPrefix string, controller interface{}) {
	if s == nil {
		return
	}

	controllerValue := reflect.ValueOf(controller)
	controllerType := controllerValue.Type()

	// 遍历控制器的所有方法
	for i := 0; i < controllerType.NumMethod(); i++ {
		methodInfo := controllerType.Method(i)

		// 验证并获取API信息
		apiInfo := apigen.BindApiInfo(methodInfo)
		if !apiInfo.Valid {
			continue
		}

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

// generateOperation 从API信息生成操作对象
func (s *Swagger) generateOperation(apiInfo *apigen.ApiInfo) Operation {
	// 提取参数信息
	params := extractParameters(apiInfo.ReqType)

	// 创建操作
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
		reqSchema := s.generateModelSchema(apiInfo.ReqType)
		schemaName := apiInfo.ReqType.Name()
		s.doc.Components.Schemas[schemaName] = reqSchema

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
	s.addResponsesToOperation(&operation, apiInfo)

	return operation
}

// addResponsesToOperation 添加响应到操作
func (s *Swagger) addResponsesToOperation(operation *Operation, apiInfo *apigen.ApiInfo) {
	// 将响应类型添加到组件schemas
	resSchema := s.generateModelSchema(apiInfo.ResType)
	resSchemaName := apiInfo.ResType.Name()
	s.doc.Components.Schemas[resSchemaName] = resSchema

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

	// 添加常见错误响应
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
}

// generateModelSchema 生成模型Schema
func (s *Swagger) generateModelSchema(typ reflect.Type) Schema {
	schema := Schema{
		Type:       "object",
		Properties: make(map[string]Schema),
	}

	// 所需字段列表
	var required []string

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
