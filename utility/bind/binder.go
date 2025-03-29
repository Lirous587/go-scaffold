package bind

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
	"net/http"
	"reflect"
	"scaffold/pkg/i18n"
	"scaffold/pkg/validator"
	"strings"
)

// Binder 负责将请求数据绑定到结构体
type Binder struct {
	ctx  *gin.Context
	lang string
}

// NewBinder 创建新的绑定器
func NewBinder(ctx *gin.Context) *Binder {
	lang := ctx.GetHeader("accept-language")

	return &Binder{
		ctx:  ctx,
		lang: lang,
	}
}

func Bind(c *gin.Context, req any) error {
	// 创建绑定器并执行绑定
	binder := NewBinder(c)
	if err := binder.Bind(req); err != nil {
		return err
	}

	return nil
}

// Bind 将请求数据智能绑定到提供的结构体指针
func (b *Binder) Bind(req interface{}) error {
	reqValue := reflect.ValueOf(req)
	if reqValue.Kind() != reflect.Ptr || reqValue.IsNil() {
		return errors.New("Invalid request object" + "Request must be a non-nil pointer")
	}

	reqElem := reqValue.Elem()
	reqType := reqElem.Type()

	// 1. 绑定URI参数
	b.bindUri(reqElem, reqType)

	// 2. 绑定请求头
	b.bindHeader(reqElem, reqType)

	// 3. 绑定查询参数
	b.bindQuery(reqElem, reqType)

	// 4. 为非GET/DELETE方法处理请求体
	method := b.ctx.Request.Method
	if method != http.MethodGet && method != http.MethodDelete {
		contentType := b.ctx.ContentType()

		var err error
		if contentType == "application/json" {
			err = b.bindJSON(req)
		} else if contentType == "application/x-www-form-urlencoded" ||
			strings.Contains(contentType, "multipart/form-data") {
			err = b.bindForm(reqElem, reqType)
		}

		if err != nil {
			return err
		}
	}

	// 验证参数
	if err := validator.V.Struct(req); err != nil {
		errMsg := i18n.TranslateValidatorError(err, b.lang)
		return errMsg
	}

	return nil
}

func (b *Binder) bindUri(reqElem reflect.Value, reqType reflect.Type) {
	for i := 0; i < reqType.NumField(); i++ {
		field := reqType.Field(i)
		fieldVal := reqElem.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		// 处理 uri 标签
		if uriTag := field.Tag.Get("uri"); uriTag != "" {
			if paramValue := b.ctx.Param(uriTag); paramValue != "" {
				setFieldValue(fieldVal, paramValue)
			}
		}
	}
}

func (b *Binder) bindHeader(reqElem reflect.Value, reqType reflect.Type) {
	for i := 0; i < reqType.NumField(); i++ {
		field := reqType.Field(i)
		fieldVal := reqElem.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		// 处理 header 标签
		headerName := field.Name
		if headerTag := field.Tag.Get("header"); headerTag != "" {
			headerName = headerTag
		}

		if headerValue := b.ctx.GetHeader(headerName); headerValue != "" {
			setFieldValue(fieldVal, headerValue)
		}
	}
}

func (b *Binder) bindQuery(reqElem reflect.Value, reqType reflect.Type) {
	queryValues := b.ctx.Request.URL.Query()

	for i := 0; i < reqType.NumField(); i++ {
		field := reqType.Field(i)
		fieldVal := reqElem.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		// 处理 query 标签
		if queryTag := field.Tag.Get("query"); queryTag != "" {
			if values, exists := queryValues[queryTag]; exists && len(values) > 0 {
				setFieldValue(fieldVal, values[0])
			}
		}
	}
}

func (b *Binder) bindJSON(req interface{}) error {
	// 检查请求体是否为空
	if b.ctx.Request.Body == nil {
		errorMsg, errorDetail := i18n.TranslateNullJSONError(b.ctx.Err(), b.lang)
		return errors.New(errorMsg + ":" + errorDetail)
	}

	bodyBytes, err := io.ReadAll(b.ctx.Request.Body)
	if err != nil {
		errorMsg, errorDetail := i18n.TranslateJSONError(err, b.lang)
		return errors.New(errorMsg + ":" + errorDetail)
	}

	// 重置请求体供后续使用
	b.ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// 检查读取的内容是否为空
	if len(bodyBytes) == 0 {
		errorMsg, errorDetail := i18n.TranslateNullJSONError(io.EOF, b.lang)
		return errors.New(errorMsg + ":" + errorDetail)
	}

	// 解析JSON
	if err := json.Unmarshal(bodyBytes, req); err != nil {
		errorMsg, errorDetail := i18n.TranslateJSONError(err, b.lang)
		return errors.New(errorMsg + ":" + errorDetail)
	}
	return nil
}

func (b *Binder) bindForm(reqElem reflect.Value, reqType reflect.Type) error {
	contentType := b.ctx.ContentType()
	// 根据Content-Type选择合适的解析方法
	if strings.Contains(contentType, "multipart/form-data") {
		// 对于multipart/form-data类型
		if err := b.ctx.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB 最大尺寸
			errorMsg, errorDetail := i18n.TranslateParseFormError(err, b.lang)
			return errors.New(errorMsg + ":" + errorDetail)
		}
	} else {
		// 对于x-www-form-urlencoded类型
		if err := b.ctx.Request.ParseForm(); err != nil {
			errorMsg, errorDetail := i18n.TranslateParseFormError(err, b.lang)
			return errors.New(errorMsg + ":" + errorDetail)
		}
	}

	// 遍历结构体字段
	for i := 0; i < reqType.NumField(); i++ {
		field := reqType.Field(i)
		fieldVal := reqElem.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		formTag := field.Tag.Get("form")

		if formTag == "" {
			continue
		}

		// 先处理文件类型字段
		if fieldVal.Type() == reflect.TypeOf((*multipart.FileHeader)(nil)) {
			if strings.Contains(contentType, "multipart/form-data") {
				file, err := b.ctx.FormFile(formTag)
				if err == nil {
					fieldVal.Set(reflect.ValueOf(file))
				}
			}
			continue
		}

		// 处理其它普通字段
		value := b.ctx.Request.FormValue(formTag)

		if value == "" {
			continue
		}

		// 设置字段值
		setFieldValue(fieldVal, value)
	}
	return nil
}
