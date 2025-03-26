package apigen

import (
	"reflect"
	"strconv"
	"strings"
)

// Meta API元数据，用于自动注册路由和生成文档
type Meta struct{}

type RouteInfoProvider interface {
	GetRouteInfo(instance interface{}) RouteInfo
}

// RouteInfo 路由信息
type RouteInfo struct {
	Method      string   // HTTP方法
	Path        string   // 路径
	Summary     string   // 摘要
	Description string   // 描述
	Tags        []string // 标签
	OperationID string   // 操作ID
	Produces    []string // 响应内容类型，如 application/json
	Consumes    []string // 请求内容类型，如 application/json
	Security    []string // 安全要求，如 JWT, OAuth2
	Headers     []string // 请求头参数名称
	Deprecated  bool     // 是否已弃用
	Version     string   // API版本
}

// GetRouteInfo 从标签中提取路由信息
func (m Meta) GetRouteInfo(parentStruct interface{}) RouteInfo {
	val := reflect.ValueOf(parentStruct)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	metaField, _ := typ.FieldByName("Meta")

	// 提取路由信息
	info := RouteInfo{
		Method:      metaField.Tag.Get("method"),
		Path:        metaField.Tag.Get("path"),
		Summary:     metaField.Tag.Get("sm"),
		Description: metaField.Tag.Get("dc"),
		Tags:        splitAndTrim(metaField.Tag.Get("tags"), ","),
		OperationID: metaField.Tag.Get("operationId"),
		Version:     metaField.Tag.Get("version"),
	}

	// 提取请求头参数
	if headers := metaField.Tag.Get("header"); headers != "" {
		info.Headers = splitAndTrim(headers, ",")
	}
	// 提取消费和生产内容类型
	if produces := metaField.Tag.Get("produces"); produces != "" {
		info.Produces = splitAndTrim(produces, ",")
	} else {
		// 默认值
		info.Produces = []string{"application/json"}
	}

	if consumes := metaField.Tag.Get("consumes"); consumes != "" {
		info.Consumes = splitAndTrim(consumes, ",")
	} else {
		// 默认值
		info.Consumes = []string{"application/json"}
	}

	// 提取安全要求
	if security := metaField.Tag.Get("security"); security != "" {
		info.Security = splitAndTrim(security, ",")
	}

	// 解析是否弃用
	if deprecated := metaField.Tag.Get("deprecated"); deprecated != "" {
		info.Deprecated, _ = strconv.ParseBool(deprecated)
	}

	return info
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

// 辅助函数，分割字符串并去除空白
func splitAndTrim(s, sep string) []string {
	if s == "" {
		return []string{}
	}

	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
