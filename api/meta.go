package api

import (
	"fmt"
	"reflect"
	"strings"
)

// Meta API元数据，用于自动注册路由和生成文档
type Meta struct{}

// GetRouteInfo 从标签中提取路由信息
func (m Meta) GetRouteInfo(parentStruct interface{}) RouteInfo {
    val := reflect.ValueOf(parentStruct)
    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }
    
    typ := val.Type()
    metaField, _ := typ.FieldByName("Meta")
    
    // 提取路由信息
    return RouteInfo{
        Method:      metaField.Tag.Get("method"),
        Path:        metaField.Tag.Get("path"),
        Summary:     metaField.Tag.Get("sm"),
        Description: metaField.Tag.Get("dc"),
        Tags:        strings.Split(metaField.Tag.Get("tags"), ","),
    }
}

// RouteInfo 路由信息
type RouteInfo struct {
    Method      string   // HTTP方法
    Path        string   // 路径
    Summary     string   // 摘要
    Description string   // 描述
    Tags        []string // 标签
}

// GetGinPath 将{param}格式转换为:param格式
func (r RouteInfo) GetGinPath() string {
    path := r.Path
    segments := strings.Split(path, "/")
    for i, segment := range segments {
        if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
            paramName := segment[1 : len(segment)-1]
            segments[i] = ":" + paramName
        }
    }
    return strings.Join(segments, "/")
}

// 参数信息
type ParamInfo struct {
    Name        string // 参数名称
    In          string // 参数位置：path, query, body, header
    Type        string // 参数类型
    Required    bool   // 是否必填
    Description string // 参数描述
}

// ExtractParameters 从结构体类型中提取参数信息
func ExtractParameters(reqType reflect.Type) []ParamInfo {
    params := []ParamInfo{}
    
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
        
        // 确定是否必填
        required := strings.Contains(field.Tag.Get("binding"), "required")
        
        // 添加参数信息
        params = append(params, ParamInfo{
            Name:        paramName,
            In:          paramIn,
            Type:        field.Type.String(),
            Required:    required,
            Description: fmt.Sprintf("%s字段", field.Name), // 默认描述
        })
    }
    
    return params
}