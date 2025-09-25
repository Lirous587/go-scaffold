package response

// 用于文档生成
type invalidParamsResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message" example:"invalid params"`
	Details map[string]interface{} `json:"details,omitempty"`
}

type errorResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message" example:"Internal server error"`
	Details map[string]interface{} `json:"details,omitempty"`
}
