package config

type Config struct {
	App     AppConfig     `mapstructure:"app"`
	Log     LogConfig     `mapstructure:"log"`
	DB      DBConfig      `mapstructure:"db"`
	Redis   RedisConfig   `mapstructure:"redis"`
	Swagger SwaggerConfig `mapstructure:"swagger"`
	Auth    AuthConfig    `mapstructure:"auth"`
}

// AppConfig 应用配置
type AppConfig struct {
	Mode      string `mapstructure:"mode"`
	Port      string `mapstructure:"port"`
	StartTime string `mapstructure:"start_time"`
}

// LogConfig 日志配置
type LogConfig struct {
	Mode       string `mapstructure:"mode"`
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

// DBConfig 数据库配置(通用)
type DBConfig struct {
	Driver     string `mapstructure:"driver"` // 数据库驱动类型
	Host       string `mapstructure:"host"`
	Username   string `mapstructure:"username"`
	Password   string `mapstructure:"password"`
	Port       string `mapstructure:"port"`
	DBName     string `mapstructure:"dbname"`
	MaxOpenCon int    `mapstructure:"max_open_con"`
	MaxIdleCon int    `mapstructure:"max_idle_con"`
	// GORM特有配置
	LogLevel      string `mapstructure:"log_level"`
	SlowThreshold int    `mapstructure:"slow_threshold"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	Password string `mapstructure:"password"`
	PoolSize int    `mapstructure:"pool_size"`
}

// SwaggerConfig swagger配置
type SwaggerConfig struct {
	// 基本配置
	Enabled  bool   `mapstructure:"enabled"`   // 是否启用Swagger
	BasePath string `mapstructure:"base_path"` // API基础路径

	// 文档位置配置
	JSONPath string `mapstructure:"json_path"` // swagger.json文件保存位置
	UIPath   string `mapstructure:"ui_path"`   // Swagger UI访问路径

	// 文档信息
	Info struct {
		Title          string `mapstructure:"title"`       // API标题
		Description    string `mapstructure:"description"` // API描述
		Version        string `mapstructure:"version"`     // API版本
		TermsOfService string `mapstructure:"terms"`       // 服务条款URL
		Contact        struct {
			Name  string `mapstructure:"name"`  // 联系人姓名
			Email string `mapstructure:"email"` // 联系人邮箱
			URL   string `mapstructure:"url"`   // 联系人网址
		} `mapstructure:"contact"`
		License struct {
			Name string `mapstructure:"name"` // 许可证名称
			URL  string `mapstructure:"url"`  // 许可证URL
		} `mapstructure:"license"`
	} `mapstructure:"info"`

	// 服务器配置
	Servers []struct {
		URL         string `mapstructure:"url"`         // 服务器URL
		Description string `mapstructure:"description"` // 服务器描述
	} `mapstructure:"servers"`

	// 安全配置
	SecurityDefinitions map[string]struct {
		Type         string `mapstructure:"type"`          // 安全类型: apiKey, http, oauth2, openIdConnect
		Name         string `mapstructure:"name"`          // 用于apiKey类型
		In           string `mapstructure:"in"`            // 用于apiKey类型: query, header, cookie
		Scheme       string `mapstructure:"scheme"`        // 用于http类型: basic, bearer, digest
		BearerFormat string `mapstructure:"bearer_format"` // 用于http类型带bearer scheme
		Description  string `mapstructure:"description"`   // 安全方案描述
	} `mapstructure:"security_definitions"`

	// 默认安全方案
	Security []string `mapstructure:"security"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	Admin AdminAuth `mapstructure:"admin"`
	User  UserAuth  `mapstructure:"user"`
}

type AdminAuth struct {
	JwtSecret       string `mapstructure:"jwt_secret"`
	JwtExpireMinute int    `mapstructure:"jwt_expire_minute"`
}

type UserAuth struct {
	JwtSecret       string `mapstructure:"jwt_secret"`
	JwtExpireMinute int    `mapstructure:"jwt_expire_minute"`
}
