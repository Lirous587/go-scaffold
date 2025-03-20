package config

type Config struct {
	App   AppConfig   `mapstructure:"app"`
	Log   LogConfig   `mapstructure:"log"`
	DB    DBConfig    `mapstructure:"db"`
	Redis RedisConfig `mapstructure:"redis"`
	Auth  AuthConfig  `mapstructure:"auth"`
}

// AppConfig 应用配置
type AppConfig struct {
	Mode      string `mapstructure:"mode"`
	Port      string `mapstructure:"port"`
	StartTime string `mapstructure:"start_time"`
	MachineID int    `mapstructure:"machine_id"`
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
