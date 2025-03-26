package swagger

import (
	"fmt"
	"reflect"
	"scaffold/pkg/apigen"
	"strconv"
	"strings"
)

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
		Responses:   make(map[string]Response),
		Security:    make([]map[string][]string, 0),
		Parameters:  make([]Parameter, 0),
	}

	// 处理所有非body参数，包括header、path、query等
	for _, param := range params {
		if param.In != "body" {
			// 路径参数始终是必需的
			if param.In == "path" {
				param.Required = true
			}

			// 创建参数对象
			paramObj := Parameter{
				Name:        param.Name,
				In:          param.In,
				Description: param.Description,
				Required:    param.Required,
				Schema: Schema{
					Type: convertGoTypeToSwagger(param.Type),
				},
			}

			// 添加示例值
			if param.Example != nil {
				paramObj.Schema.Example = param.Example
			}

			// 添加格式
			if param.Format != "" {
				paramObj.Schema.Format = param.Format
			}

			// 添加枚举值
			if len(param.Enum) > 0 {
				paramObj.Schema.Enum = param.Enum
			}

			// 添加默认值
			if param.Default != nil {
				paramObj.Schema.Default = param.Default
			}

			// 添加字符串长度限制
			if param.MinLength != nil {
				paramObj.Schema.MinLength = param.MinLength
			}

			if param.MaxLength != nil {
				paramObj.Schema.MaxLength = param.MaxLength
			}

			// 添加数值范围限制
			if param.Minimum != nil {
				paramObj.Schema.Minimum = param.Minimum
			}

			if param.Maximum != nil {
				paramObj.Schema.Maximum = param.Maximum
			}

			// 添加正则表达式模式
			if param.Pattern != "" {
				paramObj.Schema.Pattern = param.Pattern
			}

			operation.Parameters = append(operation.Parameters, paramObj)
		}
	}

	// 处理请求体参数
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

	// 添加安全要求
	if securities := apiInfo.RouteInfo.Security; len(securities) > 0 {
		for _, secName := range securities {
			secName = strings.TrimSpace(secName)
			if secName != "" {
				// 添加安全要求到操作
				operation.Security = append(operation.Security, map[string][]string{
					secName: {},
				})
			}
		}
	}

	// 添加来自RouteInfo的请求头参数（向后兼容）
	for _, headerName := range apiInfo.RouteInfo.Headers {
		// 检查是否已经添加了该请求头参数（避免重复）
		alreadyAdded := false
		for _, param := range operation.Parameters {
			if param.In == "header" && param.Name == headerName {
				alreadyAdded = true
				break
			}
		}

		if !alreadyAdded {
			// 创建请求头参数
			headerParam := Parameter{
				Name:        headerName,
				In:          "header",
				Description: headerName + " 请求头",
				Required:    true, // 默认请求头是必需的
				Schema: Schema{
					Type: "string",
				},
			}

			operation.Parameters = append(operation.Parameters, headerParam)
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

		// 跳过非请求体字段(如header、path、query参数)
		if field.Tag.Get("in") != "" || field.Tag.Get("uri") != "" ||
			field.Tag.Get("query") != "" || field.Tag.Get("header") != "" {
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
		if strings.Contains(field.Tag.Get("binding"), "required") ||
			strings.Contains(field.Tag.Get("v"), "required") {
			required = append(required, fieldName)
		}
	}

	if len(required) > 0 {
		schema.Required = required
	}

	return schema
}
