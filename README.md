# 基于 Gin 的轻量级脚手架

一个基于 Gin 框架的轻量级脚手架，集成了常用组件，帮助你快速搭建高性能的 Go Web 应用。

[![Go Version](https://img.shields.io/badge/Go-v1.18+-blue.svg)](https://golang.org/doc/devel/release.html)
[![Gin](https://img.shields.io/badge/Gin-v1.9.0+-green.svg)](https://github.com/gin-gonic/gin)
[![SQLBoiler](https://img.shields.io/badge/SQLBoiler-v4.14.0+-orange.svg)](https://github.com/volatiletech/sqlboiler)

## 🚀 特性

- 📝 完整的项目结构和最佳实践
- 🔒 JWT 认证集成
- 📊 统一的 API 响应格式
- 🔄 强大的中间件支持
- 📋 详尽的日志记录
- 🔌 多数据库支持
- 🛠️ 优雅的错误处理
- 🚦 优雅启动和关闭
- 🎇 不计其数的优雅设计

## 🔧 技术栈

- [Gin](https://github.com/gin-gonic/gin) - 高性能 HTTP Web 框架
- [Validator](https://github.com/go-playground/validator) - 参数验证库
- [SQLBoiler](https://github.com/volatiletech/sqlboiler) - 优秀的 ORM 库，基于代码生成
- [Redis](https://github.com/redis/go-redis) - Redis 客户端
- [JWT](https://github.com/golang-jwt/jwt) - JWT 鉴权管理
- [Zap](https://github.com/uber-go/zap) - 高性能、结构化日志
- [Wire](https://github.com/google/wire) - Wire 依赖注入

## 📁 项目结构

```
scaffold/
├── api/
│   ├── openapi             #swagger 文件
│   └── protobuf            # pb文件 (普通项目规模不推荐使用微服务)
├── docker/
├── internal/
│   ├── common
│   │    ├── email/         # email相关
│   │    ├── jwt/           # jwt相关
│   │    ├── logger/        # 日志配置
│   │    ├── metrics/       # 指数收集
│   │    ├── middleware/    # 中间件
│   │    ├── orm/           # SQLBoiler生成代码
│   │    ├── reqkit         # 自封装请求工具集
│   │    ├── reskit/        # 自封装响应工具集
│   │    ├── server/        # 服务配置
│   │    ├── uid/           # 唯一键生成
│   │    ├── utils/         # utils工具函数
│   │    └── validator/     # validator管理
│   └── captcha             # 验证码模块
│   └── img                 # 图库模块
│   └── user                # 用户模块
│   └── ...                 # 其余模块
├── logs/                   # 日志文件
├── tool/                   # 工具脚本
├── .air.conf               # air配置
├── .env                    # 开发环境变量(由下方得来)
├── .env.copy               # 开发环境变量
├── .env.docker             # 生产环境变量(由下方得来)
├── .env.docker_copy        # 生产环境变量
├── ....
├── main.go                 # 主入口
└── README.md
└── sqlboiler.toml          # sqlboiler配置
```

## ⚡ 快速开始

### 前置要求

- Go 1.18+
- PostgreSQL 10+
- Redis 6.0+

### 快速开始

> 以下演示以Windows作为示例

1. 新建目录
```bash
mkdir demo
cd demo
```

2. 克隆项目

```bash
git clone https://github.com/Lirous587/go-scaffold.git
```

3. 移动目录 并删除git记录
```bash
robocopy go-scaffold . /E /XD .git
```

4. 删除clone目录
```bash
Remove-Item go-scaffold -Recurse -Force
```

5. 编写并运行replace脚本

```bash
cd ./tool/replace
go build
./replace.exe demo # 填写想要的实际module名
```

6. 删除replace
```bash
cd ..
rm ./replace
```

7. 安装依赖
```bash
cd ..
go mod tidy
```

8. 使用docker构建依赖
> 开发阶段无需使用docker去构建应用程序，故可将 `docker-compose.yml` 中 `service` 下的 `go-app`注释，以此来快速构建依赖
- 修改 `docker` 目录下的 `.env` 配置
- 在 `docker` 目录下运行 `docker compose up -d --build`

9. 修改配置
- 开发环境:将 `.copy.env` 重命名为 `.env`，配置 `.env`
- 生产环境:将 `.copy.docker_copy` 重命名为 `.env.docker`，配置 `.env.docker`

10. 使用gen工具(可选)
- 根路径下运行
```bash
go run ./tool/gen/gen.go -m mock
```
- 修改入口文件main函数的 `server.RunHttpServer`
```go
server.RunHttpServer(os.Getenv("SERVER_PORT"), metricsClient, func(r *gin.RouterGroup) {
    // ......
    // 新增
    mock.InitV1(r)
})
```

11. 运行服务
```bash
# 或者运行 air
go run main.go
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
- [Validator](https://github.com/go-playground/validator)
- [SQLBoiler](https://github.com/volatiletech/sqlboiler)
- [Redis](https://github.com/redis/go-redis)
- [JWT](https://github.com/golang-jwt/jwt)
- [Zap](https://github.com/uber-go/zap)
- [Wire](https://github.com/google/wire)

---

⭐️ 如果这个项目对你有帮助，请给它一个 start！
