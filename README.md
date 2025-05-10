# 基于 Gin 的轻量级脚手架

一个基于 Gin 框架的轻量级脚手架，集成了常用组件，帮助你快速搭建高性能的 Go Web 应用。

[![Go Version](https://img.shields.io/badge/Go-v1.18+-blue.svg)](https://golang.org/doc/devel/release.html)
[![Gin](https://img.shields.io/badge/Gin-v1.9.0+-green.svg)](https://github.com/gin-gonic/gin)
[![GORM](https://img.shields.io/badge/GORM-v1.25.0+-lightblue.svg)](https://gorm.io/)

## 🚀 特性

- 📝 完整的项目结构和最佳实践
- 🔒 JWT 认证集成
- 📊 统一的 API 响应格式
- 🔄 强大的中间件支持
- 📋 详尽的日志记录
- 🔌 多数据库支持
- 🛠️ 优雅的错误处理
- 🚦 优雅启动和关闭

## 🔧 技术栈

- [Gin](https://github.com/gin-gonic/gin) - 高性能 HTTP Web 框架
- [GORM](https://gorm.io/) - 优秀的 ORM 库，支持 MySQL、PostgreSQL 等
- [Redis](https://github.com/redis/go-redis) - Redis 客户端
- [Viper](https://github.com/spf13/viper) - 完整的配置解决方案
- [Zap](https://github.com/uber-go/zap) - 高性能、结构化日志
- [Wire](https://github.com/google/wire) - Wire 依赖注入
- [JWT](https://github.com/golang-jwt/jwt) - JWT 鉴权管理
- [go-gomail](https://github.com/go-gomail/gomail) - go-gomail 邮箱库

## 📁 项目结构

```
scaffold/
├── internal/         # 实际业务逻辑
│   └── cmd           
│   └── domain        # 领域模型  
│   └── middleware    # 中间件 
├── logs/             # 日志文件
├── manifest/         # 配置目录
│   └── config/
│       └── config.yaml 
│   └── docker/
│       └── Dockerfile
├── pkg/              # 依赖项
│   ├── config/       # 配置结构化管理
│   ├── email/        # email相关
│   ├── httpserver/   # gin引擎初始化
│   ├── jwt/          # jwt相关
│   ├── logger/       # 日志配置
│   ├── repository/   # 数据存储层
│   │   ├── gorm.go   # 数据库单例
│   │   └── redis.go  # Redis单例
│   ├── response/     # 响应管理
│   ├── validator/    # validator管理
├── resource/         # 静态资源
├── utility/          # 工具函数
├── .air.conf         # air配置
├── .gitignore
├── main.go           # 主入口
└── README.md
```

## ⚡ 快速开始

### 前置要求

- Go 1.18+
- MySQL 5.7+ 或 PostgreSQL 10+
- Redis 6.0+

### 快速开始

> 以下演示以Windows作为示例

1. 新建目录

```bash
mkdir demo
```

2. 克隆项目

```bash
git clone https://github.com/Lirou587/go-scaffold.git
```

3. 移动目录 并删除git记录

```bash
robocopy go-scaffold . /E /XD .git
```

4. 删除clone目录

```bash
Remove-Item go-scaffold -Recurse -Force
```

5. 编写并运行ast脚本

```bash
go build ./ast.go 
# 填写实际的项目名称
./ast.exe demo
```

6. 删除ast

```bash
rm ./ast.go
rm ./ast.exe
```

7. 安装依赖

```bash
go mod tidy
```

8. 修改配置

编辑 `manifest/config/config.yaml` 配置文件

9. 运行服务

```bash
go run main.go
# 或者运行 air
```

## ⚙️ 配置项

配置文件路径：`manifest/config/config.yaml`

```yaml
server:
  - port: 8080
    mode: "dev"

log:
  mode: "dev"
  level: "info"
  filename: "logs/scaffold.log"
  max_size: 1
  max_age: 30
  max_backups: 7

db:
  driver: "mysql"  # 增加驱动类型字段，方便将来切换数据库
  host: "127.0.0.1"
  username: "root"
  password: "123"
  port: "3306"
  dbname: "test"
  max_open_con: 100
  max_idle_con: 50
  # 可以增加GORM特有配置
  log_level: "info"
  slow_threshold: 200  # 慢SQL阈值(ms)

redis:
  host: "127.0.0.1"
  port: "6379"
  db: 0
  #  password:
  pool_size: 200

jwt:
  issuer: "lirous"
  secret: "https://lirous.com"
  expire_minute: 120

email:
  host: "smtp.qq.com" //可以换成对应的平台
  port: 465
  username: "xxxx@xx.xx"
  password: "xxxxxx"
  from: "xxxx@xx.xx"
  fromName: "xxx"
  cc: "xxxxx.com"
```

## 🔌 主要组件

### httpserver - Web 引擎

基于 Gin 封装，支持优雅重启和关闭:

```go
// 初始化路由
s := httpserver.New(8080)

// 启动服务
s.Run()

// 支持多个端口启动 类似与gin官网的做法
{
    s1 := httpserver.New(8080)
    s2 := httpserver.New(8081)
    
    // 启动服务
    go s1.Run()
    s2.Run()
}
```

### Logger - 日志系统

基于 Zap，支持分级、轮转和多输出:

```go
// 记录信息
zap.L().Info("操作成功",
    zap.String("user", "admin"),
    zap.Int("items", 10))

// 记录错误
zap.L().Error("数据库错误",
    zap.Error(err),
    zap.String("query", "SELECT * FROM users"))
```

### 数据库 - GORM

支持 MySQL 和 PostgreSQL，自动迁移:

### Redis - 缓存

简化的 Redis 操作:

```go
// 设置缓存
err := redis.Client().Set(ctx, "key", "value", time.Minute).Err()

// 获取缓存
val, err := redis.Client().Get(ctx, "key").Result()
```

## 📝 最佳实践

1. **配置验证** - 启动时自动验证必要配置项
2. **错误处理** - 使用 `github.com/pkg/errors` 提供完整错误栈
3. **优雅关机** - 处理 SIGTERM 等信号，平滑关闭服务
4. **热重启** - 支持不停机更新应用程序

## 🤝 贡献

欢迎贡献代码或提出建议！请遵循以下步骤：

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 详情参见 [LICENSE](LICENSE) 文件

## 🙏 致谢
> 以下排名不分先后

- [Gin](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [Redis](https://github.com/redis/go-redis)
- [Viper](https://github.com/spf13/viper)
- [Zap](https://github.com/uber-go/zap)
- [Wire](https://github.com/google/wire)
- [JWT](https://github.com/golang-jwt/jwt)
- [go-gomail](https://github.com/go-gomail/gomail)

---

⭐️ 如果这个项目对你有帮助，请给它一个 start！
