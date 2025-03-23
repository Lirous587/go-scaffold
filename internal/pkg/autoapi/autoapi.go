package autoapi

import (
	"fmt"
	"reflect"
	"scaffold/api"
	"scaffold/internal/middleware"
	"scaffold/internal/pkg/swagger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterAPI 自动注册API
func RegisterAPI(router *gin.Engine, controller interface{}) {
	// 1. 为控制器注册路由
	registerControllerRoutes(router, controller)

	// 2. 生成Swagger文档
	err := swagger.GenerateSwaggerDocs(controller)
	if err != nil {
		fmt.Printf("生成Swagger文档失败: %v\n", err)
	}

	// 3. 添加Swagger UI路由
	router.StaticFile("/swagger-docs/swagger.json", "./swagger/swagger.json")
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

		// 检查方法参数是否符合要求
		// 3 个入参：接收者 + context + 请求指针
		// 2 个返回值：响应指针 + error
		if methodInfo.Type.NumIn() != 3 || methodInfo.Type.NumOut() != 2 {
			continue
		}

		// 获取请求参数类型
		// 方法的第 3 个参数(索引 2)的类型
		reqType := methodInfo.Type.In(2)
		if reqType.Kind() != reflect.Ptr {
			continue
		}

		// 获取指针指向的类型
		// 如果 *xxReq 是指针类型，Elem() 则返回 xxReq 结构体类型
		reqType = reqType.Elem()

		// 尝试获取Meta字段
		_, found := reqType.FieldByName("Meta")
		if !found {
			continue
		}

		// 创建请求实例
		// 创建指定类型的新实例，返回指向它的指针值;类似于 new(Type) 但使用反射动态创建
		reqInstance := reflect.New(reqType).Interface()

		// 获取路由信息
		// reflect.ValueOf(reqInstance)：获取 reqInstance 的反射值
		// .Elem()：获取指针指向的值（解引用） 返回结构体类型
		// .FieldByName("Meta")：获取名为 "Meta" 的字段值
		// .Interface()：将字段值转回普通 Go 值
		// .(api.Meta)：类型断言，将接口值转换为具体类型 api.Meta
		// metaVal.GetRouteInfo(reqInstance)：调用 Meta 结构的方法获取路由信息
		metaField := reflect.ValueOf(reqInstance).Elem().FieldByName("Meta")
		iface := metaField.Interface()
		metaVal, ok := iface.(api.Meta)
		if !ok {
			fmt.Printf("警告：%s.Meta 字段类型 %T 不是 api.Meta\n",
				reqType.Name(), iface)
			continue
		}
		routeInfo := metaVal.GetRouteInfo(reqInstance)

		// 获取Gin路径
		// 调用 routeInfo 的方法将 OpenAPI 风格的路径转换为 Gin 风格
		// 例如：/xx/{id} → /xx/:id
		ginPath := routeInfo.GetGinPath()

		// 注册路由
		method := controllerValue.Method(i).Interface()
		switch routeInfo.Method {
		case "get", "GET":
			router.GET(ginPath, middleware.AutoBind(method))
		case "post", "POST":
			router.POST(ginPath, middleware.AutoBind(method))
		case "put", "PUT":
			router.PUT(ginPath, middleware.AutoBind(method))
		case "delete", "DELETE":
			router.DELETE(ginPath, middleware.AutoBind(method))
		default:
			fmt.Printf("不支持的HTTP方法: %s\n", routeInfo.Method)
			continue
		}
	}
}
