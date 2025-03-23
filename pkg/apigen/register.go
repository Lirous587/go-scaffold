package apigen

import (
	"fmt"
	"reflect"
	"scaffold/pkg/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterAPI 自动注册API
func RegisterAPI(router *gin.Engine, controller interface{}) {
	// 1. 为控制器注册路由
	registerControllerRoutes(router, controller)

	// 2. 生成Swagger文档
	err := GenerateSwaggerDocs(controller)
	if err != nil {
		fmt.Printf("生成Swagger文档失败: %v\n", err)
	}

	// 3. 添加Swagger UI路由
	router.StaticFile("/swagger-docs/swagger.json", "./docs/swagger.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/swagger-docs/swagger.json")))
}

// 注册控制器的路由
func registerControllerRoutes(router *gin.Engine, controller interface{}) {
	// 获取控制器的反射值和类型
	controllerValue := reflect.ValueOf(controller)
	controllerType := controllerValue.Type()

	// 遍历控制器的所有方法
	for i := 0; i < controllerType.NumMethod(); i++ {
		methodInfo := controllerType.Method(i)

		// 验证并获取API信息
		apiInfo := ValidateMethod(controller, methodInfo)
		if !apiInfo.Valid {
			continue
		}

		// 获取Gin路径
		ginPath := apiInfo.RouteInfo.GetGinPath()

		// 注册路由
		method := controllerValue.Method(i).Interface()
		switch apiInfo.RouteInfo.Method {
		case "get", "GET":
			router.GET(ginPath, middleware.AutoBind(method))
		case "post", "POST":
			router.POST(ginPath, middleware.AutoBind(method))
		case "put", "PUT":
			router.PUT(ginPath, middleware.AutoBind(method))
		case "delete", "DELETE":
			router.DELETE(ginPath, middleware.AutoBind(method))
		default:
			fmt.Printf("不支持的HTTP方法: %s\n", apiInfo.RouteInfo.Method)
			continue
		}
	}
}
