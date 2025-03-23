package apigen

import (
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

