package swagger

import (
    "fmt"
    "reflect"
    "strconv"
    "strings"
)

// extractParameters 从结构体类型中提取参数信息
func extractParameters(reqType reflect.Type) []ParamInfo {
    var params []ParamInfo

    // 遍历所有字段
    for i := 0; i < reqType.NumField(); i++ {
        field := reqType.Field(i)

        // 跳过Meta字段
        if field.Name == "Meta" {
            continue
        }

        // 确定参数位置和名称
        paramIn := "body" // 默认在请求体
        paramName := field.Tag.Get("json")

        // 检查in标签
        if inTag := field.Tag.Get("in"); inTag != "" {
            paramIn = inTag
            // 对于header, path, query参数，优先使用字段名作为参数名
            paramName = strings.ToLower(field.Name)
        }

        // 如果仍有json标签，则使用json标签的名称
        if paramName != "" {
            // 处理json标签中的附加选项, 如 `json:"name,omitempty"`
            parts := strings.Split(paramName, ",")
            paramName = parts[0]
        }

        // 检查uri标签
        if uriTag := field.Tag.Get("uri"); uriTag != "" {
            paramIn = "path"
            paramName = uriTag
        } else if queryTag := field.Tag.Get("query"); queryTag != "" {
            paramIn = "query"
            paramName = queryTag
        } else if headerTag := field.Tag.Get("header"); headerTag != "" {
            paramIn = "header"
            paramName = headerTag
        }

        // 如果没有找到任何名称，使用字段名
        if paramName == "" {
            paramName = strings.ToLower(field.Name)
        }

        // 创建参数信息
        param := ParamInfo{
            Name: paramName,
            In:   paramIn,
            Type: field.Type.String(),
        }

        // 解析描述标签
        if dcTag := field.Tag.Get("dc"); dcTag != "" {
            param.Description = dcTag
        } else {
            param.Description = fmt.Sprintf("%s字段", field.Name) // 默认描述
        }

        // 确定是否必填
        param.Required = strings.Contains(field.Tag.Get("binding"), "required") 
                         

        // 解析示例值标签
        if exampleTag := field.Tag.Get("example"); exampleTag != "" {
            param.Example = parseExampleValue(exampleTag, field.Type)
        }

        // 解析格式标签
        if formatTag := field.Tag.Get("format"); formatTag != "" {
            param.Format = formatTag
        }

        // 解析枚举标签
        if enumTag := field.Tag.Get("enum"); enumTag != "" {
            param.Enum = strings.Split(enumTag, ",")
        }

        // 解析默认值标签
        if defaultTag := field.Tag.Get("default"); defaultTag != "" {
            param.Default = parseExampleValue(defaultTag, field.Type)
        }

        // 解析字符串长度限制
        if minLengthTag := field.Tag.Get("minLength"); minLengthTag != "" {
            if val, err := strconv.Atoi(minLengthTag); err == nil {
                param.MinLength = &val
            }
        }

        if maxLengthTag := field.Tag.Get("maxLength"); maxLengthTag != "" {
            if val, err := strconv.Atoi(maxLengthTag); err == nil {
                param.MaxLength = &val
            }
        }

        // 解析数值范围限制
        if minTag := field.Tag.Get("min"); minTag != "" {
            if val, err := strconv.ParseFloat(minTag, 64); err == nil {
                param.Minimum = &val
            }
        }

        if maxTag := field.Tag.Get("max"); maxTag != "" {
            if val, err := strconv.ParseFloat(maxTag, 64); err == nil {
                param.Maximum = &val
            }
        }

        // 解析正则表达式模式
        if patternTag := field.Tag.Get("pattern"); patternTag != "" {
            param.Pattern = patternTag
        }

        params = append(params, param)
    }

    return params
}

// parseExampleValue 根据字段类型解析示例值
func parseExampleValue(value string, fieldType reflect.Type) interface{} {
    // 处理模板变量，例如 {{token}}
    if strings.HasPrefix(value, "{{") && strings.HasSuffix(value, "}}") {
        // 保留原样，让前端处理模板变量
        return value
    }

    // 根据字段类型进行转换
    switch fieldType.Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        if val, err := strconv.ParseInt(value, 10, 64); err == nil {
            return val
        }
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        if val, err := strconv.ParseUint(value, 10, 64); err == nil {
            return val
        }
    case reflect.Float32, reflect.Float64:
        if val, err := strconv.ParseFloat(value, 64); err == nil {
            return val
        }
    case reflect.Bool:
        if val, err := strconv.ParseBool(value); err == nil {
            return val
        }
    default:
        return value
    }
    return value
}

// convertGoTypeToSwagger 将Go类型转换为Swagger类型
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