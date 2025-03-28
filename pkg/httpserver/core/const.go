package core

import "scaffold/response"

// Response 统一API响应结构
type Response struct {
	Code    response.Code `json:"code"`
	Message interface{}   `json:"message"`
	Data    interface{}   `json:"data"`
}
