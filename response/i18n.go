package response

//// 全局多语言响应映射表
//var messageI18nMap map[string]map[Code]string
//
//// GetMessageI18n 根据Code和语言获取预定义消息
//func GetMessageI18n(code Code, lang string) (string, bool) {
//	if langMap, ok := messageI18nMap[lang]; ok {
//		if message, exists := langMap[code]; exists {
//			return message, true
//		}
//	}
//
//	// 回退到默认语言
//	return GetMessage(code)
//}
//
//// RegisterI18nMessage 注册多语言消息
//func RegisterI18nMessage(lang string, code Code, message string) {
//	if _, ok := messageI18nMap[lang]; !ok {
//		messageI18nMap[lang] = make(map[Code]string)
//	}
//	messageI18nMap[lang][code] = message
//}
//
//func init() {
//	messageI18nMap = make(map[string]map[Code]string)
//
//	// 示例：注册英文消息
//	enMessages := map[Code]string{
//		CodeSuccess:         "Operation successful",
//		CodeValidationError: "Validation failed",
//		CodeJSONError:       "JSON parsing error",
//		CodeEmptyBodyError:  "Empty request body",
//		CodeServerError:     "Internal server error",
//	}
//
//	for code, message := range enMessages {
//		RegisterI18nMessage("en", code, message)
//	}
//}
