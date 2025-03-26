package apigen

import (
	"fmt"
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
func RegisterAPI(router IRoutes, apiInfos []ApiInfo) {
	for _, apiInfo := range apiInfos {
		// 获取Gin路径
		ginPath := apiInfo.RouteInfo.GetGinPath()

		// 注册路由

		switch apiInfo.RouteInfo.Method {
		case "get", "GET":
			router.GET(ginPath, bind.AutoBind(apiInfo.Method))
		case "post", "POST":
			router.POST(ginPath, bind.AutoBind(apiInfo.Method))
		case "put", "PUT":
			router.PUT(ginPath, bind.AutoBind(apiInfo.Method))
		case "delete", "DELETE":
			router.DELETE(ginPath, bind.AutoBind(apiInfo.Method))
		default:
			fmt.Printf("不支持的HTTP方法: %s\n", apiInfo.RouteInfo.Method)
			continue
		}
	}
}
