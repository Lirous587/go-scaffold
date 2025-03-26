package swagger

import (
	"scaffold/pkg/httpserver/core"
	"scaffold/pkg/httpserver/core/apigen"
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
		schemaName := getUniqueSchemaName(apiInfo.ReqType)
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

	// 简化的响应处理 - 添加标准格式的成功响应
	schema := s.generateModelSchema(apiInfo.ResType)
	schemaName := getUniqueSchemaName(apiInfo.ResType)
	s.doc.Components.Schemas[schemaName] = schema

	// 创建标准响应格式
	operation.Responses["200"] = Response{
		Description: "成功响应",
		Content: map[string]MediaType{
			"application/json": {
				Schema: Schema{
					Type: "object",
					Properties: map[string]Schema{
						"code": {
							Type:    "integer",
							Example: core.CodeSuccess,
						},
						"message": {
							Type:    "string",
							Example: "success",
						},
						"data": {
							Ref: "#/components/schemas/" + schemaName,
						},
					},
					Required: []string{"code", "message", "data"},
				},
			},
		},
	}

	return operation
}
