package dbkit

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

func ComputePages(page int, size int, count int64) (int, error) {
	if count <= 0 {
		return 1, nil
	}
	if size <= 0 {
		return 1, errors.New("无效的页码大小")
	}
	sizeInt64 := int64(size)
	pages := int((count + sizeInt64 - 1) / sizeInt64)

	if page > pages {
		return 1, errors.New("请求页码超出范围")
	}
	return pages, nil
}
