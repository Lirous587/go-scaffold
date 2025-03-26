package apigen

import (
	"context"
	"fmt"
	"reflect"
)

type ApiInfo struct {
	RouteInfo RouteInfo    // 路由信息
	ReqType   reflect.Type // 请求类型
	ResType   reflect.Type // 响应类型
	Method    interface{}  // 方法
}

// GetApiInfo 从定义的api接口中提取数据
func GetApiInfo(apiInterface interface{}) []ApiInfo {
	// 获取控制器的反射值和类型
	apiInterfaceValue := reflect.ValueOf(apiInterface)
	apiInterfaceType := apiInterfaceValue.Type()

	// 确保提供的是接口类型
	if apiInterfaceType.Kind() != reflect.Ptr && apiInterfaceType.Kind() != reflect.Interface {
		panic(fmt.Sprintf("apiInterface 必须是接口或接口指针类型，实际为 %s 类型", apiInterfaceType.Kind()))
	}

	apiInfos := make([]ApiInfo, 0, apiInterfaceType.NumMethod())

	// 遍历控制器的所有方法
	for i := 0; i < apiInterfaceType.NumMethod(); i++ {
		methodInfo := apiInterfaceType.Method(i)
		funcType := methodInfo.Type

		// 1. 检查方法签名
		// 3个入参：方法接受者 + context + 请求指针
		// 2个返回值：响应指针 + error
		if funcType.NumIn() != 3 || funcType.NumOut() != 2 {
			panic(fmt.Sprintf("方法 %s 签名错误: 需要3个入参和2个返回值，实际有 %d 个入参和 %d 个返回值",
				methodInfo.Name, funcType.NumIn(), funcType.NumOut()))
		}

		// 2. 检查第二个参数是否为context.Context
		contextInterface := reflect.TypeOf((*context.Context)(nil)).Elem()
		if !funcType.In(1).Implements(contextInterface) {
			panic(fmt.Sprintf("方法 %s 的第二个参数必须是 context.Context 类型，实际为 %s 类型",
				methodInfo.Name, funcType.In(1)))
		}

		// 3. 获取请求参数类型
		reqType := funcType.In(2)
		if reqType.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("方法 %s 的第三个参数必须是指针类型，实际为 %s 类型",
				methodInfo.Name, reqType.Kind()))
		}

		// 4. 检查返回值类型
		resType := funcType.Out(0)
		if resType.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("方法 %s 的第一个返回值必须是指针类型，实际为 %s 类型",
				methodInfo.Name, resType.Kind()))
		}

		// 5. 检查第二个返回值是否为error
		errorInterface := reflect.TypeOf((*error)(nil)).Elem()
		if !funcType.Out(1).Implements(errorInterface) {
			panic(fmt.Sprintf("方法 %s 的第二个返回值必须是 error 类型，实际为 %s 类型",
				methodInfo.Name, funcType.Out(1)))
		}

		// 6. 获取指针指向的类型
		reqElemType := reqType.Elem()
		resElemType := resType.Elem()

		// 7. 尝试获取Meta字段
		_, found := reqElemType.FieldByName("Meta")
		if !found {
			panic(fmt.Sprintf("方法 %s 的请求结构体 %s 缺少必需的 Meta 字段",
				methodInfo.Name, reqElemType.Name()))
		}

		// 8. 创建请求实例
		reqInstance := reflect.New(reqElemType).Interface()

		// 9. 获取路由信息
		metaField := reflect.ValueOf(reqInstance).Elem().FieldByName("Meta")
		iface := metaField.Interface()
		metaVal, ok := iface.(Meta)
		if !ok {
			panic(fmt.Sprintf("方法 %s 的请求结构体 %s 中的 Meta 字段类型 %T 不是 api.Meta",
				methodInfo.Name, reqElemType.Name(), iface))
		}

		// 10. 获取路由信息并保存所有结果
		routeInfo := metaVal.GetRouteInfo(reqInstance)

		// 11. 设置结果
		apiInfo := ApiInfo{
			RouteInfo: routeInfo,
			ReqType:   reqElemType,
			ResType:   resElemType,
			Method:    apiInterfaceValue.Method(i).Interface(),
		}
		apiInfos = append(apiInfos, apiInfo)
	}

	return apiInfos
}
