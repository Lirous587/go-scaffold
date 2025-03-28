package swagger

// OpenAPI文档模型定义

// OpenAPIDoc OpenAPI文档结构
type OpenAPIDoc struct {
	OpenAPI    string                `json:"openapi"`
	Info       Info                  `json:"info"`
	Servers    []Server              `json:"servers,omitempty"`
	Paths      map[string]PathItem   `json:"paths"`
	Components Components            `json:"components"`
	Security   []map[string][]string `json:"security,omitempty"`
}

// Info API文档基本信息
type Info struct {
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	Version        string      `json:"version"`
	TermsOfService string      `json:"termsOfService,omitempty"`
	Contact        interface{} `json:"contact,omitempty"`
	License        interface{} `json:"license,omitempty"`
}

// Server API服务器信息
type Server struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

// PathItem API路径项，包含不同HTTP方法的操作
type PathItem map[string]Operation

// Operation API操作定义
type Operation struct {
	Tags        []string              `json:"tags,omitempty"`
	Summary     string                `json:"summary,omitempty"`
	Description string                `json:"description,omitempty"`
	OperationID string                `json:"operationId,omitempty"`
	Deprecated  bool                  `json:"deprecated,omitempty"`
	Parameters  []Parameter           `json:"parameters,omitempty"`
	RequestBody *RequestBody          `json:"requestBody,omitempty"`
	Responses   map[string]Response   `json:"responses"`
	Security    []map[string][]string `json:"security,omitempty"`
}

// Parameter API参数定义
type Parameter struct {
	Name        string      `json:"name"`
	In          string      `json:"in"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required"`
	Schema      Schema      `json:"schema"`
	Example     interface{} `json:"example,omitempty"`
}

// RequestBody 请求体定义
type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Required    bool                 `json:"required"`
	Content     map[string]MediaType `json:"content"`
}

// MediaType 描述请求或响应的媒体类型
type MediaType struct {
	Schema Schema `json:"schema"`
}

// Response API响应定义
type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

// Schema 描述数据结构
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

// Components OpenAPI组件，包含可重用的定义
type Components struct {
	Schemas         map[string]Schema      `json:"schemas"`
	SecuritySchemes map[string]interface{} `json:"securitySchemes,omitempty"`
}

// ParamInfo 从Go结构体解析的参数信息
type ParamInfo struct {
	Name        string      // 参数名称
	In          string      // 参数位置：path, query, body, header
	Type        string      // 参数类型
	Required    bool        // 是否必填
	Description string      // 参数描述
	Example     interface{} // 示例值
	Format      string      // 格式
	Enum        []string    // 枚举值
	Default     interface{} // 默认值
	MinLength   *int        // 最小长度
	MaxLength   *int        // 最大长度
	Minimum     *float64    // 最小值
	Maximum     *float64    // 最大值
	Pattern     string      // 正则表达式模式
}
