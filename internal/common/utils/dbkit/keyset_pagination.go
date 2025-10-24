package dbkit

import (
	"fmt"
	"slices"

	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

type keyset[T any] struct {
	BeforeID int64
	AfterID  int64

	PageSize int
}

func NewKeyset[T any](pageSize int, beforeID, afterID int64) *keyset[T] {
	if pageSize <= 0 {
		pageSize = 5
	}
	return &keyset[T]{
		PageSize: pageSize,
		BeforeID: beforeID,
		AfterID:  afterID,
	}
}

// 把 Keyset 参数应用到查询 Mod 上
func (k *keyset[T]) ApplyKeysetMods(base []qm.QueryMod, IDCol string) []qm.QueryMod {
	limit := k.PageSize
	if limit <= 0 {
		limit = 5
	}

	// 游标与排序
	if k.AfterID > 0 {
		base = append(base, qm.Where(fmt.Sprintf("%s > ?", IDCol), k.AfterID))
		base = append(base, qm.OrderBy(IDCol+" ASC"))
	} else if k.BeforeID > 0 {
		base = append(base, qm.Where(fmt.Sprintf("%s < ?", IDCol), k.BeforeID))
		// 为了取到"上一页"的正确数据，先按 DESC 取最新的 N 条，然后反转切片返回给调用方（保持 ASC 展示）
		base = append(base, qm.OrderBy(IDCol+" DESC"))
	} else {
		base = append(base, qm.OrderBy(IDCol+" ASC"))
	}

	// limit
	fetchLimit := limit + 1
	base = append(base, qm.Limit(fetchLimit))

	return base
}

type paginationResult[T any] struct {
	Items   []*T
	HasNext bool
	HasPrev bool
}

// BuildPaginationResult 从查询结果构建 PaginationResult
func (k *keyset[T]) BuildPaginationResult(domainSlice []*T) *paginationResult[T] {
	limit := k.PageSize
	if limit <= 0 {
		limit = 5
	}

	hasMore := len(domainSlice) > limit

	// 截取
	if hasMore {
		domainSlice = domainSlice[:k.PageSize]
	}

	// 如果是 Before 分页，先按 DESC 取，结果需要反转为 ASC 展示
	if k.BeforeID > 0 && len(domainSlice) > 0 {
		slices.Reverse(domainSlice)
	}

	// 基于请求方向和是否多取一条来推断 hasPrev/hasNext，避免额外 DB 查询
	isAfter := k.AfterID > 0
	isBefore := k.BeforeID > 0

	var hasPrev, hasNext bool
	switch {
	case isAfter:
		// 请求为 after（向后翻页），代表存在上一页（客户端传了游标）
		hasPrev = true
		hasNext = hasMore
	case isBefore:
		// 请求为 before（向前翻页），代表存在下一页（客户端传了游标）
		hasPrev = hasMore
		hasNext = true
	default:
		// 首页
		hasPrev = false
		hasNext = hasMore
	}

	return &paginationResult[T]{
		Items:   domainSlice,
		HasPrev: hasPrev,
		HasNext: hasNext,
	}
}
