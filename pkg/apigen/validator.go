package apigen

import (
	"fmt"
	"reflect"
)

// APIMethodInfo 保存API方法的所有验证结果
type APIMethodInfo struct {
	RouteInfo   RouteInfo    // 路由信息
	ReqType     reflect.Type // 请求类型
	ResType     reflect.Type // 响应类型
	MethodValue interface{}  // 方法值
	Valid       bool         // 是否有效
}

// ValidateMethod 验证方法并提取所有需要的信息
func ValidateMethod(controller interface{}, methodInfo reflect.Method) APIMethodInfo {
	result := APIMethodInfo{
		Valid: false,
	}

	funcType := methodInfo.Type

	// 1. 检查方法签名
	// 3 个入参：接收者 + context + 请求指针
	// 2 个返回值：响应指针 + error
	if funcType.NumIn() != 3 || funcType.NumOut() != 2 {
		return result
	}

	// 2. 获取请求参数类型
	reqType := funcType.In(2)
	if reqType.Kind() != reflect.Ptr {
		return result
	}

	// 3. 检查返回值类型
	resType := funcType.Out(0)
	if resType.Kind() != reflect.Ptr {
		return result
	}

	// 4. 检查第二个返回值是否为error
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	if !funcType.Out(1).Implements(errorInterface) {
		return result
	}

	// 5. 获取指针指向的类型
	reqElemType := reqType.Elem()
	resElemType := resType.Elem()

	// 6. 尝试获取Meta字段
	_, found := reqElemType.FieldByName("Meta")
	if !found {
		return result
	}

	// 7. 创建请求实例
	reqInstance := reflect.New(reqElemType).Interface()

	// 8. 获取路由信息
	metaField := reflect.ValueOf(reqInstance).Elem().FieldByName("Meta")
	iface := metaField.Interface()
	metaVal, ok := iface.(Meta)
	if !ok {
		fmt.Printf("警告：%s.Meta 字段类型 %T 不是 api.Meta\n",
			reqElemType.Name(), iface)
		return result
	}

	// 9. 获取路由信息并保存所有结果
	routeInfo := metaVal.GetRouteInfo(reqInstance)

	// 10. 设置结果
	result.RouteInfo = routeInfo
	result.ReqType = reqElemType
	result.ResType = resElemType
	result.Valid = true

	return result
}
