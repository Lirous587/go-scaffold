package apigen

import (
	"fmt"
	"reflect"
	"scaffold/pkg/httpserver/middleware/bind"

	"github.com/gin-gonic/gin"
)

// IRoutes 接口同时被gin.Engine和gin.RouterGroup实现
type IRoutes interface {
	GET(string, ...gin.HandlerFunc) gin.IRoutes
	POST(string, ...gin.HandlerFunc) gin.IRoutes
	PUT(string, ...gin.HandlerFunc) gin.IRoutes
	DELETE(string, ...gin.HandlerFunc) gin.IRoutes
	// 可以根据需要添加其他HTTP方法
}

// RegisterAPI 自动注册API
func RegisterAPI(router IRoutes, controller interface{}) {
	// 获取控制器的反射值和类型
	controllerValue := reflect.ValueOf(controller)
	controllerType := controllerValue.Type()

	// 遍历控制器的所有方法
	for i := 0; i < controllerType.NumMethod(); i++ {
		methodInfo := controllerType.Method(i)

		// 验证并获取API信息
		apiInfo := BindApiInfo(methodInfo)
		if !apiInfo.Valid {
			continue
		}

		// 获取Gin路径
		ginPath := apiInfo.RouteInfo.GetGinPath()

		// 注册路由
		method := controllerValue.Method(i).Interface()
		switch apiInfo.RouteInfo.Method {
		case "get", "GET":
			router.GET(ginPath, bind.AutoBind(method))
		case "post", "POST":
			router.POST(ginPath, bind.AutoBind(method))
		case "put", "PUT":
			router.PUT(ginPath, bind.AutoBind(method))
		case "delete", "DELETE":
			router.DELETE(ginPath, bind.AutoBind(method))
		default:
			fmt.Printf("不支持的HTTP方法: %s\n", apiInfo.RouteInfo.Method)
			continue
		}
	}
}
