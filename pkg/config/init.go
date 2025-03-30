package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var Cfg Config

func Init() error {
	// 初始化viper
	if err := initViper(); err != nil {
		return errors.WithMessage(err, "initViper failed")
	}

	// 解析配置到结构体
	if err := loadConfig(); err != nil {
		return errors.WithMessage(err, "loadConfig failed")
	}

	if err := validationConfig(); err != nil {
		return errors.WithMessage(err, "validationConfig failed")
	}

	return nil
}

// 加载配置到结构体
func loadConfig() error {
	if err := viper.Unmarshal(&Cfg); err != nil {
		return errors.Wrap(err, "无法解析配置到结构体")
	}

	//// 调试输出
	//fmt.Printf("读取到的配置: %+v\n", Cfg)
	//if len(Cfg.Server) > 0 {
	//	fmt.Printf("Server配置: %+v\n", Cfg.Server[0])
	//} else {
	//	fmt.Println("Server配置为空")
	//}
	return nil
}

func validationConfig() error {
	// 使用反射验证配置
	return validateStruct("", Cfg)
}

// validateStruct 使用反射递归验证结构体
func validateStruct(prefix string, s interface{}) error {
	val := reflect.ValueOf(s)
	typ := val.Type()

	// 如果是指针，获取其指向的实际值
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return errors.New(fmt.Sprintf("配置项 %s 不能为 nil", prefix))
		}
		val = val.Elem()
		typ = val.Type()
	}

	// 确保我们处理的是结构体
	if val.Kind() != reflect.Struct {
		return nil // 不是结构体，跳过验证
	}

	// 遍历结构体字段
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// 获取字段名称（使用mapstructure标签或字段名）
		fieldName := fieldType.Tag.Get("mapstructure")
		if fieldName == "" || fieldName == "-" {
			fieldName = strings.ToLower(fieldType.Name)
		}

		// 构建完整字段路径
		fieldPath := fieldName
		if prefix != "" {
			fieldPath = prefix + "." + fieldName
		}

		// 根据字段类型进行验证
		switch field.Kind() {
		case reflect.Struct:
			// 递归验证嵌套结构体
			if err := validateStruct(fieldPath, field.Interface()); err != nil {
				return errors.WithMessage(err, fmt.Sprintf("验证子配置 %s 失败", fieldPath))
			}

		case reflect.Ptr:
			// 递归验证指针类型
			if !field.IsNil() {
				if err := validateStruct(fieldPath, field.Interface()); err != nil {
					return errors.WithMessage(err, fmt.Sprintf("验证子配置 %s 失败", fieldPath))
				}
			}

		case reflect.Slice, reflect.Array:
			// 处理切片或数组类型
			for j := 0; j < field.Len(); j++ {
				sliceElem := field.Index(j)
				elemPath := fmt.Sprintf("%s[%d]", fieldPath, j)

				if sliceElem.Kind() == reflect.Struct || sliceElem.Kind() == reflect.Ptr {
					if err := validateStruct(elemPath, sliceElem.Interface()); err != nil {
						return errors.WithMessage(err, fmt.Sprintf("验证切片元素 %s 失败", elemPath))
					}
				}
			}

		case reflect.String:
			// 验证必要的字符串字段不为空
			if !isPermitEmpty(fieldPath) && field.String() == "" {
				return errors.New(fmt.Sprintf("缺少配置: %s", fieldPath))
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// 验证必要的数字字段大于0
			if !isPermitEmpty(fieldPath) && field.Int() <= 0 {
				return errors.New(fmt.Sprintf("无效配置: %s 必须大于0", fieldPath))
			}

		}
	}
	return nil
}

func isPermitEmpty(fieldPath string) bool {
	permitEmptyFields := map[string]bool{
		"redis.password": true,
		"redis.db":       true,
	}
	return permitEmptyFields[fieldPath]
}
