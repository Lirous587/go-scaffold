package bind

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type TestURI struct {
	ID     string `uri:"id" validate:"required"`
	Action string `uri:"action" validate:"required"`
}

func TestBindURI(t *testing.T) {
	// 设置Gin路由
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/users/:id/:action", func(c *gin.Context) {
		var req TestURI
		err := Bind(c, &req)
		assert.NoError(t, err)
		assert.Equal(t, "123", req.ID)
		assert.Equal(t, "fuck", req.Action)
		c.String(200, "ok")
	})

	// 创建测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/123/fuck", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

type TestJSON struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age"`
}

func TestBindJSON(t *testing.T) {
	// 设置Gin路由
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/users", func(c *gin.Context) {
		var req TestJSON
		err := Bind(c, &req)
		assert.NoError(t, err)
		assert.Equal(t, "张三", req.Name)
		assert.Equal(t, 25, req.Age)
		c.String(200, "ok")
	})

	// 创建测试请求
	w := httptest.NewRecorder()
	jsonData := `{"name":"lirous","age":19}`
	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestBindJSON1(t *testing.T) {
	// 设置Gin路由
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/users", func(c *gin.Context) {
		var req TestJSON
		err := Bind(c, &req)
		assert.NoError(t, err)
		assert.Equal(t, "lirous", req.Name)
		assert.Equal(t, 0, req.Age)
		c.String(200, "ok")
	})

	// 创建测试请求
	w := httptest.NewRecorder()
	jsonData := `{"name":"张三"}`
	req, _ := http.NewRequest("POST", "/users", bytes.NewBufferString(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

type TestForm struct {
	Name string `form:"name" validate:"required"`
	Age  int    `form:"age"`
}

func TestBindForm(t *testing.T) {
	// 设置Gin路由
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/users", func(c *gin.Context) {
		var req TestForm
		err := Bind(c, &req)
		assert.NoError(t, err)
		assert.Equal(t, "lirous", req.Name)
		assert.Equal(t, 19, req.Age)
		c.String(200, "ok")
	})

	// 创建表单数据
	form := url.Values{}
	form.Add("name", "lirous")
	form.Add("age", "30")

	// 创建测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestBindForm2(t *testing.T) {
	// 设置Gin路由
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/users", func(c *gin.Context) {
		var req TestForm
		err := Bind(c, &req)
		assert.NoError(t, err)
		assert.Equal(t, "lirous", req.Name)
		assert.Equal(t, 0, req.Age)
		c.String(200, "ok")
	})

	// 创建表单数据
	form := url.Values{}
	form.Add("name", "lirous")
	form.Add("age", "0")

	// 创建测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

type TestMultipartForm struct {
	Name     string                `form:"name" validate:"required"`
	Avatar   *multipart.FileHeader `form:"avatar"`
	Document *multipart.FileHeader `form:"document"`
}

func TestBindMultipartForm(t *testing.T) {
	// 设置Gin路由
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/upload", func(c *gin.Context) {
		var req TestMultipartForm
		err := Bind(c, &req)
		assert.NoError(t, err)
		assert.Equal(t, "lirous", req.Name)

		// 检查文件是否成功绑定
		assert.NotNil(t, req.Avatar)
		assert.Equal(t, "avatar.png", req.Avatar.Filename)
		assert.NotNil(t, req.Document)
		assert.Equal(t, "doc.pdf", req.Document.Filename)

		c.String(200, "ok")
	})

	// 创建multipart表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加普通字段
	_ = writer.WriteField("name", "lirous")

	// 添加文件
	avatarPart, _ := writer.CreateFormFile("avatar", "avatar.png")
	avatarPart.Write([]byte("模拟图片内容"))
	docPart, _ := writer.CreateFormFile("document", "doc.pdf")
	docPart.Write([]byte("模拟PDF内容"))

	writer.Close()

	// 创建测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

type TestQuery struct {
	Page     int    `query:"page" validate:"required"`
	Category string `query:"category" validate:"required"`
}

func TestBindQuery(t *testing.T) {
	// 设置Gin路由
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/users", func(c *gin.Context) {
		var req TestQuery
		err := Bind(c, &req)
		assert.NoError(t, err)
		assert.Equal(t, 2, req.Page)
		assert.Equal(t, "go", req.Category)
		c.String(200, "ok")
	})

	// 创建测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users?page=2&&category=go", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

type TestHeader struct {
	Token     string `header:"Authorization"`
	UserAgent string `header:"User-Agent"`
}

func TestBindHeader(t *testing.T) {
	// 设置Gin路由
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/users", func(c *gin.Context) {
		var req TestHeader
		err := Bind(c, &req)
		assert.NoError(t, err)
		assert.Equal(t, "Bearer token123", req.Token)
		assert.Equal(t, "TestAgent", req.UserAgent)
		c.String(200, "ok")
	})

	// 创建测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	req.Header.Set("Authorization", "Bearer token123")
	req.Header.Set("User-Agent", "TestAgent")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

type TestAllBindings struct {
	// URI 参数
	ID     string `uri:"id" validate:"required"`
	Action string `uri:"action" validate:"required"`

	// 查询参数
	Page     int    `query:"page" validate:"required"`
	Category string `query:"category" validate:"required"`

	// 请求头
	Token     string `header:"Authorization"`
	UserAgent string `header:"User-Agent"`

	// 表单字段
	Name string `form:"name" validate:"required"`
	Age  int    `form:"age"`

	// 文件字段
	Avatar   *multipart.FileHeader `form:"avatar"`
	Document *multipart.FileHeader `form:"document"`
}

func TestBindAll(t *testing.T) {
	// 设置Gin路由
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/users/:id/:action", func(c *gin.Context) {
		var req TestAllBindings
		err := Bind(c, &req)
		assert.NoError(t, err)

		// 验证URI参数
		assert.Equal(t, "123", req.ID)
		assert.Equal(t, "update", req.Action)

		// 验证查询参数
		assert.Equal(t, 2, req.Page)
		assert.Equal(t, "go", req.Category)

		// 验证请求头
		assert.Equal(t, "Bearer token123", req.Token)
		assert.Equal(t, "TestAgent", req.UserAgent)

		// 验证表单字段
		assert.Equal(t, "lirous", req.Name)
		assert.Equal(t, 25, req.Age)

		// 验证文件字段
		assert.NotNil(t, req.Avatar)
		assert.Equal(t, "avatar.png", req.Avatar.Filename)
		assert.NotNil(t, req.Document)
		assert.Equal(t, "doc.pdf", req.Document.Filename)

		c.String(200, "ok")
	})

	// 创建multipart表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加普通字段
	_ = writer.WriteField("name", "lirous")
	_ = writer.WriteField("age", "25")

	// 添加文件
	avatarPart, _ := writer.CreateFormFile("avatar", "avatar.png")
	avatarPart.Write([]byte("模拟图片内容"))
	docPart, _ := writer.CreateFormFile("document", "doc.pdf")
	docPart.Write([]byte("模拟PDF内容"))

	writer.Close()

	// 创建测试请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/123/update?page=2&category=go", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer token123")
	req.Header.Set("User-Agent", "TestAgent")
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
