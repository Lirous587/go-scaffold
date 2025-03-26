package bind

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"scaffold/pkg/i18n"
)

func (b *bind) bindJson(req interface{}) {
	// 检查请求体是否为空
	if b.ctx.Request.Body == nil {
		errorMsg, errorDetail := i18n.TranslateNullJSONError(b.ctx.Err(), b.lang)
		b.ctx.Error(errors.New(errorMsg + ":" + errorDetail))
		b.ctx.Abort()
		return
	}

	bodyBytes, err := io.ReadAll(b.ctx.Request.Body)
	if err != nil {
		errorMsg, errorDetail := i18n.TranslateJSONError(err, b.lang)
		b.ctx.Error(errors.New(errorMsg + ":" + errorDetail))
		b.ctx.Abort()
		return
	}

	// 重置请求体供后续使用
	b.ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// 检查读取的内容是否为空
	if len(bodyBytes) == 0 {
		errorMsg, errorDetail := i18n.TranslateNullJSONError(io.EOF, b.lang)
		b.ctx.Error(errors.New(errorMsg + ":" + errorDetail))
		b.ctx.Abort()
		return
	}

	// 解析JSON
	if err := json.Unmarshal(bodyBytes, req); err != nil {
		errorMsg, errorDetail := i18n.TranslateJSONError(err, b.lang)
		b.ctx.Error(errors.New(errorMsg + ":" + errorDetail))
		b.ctx.Abort()
		return
	}
}
