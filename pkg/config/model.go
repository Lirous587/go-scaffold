package config

type Config struct {
	Server []ServerConfig `mapstructure:"server"`
	Log    LogConfig      `mapstructure:"log"`
	DB     DBConfig       `mapstructure:"db"`
	Redis  RedisConfig    `mapstructure:"redis"`
	JWT    JWTConfig      `mapstructure:"jwt"`
	Email  EmailConfig    `mapstructure:"email"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

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

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	Password string `mapstructure:"password"`
	PoolSize int    `mapstructure:"pool_size"`
}

type JWTConfig struct {
	Issuer       string `mapstructure:"issuer"`
	Secret       string `mapstructure:"secret"`
	ExpireMinute uint   `mapstructure:"expire_minute"`
}

type EmailConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
	FromName string `mapstructure:"fromName"`
	CC       string `mapstructure:"cc"`
}
