package utils

import "github.com/pkg/errors"

func ComputeOffset(page int, size int) (int, error) {
	if page < 1 {
		return 0, errors.New("页码必须大于0")
	}
	if size < 0 {
		return 0, errors.New("页面大小不能为负数")
	}
	return (page - 1) * size, nil
}

func ComputePages(count int64, pageSize int, currentPage int) (int, error) {
	if count <= 0 {
		return 1, nil
	}
	if pageSize <= 0 {
		return 1, errors.New("无效的页码大小")
	}
	sizeInt64 := int64(pageSize)
	pages := int((count + sizeInt64 - 1) / sizeInt64)

	if currentPage > pages {
		return 1, errors.New("请求页码超出范围")
	}
	return pages, nil
}

// BuildLikeQuery 构建用于SQL LIKE查询的模式字符串
// 参数:
//   - keyword: 要查询的关键词，如为空则返回"%"（匹配所有）
//   - matchType: 可选的匹配类型，支持以下值:
//   - "start": 前缀匹配，返回"keyword%"
//   - "end": 后缀匹配，返回"%keyword"
//   - "exact": 精确匹配，返回"keyword"
//   - 其他或不提供: 包含匹配(默认)，返回"%keyword%"
//
// 返回值:
//
//	格式化后的LIKE查询模式字符串
//
// 示例:
//
//	BuildLikeQuery("test")         // 返回 "%test%"
//	BuildLikeQuery("test", "start") // 返回 "test%"
//	BuildLikeQuery("", "exact")    // 返回 "%"
func BuildLikeQuery(keyword string, matchType ...string) string {
	if keyword == "" {
		return "%"
	}

	var match string
	if len(matchType) > 0 {
		match = matchType[0]
	}

	switch match {
	case "start":
		return keyword + "%"
	case "end":
		return "%" + keyword
	case "exact":
		return keyword
	default: // "contains"
		return "%" + keyword + "%"
	}
}
